package analysis

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// access : アクセス情報を表す構造体
type access struct {
	endpoint string
	count    int
}

// accessList : アクセス情報を表す構造体のリスト
type accessList []access

// Len : accessList の配列数を取得
func (l accessList) Len() int {
	return len(l)
}

// Swap : 指定要素の位置を入れ替える
func (l accessList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less : 指定要素の大小関係を判定
func (l accessList) Less(i, j int) bool {
	return l[i].count < l[j].count
}

// GetAccessRanking : アクセスランキングを取得する関数
func GetAccessRanking(ctx context.Context) (rankingMsg string, notifiedNum int, err error) {
	const (
		IndexPrefix     = 2
		IndexMethod     = 3
		IndexEndpoint   = 4
		AccessLogPrefix = "[AccessLog]"
	)

	// TODO: /api/v1/blogs/*/like というように正規表現で ignore 指定できるようにする
	var ignoreEndpoints = []string{"/api/v1/rankings/access", "/api/v1/tasks", "/api/v1/blogs/good_api/like"}
	var httpMethodOption = "OPTIONS"

	// app.log からアクセス記録を解析
	fp, err := os.Open(os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = fp.Close() }()

	accessCountPerEndpoint := map[string]int{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		req := scanner.Text()
		reqSlice := strings.Split(req, " ")

		if reqSlice[IndexPrefix] == AccessLogPrefix && !isContain(ignoreEndpoints, reqSlice[IndexEndpoint]) && reqSlice[IndexMethod] != httpMethodOption {
			key := reqSlice[IndexMethod] + "_" + reqSlice[IndexEndpoint]

			_, exist := accessCountPerEndpoint[key]
			if exist {
				accessCountPerEndpoint[key]++
			} else {
				accessCountPerEndpoint[key] = 1
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return "", 0, err
	}

	// アクセス数が多い順にソート
	accessList := accessList{}
	for endpoint, count := range accessCountPerEndpoint {
		accessList = append(accessList, access{
			endpoint: endpoint,
			count:    count,
		})
	}
	sort.Sort(sort.Reverse(accessList))

	// Slack 通知用のメッセージを作成
	rankingMsg = "\n【アクセスランキング】"
	for i, req := range accessList {
		// Slack 通知では10位まで表示
		if i >= 10 {
			break
		}

		rankingMsg += "\n" + strconv.Itoa(i+1) + "位 " + strconv.Itoa(req.count) + "回： " + req.endpoint
	}

	return rankingMsg, len(accessList), nil
}

// isContain : 文字列の配列に指定文字列が存在するか確認
func isContain(arr []string, str string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}
