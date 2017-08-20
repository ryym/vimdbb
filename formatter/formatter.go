package formatter

import (
	"github.com/ryym/vimdbb/mysql"
)

func ResultToString(ret *mysql.Result) string {
	lens := getMaxLengths(ret)
	s := formatRow(lens, ret.Columns) + "\n"

	for _, l := range lens {
		for i := 0; i < l+3; i++ {
			s += "-"
		}
	}
	s += "\n"

	for _, row := range ret.Rows {
		s += formatRow(lens, row) + "\n"
	}

	return s
}

func formatRow(lens []int, values []string) string {
	s := ""
	for iv, v := range values {
		s += "| " + v
		for i := 0; i <= lens[iv]-len(v); i++ {
			s += " "
		}
	}
	return s
}

func getMaxLengths(ret *mysql.Result) []int {
	lens := make([]int, len(ret.Columns))

	for _, row := range ret.Rows {
		for i, v := range row {
			col := ret.Columns[i]
			if lens[i] < len(col) {
				lens[i] = len(col)
			}
			if lens[i] < len(v) {
				lens[i] = len(v)
			}
		}
	}

	return lens
}
