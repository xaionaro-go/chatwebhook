package structs

type GiftV1 struct {
	Amount  int64  `json:"amount"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Tier    string `json:"tier"`
	Message string `json:"message"`
}
