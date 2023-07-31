package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	urls := []string{"https://google.com", "https://facebook.com", "https://twitter.com"}

	t.Run("it should return 'true' because google is contains from array", func(t *testing.T) {
		c := Contains(urls, "https://google.com")
		assert.True(t, c)
	})

	t.Run("it should return 'false' because linkedin is not contains from array", func(t *testing.T) {
		c := Contains(urls, "https://linkedin.com")
		assert.False(t, c)
	})
}
