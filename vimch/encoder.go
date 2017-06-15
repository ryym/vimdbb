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

// m: [ID, Payload] (e.g. [1, { "Query": "select 1" }])
func DecodeMessage(m string) vimdbb.Action {
	ret := strings.SplitN(m, ",", 2)
	id, _ := strconv.Atoi(ret[0][1:])

	payloadStr := ret[1][:len(ret[1])-1]
	payload := vimdbb.Payload{}
	json.Unmarshal([]byte(payloadStr), &payload)

	return vimdbb.Action{Id: float64(id), Payload: payload}
}
