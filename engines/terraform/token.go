package terraform

import "regexp"

// extractTokenContents extracts the contents between ${} from a token string
func extractTokenContents(token string) (string, bool) {
	if matches := tokenPattern.FindStringSubmatch(token); len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

var tokenPattern = regexp.MustCompile(`^\${([^}]+)}$`)
