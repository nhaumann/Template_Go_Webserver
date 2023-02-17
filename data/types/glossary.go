package glossary

import "time"

type GlossaryItem struct {
	ID         int       `json:"id"`
	Term       string    `json:"term"`
	Definition string    `json:"definition"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
