package sched

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

/**
Cron Parser

Field		Allowed Values		Allowed Special Chars
Minute		0-59				, - * /
Hour		0-23				, - * /
DayOfMonth	1-31				, - * ? / L W
Month		1-12 or Jan-Dec		, - * /
DayOfWeek	1-7 or SUN-SAT		, - * ? / L #

Special Char Vals
,	value list separator		i.e. 1,5,10
- 	range of values				i.e. 5-10
*	any value					i.e. ALL VALUES for this field
/	step values					i.e. *\/5	every 5 of field
*/

type CronField int

const (
	Minute CronField = iota + 1
	Hour
	DoM
	Month
	DoW
)

var locations = [...]string{
	"Minute", "Hour", "Day of Month", "Month", "Day of Week",
}
var maxRanges = [...]int{
	59, 23, 31, 12, 7,
}
var minRanges = [...]int{
	0, 0, 1, 1, 1,
}

func (f CronField) String() string {
	return locations[f-1]
}
func (f CronField) Max() int {
	return maxRanges[f-1]
}
func (f CronField) Min() int {
	return minRanges[f-1]
}

type CronVal struct {
	expr string
	fld  CronField
	vals []int
}
type CronEntry struct {
	vals []*CronVal
	next []time.Time // after processed move to next time.Time
}

type CronSched struct {
	minute     string
	hour       string
	dayOfMonth string
	month      string
	dayOfWeek  string
	entry      *CronEntry
}

func ParseCronSched(schedule string) *CronSched {
	t := new(CronSched)
	t.parse(schedule)
	return t
}

func (c *CronSched) parse(schedule string) error {
	if len(schedule) == 0 {
		return fmt.Errorf("Empty spec string")
	}

	fields := strings.Fields(schedule)

	// Fill in missing fields
	//fields = expandFields(fields, c.options)

	cvs := make([]*CronVal, 0)
	for i, d := range fields {
		//log.Printf("%d %s val %s", (i + 1), CronField(i+1).String(), d)
		t := &CronVal{fld: CronField(i + 1), expr: d}
		cvs = append(cvs, t)
	}
	ce := &CronEntry{vals: cvs}
	c.entry = ce
	c.expandFields(ce)
	log.Println()
	//spew.Dump(ce)
	return nil
}

// split on commas
// split on /		(set minute past hour)
// split on -		(range of minutes past hour inclusive of declared minutes)
// "*/5,*/8,10-12 2 * * *"
func (c *CronSched) expandFields(ce *CronEntry) {
	now := time.Now()
	var (
		min, max, step int
	)
	for _, d := range c.entry.vals {
		min = CronField(d.fld).Min()
		max = CronField(d.fld).Max()
		step = 0
		sects := strings.Split(d.expr, ",")
		//log.Println(sects)
		vals := make([]int, 0)
		for _, j := range sects {

			//log.Println("Processing:", j)
			slsh := strings.Split(j, "/")
			if len(slsh) > 1 { // */5 (every 5th min) or 10/5 (every 5th minute from 10 to 59)
				dsh := strings.Split(slsh[0], "-")
				if len(dsh) > 1 {
					// max is the second part and min the first
					min, _ = strconv.Atoi(dsh[0])
					max, _ = strconv.Atoi(dsh[1])
					step, _ = strconv.Atoi(slsh[1])
				} else if slsh[0] != "*" {
					// no dash so single number or *
					min, _ = strconv.Atoi(slsh[0])
					max = CronField(d.fld).Max()
					if d.fld == DoM {
						max = c.daysIn(now.Month(), now.Year())
					}
					step, _ = strconv.Atoi(slsh[1])
				} else if slsh[0] == "*" {
					// no dash so single number or *
					min = CronField(d.fld).Min()
					max = CronField(d.fld).Max()
					if d.fld == DoM {
						max = c.daysIn(now.Month(), now.Year())
					}
					step, _ = strconv.Atoi(slsh[1])
				}
			} else {
				// no slash check for dash
				dsh := strings.Split(slsh[0], "-")
				if len(dsh) > 1 {
					// max is the second part and min the first
					min, _ = strconv.Atoi(dsh[0])
					max, _ = strconv.Atoi(dsh[1])
					step = 0
				} else if slsh[0] != "*" {
					// no dash so single number or *
					step = 0
					min, _ = strconv.Atoi(slsh[0])
					if d.fld == DoM {
						dys := c.daysIn(now.Month(), now.Year())
						if min-dys == 1 {
							min = dys
						}
					}
					max = min
				}
			}
			//log.Printf("%s Expr %s Min %d Max %d Step %d", CronField(d.fld).String(), j, min, max, step)
			// vals example:

			//Build up values for each CronVal
			//d.vals
			if strings.Index(j, "*") < 0 && strings.Index(j, "/") < 0 && strings.Index(j, "-") < 0 {
				exp, _ := strconv.Atoi(j)
				if d.fld == DoM {
					dys := c.daysIn(now.Month(), now.Year())
					if exp-dys == 1 {
						exp = dys
					}
				}
				if exp == min && exp == max {
					//use original value there is no min/max for this field
					vals = append(vals, exp)
					sort.Ints(vals)
					d.vals = vals
				}
			} else {

				if step == 0 {
					step = 1
				}
				for i := min; i <= max; i += step {
					vals = append(vals, i)
				}
				sort.Ints(vals)
				d.vals = vals
			}
		}
	}
}

func (c *CronSched) nextExecution() {
	periods := make([]time.Time, 0)
	now := time.Now().UTC()

	minutes := c.entry.vals[0]
	hours := c.entry.vals[1]
	doms := c.entry.vals[2]
	months := c.entry.vals[3]
	dows := c.entry.vals[4]

	var byt bytes.Buffer
	for _, month := range months.vals {

		maxDays := c.daysIn(time.Month(month), now.Year())
		for _, dom := range doms.vals {
			for _, hour := range hours.vals {
				for _, min := range minutes.vals {
					byt.WriteString(strconv.Itoa(now.Year()))
					byt.WriteString("-")
					byt.WriteString(PrefixDigit(month))
					byt.WriteString("-")
					byt.WriteString(PrefixDigit(dom))
					byt.WriteString("T")
					byt.WriteString(PrefixDigit(hour))
					byt.WriteString(":")
					byt.WriteString(PrefixDigit(min))
					byt.WriteString(":00Z")

					if dom > maxDays {
						byt.Reset()
						continue
					}
					tmp, err := time.Parse(time.RFC3339, byt.String())
					if err != nil {
						log.Fatalf("issue parsing time via string\n%+v",err)
					}
					periods = append(periods, tmp)
					byt.Reset()
				}
			}
		}
	}

	b, _ := in_array(int(now.Weekday()), dows.vals)
	//fmt.Printf("Weekday: %d, is in %v : %v \n", now.Weekday(), dows.vals, b)
	for i, d := range periods {
		// Narrow down to current month and day
		if i < len(periods)-1 && d.Month() >= now.Month() && d.Day() >= now.Day() {
			//fmt.Printf("%v before %v\n", now, d)
			// Equal Time or After Current Time and before Next Time
			if d == now || now.Before(d) {
				if b {
					fmt.Printf("Next Trigger: %v on %v\n", d, now.Weekday())
					break
				}
			}
		} else {
			// if after now and there are no more than calc for next year
		}
	}

	//spew.Dump(periods)
}

func PrefixDigit(dig int) string {
	if dig <= 9 {
		return "0" + strconv.Itoa(dig)
	}
	return strconv.Itoa(dig)
}

func (c *CronSched) Trigger() bool {
	return false
}

func (c *CronSched) daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
func in_array(val interface{}, array interface{}) (exists bool, index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
