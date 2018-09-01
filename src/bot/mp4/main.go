package main

import (
	"bot/mp4/action"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/mp4_tts", action.GetMp4TTS)
	log.Fatal(http.ListenAndServe("0.0.0.0:5438", nil))
}
