package auth

func JoinURL(prefix, suffix string) string {
	prefixLength := len(prefix)
	suffixLength := len(suffix)

	if prefixLength > 0 && prefix[prefixLength-1] == '/' {
		prefix = prefix[:prefixLength-1]
	}

	if suffixLength > 0 && suffix[0] == '/' {
		suffix = suffix[1:]
	}

	return prefix + "/" + suffix
}
