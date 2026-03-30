package user_logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
	wxconfig "github.com/MagicPig9898/familychefassistant_server/config/wx_config"
	userrepo "github.com/MagicPig9898/familychefassistant_server/repo/user_repo"
)

type userLogicImpl struct {
	repo userrepo.UserRepo
}

func newUserLogicImpl() *userLogicImpl {
	return &userLogicImpl{repo: userrepo.NewUserRepo()}
}

func (l *userLogicImpl) GetUserInfo(ctx context.Context, id int64) (string, error) {
	return l.repo.GetUserByID(ctx, id)
}

func (l *userLogicImpl) WXLogin(ctx context.Context, userLoginDto *user_entity.UserLoginDto) (*user_entity.UserLoginDto, error) {
	// 1. 用 code 请求微信 jscode2session 接口换取 openid
	openid, err := l.code2Session(ctx, userLoginDto.Code)
	if err != nil {
		return nil, err
	}

	// 2. 把 openid 作为 token 回传（后续可替换为 JWT 等）
	userLoginDto.Token = openid
	return userLoginDto, nil
}

// code2Session 调用微信 jscode2session 接口，用 code 换 openid
func (l *userLogicImpl) code2Session(ctx context.Context, code string) (string, error) {
	params := url.Values{}
	params.Set("appid", wxconfig.AppID)
	params.Set("secret", wxconfig.AppSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := fmt.Sprintf("%s?%s", wxconfig.Code2SessionURL, params.Encode())

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
