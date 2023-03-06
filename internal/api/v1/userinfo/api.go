package userinfo

type UserInfoDto struct {
	Audiences     []string `json:"audiences"`
	Subject       string   `json:"subject"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Groups        []string `json:"groups"`
}
