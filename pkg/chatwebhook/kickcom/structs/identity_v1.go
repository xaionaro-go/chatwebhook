package structs

type IdentityV1 struct {
	UsernameColor string    `json:"username_color"`
	Badges        []BadgeV1 `json:"badges"`
}
