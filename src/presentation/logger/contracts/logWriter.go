package contracts

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type LogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w LogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
