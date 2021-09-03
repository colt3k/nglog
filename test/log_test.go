package test

/*
Test project
*/
import (
	"errors"
	"fmt"
	"github.com/colt3k/nglog/ers/bserr"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/colt3k/nglog/internal/pkg/util"
	log "github.com/colt3k/nglog/ng"
)
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Obj struct {
	Name	string `json:"name"`
	Value 	string `json:"value"`
	Id 		int `json:"id"`
}
func TestDump(t *testing.T) {
	//o := Obj{Name: "fred", Value: "testxxx", Id: 90}

	log.InitLevel(log.DBGL3)
}

func Test(t *testing.T) {

	log.InitFileAndLevel(filepath.Join(util.HomeFolder(), "test.log"),log.DBGL3)

	//log.ShowConfig()
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

	log.InitRollingFileAndLevel(filepath.Join("logtest", "roll_test.log"), log.DEBUG)
	_ = log.ShowAppenders()

	// Log several MB of data to log so it will roll
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Create a line of text (813 is 100k)
	for i := 0; i < 7000; i++ {
		log.Logf(log.INFO, "%s", StringWithCharset(80, charset, seededRand))
	}
}
func StringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestLogFile(t *testing.T) {
	var val = "HelloWorld"
	var val2 = 0.456789

	log.InitFileFilterAndLevelFailoverFormatter( "myfile.log","github.com/colt3k/nglog/", log.DBGL3, new(log.JSONLayout))

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

	log.InitFileAndLevelFormatter("output.txt", log.DBGL3, new(log.JSONLayout))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestXML(t *testing.T) {

	log.InitFileAndLevelFormatter("output.txt", log.DBGL3, new(log.XMLLayout))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

/*
TestMail in order to test install MailHog
	brew install mailhog
	https://github.com/mailhog/MailHog
	Start via: MailHog
	HTTP Port: 8025
	SMTP Port 1025, non-ssl
 */
func TestMail(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	ma, err := log.NewMailAppender("*","localhost","youruser","yourpass",
		"from@somewhere.com", "user@to.com", "Test Message", 1025, false)
	if err != nil {
		log.Logf(log.FATAL, "issue creating mail appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca, ma), log.Formatter(new(log.TextLayout)))

	log.Logln(log.DBGL3, "debug level 3 message")
	log.Logln(log.DBGL2, "debug level 2 message")
	log.Logln(log.DEBUG, "debug message")
	log.Logln(log.ERROR, "error message")
	log.Logln(log.INFO, "info message")
	log.Logln(log.WARN, "warn message")
	log.Logln(log.FATALNE, "fatal message")
}

func TestHTTP(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	ha, err := log.NewHTTPAppender("*", "http://localhost:8081", "", "")
	if err != nil {
		log.Logf(log.FATAL, "issue creating http appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.TextLayout)),log.Appenders(ca, ha))

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
	sa, err := log.NewSyslogAppender("*", "myapp")
	if err != nil {
		log.Logf(log.FATAL, "issue creating syslog appender\n%+v", err)
	}
	log.Modify(log.LogLevel(log.DBGL3), log.Formatter(new(log.XMLLayout)), log.Appenders(ca, sa))

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

type TestType struct {
	Name	string
	Value	int
}
func TestTypeOfValue(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca))

	log.PrintTypeOfValue("hello")
}
func TestStructWithFieldNames(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca))
	log.PrintStructWithFieldNames(TestType{Name: "John Doe", Value: 2})
	log.PrintStructWithFieldNamesIndent(TestType{Name: "John Doe", Value: 2}, true)
}
func TestGoSyntaxOfValue(t *testing.T) {

	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca))

	log.PrintGoSyntaxOfValue("test 1")
	log.PrintGoSyntaxOfValue(1)
	log.PrintGoSyntaxOfValue(TestType{Name: "John Doe", Value: 2})
}
