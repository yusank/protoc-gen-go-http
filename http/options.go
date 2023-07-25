package http

import (
	restyv2 "github.com/go-resty/resty/v2"
)

type CallOption interface {
	Before(req *restyv2.Request) error
	After(rsp *restyv2.Response) error
}

type baseCallOption struct {
	before func(req *restyv2.Request) error
	after  func(rsp *restyv2.Response) error
}

func (b baseCallOption) Before(req *restyv2.Request) error {
	if b.before != nil {
		return b.before(req)
	}
	return nil
}

func (b baseCallOption) After(rsp *restyv2.Response) error {
	if b.after != nil {
		return b.after(rsp)
	}
	return nil
}

// Before and After is a shortcut for baseCallOption

// Before is a shortcut for baseCallOption with only before func
func Before(f func(req *restyv2.Request) error) CallOption {
	return baseCallOption{before: f}
}

// After is a shortcut for baseCallOption with only after func
func After(f func(rsp *restyv2.Response) error) CallOption {
	return baseCallOption{after: f}
}

type ClientOption interface {
	Apply(c *restyv2.Client) error
}

type baseClientOption struct {
	apply func(c *restyv2.Client) error
}

func (b baseClientOption) Apply(c *restyv2.Client) error {
	if b.apply != nil {
		return b.apply(c)
	}
	return nil
}

// ApplyToClient is a shortcut for baseClientOption
func ApplyToClient(f func(c *restyv2.Client) error) ClientOption {
	return baseClientOption{apply: f}
}
