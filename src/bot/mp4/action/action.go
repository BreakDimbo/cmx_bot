package action

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetMp4TTS(res http.ResponseWriter, req *http.Request) {
	secrete := req.URL.Query().Get("secrete")
	t := req.URL.Query().Get("time")
	if secrete != calSig(t) {
		checkErr(fmt.Errorf("seceret err"), res)
	}

	timestamp, err := strconv.ParseInt(t, 0, 64)
	checkErr(err, res)

	if timestamp-time.Now().Unix() > 10 {
		checkErr(fmt.Errorf("time now right error"), res)
	}

	result, err := ioutil.ReadAll(req.Body)
	checkErr(err, res)

	defer req.Body.Close()

	content := string(result)

	filename, err := genTTS(content)
	checkErr(err, res)

	f, err := convertToMp4(filename)
	checkErr(err, res)

	buf, err := ioutil.ReadAll(f)
	checkErr(err, res)

	res.Write(buf)
}
