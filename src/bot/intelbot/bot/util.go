package bot

import (
	zlog "bot/log"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/sethgrid/pester"
)

func downloadPic(remoteURL string) (localURL string) {
	// Create the file
	filenameSlice := strings.Split(remoteURL, "/")
	filename := filenameSlice[len(filenameSlice)-1]
	filename = strings.Split(filename, "?")[0]
	filepath := fmt.Sprintf("/Users/break/cmx_pic/%s", filename)
	out, err := os.Create(filepath)
	if err != nil {
		zlog.SLogger.Error(err)
	}
	defer out.Close()

	// Get the data

	client := pester.New()
	client.Concurrency = 3
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true

	resp, err := client.Get(remoteURL)
	if err != nil {
		zlog.SLogger.Error(err)
	}

	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		zlog.SLogger.Error(err)
	}

	return filepath
}

func validLengthFilter(content string, validLength int) string {
	runeCount := utf8.RuneCountInString(content)
	if runeCount > validLength {
		return string([]rune(content)[:validLength])
	}
	return content
}
