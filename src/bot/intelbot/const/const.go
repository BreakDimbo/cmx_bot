package con

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const RFC3339local = "2006-01-02T15:04:05Z"

const SingleChineseByte = 3

const (
	ScopeTypePublic = "public"
	ScopeTypeLocal  = "local"
)

var Emoji []string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	dir, err := os.Getwd()
	check(err)
	emojFile := fmt.Sprintf("%s/config/emoj.txt", dir)
	dat, err := ioutil.ReadFile(emojFile)
	check(err)
	Emoji = strings.Split(string(dat), "\n")
}
