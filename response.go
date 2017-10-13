package chttp

import (
	"encoding/json"
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

func (c Response) ReadJSON(target interface{}) error {
	defer c.Raw.Body.Close()
	return json.NewDecoder(c.Raw.Body).Decode(target)
}
