package qnamegen

import "math/rand"

var tldList WeightedList

func init() {
	tldList = DefaultTLDList.ToWeightedList()
}

type TLDList map[string]uint
type WeightedList []string

func (w WeightedList) Shuffle(rounds int) {
	for i := 0; i < rounds; i++ {
		pos := rand.Intn(len(w))
		w[i], w[pos] = w[pos], w[i]
	}
}

func (t TLDList) ToWeightedList() WeightedList {
	s := make(WeightedList, 0)
	for tld, count := range t {
		for i := uint(0); i < count; i++ {
			s = append(s, tld)
		}
	}
	s.Shuffle(25)
	return s
}

var DefaultTLDList = TLDList{
	"com": 10,
	"co":  1,
	"de":  10,
	"eu":  1,
	"net": 1,
	"org": 1,
	"xyz": 3,
	"me":  5,
	"io":  5,
}
