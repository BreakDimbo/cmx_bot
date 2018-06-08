package mastodon

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var stopwords map[string]bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadStopWord() {
	dir, err := os.Getwd()
	check(err)

	stopwordFileDir := fmt.Sprintf("%s/config/stopword.txt", dir)
	dat, err := ioutil.ReadFile(stopwordFileDir)
	check(err)

	wordsSlice := strings.Split(string(dat), "\n")
	stopwords = make(map[string]bool)
	for _, word := range wordsSlice {
		stopwords[word] = true
	}
}
