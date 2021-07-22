package rule

import (
    "log"
    "reflect"
)

type RuleCause string

const (
    RuleAny RuleCause = "any"
    RuleEq  RuleCause = "eq"
    RuleGt  RuleCause = "gt"
    RuleLt  RuleCause = "lt"
)

type Rule struct {
    Cause      RuleCause `yaml:"cause"`
    Values     []string  `yaml:"values"`
    Msg        string    `yaml:"msg"`
    Seperator  string    `yaml:"seperator"`
    Expression string    `yaml:"expression"`
    catcher    *TextCatcher
    matcher    RuleMatcher
    inited     bool
}

var matchers map[RuleCause]RuleMatcher

func init() {
    matchers = make(map[RuleCause]RuleMatcher, 4)
    matchers[RuleAny] = AnyMatcher{}
    matchers[RuleEq] = EqMatcher{}
}

type RuleMatcher interface {
    IsMatch(text string, values []string) bool
}

func (self *Rule) Init() {
    if matcher, ok := matchers[self.Cause]; ok {
        self.matcher = matcher
    }

    self.catcher = &TextCatcher{
        Seperator:  self.Seperator,
        Expression: self.Expression,
    }
    self.catcher.Init()
    self.inited = true
}

// 文字是否与指定的规则匹配
func (self *Rule) Match(text string) (result string, match bool) {
    if self.Cause == "any" || self.Cause == "" {
        return result, true
    }
    if !self.inited {
        self.Init()
    }
    var err error
    result, err = self.catcher.Capture(text)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("catcher result:", result)
    if result == "" {
        return
    }

    match = self.matcher.IsMatch(result, self.Values)
    log.Printf("type:%s, match:%s", reflect.TypeOf(self.matcher).Name(), match)
    return
}
