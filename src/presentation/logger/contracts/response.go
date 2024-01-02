package contracts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strconv"
)

var regexUUID = regexp.MustCompile(`\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

type LoggerResponse struct {
	Fields logrus.Fields
}

func (ref LoggerResponse) FromDomain(c *gin.Context, timestamp string, latency int64, requestBody, responseBody interface{}) LoggerResponse {
	path, _ := url.PathUnescape(c.Request.URL.Path)
	ref.Fields = logrus.Fields{}
	ref.Fields["client_ip"] = c.ClientIP()
	ref.Fields["timestamp"] = timestamp
	ref.Fields["method"] = c.Request.Method
	ref.Fields["params"] = c.Request.URL.Query()
	ref.Fields["url"] = path
	ref.Fields["headers"] = c.Request.Header
	ref.Fields["path"] = regexUUID.ReplaceAllString(path, "entity_id")
	ref.Fields["status_code"] = strconv.Itoa(c.Writer.Status())
	ref.Fields["latency"] = latency
	ref.Fields["user_agent"] = c.Request.UserAgent()
	if c.Errors.Last() != nil {
		ref.Fields["error_message"] = c.Errors.Last().Err.Error()
	}
	ref.Fields["request_body"] = requestBody
	ref.Fields["response_body"] = responseBody
	return ref
}

func (ref LoggerResponse) ToString() string {
	out, _ := json.Marshal(ref.Fields)
	return string(out)
}
