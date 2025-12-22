package feeds

import (
	"crypto/sha256"
	"encoding/hex"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var strongRe = regexp.MustCompile(`(?is)<strong>\s*([^<]+?)\s*</strong>`)

func EventKeyFromLink (link string) (string, bool) {
	// Parse URL
	parsed, err := url.Parse(link)
	if err != nil {
		return "", false
	}
	// take last path segment
	segments := strings.Split(parsed.Path, "/")
	for i := len(segments) -1; i >= 0; i-- {
		if segments[i] != "" {
			return segments[i], true
		}
	}
	// return it if non-empty
	return "", false
}

func LatestStatusFromDescription (desc string) string {
	// html.UnescapeString(desc)
	s := html.UnescapeString(desc)
	m := strongRe.FindStringSubmatch(s)
	if len(m) < 2 {
		return "unknown"
	}
	// find first <strong>...</strong>
	// map:
	kw := strings.ToLower(strings.TrimSpace(m[1]))

	switch kw {
	case "resolved":
		return "resolved"
	case "investigating", "identified", "monitoring", "update":
		return "open"
	default:
		return "unknown"
	}
	//    Resolved -> resolved
	//	  Investigating|Identified|Monitoring|Update -> open
	// else -> unknown
}

func ContentHash (desc string) string {
	// hash the unescaped description (sha256 hex)
	s := html.UnescapeString(desc)
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}