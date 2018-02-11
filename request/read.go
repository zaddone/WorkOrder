package request

import (
	//	"fmt"
	"time"
	//	"math"
)

var (
	GranularityMap map[string]int64
	CacheList      []*Cache
)

func init() {

	GranularityMap = make(map[string]int64)
	GranularityMap["S5"] = 5
	GranularityMap["S10"] = 10
	GranularityMap["S15"] = 15
	GranularityMap["S30"] = 30
	GranularityMap["M1"] = 60
	GranularityMap["M2"] = 60 * 2
	GranularityMap["M4"] = 60 * 4
	GranularityMap["M5"] = 60 * 5
	GranularityMap["M10"] = 600
	GranularityMap["M15"] = 60 * 15
	GranularityMap["M30"] = 60 * 30
	GranularityMap["H1"] = 3600
	GranularityMap["H2"] = 3600 * 2
	GranularityMap["H3"] = 3600 * 3
	GranularityMap["H4"] = 3600 * 4
	GranularityMap["H6"] = 3600 * 6
	GranularityMap["H8"] = 3600 * 8
	GranularityMap["H12"] = 3600 * 12
	GranularityMap["D"] = 3600 * 24
	GranularityMap["W"] = 3600 * 24 * 7
	//GranularityMap["M"]= 3600*24*30
	i := 0
	for k, v := range GranularityMap {
		ca := new(Cache)
		ca.Init(k, v)
		CacheList = append(CacheList, ca)
		SortCacheList(i)
		i++
	}

}
func SortCacheList(i int) {
	if i == 0 {
		return
	}
	I := i - 1
	if CacheList[I].Scale > CacheList[i].Scale {
		CacheList[I], CacheList[i] = CacheList[i], CacheList[I]
		SortCacheList(I)
	}
}
func ReadTemplate() {

	end := len(CacheList) - 1
	for _, _ca := range CacheList[1:end] {
		go _ca.SyncRun(_ca.UpdateSuccessive)
	}
	endCa := CacheList[end]
	go endCa.SyncRun(endCa.UpdateTmpCan)
	time.Sleep(time.Second)

	CacheList[0].Sensor(1)
	//	ca.Sensor(1)

}

//func ReadTemplateExt(){
//
//	for _,_ca := range CacheList[1:]{
//		go _ca.SyncSuccessive()
//	}
//	CacheList[0].Sensor(1)
//
//}
