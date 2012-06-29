package gcm

import (
	"testing"
)

func TestMessage_Send(t *testing.T) {
	message := &Message{
		RegistrationIDs: []RegistrationID{
			RegistrationID("TEST ID"),
		},
		CollapseKey: "Testing",
		Data: map[string]interface{}{
			"mykey":      "myval",
			"otherthing": "otherval",
		},
	}
	resp, err := message.Send("mytestkey")
	if err != nil {
		t.Errorf(err.Error())
	}
	if resp.Failure != 0 {
		t.Errorf("%v messages could not be processed.", resp.Failure)
		for _, msg := range resp.Results {
			if msg.Error != "" {
				t.Logf(msg.Error)
			}
		}
	}
}
