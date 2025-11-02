package structs

type RepliesToV1 struct {
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
	Sender    UserV1 `json:"sender"`
}
