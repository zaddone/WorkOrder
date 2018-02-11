package request

import (
	"fmt"
	"os"
	//	"math"
	"io/ioutil"
	"net/http"
	"net/url"
	//	"net"
	//	"golang.org/x/net/proxy"
	"encoding/json"
	"flag"
	"path/filepath"
	"time"

	"strconv"
	"strings"
	//	"bytes"
	"io"
)

var (
	Accounts []*Account
	Client   *http.Client
	Header   http.Header
	Instr    *Instrument

	Account_ID    = flag.String("accountid", "101-011-2471429-001", "Account ID")
	Authorization = flag.String("auth", "07b98f22c1eafe0359d287f68189a6db-2ed5b9d4ee662581965275d7a6dbcb58", "Auth")
	//	Proxy = flag.String("p","192.168.1.70:1081","proxy")
	Proxy     = flag.String("proxy", "", "proxy")
	LogFile   = flag.String("L", "LogInfo.log", "LogInfo")
	Host      = flag.String("h", "https://api-fxpractice.oanda.com/v3", "host")
	InsName   = flag.String("n", "EUR_JPY", "INS NAME")
	BEGINTIME = flag.String("begintime", "2009-01-01T00:00:00", "2009-01-01T00:00:00")
)

func init() {

	flag.Parse()
	Header = make(http.Header)
	Header.Add("Authorization", fmt.Sprintf("Bearer %s", *Authorization))
	Header.Add("Connection", "Keep-Alive")
	Header.Add("Accept-Datetime-Format", "UNIX")
	Header.Add("Content-type", "application/json")

	if *Proxy == "" {
		Client = new(http.Client)
	} else {
		panic(0)
		//	dialer, err := proxy.SOCKS5("tcp",*Proxy,
		//	    nil,
		//	    &net.Dialer {
		//	        Timeout: 30 * time.Second,
		//	        KeepAlive: 30 * time.Second,
		//	    },
		//	)
		//	if err != nil {
		//		panic(err)
		//	}
		//	transport := &http.Transport{
		//	    Proxy: nil,
		//	    Dial: dialer.Dial,
		//	    TLSHandshakeTimeout: 10 * time.Second,
		//	}
		//	Client = &http.Client{Transport:transport}
	}
	err := InitAccounts(false)
	if err != nil {
		panic(err)
	}
	Ins, err := GetInstruments()
	if err != nil {
		panic(err)
	}
	Instr = Ins[*InsName]
	if Instr == nil {
		panic("instr == nil")
	}

}
func ClientPut(path string, val io.Reader, da interface{}) error {

	Req, err := http.NewRequest("PUT", path, val)
	if err != nil {
		return err
	}
	Req.Header = Header
	res, err := Client.Do(Req)
	if err != nil {
		if res != nil {
			b, err := ioutil.ReadAll(res.Body)
			fmt.Println(string(b), err)
		}
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(da)

}
func ClientPost(path string, val io.Reader, da interface{}) error {

	Req, err := http.NewRequest("POST", path, val)
	if err != nil {
		return err
	}
	Req.Header = Header
	res, err := Client.Do(Req)
	if err != nil {
		if res != nil {
			b, err := ioutil.ReadAll(res.Body)
			fmt.Println(string(b), err)
		}
		return err
	}
	defer res.Body.Close()

	//	b,err:=ioutil.ReadAll(res.Body)
	//	fmt.Println(string(b),err)
	//	return json.Unmarshal(b,da)

	return json.NewDecoder(res.Body).Decode(da)

}

func ClientDO(path string, da interface{}) error {
	Req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	Req.Header = Header
	res, err := Client.Do(Req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		//	b,err:=ioutil.ReadAll(res.Body)
		//	fmt.Println(string(b),err)
		return fmt.Errorf("status code %d %s", res.StatusCode, path)
	}
	return json.NewDecoder(res.Body).Decode(da)

	//	b,err:=ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		return err
	//	}
	//	return json.Unmarshal(b,da)

}

func GetInstruments() (ins map[string]*Instrument, err error) {
	//func GetInstruments(n int) (ins map[string]*Instrument,err error) {

	var Nacc *Account
	for _, acc := range Accounts {
		if acc.Id == *Account_ID {
			Nacc = acc
			if len(acc.Instruments) != 0 {
				return acc.Instruments, nil
			}
			break
		}
	}
	if Nacc == nil {
		err := InitAccounts(true)
		if err != nil {
			panic(err)
		}
		return GetInstruments()
	}
	//	if len(Accounts[n].Instruments) != 0 {
	//		return Accounts[n].Instruments,nil
	//	}

	path := GetAccountPath()
	path += "/instruments"
	da := make(map[string]interface{})
	err = ClientDO(path, &da)
	if err != nil {
		return nil, err
	}
	in := da["instruments"].([]interface{})
	ins = make(map[string]*Instrument)
	for _, n := range in {
		//		fmt.Println(n.(InstrumentTmp).Name)
		in := new(Instrument)
		in.Init(n.(map[string]interface{}))

		ins[in.Name] = in
	}
	Nacc.Instruments = ins
	err = SaveAccounts()
	return ins, err

}
func GetAccountPath() string {

	return fmt.Sprintf("%s/accounts/%s", *Host, *Account_ID)

}
func InitAccounts(isU bool) (err error) {
	if !isU {
		err = ReadAccounts()
		if err == nil {
			return nil
		}
	}
	da := make(map[string][]*Account)
	err = ClientDO(*Host+"/accounts", &da)
	if err != nil {
		return err
	}
	Accounts = da["accounts"]
	//	L := len(self.Accounts)
	if len(Accounts) == 0 {
		fmt.Errorf("accounts == nil")
	}
	return SaveAccounts()

}
func SaveAccounts() error {
	f, err := os.OpenFile(*LogFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_SYNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := json.Marshal(Accounts)
	if err != nil {
		return err
	}
	_, err = f.Write(d)
	if err != nil {
		return err
	}
	return nil
}
func ReadAccounts() error {
	fi, err := os.Stat(*LogFile)
	if err != nil {
		return err
	}
	data := make([]byte, fi.Size())
	f, err := os.Open(*LogFile)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Read(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("%d %d", n, len(data))
	}
	return json.Unmarshal(data, &(Accounts))
}

func Down(from, to int64, gran string, Handle func(*CandlesMin)) {

	gr := GranularityMap[gran]
	var file *os.File = nil
	var LogFile string = ""
	var err error
	var Begin time.Time
	for {
		//		fmt.Println(time.Unix(from,0).UTC(),gran)
		err = GetCandlesHandle(Instr.Name, gran, from, 500, func(c interface{}) error {
			can := new(CandlesMin)
			can.Init(c.(map[string]interface{}))
			Begin = can.GetTime()
			Handle(can)
			path := filepath.Join(Instr.Name, gran, fmt.Sprintf("%d", Begin.Year()))
			_, err := os.Stat(path)
			if err != nil {
				os.MkdirAll(path, 0777)
			}
			path = filepath.Join(path, Begin.Format("20060102"))
			if file == nil {
				LogFile = path
				file, err = os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
				if err != nil {
					panic(err)
				}
			} else if LogFile != path {
				//	fmt.Println(path)
				file.Close()
				LogFile = path
				file, err = os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
				if err != nil {
					panic(err)
				}
			}
			file.WriteString(can.String())
			return nil
		})
		if err != nil {
			//	fmt.Println("err:",err)
			aft := from - time.Now().Unix()
			if aft > 0 {
				<-time.After(time.Second * time.Duration(aft))
			} else {
				<-time.After(time.Second * 3)
			}
		} else {
			from = Begin.Unix() + gr
			if to > 0 && from >= to {
				break
			}
			aft := from - time.Now().Unix()
			if aft > 0 {
				<-time.After(time.Second * time.Duration(aft))
			}
		}
	}

}
func GetCandlesHandle(Ins_name, granularity string, from, count int64, f func(c interface{}) error) (err error) {

	path := fmt.Sprintf("%s/instruments/%s/candles?", *Host, Ins_name)
	uv := url.Values{}
	uv.Add("granularity", granularity)
	uv.Add("price", "M")
	uv.Add("from", fmt.Sprintf("%d", from))
	uv.Add("count", fmt.Sprintf("%d", count))
	path += uv.Encode()
	//	fmt.Println(path)
	da := make(map[string]interface{})
	err = ClientDO(path, &da)
	if err != nil {
		return err
	}
	ca := da["candles"].([]interface{})
	lc := len(ca)
	if lc == 0 {
		return fmt.Errorf("candles len = 0")
	}
	//	can = make([]*CandlesMin,lc)
	for _, c := range ca {
		er := f(c)
		if er != nil {
			fmt.Println(er)
			break
		}
		//		can[i] = new(CandlesMin)
		//		can[i].Init(c.(map[string]interface{}))
	}
	return nil

}

type Account struct {
	Id          string                 `json:"id"`
	Tags        []string               `json:"tags"`
	Instruments map[string]*Instrument `json:"instruments"`
}

//func GetAccountID() string {
//	return Accounts[*Account_ID].Id
//}
type Instrument struct {
	Name string

	DisplayPrecision float64
	MarginRate       float64

	MaximumOrderUnits           float64
	MaximumPositionSize         float64
	MaximumTrailingStopDistance float64

	MinimumTradeSize            float64
	MinimumTrailingStopDistance float64

	PipLocation         float64
	TradeUnitsPrecision float64
	Type                string
}

func (self *Instrument) Init(tmp map[string]interface{}) (err error) {
	self.Name = tmp["name"].(string)
	self.PipLocation = tmp["pipLocation"].(float64)
	self.TradeUnitsPrecision = tmp["tradeUnitsPrecision"].(float64)
	self.Type = tmp["type"].(string)
	self.DisplayPrecision = tmp["displayPrecision"].(float64)

	self.MarginRate, err = strconv.ParseFloat(tmp["marginRate"].(string), 64)
	if err != nil {
		return err
	}
	self.MaximumOrderUnits, err = strconv.ParseFloat(tmp["maximumOrderUnits"].(string), 64)
	if err != nil {
		return err
	}
	self.MaximumPositionSize, err = strconv.ParseFloat(tmp["maximumPositionSize"].(string), 64)
	if err != nil {
		return err
	}
	self.MaximumTrailingStopDistance, err = strconv.ParseFloat(tmp["maximumTrailingStopDistance"].(string), 64)
	if err != nil {
		return err
	}
	self.MinimumTradeSize, err = strconv.ParseFloat(tmp["minimumTradeSize"].(string), 64)
	if err != nil {
		return err
	}
	self.MinimumTrailingStopDistance, err = strconv.ParseFloat(tmp["minimumTrailingStopDistance"].(string), 64)
	if err != nil {
		return err
	}

	return nil
}

type CandlesMin struct {
	Mid    [4]float64
	Time   int64
	Volume float64
	Val    float64
}

func (self *CandlesMin) GetMidLong() float64 {
	return self.Mid[2] - self.Mid[3]
}

func (self *CandlesMin) GetMidAverage() float64 {
	if self.Val == 0 {
		var sum float64 = 0
		for _, m := range self.Mid {
			sum += m
		}
		self.Val = sum / 4
	}
	return self.Val
}

func (self *CandlesMin) GetInOut() bool {
	return self.Mid[0] < self.Mid[1]
}

func (self *CandlesMin) GetTime() time.Time {
	return time.Unix(self.Time, 0).UTC()

}
func (self *CandlesMin) Show() {
	fmt.Printf("%.6f %.6f %.6f %.6f %s %.6f\r\n", self.Mid[0], self.Mid[1], self.Mid[2], self.Mid[3], time.Unix(self.Time, 0).String(), self.Volume)
}
func (self *CandlesMin) Load(str string) (err error) {
	strl := strings.Split(str, " ")
	self.Mid[0], err = strconv.ParseFloat(strl[0], 64)
	if err != nil {
		return err
	}
	self.Mid[1], err = strconv.ParseFloat(strl[1], 64)
	if err != nil {
		return err
	}
	self.Mid[2], err = strconv.ParseFloat(strl[2], 64)
	if err != nil {
		return err
	}
	self.Mid[3], err = strconv.ParseFloat(strl[3], 64)
	if err != nil {
		return err
	}
	Ti, err := strconv.Atoi(strl[4])
	if err != nil {
		return err
	}
	self.Time = int64(Ti)
	self.Volume, err = strconv.ParseFloat(strl[5], 64)
	if err != nil {
		return err
	}
	return nil
}
func (self *CandlesMin) String() string {
	return fmt.Sprintf("%.5f %.5f %.5f %.5f %d %.5f\r\n", self.Mid[0], self.Mid[1], self.Mid[2], self.Mid[3], self.Time, self.Volume)
}
func (self *CandlesMin) Init(tmp map[string]interface{}) (err error) {
	Mid := tmp["mid"].(map[string]interface{})
	if Mid != nil {
		self.Mid[0], err = strconv.ParseFloat(Mid["o"].(string), 64)
		if err != nil {
			return err
		}
		self.Mid[1], err = strconv.ParseFloat(Mid["c"].(string), 64)
		if err != nil {
			return err
		}
		self.Mid[2], err = strconv.ParseFloat(Mid["h"].(string), 64)
		if err != nil {
			return err
		}
		self.Mid[3], err = strconv.ParseFloat(Mid["l"].(string), 64)
		if err != nil {
			return err
		}
	}
	self.Volume = tmp["volume"].(float64)
	ti, err := strconv.ParseFloat(tmp["time"].(string), 64)
	if err != nil {
		return err
	}
	self.Time = int64(ti)
	return nil
}
