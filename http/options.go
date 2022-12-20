package http

import (
	"net/http"
	"strings"
	"time"

	restyv2 "github.com/go-resty/resty/v2"
)

type CallOption interface {
	Before(req *restyv2.Request) error
	After(rsp *restyv2.Response) error
}

func Header(h http.Header) CallOption {
	return HeaderCallOption{h: h}
}

type HeaderCallOption struct {
	onReq bool
	h     http.Header
}

func (h HeaderCallOption) Before(req *restyv2.Request) error {
	if h.onReq {
		req.SetHeaderMultiValues(h.h)
	}

	return nil
}

func (h HeaderCallOption) After(rsp *restyv2.Response) error {
	if !h.onReq {
		for s, ss := range h.h {
			rsp.Header().Add(s, strings.Join(ss, ","))
		}
	}
	return nil
}

type ClientOption interface {
	Apply(c *restyv2.Client) error
}

func Retry(times int, wait time.Duration) ClientOption {
	return RetryClientOption{
		times:        times,
		waitDuration: wait,
	}
}

type RetryClientOption struct {
	times        int
	waitDuration time.Duration
}

func (r RetryClientOption) Apply(c *restyv2.Client) error {

	c.SetRetryCount(r.times).SetRetryWaitTime(r.waitDuration)
	return nil
}
