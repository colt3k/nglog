package test

/*
Test project
*/
import (
	"errors"
	"fmt"
	"github.com/colt3k/nglog/ers/bserr"
	"path/filepath"
	"testing"

	"github.com/colt3k/nglog/internal/pkg/util"
	log "github.com/colt3k/nglog/ng"
)

func Test(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	fa, err := log.NewFileAppender("*", filepath.Join(util.HomeFolder(), "test.log"), "FileAppender", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.ColorsOn(), log.Appenders(ca, fa))

	log.ShowConfig()
	var val = "HelloWorld"
	var val2 = 0.456789
	log.Println("---FIRST LINE---", val2, val)
	log.Printf("Printf Test some text %s %v ", val, val2)
	log.Printf("Printf Test some text %v %s ", val2, val)
	log.Logf(log.INFO, "Logf Test some text %s %f", val, val2)
	log.Logf(log.INFO, "Throw Error %v", errors.New("test1"))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")

}

func TestRollingFileAppender(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	rfa, err := log.NewRollingFileAppenderWithTriggerAndStrategy("*", filepath.Join("logtest", "roll_test.log"),
		"RollingFileAppender", -1, log.NewSizeTriggerPolicy(0.000100), log.DefaultFileStrategy())
	if err != nil {
		log.Logf(log.FATAL, "issue creating rolling file appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DEBUG), log.HiColorsOn(), log.Appenders(ca, rfa))
	_ = log.ShowAppenders()

	log.Logln(log.INFO, "12345")
	log.Logln(log.INFO, "6789")
	log.Logln(log.INFO, "9876")
	log.Logln(log.INFO, "54321")
	log.Logln(log.INFO, "2468")
	log.Logln(log.INFO, "8642")
	log.Logln(log.INFO, "139")
	log.Logln(log.INFO, "931")
}

func TestLogFile(t *testing.T) {
	var val = "HelloWorld"
	var val2 = 0.456789

	ca := log.NewConsoleAppender("github.com/colt3k/nglog/")
	fileAp, err := log.NewFileAppender("github.com/colt3k/nglog/", "myfile.log", "FileAppender", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	appenders := []log.Appender{ca}
	fa := log.NewFailoverAppender(fileAp, appenders)

	//log.Modify(log.LogLevel(log.DEBUG), log.LogOut(f), log.Formatr(new(log.JSONLayout)))
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.JSONLayout)), log.Appenders(fa))
	//log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn())
	log.Println("---FIRST LINE---", val2, val)
	log.Printf("Printf Test some text %s %v ", val, val2)
	log.Printf("Printf Test some text %v %s ", val2, val)
	log.Logf(log.INFO, "Logf Test some text %s %f (should not include newline)", val, val2)
	log.Logf(log.INFO, "2nd Logf text %s %f (should not include newline)", val, val2)
	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestJSON(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	fa, err := log.NewFileAppender("*", "output.txt", "", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.JSONLayout)), log.Appenders(ca, fa))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestMail(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	ma, err := log.NewMailAppender("*","my.mailserver.com","youruser","yourpass",
		"from@somewhere.com", "user@to.com", "Test message", 25)
	if err != nil {
		log.Logf(log.FATAL, "issue creating mail appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca, ma))

	log.Logln(log.DEBUG, "Mail Message")
}

func TestHTTP(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	ha, err := log.NewHTTPAppender("*", "http://localhost:8080", "", "")
	if err != nil {
		log.Logf(log.FATAL, "issue creating http appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.JSONLayout)),log.Appenders(ca, ha))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestTCPSocket(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	ha, err := log.NewSocketAppender("*", "localhost", "9090")
	if err != nil {
		log.Logf(log.FATAL, "issue creating socket appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.JSONLayout)),log.Appenders(ca, ha))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

// MAC goes to /var/log/system.log
func TestSyslog(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	sa, err := log.NewSyslogAppender("*", "myapp")
	if err != nil {
		log.Logf(log.FATAL, "issue creating syslog appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Appenders(ca, sa))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestSyslogXML(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	fa, err := log.NewFileAppender("*", "output.txt", "", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.XMLLayout)), log.Appenders(ca, fa))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestErrStacks(t *testing.T) {
	bserr.StopErr(fmt.Errorf("some error one"), nil...)
}
