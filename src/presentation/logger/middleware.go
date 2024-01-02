package logger

import (
	"bitbucket.org/bexstech/temis-compliance/src/presentation/logger/contracts"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"time"
)

func Middleware(c *gin.Context) {

	response := &contracts.LogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = response

	requestBuffer, _ := ioutil.ReadAll(c.Request.Body)
	bufferBody := ioutil.NopCloser(bytes.NewBuffer(requestBuffer))
	bufferBodyCopy := ioutil.NopCloser(bytes.NewBuffer(requestBuffer))

	c.Request.Body = bufferBodyCopy

	start := time.Now()
	c.Next()
	end := time.Now()

	var responseBody interface{}
	var requestBody interface{}

	json.Unmarshal(response.Body.Bytes(), &responseBody)
	json.Unmarshal(readBody(bufferBody), &requestBody)

	contentLog := contracts.LoggerResponse{}.FromDomain(c, start.String(), end.Sub(start).Milliseconds(), requestBody, responseBody)
	if c.Writer.Status() >= 500 {
		logrus.WithFields(contentLog.Fields).Error()
	} else {
		logrus.WithFields(contentLog.Fields).Info()
	}

}

func readBody(reader io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.Bytes()
	return s
}
