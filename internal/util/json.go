package util

import "strings"

func FlattenJSONString(json string) string {
	return strings.ReplaceAll(json, "\n", "")
}
