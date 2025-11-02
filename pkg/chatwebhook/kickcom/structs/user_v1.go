package structs

type UserV1 struct {
	IsAnonymous    bool        `json:"is_anonymous"`
	UserID         int64       `json:"user_id"`
	Username       string      `json:"username"`
	IsVerified     bool        `json:"is_verified"`
	ProfilePicture string      `json:"profile_picture"`
	ChannelSlug    string      `json:"channel_slug"`
	Identity       *IdentityV1 `json:"identity"`
}
