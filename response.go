package main

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	Raw *http.Response
	Err error
}

func (c Response) ReadAll() (body []byte, err error) {
	defer c.Raw.Body.Close()
	body, err = ioutil.ReadAll(c.Raw.Body)
	return
}
