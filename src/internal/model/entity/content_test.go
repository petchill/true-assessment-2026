package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContentGetRecencyFactory(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name      string
		createdAt time.Time
	}{
		{
			name:      "new content should have score close to 1",
			createdAt: now,
		},
		{
			name:      "content from 30 days ago should decay slightly",
			createdAt: now.AddDate(0, 0, -30),
		},
		{
			name:      "content from one year ago should decay to 0.5",
			createdAt: now.AddDate(-1, 0, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := Content{CreatedAt: tc.createdAt}

			actual := content.GetRecencyFactory()
			expected := 1 / (1 + (time.Since(tc.createdAt).Hours()/24)/365)

			assert.InDelta(t, expected, actual, 0.0001)
		})
	}
}
