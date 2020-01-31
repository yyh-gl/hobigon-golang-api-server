package blog

import "strconv"

// initCount : Count のデフォルト初期値
const initCount = 0

// Count : ブログのいいね数を表す値オブジェクト
type Count int

// NewCount : Count を生成
func NewCount() (*Count, error) {
	n := Count(initCount)
	return &n, nil
}

// NewCount : 引数の値をもとに Count を生成
func NewCountWithArg(val int) (*Count, error) {
	n := Count(val)
	return &n, nil
}

// Int : Countの 値を Int として返却
func (n Count) Int() int {
	return int(n)
}

// String : Countの 値を String として返却
func (n Count) String() string {
	return strconv.Itoa(n.Int())
}
