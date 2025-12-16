package log

import (
	"testing"
)

func TestLoggerFunctions(t *testing.T) {
	Init()

	t.Run("Info", func(t *testing.T) {
		LogInfo("This is an info message for testing.")
	})

	t.Run("Warning", func(t *testing.T) {
		LogWarning("This is a warning message for testing.")
	})

	t.Run("Error", func(t *testing.T) {
		LogError("This is an error message for testing.")
	})
}
