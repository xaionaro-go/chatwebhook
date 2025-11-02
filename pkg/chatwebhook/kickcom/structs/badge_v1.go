package structs

type BadgeV1 struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Count *int64 `json:"count,omitempty"`
}
