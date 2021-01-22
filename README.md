# nglog

[breaking changes](#breaking-changes")

[changes](#changes)

# GOAL

- establish a Great logging solution for Go 
- support most of what Log4j and others provide

# Logging

### Setup Logging for project (Basic) (log to Console Only)

    ca := log.NewConsoleAppender("*")
    log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca))
    
### Setup Logging with file output (log to Console (using bright-colors) and File)

    ca := log.NewConsoleAppender("*")
    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.INFO), log.HiColorsOn(), log.Appenders(ca, fa))
    
### Setup Logging with JSON output  (log to Console and File using JSON layout)

    ca := log.NewConsoleAppender("*")
    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(ca,fa))
    
    Examples: 
    {"level":"DEBUG","msg":"test message","time":"2021-01-22T10:53:54-07:00"}
    {"level":"ERROR","msg":"error message","time":"2021-01-22T10:53:54-07:00"}
    {"level":"INFO","msg":"info message","time":"2021-01-22T10:53:54-07:00"}
    {"level":"WARN","msg":"warn message","time":"2021-01-22T10:53:54-07:00"}
    {"level":"FATAL","msg":"fatal message","time":"2021-01-22T10:53:54-07:00"}

    
### Setup Logging with Failover (log to File then Console using Failover)

    // NOTE: If primary logger fails alternative appends are used

    ca := log.NewConsoleAppender("*")
    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fileAp, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    fa := log.NewFailoverAppender(fileAp, []log.Appender{ca})
    log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(fa))
    
### Setup Logging with Failover and Filtering (log to File, log to Console upon Failover, only log for filtered path)

    // INFO: Setup logging by Filter i.e. 'github.com/colt3k/nglog/'
    //  The filter is a path in your source code 
    //    that files are permitted to log for the specified appender.
    //  One to many log appenders with different filter's can be created. 

    ca := log.NewConsoleAppender("*")
    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fileAp, err := log.NewFileAppender("github.com/colt3k/nglog/", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    fa := log.NewFailoverAppender(fileAp, []log.Appender{ca})
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(fa))

### Setup Logging with Rolling File output (log to Console and RollingFile)

    ca := log.NewConsoleAppender("*")
    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    rfa, err := log.NewRollingFileAppender("*", filepath.Join("logtest", "roll_test.log"),
		"RollingFileAppender", -1, log.NewSizeTriggerPolicy(0.000100),
		log.NewDefaultStrategy(4, flate.BestCompression))
	 if err != nil {
		log.Logf(log.FATAL, "issue creating rolling file appender\n%+v", err)
	 }
    log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca, rfa))
        
### Setup Logging with Mailer Appender (log to Console and Mailer)

    ca := log.NewConsoleAppender("*")    
    ma, err := log.NewMailAppender("*", "my.mailserver.com", "youruser", "yourpass", "from@somewhere.com", "user@to.com", "Test message", 25)
    if err != nil {
        log.Logf(log.FATAL, "issue creating mail appender\n%+v", err)
    }
    tl := new(log.TextLayout)
    tl.DisableColors = true
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(tl), log.Appenders(ca, ma))        
    
### Setup Logging with HTTP Appender using JSON (log to Console and Push to HTTP)

    ca := log.NewConsoleAppender("*")
    ha, err := log.NewHTTPAppender("*", "http://localhost:8080", "", "")
    if err != nil {
        log.Logf(log.FATAL, "issue creating http appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(ca, ha))

### Setup Logging to TCP Socket and Console

    ca := log.NewConsoleAppender("*")
    ha, err := log.NewSocketAppender("*", "localhost", "9090")
    if err != nil {
        log.Logf(log.FATAL, "issue creating socket appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.DBGL3), log.Appenders(ca, ha))

### Setup Syslog Logging (log to Console and Syslog)

    ca := log.NewConsoleAppender("*")
    sa, err := log.NewSyslogAppender("*", "myapp")
    if err != nil {
		log.Logf(log.FATAL, "issue creating syslog appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.DEBUG), log.Appenders(ca, sa))

### Setup Logging with XML output (log to Console and File using XML layout)

    ca := log.NewConsoleAppender("*")
    fa, err := log.NewFileAppender("*", "output.txt", "", 0)
    if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.XMLLayout)), log.Appenders(ca, fa))
	        
### Appenders

    - [x] ConsoleAppender       : sends output to your console
    - [x] FailoverAppender      : provides a way to provide multiple appenders with a primary and alts
    - [x] FileAppender          : sends output to a defined file
    - [x] HTTPAppender          : sends output to an HTTP URL
    - [x] MailAppender          : sends output to Mail
    - [x] RollingFileAppender   : like a FileAppender but allows a Trigger and Strategy to be set
        - TriggerPolicies
            - SizeTriggerPolicy
                - [x] DefaultSizeTriggerPolicy()    maxSizeMB: 4
                - [x] NewSizeTriggerPolicy   (maxSizeMB  float64)
            - [] CronTriggerPolicy
            - [] TimeTriggerPolicy
        - Strategies
            - [x] DefaultFileStrategy()         maxKeep: 4, compressionLevel: flate.BestCompression
            - [x] NewDefaultFileStrategy(maxKeep, compressionLevel int)
    - [x] SocketAppender        : TCP socket
    - [x] SyslogAppender        : syslog
    - [] Database storage (db, nosql(mongodb, couchdb))
    - [] MQ Apps (ZeroMQ, JeroMQ, RabbitMQ)
    - [] Rewrite
    - [] JPA, JMS, Cassandra, Async

    
### Modifiers
    - [x] LogLevel(log.LogLevel)    : log.Level can be NONE, DBGL3, DBGL2, DEBUG, INFO, WARN, ERROR, FATAL
        -   this is the lowest level that will be output
    - [x] LogOut(io.Writer)         : Default os.Stdout
    - [x] CallDepth(int)            : Default 5, how deep to go in order to find caller
    - [x] SetFlgs(log.Flags)        : Default LstdFlags 
        - [x] Ldate               : date in the local tz (2009/01/23)
        - [x] Ltime               : time in the local tz (01:23:23)
        - [x] Lmicroseconds       : microsecond resolution (01:23:23.123123) based on Ltime
        - [x] Llongfile           : full path of file name and line number: /a/b/c/d.go:23
        - [x] Lshortfile          : file name and line number: d.go:23
        - [x] LUTC                : if Ldate or Ltime is set, use UTC rather than the local time zone
    - [x] ErrColor(log.ColorAttr)   : default FgRed
    - [x] WarnColor(log.ColorAttr)  : default FgYellow
    - [x] InfoColor(log.ColorAttr)  : default FgBlue
    - [x] DebugColor(log.ColorAttr) : default FgWhite
    - [x] DebugLvl2Color(log.ColorAttr) : default FgWhite
    - [x] DebugLvl3Color(log.ColorAttr) : default FgWhite
        -   Color Options
            -   FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite
        -   Bright Color Options
            -   HiBlack, HiRed, HiGreen, HiYellow, HiBlue, HiMagenta, HiCyan, HiWhite
    - [x] ColorsOn()
        -   Non Bright "Color Options" for 
            -   Default Color, Error, Warn, Info, Debug, DebugL2 and DebugL3
    - [x] HiColorsOn()
        -   Bright "Color Options" for 
            -   Default Color, Error, Warn, Info, Debug, DebugL2 and DebugL3
    - [x] Formatter(log.Layout)     : Default TextLayout{ForceColor: true}
        -   log.JSONLayout
        -   log.XMLLayout
    - [x] Appenders(log.Appender)   : default ConsoleAppender, SEE Appenders above
    
### Layouts

    - [x] TextLayout
    - [x] JSONLayout
    - [x] XMLLayout
        
## Influencers    
* originally influenced by 
    * logrus
    * log4j 
    * other great loggers

#changes
----------
### v0.0.13 -> v0.0.14
- log.DBGL3 was added for Debug Level 3
- log.FATALNE was added, provides a logger that won't peform an exit
- Test examples updated and new ones added

#breaking-changes 
----------
### v0.0.13 -> v0.0.14
- log.DEBUGX2 is now log.DBGL2  (Debug Level 2)
- log.NewSizeTriggerPolicy(maxsizeMB int64, evalOnStart bool)   -> log.NewSizeTriggerPolicy(maxSizeMB float64)
    -   maxsizeMB named changed case to maxSizeMB
    -   maxsizeMB type change to handle smaller size then 1MB
    -   evalOnStart was never used, removed
- log.Formatr changed to proper spelling log.Formatter
- Hi<Color> Constants changed to proper case
    -   HiBLACK -> HiBlack, etc...
- HiBlack function changed to BrightBlack
