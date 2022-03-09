package git

import (
	"strings"
)

func ExplicitRefOrHead(refName string) string {
	if refName == "" {
		return "HEAD"
	}
	branchRefPrefix := "refs/heads/"
	if strings.HasPrefix(refName, branchRefPrefix) {
		return refName
	}
	return branchRefPrefix + refName
}
