package vimch

import (
	"encoding/json"
	"github.com/ryym/vimdbb"
	"strconv"
	"strings"
)

func EncodeMessage(id float64, payload interface{}) ([]byte, error) {
	message := [2]interface{}{}
	message[0] = id
	message[1] = payload
	return json.Marshal(message)
}

// m: [ID, Type, Payload] (e.g. [1, "Query", { "Query": "select 1" }])
func DecodeMessage(m string) *vimdbb.Message {
	ret := strings.SplitN(m, ",", 2)
	id, _ := strconv.Atoi(ret[0][1:])

	ret = strings.SplitN(ret[1], ",", 2)
	var command string
	json.Unmarshal([]byte(ret[0][1:]), &command)

	payloadStr := ret[1][:len(ret[1])-2]

	return &vimdbb.Message{
		Id:      float64(id),
		Command: command,
		Payload: payloadStr,
	}
}

func DecodePayload(payload string, v interface{}) error {
	return json.Unmarshal([]byte(payload), v)
}
