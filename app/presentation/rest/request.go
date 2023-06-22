package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

// bindReqWithValidate : リクエスト内容を構造体にマッピングし、バリデーションを実施
func bindReqWithValidate(ctx context.Context, src, dist any) error {
	switch reflect.TypeOf(src) {
	case reflect.TypeOf(&http.Request{}):
		if err := bindFromHTTPBody(ctx, src.(*http.Request), dist); err != nil {
			return fmt.Errorf("bindFromHTTPBody() > %w", err)
		}
	default:
		if err := bindFromPathParams(ctx, src.(map[string]string), dist); err != nil {
			return fmt.Errorf("bindFromPathParams() > %w", err)
		}
	}

	// TODO: バリデーションエラー専用のエラーレスポンスを用意
	v := validator.New()
	err := v.Struct(dist)
	if err != nil {
		return fmt.Errorf("validation error > %w", err)
	}
	return nil
}

// bindFromHTTPBody: リクエストボディの内容を構造体にマッピング
func bindFromHTTPBody(ctx context.Context, r *http.Request, dist any) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll() > %w", err)
	}
	defer func() { _ = r.Body.Close() }()

	err = json.Unmarshal(body, dist)
	if err != nil {
		return fmt.Errorf("json.Unmarshal() > %w", err)
	}

	return nil
}

// bindFromPathParams : パスパラメータの内容を構造体にマッピング（bindReqWithValidate()からしか呼び出さない）
func bindFromPathParams(ctx context.Context, src map[string]string, dist any) error {
	if err := mapstructure.Decode(src, &dist); err != nil {
		return fmt.Errorf("mapstructure.Decode() > %w", err)
	}

	// TODO: バリデーションエラー専用のエラーレスポンスを用意
	v := validator.New()
	err := v.Struct(dist)
	if err != nil {
		return fmt.Errorf("validate error > %w", err)
	}

	return nil
}
