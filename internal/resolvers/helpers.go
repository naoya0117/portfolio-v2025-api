package resolvers

import (
	"fmt"
	"regexp"
	"time"
)

// Helper functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func detectURLInText(text string) string {
	urlRegex := `https?://[^\s]+`
	re := regexp.MustCompile(urlRegex)
	match := re.FindString(text)
	return match
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

