package chttp

import (
	"net/http"
)

type Req struct {
	Client *http.Client
	Raw    *http.Request
}

func (r Req) Response() chan Response {
	response := make(chan Response)

	go func() {
		client := r.Client
		res, err := client.Do(r.Raw)
		response <- Response{res, err}
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
			return Req{}, err
		}
	}

	return Req{client, raw}, nil
}
