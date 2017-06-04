package main

import (
    "database/sql"
	"encoding/json"
	"strings"
	"strconv"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	// "io"
	// "bytes"
	"bufio"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Println("server started")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		defer conn.Close()

		fmt.Println("accepted")
		go handle(conn)
	}

}

type Payload struct {
	Query string
}

type Action struct {
	Id float64
	Payload Payload
}
type Result struct {
	Rows string
}

func handle(conn net.Conn) {
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		fmt.Println("scanned")
		message := sc.Text()

		// conn.Write([]byte(message + "\n"))
		fmt.Println(message)

		action := decodeMessage(message)
		// fmt.Println(action.Payload.Query)

		result := run(action.Payload.Query)
		// resJson, _ := json.Marshal(Result{ Rows: result })

		results := [2]interface{}{}
		results[0] = action.Id
		results[1] = Result{ Rows: result }

		retJson, err := json.Marshal(results)
		if err != nil {
			panic(err.Error())
		}
		conn.Write(retJson)
		// fmt.Println(result)
	}

	fmt.Println("END!")
}

func decodeMessage(m string) Action {
	ret := strings.SplitN(m, ",", 2)
	id, _ := strconv.Atoi(ret[0][1:])

	payloadStr := ret[1][:len(ret[1])-1]
	payload := Payload{}
	json.Unmarshal([]byte(payloadStr), &payload)

	return Action{ Id: float64(id), Payload: payload }
}

func run(queryStr string) string {
    // Open database connection
    db, err := sql.Open("mysql", "root:root@/sample")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    // rows, err := db.Query("SELECT * FROM customer")
    rows, err := db.Query(queryStr)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Make a slice for the values
    values := make([]sql.RawBytes, len(columns))

    // rows.Scan wants '[]interface{}' as an argument, so we must copy the
    // references into such a slice
    // See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

	results := ""

	for i, col := range columns {
		if i != 0 {
			results += " | "
		}
		results += col
	}

	results += "\n---------------\n"

    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }

		var value string
		for i ,col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
			if i != 0 {
				results += " | "
			}
			results += value
		}
		results += "\n"
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

	return results
}

