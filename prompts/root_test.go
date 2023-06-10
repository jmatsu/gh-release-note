package prompts

import (
	"fmt"
	"testing"
)

func Test_embeds(t *testing.T) {
	prompts := []string{
		SimpleTxt,
	}

	for i, prompt := range prompts {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			if prompt == "" {
				t.Errorf("prompt is empty")
			}
		})
	}
}
