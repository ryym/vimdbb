package vimch

import (
	"encoding/json"
	"strconv"
	"strings"
)

func EncodeMessage(id float64, payload interface{}) ([]byte, error) {
	message := [2]interface{}{}
	message[0] = id
	message[1] = payload
	return json.Marshal(message)
}

// m: [ID, Action, Payload] (e.g. [1, "Query", { "Query": "select 1" }])
func DecodeMessage(m string) (float64, string, string) {
	ret := strings.SplitN(m, ",", 2)
	id, _ := strconv.Atoi(ret[0][1:])

	ret = strings.SplitN(ret[1], ",", 2)
	var action string
	json.Unmarshal([]byte(ret[0][1:]), &action)

	payloadStr := ret[1][:len(ret[1])-2]

	return float64(id), action, payloadStr
}

func DecodePayload(payload string, v interface{}) error {
	return json.Unmarshal([]byte(payload), v)
}
