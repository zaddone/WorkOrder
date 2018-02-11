package request

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache struct {
	Scale       int64
	Name        string
	TmpCan      []*CandlesMin
	EndtimeChan chan int64
	Successive  *Successive
	Stop        chan bool

	FutureChan chan *CacheFile

	TemplagesLib []*Templages
	Templages    *Templages
	DiffLong     float64
}

func (self *Cache) Init(name string, scale int64) {

	self.Scale = scale
	self.Name = name
	self.EndtimeChan = make(chan int64, 10)
	self.Successive = new(Successive)
	self.Stop = make(chan bool, 1)
	self.FutureChan = make(chan *CacheFile, 10)
	go self.ReadToCache(filepath.Join(*InsName, name))

}

func (self *Cache) ShowTemplageLib() {
	count := 0
	var winning [2]float64
	for _, tes := range self.TemplagesLib {
		count += len(tes.TemplageList)
		winning[0] += tes.Winning[0]
		winning[1] += tes.Winning[1]
	}
	fmt.Println("Show", self.Name, count, winning)
}

func FindCandlesBack(can []*CandlesMin, t int64, M *int, L, R int) {

	if can[*M].Time == t {
		return
	}
	var le int
	if can[*M].Time > t {
		R = *M
		le = (*M - L)
		if le == 0 {
			*M--
			return
		}
		*M = L + le/2
	} else {
		L = *M + 1
		le = (R - L)
		if le == 0 {
			//*M++
			return
		}
		*M = L + le/2
	}
	FindCandlesBack(can, t, M, L, R)

}

func FindCandles(can []*CandlesMin, t int64, M *int, L, R int) {

	if can[*M].Time == t {
		return
	}
	var le int
	if can[*M].Time > t {
		R = *M
		le = (*M - L)
		if le == 0 {
			//*M--
			return
		}
		*M = L + le/2
	} else {
		L = *M + 1
		le = (R - L)
		if le == 0 {
			*M++
			return
		}
		*M = L + le/2
	}
	FindCandles(can, t, M, L, R)

}

func (self *Cache) SearchHandle(beginTime int64, Handle func(can *CandlesMin) error) {

	le := len(self.TmpCan)
	if le == 0 {
		return
	}
	I := le / 2
	FindCandles(self.TmpCan, beginTime, &I, 0, le)
	var err error
	for ; I < le; I++ {
		err = Handle(self.TmpCan[I])
		if err != nil {
			return
		}
	}
	return

}

func (self *Cache) Sensor(ca_id int) {

	var lastCan *CandlesMin
	//	var dir,lastdir bool
	//	var count int
	var add bool
	var tes []*Templage

	handle := func(te *Templage) *Templage {
		if len(tes) > 0 {
			for _, _te := range tes {
				if _te.Comp(te) {
					return nil
				}
			}
		}
		err := te.FittingInit()
		if err == nil {
			return te
		}
		return nil
	}
	self.Read(func(can *CandlesMin) {
		//fmt.Printf("%s\r", time.Unix(can.Time, 0).UTC())
		add, lastCan = self.CacheUpdate(can)
		if !add {
			return
		}
		if len(tes) > 0 {
			var Ntes []*Templage
			for _, te := range tes {
				if te.CheckIsF(can) {
					//SortAppendTemplageToLib(te)
				} else {
					Ntes = append(Ntes, te)
				}
			}
			tes = Ntes
		}
		_tes, lastid := self.CheckTemplages(can, ca_id)
		if len(_tes) == 0 {
			return
		}
		var tmpTes []*Templage
		//var endTe *Templage
		for _, te_ := range _tes {
			te_ = handle(te_)
			if te_ != nil {
				tmpTes = append(tmpTes, te_)
			}
		}
		if len(tmpTes) == 0 {
			return
		}
		tes = append(tes, tmpTes...)

		endTe := tmpTes[len(tmpTes)-1]

		var lastCaTe []*Templage
		var SameCaTe []*Templage
		var isPOTe *Templage = nil
		for _, te_ := range tes[:len(tes)-1] {
			if te_.IsOrder {
				isPOTe = te_
			}
			if te_.Duration >= endTe.Duration || te_.CaRight == endTe.CaRight {
				//if te_.CaRight == endTe.CaRight {
				if te_.Direction != endTe.Direction {
					return
				}
				//SameCaTe = append(SameCaTe, te_)
				if te_.CaRight == endTe.CaRight {
					SameCaTe = append(SameCaTe, te_)
				} else {
					lastCaTe = append(lastCaTe, te_)
				}
			}
		}

		if len(SameCaTe) == 0 {
			return
		}
		if len(lastCaTe) == 0 {
			if (lastid == len(CacheList)-1) || (CacheList[lastid].Successive.Diff > 0) == endTe.Direction {
				return
			}
		}

		//if !endTe.SetFarTes() {
		//	return
		//}
		go endTe.PostOrderCheck(isPOTe)
		return

	})

}

func (self *Cache) SyncRun(hand func(can *CandlesMin)) {
	endTime := <-self.EndtimeChan
	var h func(can *CandlesMin)
	h = func(can *CandlesMin) {
		if can.Time+self.Scale <= endTime {
			hand(can)
			return
		}
		if len(self.EndtimeChan) == 0 {
			if len(self.Stop) == 0 {
				self.Stop <- true
			}
		}
		endTime = <-self.EndtimeChan
		h(can)
		return
	}
	self.Read(h)
}

func (self *Cache) UpdateTmpCan(can *CandlesMin) {

	le := len(self.TmpCan)
	if le == 0 {
		self.TmpCan = []*CandlesMin{can}
		return
	}
	lastCan := self.TmpCan[le-1]
	if lastCan.Time > can.Time {
		//fmt.Println(time.Unix(can.Time,0))
		return
	}
	self.TmpCan = append(self.TmpCan, can)

	var sumdifflong float64 = can.GetMidLong()
	var count float64 = 1
	var ave float64 = 1
	for i := le - 1; i >= 0; i-- {
		sumdifflong += self.TmpCan[i].GetMidLong()
		count += 1
		self.DiffLong = can.GetMidAverage() - self.TmpCan[i].GetMidAverage()
		if math.Abs(self.DiffLong) > (sumdifflong/count)*ave {
			self.TmpCan = self.TmpCan[i:]
			return
		}
	}
	return
}

func (self *Cache) UpdateSuccessive(can *CandlesMin) {

	add, lastCan := self.CacheUpdate(can)
	if !add {
		return
	}
	if lastCan == nil {
		return
	}
	self.Successive.Update(can)

}

func (self *Cache) CacheUpdate(can *CandlesMin) (add bool, last *CandlesMin) {

	le := len(self.TmpCan)
	if le == 0 {
		self.TmpCan = []*CandlesMin{can}
		return true, nil
	}
	last = self.TmpCan[le-1]
	if last.Time >= can.Time {
		//fmt.Println(time.Unix(can.Time,0))
		return false, last
	}
	self.TmpCan = append(self.TmpCan, can)

	lastCa := CacheList[len(CacheList)-1]
	if self == lastCa {
		return true, last
	}
	if len(lastCa.TmpCan) == 0 {
		return true, last
	}
	be := lastCa.TmpCan[0].Time
	for i, _can := range self.TmpCan {
		if _can.Time >= be {
			self.TmpCan = self.TmpCan[i:]
			break
		}
	}
	return true, last

}

func (self *Cache) ReadToCache(path string) {
	f, err := os.Stat(path)
	var cf *CacheFile
	if err == nil && f.IsDir() {
		filepath.Walk(path, func(pa string, fi os.FileInfo, er error) error {
			if fi.IsDir() {
				return er
			}
			cf = new(CacheFile)
			MaxInt := int(86400/self.Scale+1) * 2
			er = cf.Init(pa, fi, MaxInt)
			if er == nil {
				self.FutureChan <- cf
			} else {
				fmt.Println(er)
			}
			return er
		})
	}
	var begin int64
	if cf == nil {
		beginT, err := time.Parse("2006-01-02T15:04:05", *BEGINTIME)
		if err != nil {
			panic(err)
		}
		begin = beginT.Unix()
	} else {
		begin = cf.EndCan.Time + self.Scale
	}
	cf = new(CacheFile)
	cf.Can = make(chan *CandlesMin, 1000)
	self.FutureChan <- cf
	Down(begin, 0, self.Name, func(can *CandlesMin) {
		cf.Can <- can
	})
	close(cf.Can)
}
func ReadPath(path string, f func(can *CandlesMin)) {

	filepath.Walk(path, func(pa string, fi os.FileInfo, er error) error {
		if fi.IsDir() {
			//	fmt.Println(pa,fi.Name())
			//	ReadDB(pa,f)
			return er
		}
		fe, er := os.Open(pa)
		if er != nil {
			return er
		}
		defer fe.Close()
		//		lastPath = pa
		r := bufio.NewReader(fe)
		for {
			db, _, e := r.ReadLine()
			if e != nil {
				break
			}
			can := new(CandlesMin)
			can.Load(string(db))
			f(can)
			//	if er != nil {
			//		return er
			//		//fmt.Println(er)
			//	}
		}
		return er
	})

}

func (self *Cache) Read(Handle func(can *CandlesMin)) {
	for {
		cf := <-self.FutureChan
		if cf == nil {
			break
		}
		for {
			can := <-cf.Can
			if can == nil {
				break
			}
			Handle(can)
		}
	}

}

func (self *Cache) CheckTemplages(can *CandlesMin, ca_id int) (Tmp []*Templage, lastid int) {
	endTime := can.Time + self.Scale
	cas := CacheList[ca_id:]
	w := new(sync.WaitGroup)
	for _, ca := range cas {
		w.Add(1)
		go func(_ca *Cache) {
			_ca.EndtimeChan <- endTime
			<-_ca.Stop
			w.Done()
		}(ca)
	}
	w.Wait()
	fu := func(Dur int64) (_ca *Cache, le int64) {
		if Dur > 3600 {
			return nil, 0
		}
		for j := len(CacheList) - 1; j >= 0; j-- {
			_ca = CacheList[j]
			le = Dur / _ca.Scale
			if Dur%_ca.Scale > 0 {
				le++
			}
			if le > 100 {
				return _ca, le
			}
		}
		return _ca, le
	}
	for _i, ca := range cas[:len(cas)-1] {
		lastid = _i
		if !ca.Successive.Check(can) {
			continue
		} else {
			beginCan := ca.Successive.LastSuccessive.tmpCans[0]
			_ca, le := fu(endTime - beginCan.Time)
			if le < 20 {
				continue
			}
			te := new(Templage)
			te.Init(beginCan, can, _ca, ca, le)
			Tmp = append(Tmp, te)
		}
	}
	return Tmp, lastid + ca_id

}
func CompUint8(ke uint64) (u int) {
	u = 8
	if ke == 0 {
		return u
	}
	var i uint
	var k uint64
	for i = 0; i < 64; i += 8 {
		k = ke >> i
		if (k | 1) == k {
			u--
		}
	}
	return u
}
