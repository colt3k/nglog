# nglog

## GOAL

- Log4j equivalent for Go

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
    log.Modify(log.LogLevel(log.DEBUG), log.Formatr(new(log.JSONLayout)), log.Appenders(fa))
    
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
    
### Appenders

    Console
    Failover
    File
    
### Layouts

    Text
    JSON    
        
# Coming    
   
## Appenders    
    - mailer SMTP
    - file appender 
    - storage (db, file, nosql(mongodb, couchdb))
    - random access file
    - socket
    - ssl
    - syslog (possibly in golang already)
    - MQ Apps (ZeroMQ, JeroMQ, RabbitMQ)
    - Rewrite
    - HTTP
    - JPA, JMS, Cassandra, Async
    - Console (default) StdOut can change to StdErr
        - Add template formatting ability
        - filter ability for routing events to log appenders
        - direct to file, skip StdOut/StdErr
        
    
