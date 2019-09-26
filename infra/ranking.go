package infra

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func GetAccessRanking() (ranking string, err error) {
	const (
		IndexPrefix     = 2
		IndexMethod     = 3
		IndexEndpoint   = 4
		AccessLogPrefix = "[AccessLog]"
		IgnoreEndpoint  = "/api/v1/rankings/access"
	)

	// app.log からアクセス記録を解析
	fp, err := os.Open(os.Getenv("LOG_PATH") + "/app.log")
	if err != nil {
		return "", err
	}
	defer fp.Close()

	accessCountPerEndpoint := map[string]int{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		req := scanner.Text()
		reqSlice := strings.Split(req, " ")

		if reqSlice[IndexPrefix] == AccessLogPrefix && reqSlice[IndexEndpoint] != IgnoreEndpoint {
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
		return "", err
	}

	// アクセス数が多い順にソート
	//sort.Slice(accessCountPerEndpoint, func(i, j int) bool {
	//	return accessCountPerEndpoint[i].Price < foods[j].Price
	//})

	rankingStr := "\nアクセスランキング"
	rank := 1
	for request, count := range accessCountPerEndpoint {
		rankingStr += "\n" + strconv.Itoa(rank) + "位 " + strconv.Itoa(count) + "回： " + request
		rank++
	}

	return rankingStr, nil
}
