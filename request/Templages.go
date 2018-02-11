package request

import (
	"sync"
)

type Templages struct {
	TemplageList []*Templage
	Key          uint64
	Winning      [2]float64
	Wlen         int
}

func (self *Templages) ContrastRate(te *Templage) bool {

	if (self.Key|1 == self.Key) != te.Direction {
		return true
	}
	TsCa := self.TemplageList[0].CaLeft
	num := int64(self.Key >> 1)
	isDiffCa := (TsCa == te.CaLeft)

	isfChan := make(chan bool, 10)
	w := new(sync.WaitGroup)
	defer func() {
		w.Wait()
		close(isfChan)
	}()
	go func() {
		for {
			isf, ok := <-isfChan
			if !ok {
				return
			}
			te.UpdateWinning(isf)
		}

	}()
	if num < te.Duration {
		les := float64(te.Duration / TsCa.Scale)
		if te.Duration%TsCa.Scale > 0 {
			les += 1
		}
		if (les-float64(self.TemplageList[0].Len))/les > te.Rate/2 {
			return false
		}
		te_ := te
		if !isDiffCa {
			te_ = new(Templage)
			te_.Init(te.Begin, te.End, TsCa, te.CaRight, int64(les))
			if te_.FittingInit() != nil {
				return false
			}
		}
		for _, _te := range self.TemplageList {
			w.Add(1)
			go func(_te_ *Templage) {
				if te_.ContrastDiff(*_te_) {
					isfChan <- _te_.IsF
				}
				w.Done()
			}(_te)
		}
	} else if num > te.Duration {

		les := float64(num / te.CaLeft.Scale)
		if num%te.CaLeft.Scale > 0 {
			les += 1
		}
		if (les-float64(te.Len))/les > te.Rate/2 {
			return false
		}

		for _, te_ := range self.TemplageList {
			w.Add(1)
			go func(te2 *Templage) {
				_te_ := te2
				if !isDiffCa {
					_te_ := new(Templage)
					_te_.Init(te2.Begin, te2.End, te.CaLeft, te2.CaRight, int64(les))
					if _te_.FittingInit() != nil {
						w.Done()
						return
					}
				}
				if _te_.ContrastDiff(*te) {
					isfChan <- _te_.IsF
				}
				w.Done()
			}(te_)
		}

	} else {
		for _, te_ := range self.TemplageList {

			w.Add(1)
			go func(_te_ *Templage) {
				if te.Contrast(*_te_) {
					isfChan <- _te_.IsF
				}
				w.Done()
			}(te_)
		}

	}
	return true
}

func (self *Templages) Inits(te *Templage, key uint64) {

	self.Key = key
	self.Wlen = len(te.Weight)
}
func SortAppendTemplageToLib(te *Templage) {

	le := len(TemplagesLib)
	if le == 0 {
		tes := new(Templages)
		tes.Inits(te, te.GetKey())
		te.farTes = tes
		tes.Append(te)
		TemplagesLib = []*Templages{tes}

		return
	}
	if te.farTes != nil {
		te.farTes.Append(te)
		return
	}
	M := le / 2
	tes := FindTemplages(TemplagesLib, te.GetKey(), &M, 0, le)
	if tes == nil {
		tes = new(Templages)
		tes.Inits(te, te.GetKey())
		te.farTes = tes
		tes.Append(te)
		//	te.CaRight.Templages = append(te.CaRight.Templages,tes)
		TemplagesLib = InsertTemplagesToLib(TemplagesLib, tes, M)

	} else {
		tes.Append(te)
		te.farTes = tes
	}

}
func FindTemplages(TLib []*Templages, t uint64, M *int, L, R int) *Templages {

	if TLib[*M].Key == t {
		return TLib[*M]
	}
	var le int
	if TLib[*M].Key > t {
		R = *M
		le = (*M - L)
		if le == 0 {
			//	*M--
			return nil
		}
		*M = L + le/2
	} else {
		L = *M + 1
		le = (R - L)
		if le == 0 {
			*M++
			//	if *M < len(TLib) && TLib[*M].Key == t {
			//		return TLib[*M]
			//	}
			return nil
		}
		*M = L + le/2
	}
	return FindTemplages(TLib, t, M, L, R)

}

func (self *Templages) Append(te *Templage) {

	self.TemplageList = append(self.TemplageList, te)
	if te.IsF {
		self.Winning[1]++
	} else {
		self.Winning[0]++
	}
	lastCa := CacheList[len(CacheList)-1]
	if len(lastCa.TmpCan) == 0 {
		return
	}
	be := lastCa.TmpCan[0].Time
	for {
		if len(self.TemplageList) < 2 || self.TemplageList[0].Begin.Time >= be {
			break
		}
		if self.TemplageList[0].IsF {
			self.Winning[1]--
		} else {
			self.Winning[0]--
		}
		self.TemplageList = self.TemplageList[1:]

	}

}
func InsertTemplagesToLib(TLib []*Templages, tes *Templages, M int) (blib []*Templages) {

	le := len(TLib)
	blib = make([]*Templages, le+1)
	copy(blib, TLib[:M])
	blib[M] = tes
	copy(blib[M+1:], TLib[M:])
	return blib

}
