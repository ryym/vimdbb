package server

import (
	"errors"
	"github.com/ryym/vimdbb"
	"github.com/ryym/vimdbb/formatter"
	"github.com/ryym/vimdbb/mysql"
	"github.com/ryym/vimdbb/vimch"
)

func handleUserMessage(mes string) ([]byte, error) {
	m := vimch.DecodeMessage(mes)
	result, err := handleUserCommand(m)

	if err != nil {
		encm, encErr := encodeResp(m.ID, "ERR", err.Error())
		if encErr != nil {
			return nil, encErr
		}
		return encm, err
	}

	if result != nil {
		return encodeResp(m.ID, m.Command, result)
	}

	return nil, nil
}

func encodeResp(id float64, command string, result interface{}) ([]byte, error) {
	res := make(map[string]interface{})
	res["Command"] = command
	res["Result"] = result
	return vimch.EncodeMessage(id, res)
}

func handleUserCommand(m *vimdbb.Message) (interface{}, error) {
	switch m.Command {
	case "Query":
		queryP := vimdbb.QueryPayload{}
		vimch.DecodePayload(m.Payload, &queryP)
		return handleQuery(queryP)
	}

	return nil, errors.New("Unknown command " + m.Command)
}

func handleQuery(p vimdbb.QueryPayload) (*vimdbb.Result, error) {
	db, err := mysql.Open(p.ConnectionURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	queryResult, err := db.Query(p.Query)
	if err != nil {
		return nil, err
	}

	rows := formatter.ResultToString(queryResult)
	result := vimdbb.Result{QueryID: p.QueryID, Rows: rows}
	return &result, nil
}
