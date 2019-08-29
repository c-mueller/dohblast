package qnamegen

import (
	"fmt"
	"math/rand"
)

func GenerateRandomQName() string {
	tld := tldList[rand.Intn(len(tldList))]
	body := GenerateRandomString(4)
	return fmt.Sprintf("%s.%s.", body, tld)
}
