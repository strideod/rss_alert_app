package feeds

import "strings"

func DeriveSeverity(title, desc string) string {
	s := strings.ToLower(title + " " + desc)

	// Order matters: match strongest first
	switch {
	case strings.Contains(s, "major outage") || strings.Contains(s, "outage"):
		return "outage"
	case strings.Contains(s, "degraded") || strings.Contains(s, "partial outage") ||
		strings.Contains(s, "elevated error") || strings.Contains(s, "elevated errors"):
		return "degraded"
	case strings.Contains(s, "maintenance") || strings.Contains(s, "scheduled maintenance"):
		return "maintenance"
	default:
		return "info"
	}
}