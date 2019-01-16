package cathand

import "strings"

func NormalizeLineFeed(raw []byte) string {
	return strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(string(raw))
}

func NormalizeLineFeedBytes(raw []byte) []byte {
	return []byte(NormalizeLineFeed(raw))
}

