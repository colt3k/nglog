package sched

import (
	"testing"
)

// Min, Hour, Day of Month, Month, Day of Week
var data = []string{
	"5 11,19 * * *",
	"20/5,*/8,10-12 20,21 * * *",
	"20/5,*/8,10-12 20,21 1,31 10-12 *",
}

func TestCronSched_Parse(t *testing.T) {
	// ALL UTC TIMES
	for _, d := range data {
		tmp := ParseCronSched(d)
		tmp.nextExecution()
	}
}
