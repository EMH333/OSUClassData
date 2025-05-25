package util

import "strings"

func ClassToSubject(class string) string {
	return strings.TrimSuffix(class, "H")[:strings.IndexAny(class, "0123456789")]
}
