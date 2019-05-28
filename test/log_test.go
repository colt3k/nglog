package test

/*
Test project
*/
import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/colt3k/nglog/internal/pkg/util"
	log "github.com/colt3k/nglog/ng"
)

func Test(t *testing.T) {
	//f, err := os.Create("test.log")
	//if err != nil {
	//	log.Println(err)
	//}

	ca := log.NewConsoleAppender("*")
	fa, err := log.NewFileAppender("*", filepath.Join(util.HomeFolder(), "test.log"), "FileAppender", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca, fa))

	log.ShowConfig()
	var val = "HelloWorld"
	var val2 = 0.456789
	//ca := log.NewConsoleAppender("github.com/colt3k/nglog/")
	//fa := log.NewFileAppender("*", "test.log", 0)
	//log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca, fa))
	//log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn())
	log.Println("---FIRST LINE---", val2, val)
	log.Printf("Printf Test some text %s %v ", val, val2)
	log.Printf("Printf Test some text %v %s ", val2, val)
	log.Logf(log.INFO, "Logf Test some text %s %f", val, val2)
	log.Logf(log.INFO, "Throw Error %v", errors.New("test1"))

	log.Logln(log.DEBUG, "Logln Teest")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
	log.Logln(log.FATAL, "---LAST LINE OF TEST---")

}

func TestRollingFileAppender(t *testing.T) {

	log.NewRollingFileAppender("*", filepath.Join(util.HomeFolder(), "roll_test.log"),
		"RollingFileAppender", -1, log.NewTimeTriggerPolicy(5, 1),
		log.NewDefaultStrategy("", "log/", 1, 4, 9))
}
func TestLogFile(t *testing.T) {
	var val = "HelloWorld"
	var val2 = 0.456789

	fileAp, err := log.NewFileAppender("github.com/colt3k/nglog/", "myfile.log", "FileAppender", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("github.com/colt3k/nglog/")
	appenders := []log.Appender{ca}
	fa := log.NewFailoverAppender(fileAp, appenders)

	//log.Modify(log.LogLevel(log.DEBUG), log.LogOut(f), log.Formatr(new(log.JSONLayout)))
	log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(fa))
	//log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn())
	log.Println("---FIRST LINE---", val2, val)
	log.Printf("Printf Test some text %s %v ", val, val2)
	log.Printf("Printf Test some text %v %s ", val2, val)
	log.Logf(log.INFO, "Logf Test some text %s %f", val, val2)
	log.Logln(log.DEBUG, "Logln Test")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
	log.Logln(log.FATAL, "---LAST LINE OF TEST---")
}
