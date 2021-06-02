package tailer

import (
	"os"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestIsNewLine(t *testing.T) {
	pwd, _ := os.Getwd()

	handler := NewFileTailer(pwd + "/hello.log", "^\\d{4}\\/\\d{2}\\/\\d{2} ")

	line := []byte("2021-06-02 16:07:00,329 [http-nio-28083-exec-2] ERROR [com.legendshop.base.log.PaymentLog] PaymentLog.java:34 - pay notify , error {} ")
	newLine := handler.isNewLine(line)
	assert.False(t, newLine)

	line = []byte("2021/06/02 sfsdfsdf")
	newLine = handler.isNewLine(line)
	assert.True(t, newLine)
}

func TestProdExp(t *testing.T) {
	handler := NewFileTailer("/hello.log", "^\\d{4}\\-\\d{2}\\-\\d{2} \\d{2}\\:\\d{2}\\:\\d{2},\\d{3} ")

	line := []byte("Param [ e314e51f-0d9f-4d06-8ca1-0392b7a48055,H5,user], result list size is 1, cost time 14")
	newLine := handler.isNewLine(line)
	assert.False(t, newLine)
}