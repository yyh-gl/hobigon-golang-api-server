package entity

// Access : アクセス情報の構造体
type Access struct {
	Endpoint string
	Count    int
}

// AccessList : アクセス情報の配列構造体
type AccessList []Access

// Len : AccessList の配列数を取得
func (al AccessList) Len() int {
	return len(al)
}

// Swap : 指定要素の位置を入れ替える
func (al AccessList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

// Less : 指定要素の大小関係を判定
func (al AccessList) Less(i, j int) bool {
	return al[i].Count > al[j].Count
}
