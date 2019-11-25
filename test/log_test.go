package test

/*
Test project
*/
import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/colt3k/nglog/ers/bserr"
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
	log.Modify(log.LogLevel(log.DEBUGX2), log.ColorsOn(), log.Appenders(ca, fa))

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

	log.Logln(log.DEBUGX2, "Logln X2")
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

func TestErrStacks(t *testing.T) {
	bserr.StopErr(fmt.Errorf("some error one"), nil...)
}

func TestJSON(t *testing.T) {
	fa, err := log.NewFileAppender("*", "output.txt", "", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(ca, fa))

	log.Logln(log.DEBUGX2, "Logln X2")
	log.Logln(log.DEBUG, "Logln Teest")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
	log.Logln(log.FATAL, "---LAST LINE OF TEST---")
}

func TestMail(t *testing.T) {
	ma, err := log.NewMailAppender("*", "my.mailserver.com", "youruser", "yourpass", "from@somewhere.com", "user@to.com", "Test message", 25)
	if err != nil {
		log.Logf(log.FATAL, "issue creating mail appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca, ma))

	log.Logln(log.DEBUG, "Mail Message")
}

func TestHTTP(t *testing.T) {
	ha, err := log.NewHTTPAppender("*", "http://localhost:8080", "", "")
	if err != nil {
		log.Logf(log.FATAL, "issue creating http appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)),log.Appenders(ca, ha))

	log.Logln(log.DEBUGX2, "Logln X2")
	log.Logln(log.DEBUG, "Logln Teest")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
}

func TestTCPSocket(t *testing.T) {
	ha, err := log.NewSocketAppender("*", "localhost", "9090")
	if err != nil {
		log.Logf(log.FATAL, "issue creating socket appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)),log.Appenders(ca, ha))

	log.Logln(log.DEBUGX2, "Logln X2")
	log.Logln(log.DEBUG, "Logln Teest")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
}

func TestSyslog(t *testing.T) {

	// MAC goes to /var/log/system.log
	sa, err := log.NewSyslogAppender("*", "myapp")
	if err != nil {
		log.Logf(log.FATAL, "issue creating syslog appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca, sa))

	log.Logln(log.DEBUGX2, "Logln X2")
	log.Logln(log.DEBUG, "Logln Teest")
	log.Logln(log.ERROR, "Logln Er")
	log.Logln(log.INFO, "Logln inf o msg")
	log.Logln(log.WARN, "Logln wrn msg")
	log.Logln(log.DEBUG, "Logln ftl msg")
}