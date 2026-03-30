package user_entity

// UserLoginDto 用户登录Dto
type UserLoginDto struct {
	Code      string `json:"code"`
	AvatarUrl string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	NickName  string `json:"nick_name"`
	Token     string `json:"token"`
}

// Tb_user 用户表
type Tb_user struct {
	ID        string `db:"id"` // 就是openid
	NickName  string `db:"nick_name"`
	AvatarUrl string `db:"avatar_url"`
	City      string `db:"city"`
	Country   string `db:"country"`
	Gender    int    `db:"gender"`
}

// WXLoginResp 微信 jscode2session 接口返回结构
type WXLoginResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
