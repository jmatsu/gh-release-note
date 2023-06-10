package queries

import (
	"fmt"
	"testing"
)

func Test_embeds(t *testing.T) {
	queries := []string{
		ListTagGraphQL,
	}

	for i, query := range queries {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			if query == "" {
				t.Errorf("query is empty")
			}
		})
	}
}
