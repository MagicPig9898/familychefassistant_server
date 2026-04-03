package user_logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	conf "github.com/MagicPig9898/familychefassistant_server/conf"
	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
)

type UserLogic interface {
	GetUserInfo(ctx context.Context, id string) (*user_entity.TbUser, error)

	WXLogin(ctx context.Context, userLoginDto *user_entity.UserLoginDto) (*user_entity.UserLoginDto, error)

	ValidToken(ctx context.Context, token string) (string, error)
}

func NewUserLogic() UserLogic {
	return newUserLogicImpl()
}

// code2Session 调用微信 jscode2session 接口，用 code 换 openid
// 在 微信 体系中：
// 每个用户访问一个小程序时
// 微信会给这个用户生成一个 唯一的 openid
// 这个 ID 只在当前小程序内唯一
func code2Session(ctx context.Context, code string) (string, error) {
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
