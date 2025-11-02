package structs

import (
	"time"
)

type RFC3339String string

func (t RFC3339String) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}
