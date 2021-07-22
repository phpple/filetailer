package rule

type AnyMatcher struct {
}

func (self AnyMatcher) IsMatch(text string, values []string) bool {
    return true
}
