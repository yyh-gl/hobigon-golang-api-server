package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

// !! deprecate !!
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

// bindReqWithValidate : リクエスト内容を構造体にマッピングし、バリデーションを実施
func bindReqWithValidate(ctx context.Context, reqStruct interface{}, src interface{}) error {
	switch reflect.TypeOf(src) {
	case reflect.TypeOf(&http.Request{}):
		if err := bindFromHTTPBody(ctx, src.(*http.Request), reqStruct); err != nil {
			return fmt.Errorf("bindFromHTTPBody()でエラー: %w", err)
		}
	default:
		if err := bindFromPathParams(ctx, reqStruct); err != nil {
			return fmt.Errorf("bindFromPathParams()でエラー: %w", err)
		}
	}

	// TODO: バリデーションエラー専用のエラーレスポンスを用意
	v := validator.New()
	err := v.Struct(reqStruct)
	if err != nil {
		return fmt.Errorf("バリデーションエラー: %w", err)
	}

	return nil
}

// bindFromHTTPBody: リクエストボディの内容を構造体にマッピング（bindReqWithValidate()からしか呼び出さない）
func bindFromHTTPBody(ctx context.Context, r *http.Request, req interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll()でエラー: %w", err)
	}
	defer func() { _ = r.Body.Close() }()
	err = json.Unmarshal(body, req)
	if err != nil {
		return fmt.Errorf("json.Unmarshal()でエラー: %w", err)
	}

	return nil
}

// bindFromPathParams : パスパラメータの内容を構造体にマッピング（bindReqWithValidate()からしか呼び出さない）
func bindFromPathParams(ctx context.Context, reqStruct interface{}) error {
	ps := httprouter.ParamsFromContext(ctx)

	refValReqStruct := reflect.ValueOf(reqStruct).Elem()
	refTypeReqStruct := refValReqStruct.Type()
	for i := 0; i < refTypeReqStruct.NumField(); i++ {
		fieldName := strings.ToLower(refTypeReqStruct.Field(i).Name)
		psVal := ps.ByName(fieldName)
		val := reflect.ValueOf(strings.ToLower(psVal))
		refValReqStruct.Field(i).Set(val)
	}

	return nil
}

// bindReqWithValidate : リクエスト内容を構造体にマッピングし、バリデーションを実施
func bindReqWithValidate2(ctx context.Context, src, dist interface{}) error {
	switch reflect.TypeOf(src) {
	case reflect.TypeOf(&http.Request{}):
		if err := bindFromHTTPBody2(ctx, src.(*http.Request), dist); err != nil {
			return fmt.Errorf("bindFromHTTPBody() > %w", err)
		}
	default:
		if err := bindFromPathParams2(ctx, src.(map[string]string), dist); err != nil {
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
func bindFromHTTPBody2(ctx context.Context, r *http.Request, dist interface{}) error {
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
func bindFromPathParams2(ctx context.Context, src map[string]string, dist interface{}) error {
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
