package rule

import (
    "bytes"
    "github.com/benhoyt/goawk/interp"
    "github.com/benhoyt/goawk/parser"
    "strings"
)

type TextCatcher struct {
    Seperator  string
    Expression string
    prog       *parser.Program
    inited     bool
    config     *interp.Config
}

func (self *TextCatcher) Init() (err error) {
    self.prog, err = parser.ParseProgram([]byte(self.Expression), nil)
    if err != nil {
        return err
    }
    self.config = &interp.Config{
        Stdin:  nil,
        Output: nil,
        Vars:   []string{"FS", self.Seperator},
    }
    self.inited = true
    return err
}

func (self *TextCatcher) Capture(text string) (ret string, err error) {
    if !self.inited {
        err = self.Init()
        if err != nil {
            return
        }
    }
    buffer := new(bytes.Buffer)
    config := self.config
    config.Stdin = bytes.NewReader([]byte(text))
    config.Output = buffer

    _, err = interp.ExecProgram(self.prog, config)
    if err != nil {
        return
    }
    ret = strings.TrimSpace(buffer.String())
    return
}
