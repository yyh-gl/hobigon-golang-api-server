package repository

import "errors"

// ErrRecordNotFound : DBにおける"record not found"エラー
var ErrRecordNotFound = errors.New("record not found")
