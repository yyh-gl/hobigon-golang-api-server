package iservice

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/ranking"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/service"
)

//////////////////////////////////////////////////
// NewRankingService
//////////////////////////////////////////////////

type rankingService struct{}

// NewRankingService : ランキング用のサービスを取得
func NewRankingService() service.RankingService {
	return &rankingService{}
}

//////////////////////////////////////////////////
// GetAccessRanking
//////////////////////////////////////////////////

// isContain : 文字列の配列に指定文字列が存在するか確認
func isContain(arr []string, str string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

// GetAccessRanking : アクセスランキングを取得する関数
func (rs rankingService) GetAccessRanking(ctx context.Context) (rankingMsg string, accessList ranking.AccessList, err error) {
	const (
		IndexPrefix     = 2
		IndexMethod     = 3
		IndexEndpoint   = 4
		AccessLogPrefix = "[AccessLog]"
	)

	// TODO: /api/v1/blogs/*/like というように正規表現で ignore 指定できるようにする
	var IgnoreEndpoints = []string{"/api/v1/rankings/access", "/api/v1/tasks", "/api/v1/blogs/good_api/like"}

	// app.log からアクセス記録を解析
	fp, err := os.Open(os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
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
	accessList = ranking.AccessList{}
	for endpoint, count := range accessCountPerEndpoint {
		e := ranking.Access{
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
