package chttp

import (
	"net/http"
)

type Req interface {
	Client() *http.Client
	RawRequest() *http.Request
	Response() chan Response
	Request() Req
}

type SimpleReq struct {
	C   *http.Client
	Raw *http.Request
}

func (r SimpleReq) Client() *http.Client {
	return r.C
}

func (r SimpleReq) RawRequest() *http.Request {
	return r.Raw
}

func (r SimpleReq) Request() Req {
	return r
}

func (r SimpleReq) Response() chan Response {
	response := make(chan Response)

	go func() {
		res, err := r.C.Do(r.Raw)
		response <- Response{res, err, r.Request()}
	}()

	return response
}

func RequestOrPanic(opts ...RequestOpt) Req {
	r, err := Request(opts...)
	if err != nil {
		panic(err)
	}

	return r
}

func Request(opts ...RequestOpt) (Req, error) {
	var (
		raw *http.Request
		err error
	)

	client := &http.Client{}
	for _, opt := range opts {
		raw, err = opt(client, raw)
		if err != nil {
			return SimpleReq{}, err
		}
	}

	return SimpleReq{client, raw}, nil
}
