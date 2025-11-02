package structs

type EmoteV1 struct {
	EmoteID   string       `json:"emote_id"`
	Positions []PositionV1 `json:"positions"`
}
