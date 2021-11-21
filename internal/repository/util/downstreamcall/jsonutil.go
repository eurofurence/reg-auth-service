package downstreamcall

import "encoding/json"

// helper functions for dealing with json

func RenderJson(dto interface{}) (string, error) {
	representationBytes, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}
	return string(representationBytes), nil
}

// tip: dto := &whatever.WhateverDto{}
func ParseJson(body string, dto interface{}) error {
	err := json.Unmarshal([]byte(body), dto)
	return err
}

