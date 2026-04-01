package user_entity

// UserLoginDto 用户登录Dto
type UserLoginDto struct {
	Code      string `json:"code"`
	AvatarUrl string `json:"avatar_url"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	NickName  string `json:"nick_name"`
	Token     string `json:"token"`
}

// Tb_user 用户表
type TbUser struct {
	ID             string `db:"id"` // 就是openid
	NickName       string `db:"nick_name"`
	AvatarUrl      string `db:"avatar_url"`
	City           string `db:"city"`
	Country        string `db:"country"`
	Gender         int    `db:"gender"`
	FristLoginTime int64  `db:"frist_login_time"`
}

// ValidTokenDto token 校验请求
type ValidTokenDto struct {
	Token string `json:"token"`
}

// WXLoginResp 微信 jscode2session 接口返回结构
type WXLoginResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
