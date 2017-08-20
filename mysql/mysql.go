package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Result struct {
	Columns []string
	Rows    [][]string
}

type DB interface {
	Query(query string) (*Result, error)
	Close() error
}

type myDB struct {
	db *sql.DB
}

func Open(url string) (DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return &myDB{db: db}, nil
}

func (mdb *myDB) Close() error {
	return mdb.db.Close()
}

func (mdb *myDB) Query(query string) (*Result, error) {
	return fetchRows(mdb.db, query)
}

func fetchRows(db *sql.DB, query string) (*Result, error) {
	res, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, err := res.Columns()
	if err != nil {
		return nil, err
	}

	// res.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rows := make([][]string, 0)
	for res.Next() {
		err = res.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		row := make([]string, len(columns))
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			row[i] = value
		}
		rows = append(rows, row)
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return &Result{Columns: columns, Rows: rows}, nil
}
