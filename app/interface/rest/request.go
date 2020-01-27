package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// !! deprecate !1
// decodeRequest : リクエストボディの内容を構造体にマッピング
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

// bindReqWithValidate : リクエストボディの内容を構造体にマッピングし、バリデーションを実施
func bindReqWithValidate(r *http.Request, req interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll()でエラー: %w", err)
	}
	defer func() { _ = r.Body.Close() }()

	err = json.Unmarshal(body, req)
	if err != nil {
		return fmt.Errorf("json.Unmarshal()でエラー: %w", err)
	}

	// TODO: バリデーションエラー専用のレスポンスを用意
	v := validator.New()
	err = v.Struct(req)
	if err != nil {
		return fmt.Errorf("バリデーションエラー: %w", err)
	}

	return nil
}
