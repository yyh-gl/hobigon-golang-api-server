package infra

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

func GetAccessRanking() (rankingMsg string, accessList model.AccessList, err error) {
	const (
		IndexPrefix     = 2
		IndexMethod     = 3
		IndexEndpoint   = 4
		AccessLogPrefix = "[AccessLog]"
	)

	var IgnoreEndpoints = []string{"/api/v1/rankings/access", "/api/v1/tasks"}

	// app.log からアクセス記録を解析
	fp, err := os.Open(os.Getenv("LOG_PATH") + "/app.log")
	if err != nil {
		return "", nil, err
	}
	defer fp.Close()

	accessCountPerEndpoint := map[string]int{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		req := scanner.Text()
		reqSlice := strings.Split(req, " ")

		if reqSlice[IndexPrefix] == AccessLogPrefix && !isContain(IgnoreEndpoints, reqSlice[IndexEndpoint]) {
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
		return "", nil, err
	}

	// アクセス数が多い順にソート
	accessList = model.AccessList{}
	for endpoint, count := range accessCountPerEndpoint {
		e := model.Access{
			Endpoint: endpoint,
			Count:    count,
		}
		accessList = append(accessList, e)
	}
	sort.Sort(accessList)

	// Slack 通知用のメッセージを作成
	rankingMsg = "\n【アクセスランキング】"
	for i, req := range accessList {
		// Slack 通知では10位まで表示
		if i >= 10 {
			break
		}

		rankingMsg += "\n" + strconv.Itoa(i+1) + "位 " + strconv.Itoa(req.Count) + "回： " + req.Endpoint
	}

	return rankingMsg, accessList, nil
}

func isContain(arr []string, str string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}
