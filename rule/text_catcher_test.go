package rule

import (
    "github.com/stretchr/testify/assert"
    "log"
    "testing"
)

var str = `10.202.0.40 - - [02/Jul/2021:17:58:41 +0800] "POST /api/p/shopCart HTTP/1.1" 200 106 dm-h5.amusgame.net "https://dm-h5.amusgame.net/cart" "-" "210.0.159.194, 47.88.158.237, 10.202.2.21" "x-platform:-" "x-device:-" "x-version:-" "hid:- n:- c:- bimcode:-" "content-type:application/x-www-form-urlencoded data:-" "user_agent: Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1" 0.011`

func TestCapture(t *testing.T) {
    capture := &TextCatcher{
        Expression: "$NF>0.01{print $9}",
        Seperator: " ",
    }

    result, err := capture.Capture(str)
    if err != nil {
        log.Fatalln(err)
    }

    assert.Equal(t, "200", result)
}
