package userdto

type LoginByGoogleReq struct {
	IDToken   string `json:"id_token" validate:"required"`
	UserAgent string
	IPAddress string
}

type VerifyTokenReq struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type VerifyTokenRes struct {
	Sub       string `json:"sub"`
	SessionID uint64 `json:"sid"`
	JTI       string `json:"jti"`
	ExpiresAt int64  `json:"exp"`
}

type RegisterReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name,omitempty"`
	UserAgent string `json:"-"`
	IPAddress string `json:"-"`
}

type LogoutReq struct {
	AccessToken string
}

type RefreshTokenReq struct {
	RefreshToken string
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
