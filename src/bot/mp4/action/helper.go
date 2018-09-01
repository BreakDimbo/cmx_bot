package action

import (
	"bot/bredis"
	"bot/log"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func genTTS(s string) (string, error) {
	content := make(map[string]string)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	filename := hex.EncodeToString(bs)
	content["id"] = filename
	content["txt"] = s
	jsonMap, err := json.Marshal(content)
	if err != nil {
		log.SLogger.Errorf("marshal to json error: %v", err)
		return "", err
	}
	err = bredis.Client.Publish("tts", jsonMap).Err()
	if err != nil {
		log.SLogger.Errorf("pub to redis error: %v", err)
		return "", err
	}
	time.Sleep(1 * time.Second)
	return filename, nil
}

func convertToMp4(fn string) (io.Reader, error) {
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	output := strings.Replace(fn, ".wav", ".mp4", -1)
	cmd := exec.Command("ffmpeg", "-i", fn, "-vn", "-acodec", "aac", "-strict", "-2", output)
	// 运行命令
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err := os.Open(output)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func calSig(t string) string {
	const signature = "a7d040cee7bed322a188c9ec7fd9b8b8b34ac893"
	h := sha1.New()
	h.Write([]byte(signature + string(t)))
	bs := h.Sum(nil)
	bsStr := hex.EncodeToString(bs)
	return bsStr
}

func checkErr(err error, res http.ResponseWriter) {
	if err != nil {
		res.Write([]byte(err.Error()))
	}
}
