# nglog

## GOAL

- to be a Log4j equivalent for Go

# Logging

- If the log is just a record of something happening, if it could be aggregated, it's probably a metric (Prometheus).

- If the log is something important, then it's not a log it's an alert. Send it by email, pagerduty, chat-ops, etc.

- If the log is not something important then remove it.


### Setup Logging for project (Basic)

    ca := log.NewConsoleAppender("*")
    log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca, fa))
    
### Setup Logging with file output

    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca, fa))
    
### Setup Logging with JSON output

    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(ca,fa))
    
### Setup Logging with Failover

    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("*", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    fa := log.NewFailoverAppender(fileAp, []log.Appender{ca})
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(fa))
    
### Setup Logging with Failover and Filtering

    logfile = filepath.Join(file.HomeFolder(), appName+".log")
    fa, err := log.NewFileAppender("github.com/colt3k/nglog/", logfile, "", 0)
    if err != nil {
        log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    fa := log.NewFailoverAppender(fileAp, []log.Appender{ca})
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(fa))
    
### Setup Logging with Mailer

    ma, err := log.NewMailAppender("*", "my.mailserver.com", "youruser", "yourpass", "from@somewhere.com", "user@to.com", "Test message", 25)
    if err != nil {
        log.Logf(log.FATAL, "issue creating mail appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    tl := new(log.TextLayout)
    tl.DisableColors = true
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(tl), log.Appenders(ca, ma))        
    
### Setup Logging with HTTP Appender using JSON

    ha, err := log.NewHTTPAppender("*", "http://localhost:8080", "", "")
    if err != nil {
        log.Logf(log.FATAL, "issue creating http appender\n%+v", err)
    }
    ca := log.NewConsoleAppender("*")
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(ca, ha))
        
### Appenders

    - Console
    - Failover
    - File
    - Mail
    - HTTP
    - TCP socket
    - syslog
    
### Layouts

    Text
    JSON    
        
# Coming    
   
## Appenders
    - storage (db, file, nosql(mongodb, couchdb))
    - MQ Apps (ZeroMQ, JeroMQ, RabbitMQ)
    - Rewrite
    - JPA, JMS, Cassandra, Async
        
## Influencers    
originally influenced by logrus, log4j and other great loggers
