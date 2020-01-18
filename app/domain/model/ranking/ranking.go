package ranking

// TODO: ドメインモデル貧血症を治す

// Ranking : ランキングを表すドメインモデル
// TODO: ドメインモデルらしくする
type Ranking []Access

// Len : AccessList の配列数を取得
func (r Ranking) Len() int {
	return len(r)
}

// Swap : 指定要素の位置を入れ替える
func (r Ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Less : 指定要素の大小関係を判定
func (r Ranking) Less(i, j int) bool {
	return r[i].Count > r[j].Count
}
