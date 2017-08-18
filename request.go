package main

import (
	"log"
	"net/http"
	"time"
)

type Req struct {
	Client *http.Client
	Raw    *http.Request
}

func (r Req) Response() chan Response {
	log.Println("REQ STARTED.")
	response := make(chan Response)

	go func() {
		client := r.Client
		res, err := client.Do(r.Raw)
		response <- Response{res, err}
	}()

	return response
}

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
