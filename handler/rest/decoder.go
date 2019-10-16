package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func decodeRequest(r *http.Request, req interface{}) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}

	return req.(map[string]interface{}), nil
}
