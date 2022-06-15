package blog

import "strconv"

// initCount : Countのデフォルト初期値
const initCount = 0

// Count : ブログのいいね数を表す値オブジェクト
type Count int

// newCount : Countを生成
func newCount() (Count, error) {
	return Count(initCount), nil
}

// Int : Countの値をIntとして返却
func (n Count) Int() int {
	return int(n)
}

// String : Countの値をStringとして返却
func (n Count) String() string {
	return strconv.Itoa(n.Int())
}
