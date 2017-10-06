package chttp

import (
	"io"
	"net/http"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

type RequestOpt func(*http.Client, *http.Request) (*http.Request, error)

func GET(url string) RequestOpt {
	return func(client *http.Client, req *http.Request) (r *http.Request, err error) {
		r, err = http.NewRequest("GET", url, nil)
		return
	}
}

func POST(url string, reader io.Reader) RequestOpt {
	return func(client *http.Client, req *http.Request) (r *http.Request, err error) {
		r, err = http.NewRequest("POST", url, reader)
		return
	}
}

func Oauth(token, secret string) RequestOpt {
	return func(client *http.Client, req *http.Request) (*http.Request, error) {
		var ocl oauth.Client
		ocl.Credentials = oauth.Credentials{Token: token, Secret: secret}

		// Compute authorization header based on current request configuration.
		header := ocl.AuthorizationHeader(nil, req.Method, req.URL, req.PostForm)
		req.Header.Add("Authorization", header)
		return req, nil
	}
}

func Timeout(seconds int) RequestOpt {
	return func(client *http.Client, req *http.Request) (*http.Request, error) {
		client.Timeout = time.Second * time.Duration(seconds)
		return req, nil
	}
}

func Header(key, value string) RequestOpt {
	return func(client *http.Client, req *http.Request) (*http.Request, error) {
		req.Header.Add(key, value)
		return req, nil
	}
}
