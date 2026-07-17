package contact

import "time"

type Tonal string

const (
	TonalNeutral    Tonal = "neutral"
	TonalNegative   Tonal = "negative"
	TonalPositive   Tonal = "positive"
	TonalUnexpected Tonal = "unexpected"
)

type Contact struct {
	ID         int
	Name       string
	Phone      string
	Email      string
	Comment    string
	Tonal      Tonal
	DateCreate time.Time
}
