package context

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"statefulset/cmds/server/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	maxBodyBytes = 1024 * 1024 * 2
)

var defaultXMLBodyHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>`)
var (
	bytesBufferPool = utils.NewBufferPool(256)
	contextPool     = sync.Pool{New: func() interface{} { return new(Context) }}
	emptyContext    = Context{}
)

// empty response
type EmptyResponse struct{}

// Context is the http context for every request
type Context struct {
	ctx            *gin.Context
	Logger         *logrus.Logger
	RequestIDInt64 int64
	RequestIDStr   string
	Data           map[string]interface{}
	StatusCode     int
}

func NewContext(ctx *gin.Context) *Context {
	c := contextPool.Get().(*Context)
	c.ctx = ctx
	return c
}
func (c *Context) Free() {
	// Reset
	*c = emptyContext
	contextPool.Put(c)
}

func (c *Context) Discard() {
	select {
	case <-c.ctx.Request.Context().Done():
		c.Logger.Debugf("Client has be broken")
	default:
		if c.ctx.Request.Body != nil {
			w, err := io.Copy(ioutil.Discard, c.ctx.Request.Body)
			if err == io.EOF {
				return
			}
			if err != nil {
				c.Logger.Warnf("Discard unread body fail: %v", err)
				return
			}
			if w > 0 {
				c.Logger.Debugf("Discard unread body, bytes=%d", w)
			}
		}
	}
}

//Set data
func (c *Context) Set(key string, v interface{}) {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data[key] = v
}

//Get data
func (c *Context) Get(key string) interface{} {
	return c.Data[key]
}

// GetQueryParam get the url query param from URL.Query
func (c *Context) GetQueryParam(key string) (value string, exist bool) {
	var temp []string
	temp, exist = c.ctx.GetQueryArray(key)
	if exist {
		value = temp[0]
	}
	return
}

// GetQueryParam get the url query param from URL.Query
func (c *Context) GetPageParam() (int, int) {
	limit := 10
	if lim, exist := c.GetQueryParam("limit"); exist {
		li, err := strconv.Atoi(lim)
		if err == nil {
			limit = li
		}
	}
	offset := 0
	if off, exist := c.GetQueryParam("offset"); exist {
		of, err := strconv.Atoi(off)
		if err == nil {
			offset = of
		}
	}
	return limit, offset
}

// GetTimeParam get the url query param from URL.Query
func (c *Context) GetTimeParam() (time.Time, time.Time) {
	logger := c.Logger
	var startAt time.Time
	var endAt time.Time
	if start_at, exist := c.GetQueryParam("start_at"); exist {
		st, err := strconv.Atoi(start_at)
		if err != nil {
			logger.Debugf(" api: Get alerts: string to int err:%v", err)
		}
		if err == nil {
			startAt = time.Unix(int64(st), 0)
		}
	}
	if end_at, exist := c.GetQueryParam("end_at"); exist {
		et, err := strconv.Atoi(end_at)
		if err != nil {
			logger.Debugf(" api: Get alerts: string to int err:%v", err)
		}
		if err == nil {
			endAt = time.Unix(int64(et), 0)
		}
	}
	return startAt, endAt
}

// GetRawQueryParam get the specified query param value by specified key
func (c *Context) GetRawQueryParam(key string) (value string) {
	value, _ = c.GetQueryParam(key)
	return
}

// GetQueryParamDefault is wrapper for GetQueryParam; The defaultValue will be return
// if the value of key is "";
func (c *Context) GetQueryParamDefault(key string, defaultValue string) (value string) {
	value, _ = c.GetQueryParam(key)
	if value == "" {
		value = defaultValue
	}
	return
}

// GetHeader get the specified header value by specified key
func (c *Context) GetHeader(key string) (value string, exist bool) {
	var temp []string
	temp, exist = c.ctx.Request.Header[key]
	if exist {
		value = temp[0]
	}
	return
}

// GetRawHeader get the specified header value by specified key
func (c *Context) GetRawHeader(key string) (value string) {
	return c.ctx.Request.Header.Get(key)
}

// GetRawHeaderDefault is wrapper for GetHeader; The defaultValue will be return
// if the value of key is "";
func (c *Context) GetRawHeaderDefault(key string, defaultValue string) (value string) {
	value = c.ctx.Request.Header.Get(key)
	if value == "" {
		value = defaultValue
	}
	return
}

// ClientIP get the ip address of client
func (c *Context) ClientIP() string {
	var ip string
	ip, _ = c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(strings.TrimSpace(c.ctx.Request.RemoteAddr))
	}
	return ip
}

// ParseJSONBody decode the response body to specified data
// Will no allow unknown fields in request body if strict is true
func (c *Context) ParseJSONBody(v interface{}, strict bool) (err error) {
	lenstr := c.GetRawHeader("Content-Length")
	if lenstr != "" {
		c.Logger.Warn("not contain Content-Length")
		length, err := strconv.Atoi(lenstr)
		if err != nil {
			c.Logger.Warnf("%s cannot convert to int, err: %v", lenstr, err)
			return fmt.Errorf("invalid header Content-Length in request")
		}
		// Ignore nil body
		if length == 0 {
			return nil
		}
		if length > maxBodyBytes {
			return fmt.Errorf("the Content-Length can not greater than %d", maxBodyBytes)
		}
	}

	buf := bytesBufferPool.Get()
	defer bytesBufferPool.Put(buf)

	_, err = buf.ReadFrom(c.ctx.Request.Body)
	if err != nil {
		c.Logger.Warnf("Read request body content fail: %v", err)
		return fmt.Errorf("Read request body content fail: %v", err)
	}

	c.Logger.Debugf("Request Body: <%s>", buf.Bytes())

	decoder := json.NewDecoder(buf)
	if strict {
		decoder.DisallowUnknownFields()
	}
	// Decode json body
	if err = decoder.Decode(v); err != nil {
		c.Logger.Warnf("Parse json body fail: %v", err)
		return fmt.Errorf("Parse json body fail: %v", err)
	}
	return
}

// SetHeader set key and value to the response header
func (c *Context) SetHeader(key, value string) {
	c.ctx.Writer.Header().Set(key, value)
}

// WriteResponseStatusCode send a http status code to client
func (c *Context) WriteResponseStatusCode(StatusCode int) {
	// Cannot duplicate write status code to client
	if c.StatusCode != 0 {
		panic(fmt.Sprintf("duplicate write response status code, %s", c.RequestIDStr))
	}
	c.StatusCode = StatusCode
	// Set response request id
	c.SetHeader("x-qs-request-id", c.RequestIDStr)
	c.ctx.Writer.WriteHeader(StatusCode)
}

// WriteJSONResponse write the data in JSON format to the response body;
func (c *Context) WriteJSONResponse(StatusCode int, data interface{}) (err error) {
	if data == nil {
		panic("invalid data for JSON response")
	}

	var b []byte
	if b, err = json.Marshal(data); err != nil {
		c.Logger.Errorf("Encode data %v to json format fail: %v", data, err)
		return
	}

	c.SetHeader("Content-Length", strconv.Itoa(len(b)))
	c.SetHeader("Content-Type", "application/json")

	c.WriteResponseStatusCode(StatusCode)

	c.Logger.Debugf("Writing JSON body <%s>", b)

	// Can't return error of write body error after write status code to response
	if _, e := c.ctx.Writer.Write(b); e != nil {
		c.Logger.Warnf("Write JSON body to response fail: %v", e)
		return
	}
	return
}

// WriteEmptyResponse write empty data in JSON format to the response body;
func (c *Context) WriteEmptyResponse(StatusCode int) (err error) {
	var b []byte
	if b, err = json.Marshal(EmptyResponse{}); err != nil {
		c.Logger.Errorf("Encode empty struct data to json format fail: %v", err)
		return
	}

	c.SetHeader("Content-Length", strconv.Itoa(len(b)))
	c.SetHeader("Content-Type", "application/json")

	c.WriteResponseStatusCode(StatusCode)

	c.Logger.Debugf("Writing JSON body empty")

	// Can't return error of write body error after write status code to response
	if _, e := c.ctx.Writer.Write(b); e != nil {
		c.Logger.Warnf("Write JSON body to response fail: %v", e)
		return
	}
	return
}
