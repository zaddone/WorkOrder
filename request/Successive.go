package request

import (
	//"fmt"
	"math"
)

type Successive struct {
	//Count          int
	SumLong        float64
	Diff           float64
	tmpCans        []*CandlesMin
	LastSuccessive *Successive
}

func (self *Successive) Init(cans []*CandlesMin) {

	le := len(cans)
	self.tmpCans = make([]*CandlesMin, le)
	self.SumLong = 0
	for i, _can := range cans {
		self.SumLong += _can.GetMidLong()
		self.tmpCans[i] = _can
	}
	self.Diff = cans[le-1].GetMidAverage() - cans[0].GetMidAverage()

}
func (self *Successive) Update(can *CandlesMin) {

	defer func() {
		self.tmpCans = append(self.tmpCans, can)
		self.SumLong += can.GetMidLong()
		if len(self.tmpCans) > 1 {
			self.Diff = can.GetMidAverage() - self.tmpCans[0].GetMidAverage()
		}
	}()
	le := len(self.tmpCans)
	if le < 2 {
		return
	}
	var maxDif float64 = 0
	var maxId int
	ave := self.GetLongAve()
	for i, _can := range self.tmpCans {
		dif := can.GetMidAverage() - _can.GetMidAverage()
		if (dif > 0) != (self.Diff > 0) {
			dif = math.Abs(dif)
			if (dif > ave) && (dif > maxDif) {
				maxDif = dif
				maxId = i
			}
		}
	}
	if maxDif == 0 || maxId == 0 {
		return
	}
	if self.LastSuccessive == nil {
		self.LastSuccessive = new(Successive)
	}
	self.LastSuccessive.Init(self.tmpCans[:maxId])
	self.Init(self.tmpCans[maxId:])

}

func (self *Successive) GetLongAve() float64 {
	return self.SumLong / float64(len(self.tmpCans))
}

func (self *Successive) Check(can *CandlesMin) bool {

	if self.LastSuccessive == nil {
		return false
	}
	dif1 := can.GetMidAverage() - self.tmpCans[0].GetMidAverage()
	dif2 := can.GetMidAverage() - self.LastSuccessive.tmpCans[0].GetMidAverage()
	//if (dif1 > 0) == (dif2 > 0) {
	//	return false
	//}
	if ((dif1 > 0) == (dif2 > 0)) || (math.Abs(dif2) < ((self.LastSuccessive.GetLongAve() + self.GetLongAve()) / 2)) {
		return true
	}
	return false

}
