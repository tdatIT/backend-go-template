package userdto

type LoginByGoogleReq struct {
	IDToken string `json:"id_token" validate:"required"`
}

type VerifyTokenReq struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type VerifyTokenRes struct {
	SubID     uint64 `json:"sub_id"`
	JTI       string `json:"jti"`
	ExpiresAt int64  `json:"exp"`
}

type RegisterReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name,omitempty"`
}

type LogoutReq struct {
	UserID       uint64 `json:"user_id" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
