package rule

import "strconv"

type GtMatcher struct {
}

func (self GtMatcher) IsMatch(text string, values []string) bool {
    var floatText float64
    floatText, _ = strconv.ParseFloat(text, 64)
    for _, val := range values {
        var floatVal float64
        floatVal, _ = strconv.ParseFloat(val, 64)
        if floatText > floatVal {
            return true
        }
    }
    return false
}
