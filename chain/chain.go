package chain

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fr3fou/polo/polo"
)

type Entries struct {
	Commits []string `json:"commits"`
}

func New(order int, r io.Reader) (polo.Chain, error) {
	entries := Entries{}
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		return polo.Chain{}, fmt.Errorf("failed decoding entries, %w", err)
	}
	return polo.NewFromText(order, entries.Commits), nil
}
