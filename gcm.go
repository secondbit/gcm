package gcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegistrationID string

type Message struct {
	RegistrationIDs []RegistrationID       `json:"registration_ids"`
	CollapseKey     string                 `json:"collapse_key,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	DelayWhileIdle  bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive      int                    `json:"time_to_live",omitempty`
}

type Response struct {
	MulticastID  int64     `json:"multicast_id"`
	Success      int       `json:"success"`
	Failure      int       `json:"failure"`
	CanonicalIDs int       `json:"canonical_ids"`
	Results      []*Result `json:"results"`
}

type Result struct {
	MessageID      int64  `json:"message_id"`
	RegistrationID int64  `json:"registration_id"`
	Error          string `json:"error"`
}

func (m *Message) Send(key string) (resp *Response, err error) {
	var msg_body []byte
	msg_body, err = json.Marshal(m)
	if err != nil {
		return
	}
	msg_reader := bytes.NewReader(msg_body)
	var req *http.Request
	req, err = http.NewRequest("POST", "https://android.googleapis.com/gcm/send", msg_reader)
	if err != nil {
		return
	}
	req.Header.Set(http.CanonicalHeaderKey("Authorization"), "key="+key)
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")
	client := &http.Client{}
	var r *http.Response
	r, err = client.Do(req)
	if err != nil {
		return
	}
	switch r.StatusCode {
	case 200:
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		resp = &Response{}
		err = json.Unmarshal(body, resp)
		return
	case 503:
		fmt.Println("503 Error: %v", r.Header.Get("Retry-After"))
		return
	default:
		fmt.Println(r.Status)
	}
	return
}
