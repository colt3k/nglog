package ng

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/colt3k/utils/archive/gz"
)

type TriggerPolicy interface {
	Rotate(string, Strategy) (bool,error)
}

// tests using cron
type CronTriggerPolicy struct {
	schedule    string //	i.e. 0 0 * * * ?
	evalOnStart bool
}

// checks while running or on start
type SizeTriggerPolicy struct {
	maxSize     int64
	evalOnStart bool
}

// tests on interval of hours with a possible max random delay of X
type TimeTriggerPolicy struct {
	interval       int // based on most specific date pattern field i.e. hours (default)
	maxRandomDelay int // 0 no delay otherwise seconds to delay after interval
}

func NewSizeTriggerPolicy(maxsizeBytes int64, evalOnStart bool) *SizeTriggerPolicy {
	t := new(SizeTriggerPolicy)
	t.maxSize = maxsizeBytes
	t.evalOnStart = evalOnStart
	return t
}

func (s *SizeTriggerPolicy) Rotate(fileName string, strategy Strategy) (bool,error) {
	// 1. Open file and check size if over maxsize then rotate
	fi, err := os.Stat(fileName)
	if err != nil {
		return false, err
	}
	// get the size
	size := fi.Size()
	if size > s.maxSize {
		fo,err := os.Open(fileName)
		if err != nil {
			return false, err
		}
		// List existing gz files and increment to next
		renameFiles(fileName, strategy.Max())

		fw, err := os.Create(fileName+"1.gz")
		if err != nil {
			return false, err
		}
		defer fw.Close()
		gzC := gz.NewLvl(strategy.CompressionLevel())
		err = gzC.Compress(fo, fw)
		if err != nil {
			return false, err
		}
		fo.Close()
		err = os.Truncate(fileName, 0)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func NewTimeTriggerPolicy(interval, maxRandomDelay int) *TimeTriggerPolicy {
	t := new(TimeTriggerPolicy)
	t.interval = interval
	t.maxRandomDelay = maxRandomDelay
	return t
}

func (s *TimeTriggerPolicy) Rotate(fileName string, strategy Strategy) (bool,error) {
	// TODO check rotation triggers
	return false, nil
}

func NewCronTriggerPolicy(schedule string, evalOnStart bool) *CronTriggerPolicy {
	t := new(CronTriggerPolicy)
	t.schedule = schedule
	t.evalOnStart = evalOnStart
	return t
}

func (s *CronTriggerPolicy) Rotate(fileName string, strategy Strategy) (bool,error) {
	// TODO check rotation triggers
	return false, nil
}

// Rename and roll off after XX
func renameFiles(fileName string, max int) error {
	files, err := ioutil.ReadDir(path.Dir(fileName))
	if err != nil {
		return err
	}
	sort.Sort(ByNumericalFilename(files))
	found := make([]os.FileInfo, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".gz") {
			fmt.Println(file.Name())
			found = append(found, file)
		}
	}

	sort.Sort(ByNumericalFilenameRev(found))
	// Rename to 1 higher
	re := regexp.MustCompile("[0-9]+")
	for _,j := range found {
		if ar := re.FindAllString(j.Name(), 1); ar != nil {
			lastNum,err := strconv.ParseInt(ar[0], 10, 64)
			if err != nil {
				fmt.Printf("issue renaming convert found int %v\n", err)
			}
			newNum := lastNum + 1
			oldName := j.Name()
			newName := strings.ReplaceAll(oldName, strconv.Itoa(int(lastNum)), strconv.Itoa(int(newNum)))
			err = os.Rename(path.Join(path.Dir(fileName),oldName), path.Join(path.Dir(fileName),newName))
			if err != nil {
				fmt.Printf("issue renaming %v\n", err)
			}
		}

	}

	// Re-read and clean up
	files, err = ioutil.ReadDir(path.Dir(fileName))
	if err != nil {
		return err
	}
	sort.Sort(ByNumericalFilename(files))
	for i,j := range files {
		if i > max {
			err := os.Remove(path.Join(path.Dir(fileName), j.Name()))
			if err != nil {
				fmt.Printf("issue removing %v\n", err)
			}
		}
	}

	return nil
}

type ByNumericalFilename []os.FileInfo

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	// Use path names
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	// Grab integer value of each filename by parsing the string and slicing off
	// the extension
	re := regexp.MustCompile("[0-9]+")
	var a int64
	var b int64
	var err1 error
	var err2 error
	oneAr := re.FindAllString(pathA, 1)
	twoAr := re.FindAllString(pathB, 1)
	if oneAr != nil {
		a, err1 = strconv.ParseInt(oneAr[0], 10, 64)
	} else {
		err1 = fmt.Errorf("no numbers found")
	}
	if twoAr != nil {
		b, err2 = strconv.ParseInt(twoAr[0], 10, 64)
	} else {
		err2 = fmt.Errorf("no numbers found")
	}

	// If any were not numbers sort lexographically
	if err1 != nil || err2 != nil {
		return pathA < pathB
	}

	// Which integer is smaller?
	return a < b
}

type ByNumericalFilenameRev []os.FileInfo

func (nf ByNumericalFilenameRev) Len() int      { return len(nf) }
func (nf ByNumericalFilenameRev) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilenameRev) Less(i, j int) bool {

	// Use path names
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	// Grab integer value of each filename by parsing the string and slicing off
	// the extension
	re := regexp.MustCompile("[0-9]+")
	var a int64
	var b int64
	var err1 error
	var err2 error
	oneAr := re.FindAllString(pathA, 1)
	twoAr := re.FindAllString(pathB, 1)
	if oneAr != nil {
		a, err1 = strconv.ParseInt(oneAr[0], 10, 64)
	} else {
		err1 = fmt.Errorf("no numbers found")
	}
	if twoAr != nil {
		b, err2 = strconv.ParseInt(twoAr[0], 10, 64)
	} else {
		err2 = fmt.Errorf("no numbers found")
	}

	// If any were not numbers sort lexographically
	if err1 != nil || err2 != nil {
		return pathB < pathA
	}

	// Which integer is smaller?
	return b < a
}