package user_logic

import (
	"context"
	"errors"
	"time"

	jwt "github.com/MagicPig9898/familychefassistant_server/config/jwt_config"
	log "github.com/MagicPig9898/familychefassistant_server/config/log_config"
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
	openid, err := code2Session(ctx, userLoginDto.Code)
	if err != nil {
		log.Errorf("code2Session failed: %v", err)
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
		if err != nil {
			log.Errorf("%s InsertUser failed: %v", openid, err)
			return nil, err // 登录失败
		}
	} else {
		err = l.repo.UpdateUser(ctx, &user_entity.TbUser{
			ID:            openid,
			NickName:      userLoginDto.NickName,
			AvatarUrl:     userLoginDto.AvatarUrl,
			City:          userLoginDto.City,
			Province:      userLoginDto.Province,
			Country:       userLoginDto.Country,
			Gender:        userLoginDto.Gender,
			LastLoginTime: time.Now().Unix(),
		})
		if err != nil {
			// 更新失败，打印警告即可
			log.Warnf("%s UpdateUser failed: %v", openid, err)
		}
	}
	token, err := jwt.GenerateToken(openid, 24*time.Hour)
	if err != nil {
		log.Errorf("%s GenerateToken failed: %v", openid, err)
		return nil, err
	}
	userLoginDto.Token = token
	log.Infof("%s Login success: %s", openid, token)
	return userLoginDto, nil
}

func (l *userLogicImpl) ValidToken(ctx context.Context, token string) (string, error) {
	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Errorf("ParseToken failed: %v", err)
		return "", err
	}
	// 判断 token 是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		log.Errorf("%s Token expired: %v", claims.OpenID, err)
		return "", errors.New("token expired")
	}
	// 如果 token 有效，创建一个新的 token
	newtoken, err := jwt.GenerateToken(claims.OpenID, 24*time.Hour)
	if err != nil {
		log.Errorf("%s GenerateToken failed: %v", claims.OpenID, err)
		return "", err
	}
	log.Infof("%s Token valid: %s", claims.OpenID, newtoken)
	return newtoken, nil

}
