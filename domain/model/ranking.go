package model

// アクセスランキング用の構造体
// TODO: JSON のキー名をここにドメイン知識として持ってよいのか検討
type Access struct {
	Endpoint string `json:"endpoint"`
	Count    int    `json:"count"`
}

type AccessList []Access

func (al AccessList) Len() int {
	return len(al)
}

func (al AccessList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

// ソート用関数：リクエスト回数の降順でソート
func (al AccessList) Less(i, j int) bool {
	return al[i].Count > al[j].Count
}
