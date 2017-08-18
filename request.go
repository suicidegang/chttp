package main

import (
	"net/http"
	"time"
)

type RequestOpt func(*http.Client, *http.Request) (*http.Request, error)

func GET(url string) RequestOpt {
	return func(client *http.Client, req *http.Request) (r *http.Request, err error) {
		r, err = http.NewRequest("GET", url, nil)
		return
	}
}

func Timeout(seconds int) RequestOpt {
	return func(client *http.Client, req *http.Request) (*http.Request, error) {
		client.Timeout = time.Second * time.Duration(seconds)
		return req, nil
	}
}

func RequestOrPanic(opts ...RequestOpt) chan Response {
	r, err := Request(opts...)
	if err != nil {
		panic(err)
	}

	return r
}

func Request(opts ...RequestOpt) (chan Response, error) {
	var (
		req *http.Request
		err error
	)

	client := &http.Client{}
	for _, opt := range opts {
		req, err = opt(client, req)
		if err != nil {
			return nil, err
		}
	}

	ch := make(chan Response)
	go func() {
		res, err := client.Do(req)

		// Having done the request send it to the response channel.
		ch <- Response{res, err}
	}()

	return ch, nil
}
