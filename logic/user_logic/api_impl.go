package user_logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	conf "github.com/MagicPig9898/familychefassistant_server/conf"
	jwt "github.com/MagicPig9898/familychefassistant_server/config/jwt_config"
	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
	userrepo "github.com/MagicPig9898/familychefassistant_server/repo/user_repo"
)

type userLogicImpl struct {
	repo userrepo.UserRepo
}

func newUserLogicImpl() *userLogicImpl {
	return &userLogicImpl{repo: userrepo.NewUserRepo()}
}

func (l *userLogicImpl) GetUserInfo(ctx context.Context, id string) (*user_entity.TbUser, error) {
	return l.repo.GetUserByID(ctx, id)
}

func (l *userLogicImpl) WXLogin(ctx context.Context, userLoginDto *user_entity.UserLoginDto) (*user_entity.UserLoginDto, error) {
	// 1. 用 code 请求微信 jscode2session 接口换取 openid
	openid, err := l.code2Session(ctx, userLoginDto.Code)
	if err != nil {
		return nil, err
	}
	token, err := jwt.GenerateToken(openid, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	_, err = l.repo.GetUserByID(ctx, openid)
	if err != nil {
		err = l.repo.InsertUser(ctx, &user_entity.TbUser{
			ID:             openid,
			NickName:       userLoginDto.NickName,
			AvatarUrl:      userLoginDto.AvatarUrl,
			City:           userLoginDto.City,
			Country:        userLoginDto.Country,
			Gender:         userLoginDto.Gender,
			FristLoginTime: time.Now().Unix(),
		})
		return nil, err
	} else {
		err = l.repo.UpdateUser(ctx, &user_entity.TbUser{
			ID:        openid,
			NickName:  userLoginDto.NickName,
			AvatarUrl: userLoginDto.AvatarUrl,
			City:      userLoginDto.City,
			Country:   userLoginDto.Country,
			Gender:    userLoginDto.Gender,
		})
		if err != nil {
			fmt.Println("update user err:", err)
			return nil, err
		}
	}
	userLoginDto.Token = token
	return userLoginDto, nil
}

func (l *userLogicImpl) ValidToken(ctx context.Context, token string) (string, error) {
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return "", err
	}
	// 判断 token 是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return "", errors.New("token expired")
	}
	// 如果 token 有效，创建一个新的 token
	newtoken, err := jwt.GenerateToken(claims.OpenID, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return newtoken, nil

}

// code2Session 调用微信 jscode2session 接口，用 code 换 openid
// 在 微信 体系中：
// 每个用户访问一个小程序时
// 微信会给这个用户生成一个 唯一的 openid
// 这个 ID 只在当前小程序内唯一
func (l *userLogicImpl) code2Session(ctx context.Context, code string) (string, error) {
	params := url.Values{}
	params.Set("appid", conf.Cfg.WX.AppID)
	params.Set("secret", conf.Cfg.WX.AppSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := fmt.Sprintf("%s?%s", conf.Cfg.WX.Code2SessionURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求微信接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var wxResp user_entity.WXLoginResp
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if wxResp.ErrCode != 0 {
		return "", fmt.Errorf("微信登录失败: errcode=%d, errmsg=%s", wxResp.ErrCode, wxResp.ErrMsg)
	}

	return wxResp.OpenID, nil
}
