package _package

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func CalcTime(st string, cnt string) []int {
	NowTime := time.Now().Unix()
	SetTime := regexp.MustCompile("[/:]").Split(st, -1)
	t := make([]int, 5)
	for i := 0; i < len(SetTime); i++ {
		n, _ := strconv.Atoi(SetTime[i])
		t[i] = n
	}
	t1 := time.Date(t[0], time.Month(t[1]), t[2], t[3], t[4], 0, 0, Location())
	return InformCnt(NowTime, t1.Unix(), cnt)
}

// Location 東京のタイムゾーンを取得
func Location() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// InformCnt 優先度によって通知する時間を決定
func InformCnt(Now int64, info int64, cnt string) []int {
	TimeRemain := int(info - Now)
	var informTime []int
	switch cnt {
	case "1":
		informTime = append(informTime, 86400)
	case "2":
		informTime = append(informTime, 86400, 259200)
	case "3":
		informTime = append(informTime, 10800, 86400, 172800, 259200, 604800)
	case "4":
		informTime = append(informTime, 3600, 10800, 21600, 43200, 86400, 172800, 259200,
			345600, 432000, 518400, 604800, 2592000, 7776000, 15552000)
	case "5":
		informTime = append(informTime, 600, 1800, 3600, 10800, 21600, 43200, 86400, 172800, 259200,
			345600, 432000, 518400, 604800, 1209600, 1814400, 2419200, 5184000, 7776000, 10368000, 12960000, 15552000, 31536000)
	default:
		fmt.Println("error")
	}

	//返り値
	var PushTime []int
	for _, Tm := range informTime {
		if TimeRemain > Tm { //残されている時間 > 通知する候補の時間の場合，通知する時間になる
			PushTime = append(PushTime, int(info)-Tm)
		}
	}
	return PushTime
}
