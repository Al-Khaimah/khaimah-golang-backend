package users

type OAuthRequestDTO struct {
	Provider string `json:"provider"`
}

type SSOLoginResponse struct {
	Token      string       `json:"token"`
	UserExists bool         `json:"user_exists"`
	User       *UserBaseDTO `json:"user"`
}
