package userdto

type LoginByGoogleReq struct {
	IDToken   string `json:"id_token" validate:"required"`
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

type VerifyTokenReq struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type VerifyTokenRes struct {
	Sub       string `json:"sub"`
	SessionID string `json:"sid"`
	JTI       string `json:"jti"`
	ExpiresAt int64  `json:"exp"`
}

type RegisterReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

type LogoutReq struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
