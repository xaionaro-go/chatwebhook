package structs

type BanMetadataV1 struct {
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
