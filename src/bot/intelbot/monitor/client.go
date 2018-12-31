package monitor

import pusher "github.com/pusher/pusher-http-go"

var Client = pusher.Client{
	AppId:   "681531",
	Key:     "ba844c624003f02c6c0f",
	Secret:  "78d4a04ab77b874f4116",
	Cluster: "ap1",
	Secure:  true,
}

// visitsData is a struct
type VisitsData struct {
	Count int
	Time  string
}
