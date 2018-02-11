package request

import (
	"../fitting"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var (
	TemplagesLib []*Templages
	Winning      [2]float64
	//	TimeWinning  [][2]float64
	LastMonth int
	LastDay   int
	//Order_Response *OrderResponse
	//TemplagesLib2 []*Templages
	Winning2 [2]float64
	Winning3 [2]float64
	Plval    float64
	DebugLog *log.Logger
)

func init() {
	flag.Parse()
	LastMonth = -1
	fileName := fmt.Sprintf("%s_%s.log", time.Now().Format("20060102"), *InsName)
	//	logFile, err := os.Create(fileName)
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_SYNC, 0777)
	if err != nil {
		log.Fatal(err)
	}

	DebugLog = log.New(logFile, "", log.LstdFlags)
}

type Templage struct {
	CaLeft  *Cache
	CaRight *Cache

	Begin *CandlesMin
	End   *CandlesMin

	Duration  int64
	Direction bool
	Len       int64

	X    []float64
	Y    []float64
	YMin float64
	YLon float64

	WList  []float64
	Weight []float64

	Rate float64
	//Count int

	IsF     bool
	MIsF    bool
	IsBuy   bool
	IsFar   bool
	IsOrder bool

	//	tmpId []int
	tmpSame []*Templage

	farTes *Templages
	farM   int

	Winning  [2]float64
	Winning1 [2]float64

	//	BigKey uint64
	Hide bool
	//	Successive Successive
}

func (self *Templage) Show(Y []float64) {
	//	return
	var Dx int = 64
	var Dxf float64 = 63
	max := make([][]bool, len(self.WList))
	if Y != nil {
		var Dy float64 = float64(len(self.WList) - 1)
		for _i, x := range self.X {
			i := int(Rounding(x * Dy))
			max[i] = make([]bool, Dx)
			j := int(Rounding(Y[_i] * Dxf))
			max[i][j] = true
		}
	}
	for i, l := range self.WList {
		if max[i] == nil {
			max[i] = make([]bool, Dx)
		}
		j := int(Rounding(l * Dxf))
		if j < 0 {
			j = 0
		} else if j >= Dx {
			j = Dx - 1
		}
		//fmt.Println(j,l)
		max[i][j] = true
		str := ""
		for _, m := range max[i] {
			if m {
				str = fmt.Sprintf("%s%d", str, 1)
			} else {
				str = fmt.Sprintf("%s%d", str, 0)
			}
		}
		fmt.Println(str)
	}
	fmt.Println(self.Rate, self.Duration, self.Weight, len(self.X), self.Len)
	var cmd string
	fmt.Scanf("%s", &cmd)
}

func CompKey(keyA uint64, keyB uint64) (long int) {

	ke := keyA ^ keyB
	if ke == 0 {
		return 0
	}
	//	ShowUint64(ke)
	var i uint
	var k uint64
	for i = 0; i < 64; i += 8 {
		k = ke >> i
		if (k | 1) == k {
			return -1
		}
	}
	for i = 0; i < 64; i++ {
		k = ke >> i
		if (k | 1) == k {
			long++
		}
	}
	return long

}
func ShowUint64(k uint64) {

	var str string
	var _k uint64
	var i uint
	for i = 0; i < 64; i++ {
		if i%8 == 0 {
			str = " " + str
		}
		_k = k >> i
		if _k|1 == _k {
			str = "1" + str
		} else {
			str = "0" + str
		}
	}
	fmt.Println(str)
	//	fmt.Scanf("%s",&str)
	//	str = ""

}

func (self *Templage) SetMIsF() bool {

	//sum := self.Winning[0] + self.Winning[1]
	//if sum > 50 {
	//	if (self.Winning[0] / sum) < 0.6 {
	//		return false
	//	}
	//} else {
	if self.Winning[1] > self.Winning[0] {
		return false
	}
	//}

	self.MIsF = true
	return true

}

func (self *Templage) CheckIsF(can *CandlesMin) bool {
	dif := can.Mid[1] - self.End.Mid[1]
	MaxLon := self.YLon / 2
	if MaxLon < Instr.MinimumTrailingStopDistance {
		MaxLon = Instr.MinimumTrailingStopDistance
	}
	if math.Abs(dif) < MaxLon {
		//	if (can.Time - self.End.Time) < (self.CaRight.Scale*2) {
		return false
		//	}
	}
	self.Hide = false
	self.IsF = (dif > 0 == self.Direction)
	func() {
		if !self.IsBuy {
			return
		}
		if self.IsOrder {
			go func() {
				pl, err := ClosePosition()
				if err != nil {
					Plval += pl
				}
				DebugLog.Print("order ", dif, pl, Plval)
				//fmt.Println(dif, pl, Plval)
			}()
		}
		now := time.Unix(can.Time, 0).UTC()
		y, m, d := now.Date()
		if d != LastDay {
			DebugLog.Print("day ", y, m, d, Winning3)
			//fmt.Println("day", y, m, d, Winning3)
			Winning3[0] = 0
			Winning3[1] = 0
			LastDay = d
		}
		if int(m) != LastMonth {
			//TimeWinning = append(TimeWinning, Winning)
			sumW := (Winning[0] + Winning[1])
			dif := (Winning[0] - Winning[1])

			Winning2[0] += Winning[0]
			Winning2[1] += Winning[1]
			sumW2 := (Winning2[0] + Winning2[1])
			dif2 := (Winning2[0] - Winning2[1])
			DebugLog.Print("month ", dif/sumW, int(dif), int(sumW), "---", dif2/sumW2, int(dif2), int(sumW2), y, m, Plval)
			fmt.Println(dif/sumW, int(dif), int(sumW), "---", dif2/sumW2, int(dif2), int(sumW2), y, m, Plval)

			Winning[0] = 0
			Winning[1] = 0
			LastMonth = int(m)
		}
		if !self.IsF {
			Winning[0]++
			Winning3[0]++
		} else {
			Winning[1]++
			Winning3[1]++
		}
		//fmt.Printf("%f %d %d %s %d %d %d\r\n",dif/sumW,int(dif),int(sumW),time.Unix(can.Time,0).UTC(),self.Duration,int(self.Winning[0]),int(self.Winning[1]))
	}()
	return true
}
func (self *Templage) Comp(te *Templage) bool {
	//	if self.Duration != te.Duration {
	//		return false
	//	}
	if self.Begin.Time != te.Begin.Time {
		return false
	}
	return true
}

func (self *Templage) Init(begin, end *CandlesMin, CaL, CaR *Cache, Le int64) {

	//	self.EndTime = EndTime
	self.Begin = begin
	self.End = end
	self.CaLeft = CaL
	self.CaRight = CaR
	self.Len = Le

	self.Duration = Le * self.CaLeft.Scale
	//	self.BigKey = Key
	self.Hide = true
	self.IsBuy = false
	self.MIsF = false
	self.IsFar = false
	self.Direction = self.CaRight.Successive.Diff > 0

}

func (self *Templage) UpdateWinning(isf bool) {
	if isf {
		self.Winning[1]++
	} else {
		self.Winning[0]++
	}
}
func (self *Templage) SetFarTes() bool {
	self.IsFar = true
	le := len(TemplagesLib)
	if le == 0 {
		return false
	}
	M := le / 2
	self.farTes = FindTemplages(TemplagesLib, self.GetKey(), &M, 0, le)
	if self.farTes != nil {
		self.farM = M
		self.farTes.ContrastRate(self)
	}
	/**
	for i := M; i < le; i++ {
		if !TemplagesLib[i].ContrastRate(self) {
			break
		}
	}
	for i := M - 1; i >= 0; i-- {
		if !TemplagesLib[i].ContrastRate(self) {
			break
		}
	}

	**/
	return self.SetMIsF()

}

func (self *Templage) PostOrderCheck(pote *Templage) {
	self.IsBuy = true
	if (self.End.Time + CacheList[0].Scale*2) < time.Now().Unix() {
		return
	}
	unit := 100
	if pote != nil {
		unit *= 2
	}
	if self.Direction {
		unit = -unit
	}
	//if Order_Response == nil {
	MaxLon := self.YLon
	if MaxLon < Instr.MinimumTrailingStopDistance*1.5 {
		MaxLon = Instr.MinimumTrailingStopDistance * 1.5
	}
	tp := self.End.GetMidAverage()
	sl := self.End.GetMidAverage()
	if self.Direction {
		tp -= MaxLon
		sl += MaxLon
	} else {
		tp += MaxLon
		sl -= MaxLon
	}

	or, err := HandleOrder(unit, MaxLon, tp, sl)
	if err != nil {
		fmt.Println(err)
		return
	}
	if or.ErrorCode != "" {
		DebugLog.Println("Post ", MaxLon, tp, sl, or.LastTransactionID)
	} else {
		DebugLog.Print("Post ", self.End.GetMidAverage, or)
		self.IsOrder = true
		if pote != nil {
			pote.IsOrder = false

		}
	}
	//Order_Response = &or
	//}
}
func (self *Templage) ContrastDiff(te Templage) bool {
	leDif := float64(self.Len - te.Len)
	scale := float64(te.CaLeft.Scale)
	var i float64 = 0
	X := make([]float64, len(te.X))
	copy(X, te.X)
	te.X = X
	for i = 0; i < leDif; i++ {
		for j, x := range te.X {
			te.X[j] = x + scale
		}
		if self.ContrastRateDiff(te, leDif) {
			return true
		}
	}
	return false

}
func (self *Templage) ContrastRateDiff(te Templage, j float64) bool {

	le := len(te.Y)
	Y := make([]float64, le)
	copy(Y, te.Y)
	NorMalizationOEM(Y, &te.YMin, &self.YLon)

	X := make([]float64, le)
	copy(X, te.X)
	dur := float64(self.Duration)
	NorMalizationOEM(X, &te.X[0], &dur)
	//fmt.Println(Y)

	Dx := float64(self.Len - 1)
	var Err, _y float64
	for _i, _x := range X {
		if _x > 1 {
			_x = 1
			//} else if _x < 0 {
			//	_x = 0
		}
		_y = self.WList[int(Rounding(_x*Dx))]
		Err += math.Abs(_y - Y[_i])
	}
	return self.Rate > (Err+j)/float64(le)

}
func (self *Templage) ContrastRate(te Templage) bool {

	le := len(te.Y)
	Y := make([]float64, le)
	copy(Y, te.Y)
	NorMalizationOEM(Y, &te.YMin, &self.YLon)

	X := make([]float64, le)
	copy(X, te.X)
	dur := float64(self.Duration)
	NorMalizationOEM(X, &te.X[0], &dur)
	//fmt.Println(Y)

	Dx := float64(self.Len - 1)

	var Err, _y float64
	for _i, _x := range X {
		if _x > 1 {
			_x = 1
			//} else if _x < 0 {
			//	_x = 0
		}
		_y = self.WList[int(Rounding(_x*Dx))]
		Err += math.Abs(_y - Y[_i])
	}
	return self.Rate > Err/float64(le)

}

func (self *Templage) Contrast(te Templage) bool {

	//	return self.ContrastRate(te)
	//	if te.CaLeft != self.CaLeft {
	//		te.Len = te.Duration/self.CaLeft.Scale
	//		if te.Duration%self.CaLeft.Scale>0 {
	//			te.Len++
	//		}
	//		te.CaLeft = self.CaLeft
	//		err := te.FittingInit()
	//		if err != nil {
	//			return false
	//		}
	//
	//	}
	if self.ContrastRate(te) || te.ContrastRate(*self) {
		//	te.Show(nil)
		//	self.Show(nil)
		return true
	}
	return false

}
func (self *Templage) GetKey() (key uint64) {
	key = uint64(self.Duration) << 1
	if self.Direction {
		key++
	}
	return key
}
func (self *Templage) FittingInit() error {

	Le := int(self.Len)
	self.X = make([]float64, Le)
	self.Y = make([]float64, Le)

	Yh := make([]float64, Le)
	Yl := make([]float64, Le)
	//	self.YMin
	var YMax float64 = 0
	self.YMin = self.Begin.GetMidAverage()
	i := 0
	self.CaLeft.SearchHandle(self.Begin.Time, func(can *CandlesMin) error {
		if (can.Time + self.CaLeft.Scale) > self.End.Time {
			//	fmt.Println(self.CaLeft.Name,"end")
			return fmt.Errorf("end")
		}
		self.X[i] = float64(can.Time)
		self.Y[i] = can.GetMidAverage()
		Yh[i] = can.Mid[2]
		Yl[i] = can.Mid[3]
		if can.Mid[2] > YMax {
			YMax = can.Mid[2]
		}
		if can.Mid[3] < self.YMin {
			self.YMin = can.Mid[3]
		}
		i++
		return nil
	})
	if i < 10 {
		return fmt.Errorf("err %d with double < %d", i, self.Len)
	}

	self.X = self.X[:i]
	X := make([]float64, i)
	copy(X, self.X)
	dur := float64(self.Duration)
	XMin := self.X[0]
	NorMalizationOEM(X, &XMin, &dur)

	self.Y = self.Y[:i]
	Y := make([]float64, i)
	copy(Y, self.Y)
	NorMalizationOEM(Y, &self.YMin, &self.YLon)

	Wlen := 0

	liblen := len(TemplagesLib)
	if liblen > 0 {
		M := liblen / 2
		tes := FindTemplages(TemplagesLib, self.GetKey(), &M, 0, liblen)
		if tes != nil {
			Wlen = tes.Wlen
		}
	}

	if Wlen == 0 {
		//fmt.Println("X",self.X)
		//fmt.Println("Y",Y)
		lastW := fitting.GetBastCols(X, Y)
		self.WList = nil
		for _, w := range lastW {
			//	fmt.Println(w)
			vl, r := CheckWeight(w, Le)
			if r {
				self.WList = vl
				self.Weight = w
				//fmt.Println(len(w))
				break
				//		return w
			}
		}
	} else {
		//		fmt.Println(Wlen)
		W := make([]float64, Wlen)
		if fitting.GetCurveFittingWeight(X, Y, Wlen, W) {
			vl, r := CheckWeight(W, Le)
			if r {
				self.WList = vl
				self.Weight = W
			}
		}
	}

	if self.WList != nil {
		self.SetRate(X, Yh[:i], Yl[:i])
		if len(self.Weight) < 4 {
			self.WList = nil
			self.Weight = nil
			return fmt.Errorf("len(self.Weight) <4")
		}
		//self.Show(Y)
		return nil
	}
	return fmt.Errorf("self.WList == nil")

}
func (self *Templage) SetRate(X, Yh, Yl []float64) {

	NorMalizationOEM(Yh, &self.YMin, &self.YLon)
	NorMalizationOEM(Yl, &self.YMin, &self.YLon)
	Dx := float64(self.Len - 1)
	//	var ErrH,ErrL,_y float64
	var ErrH, _y float64
	for i, _x := range X {
		_y = self.WList[int(Rounding(_x*Dx))]

		dh := math.Abs(_y - Yh[i])
		dl := math.Abs(_y - Yl[i])
		if dh > dl {
			ErrH += dh
		} else {
			ErrH += dl
		}

		//	ErrH += dh*dh
		//	ErrL += dl*dl

		//	ErrH += math.Abs(_y-Yh[i])
		//	ErrL += math.Abs(_y-Yl[i])
	}
	//	self.Rate = (ErrH+ErrL)/2/float64(len(self.X))

	self.Rate = ErrH / float64(len(X))

	//	fmt.Println(self.CoverageRate)

}

func Rounding(val float64) float64 {
	x, y := math.Modf(val)
	if y < 0.5 {
		return x
	}
	return x + 1
}
func CheckWeight(w []float64, dx int) (val []float64, r bool) {
	val = make([]float64, dx)
	var i, le float64 = 1, float64(dx)
	var j int = 0
	for ; i <= le; i = i + 1 {
		_x := i / le
		val[j] = 0
		for _j, _w := range w {
			val[j] += math.Pow(_x, float64(_j)) * _w
		}
		if val[j] > 1.2 || val[j] < -0.2 {
			return nil, false
		}
		j++
	}
	return val, true
}
func NorMalizationOEM(Val []float64, ValMin *float64, ValLong *float64) {

	if *ValMin == 0 || *ValLong == 0 {
		ValMax := Val[0]
		*ValMin = Val[0]
		for _, v := range Val[1:] {
			if v > ValMax {
				ValMax = v
			}
			if *ValMin > v {
				*ValMin = v
			}
		}
		*ValLong = ValMax - *ValMin
	}
	for i, v := range Val {
		Val[i] = (v - *ValMin) / *ValLong
		if Val[i] > 1 {
			Val[i] = 1
		}
	}

}
