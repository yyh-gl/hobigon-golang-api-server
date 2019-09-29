package infra

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	listForSort := requestList{}
	for endpoint, count := range accessCountPerEndpoint {
		e := request{endpoint, count}
		listForSort = append(listForSort, e)
	}
	sort.Sort(listForSort)

	rankingStr := "\nアクセスランキング"
	rank := 1
	for r, count := range listForSort {
		fmt.Println("========================")
		fmt.Println(r)
		fmt.Println(count)
		fmt.Println("========================")
		//rankingStr += "\n" + strconv.Itoa(rank) + "位 " + strconv.Itoa(count) + "回： " + request
		rank++
	}

	return rankingStr, nil
}

// ソート用の構造体およびメソッドを用意
type request struct {
	endpoint string
	count    int
}
type requestList []request

func (rl requestList) Len() int {
	return len(rl)
}

func (rl requestList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

func (rl requestList) Less(i, j int) bool {
	if rl[i].endpoint == rl[j].endpoint {
		return rl[i].count > rl[j].count
	}

	return rl[i].endpoint > rl[j].endpoint
}
