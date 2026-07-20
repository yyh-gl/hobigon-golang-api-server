package line

// LINEMessages : メッセージキーと文言のマッピング
type LINEMessages map[string]string

// MessageFor : 指定されたキーに対応するメッセージ文言を返す
func (m LINEMessages) MessageFor(key string) (string, error) {
	return m[key], nil
}
