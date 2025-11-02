package structs

type MetadataV1 struct {
	Title            string     `json:"title"`
	Language         string     `json:"language"`
	HasMatureContent bool       `json:"has_mature_content"`
	Category         CategoryV1 `json:"category"`
}
