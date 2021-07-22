package rule

type EqMatcher struct {
}

func (self EqMatcher) IsMatch(text string, values []string) bool {
    for _, val := range values {
        if val == text {
            return true
        }
    }
    return false
}
