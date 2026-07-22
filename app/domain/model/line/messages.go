package line

import "errors"

// ErrLINEMessageKeyNotFound : 指定されたメッセージキーが設定に存在しない場合のエラー
var ErrLINEMessageKeyNotFound = errors.New("line message key not found")

// LINEMessages : メッセージキーと文言のマッピング
type LINEMessages map[string]string

// FindMessageFor : 指定されたキーに対応するメッセージ文言を返す
func (m LINEMessages) FindMessageFor(key string) (string, error) {
	msg, ok := m[key]
	if !ok {
		return "", ErrLINEMessageKeyNotFound
	}
	return msg, nil
}
