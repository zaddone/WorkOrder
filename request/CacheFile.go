package request

import (
	//	"fmt"
	"bufio"
	"os"
)

type CacheFile struct {
	Can    chan *CandlesMin
	Fi     os.FileInfo
	Path   string
	EndCan *CandlesMin
}

func (self *CacheFile) Init(path string, fi os.FileInfo, Max int) (err error) {

	self.Path = path
	self.Fi = fi
	self.Can = make(chan *CandlesMin, Max)
	var fe *os.File
	fe, err = os.Open(path)
	if err != nil {
		return err
	}
	defer fe.Close()
	r := bufio.NewReader(fe)
	for {
		db, _, e := r.ReadLine()
		if e != nil {
			break
		}
		self.EndCan = new(CandlesMin)
		self.EndCan.Load(string(db))
		self.Can <- self.EndCan
	}
	close(self.Can)
	return nil

}
