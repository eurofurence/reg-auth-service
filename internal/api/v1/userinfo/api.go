package userinfo

type UserInfoDto struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
