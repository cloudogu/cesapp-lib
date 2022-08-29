# The logging interface 'core.Logger

## Summary

This page describes background and use-cases of the logging interface [`core.Logger`](../../core.logger.go).

The interface described here handles some needs that arise in the interaction of logging that arise in the development of multiple software applications:

- Need for level logging (`logger.Debug(msg)` versus `logger.V(42).Log(msg)`)
- Need for different data stream repositories
   - syslog
   - TTY
   - Stdout for containers
- Ease of use is more important than a highly simplified log interface
   - `logger.Debug(msg)` vs `logger.Log(log.DebugLevel, "msg")`
- Client determines HOW to log, not the library above it.

At Cloudogu, we develop many applications using Golang. Among them are Golang libraries (like this one) for reasons of [DRY](https://clean-code-developer.com/grades/grade-1-red/#Don8217t_Repeat_Yourself_DRY). One common feature in each application is logging.

Logging is a cross-cutting aspect, i.e., an aspect does touch more than a single (usually domain-specific) area. Even more, this aspect touches several ranges in the same way. For example, logging should work the same way in a `storage` package as it does in an `api` package.

"Log-leveling" is the name given to the classification of messages that satisfy one or more log levels. A small number of log levels simplifies usage and assignment during application development.

## Usage

### Configuring the logger of this library

This library generates log and console output that can be delegated to its own valid logger.

This logger must implement the [logger interface](../../core/logger.go) of this library so that it can successfully output all log-
messages successfully.

To achieve this the logger interface must be implemented in the own application. For this purpose an own
struct (in the example `libraryLogger`) that contains the log solution of the own application (e.g. `logr`,
`logrus`, `log`, etc.). Interface delegates ensure that all log calls are passed to the internal logger.

After that, the logger instance in the library can be overwritten with its own instance.

### Example

The custom interface implementation with the reference to the internal application logger:

```go
package example

type libraryLogger struct {
	logger *yourInternalLogger
}

func (l *libraryLogger) Debug(args ...interface{}) { ... }
func (l *libraryLogger) Info(args ...interface{}) { ... }
func (l *libraryLogger) Warning(args ...interface{}) { ... }
func (l *libraryLogger) Error(args ...interface{}) { ... }
func (l *libraryLogger) Print(args ...interface{}) { ... }

func (l *libraryLogger) Debugf(format string, args ...interface{}) { ... }
func (l *libraryLogger) Infof(format string, args ...interface{}) { ... }
func (l *libraryLogger) Warningf(format string, args ...interface{}) { ... }
func (l *libraryLogger) Errorf(format string, args ...interface{}) { ... }
func (l *libraryLogger) Printf(format string, args ...interface{}) { ... }
```

Now a reference of the own interface implementation overwrites the logger in the library:

```go
package example

func configureLibraryLogger(applicationLogger *yourInternalLogger) {
	// assign the application logger as the internal logger for our logger wrapper.
	cesappLibLogger := libraryLogger{logger: applicationLogger}
	
	// The method core.GetLogger provides the logger instance to the whole library
	core.GetLogger = func() core.Logger {
		return &cesappLibLogger
	}
}
```

## Interface architecture

The following sections address these aspects.

### Actual log behavior

The domain-oriented application (i.e., the application that is actually executed) determines what log messages look like and where these logs should flow. This is important because the application knows better than an abstract library in which scenario log messages make the most sense.

The external appearance is influenced by two points:
1. the format of log messages
2. the output medium of the log messages.

#### Log behavior: Format of the log messages

The first log behavior refers to the structure of the log message itself.
Log formats can be roughly divided into two categories. These have opposing advantages and disadvantages:

| format      | readable by humans | processable by machines | frameworks                         |
|-------------|--------------------|-------------------------|------------------------------------|
| Traditional | good               | bad                     | Golang `log`, `glog`               |
| Structured  | bad                | good                    | `logrus`, `logr`, `log15`, `gokit` |

It depends on the particular use case of the application which format should be used.

##### Traditional format

```
log.Infof("A group of %v walrus emerges from the ocean", 10)

I0522 19:59:41.842000 11996 main.go:18] A group of 10 walrus emerges from the ocean
```

##### Structured format

```
log15.Info("A group of walrus emerges from the ocean", "animal", "walrus", "size", 10)

JSON : {"animal": "walrus", "lvl":3, "msg": "A group of walrus emerges from the ocean", "size":10, "t": "2015-05-21T20:51:53.594-04:00"}
Logfmt: t=2015-05-21T20:48:06-0400 lvl=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
```

#### Log behavior: Output medium of the log messages

In addition to the log format mentioned, the output medium also plays a role in determining the desired log behavior.

Examples of output media:
- StdOut / StdErr
- Syslog data stream
- Direct storage in files
- Combinations of these (e.g. log message with specific log level on TTY; simultaneously ALL log messages to syslog).

The same output medium can be interpreted differently in different contexts.

Contrast example for `StdOut`:
- a CLI application logs warnings to a TTY - a human should immediately read the message
- a Kubernetes container logs to StdOut - first all log messages of the cluster should be aggregated, then only a human reads the message (e.g. by ELK stack)

### Log level

In order to identify relevant log messages (e.g. for debugging) more quickly, mainly log messages with a level are to be produced. The interface therefore defines a series of functions that make these levels available in a simple manner. Here, instead of a minimal interface, it was considered that interface implementations should be easy to use.

These are divided into five groups:

- `Print`
   - is always printed, especially interesting for CLI applications
- `Error`
   - an application error has occurred, a stack trace may give information about the code path where the error occurred
- `Warn`
   - a possibly incorrect condition has been detected, which may cause the application to fail if it continues to run
   - contains the `Error` log level
- `Info`
   - a general information
   - contains the `Warn` log level and all log levels above it
- debug
   - for debugging the application outputs current execution paths with context information
   - contains the `Info` log level as well as all log levels above it

### About the API

If you write a lot of code, you will quickly get to the point of logging messages at many different code locations in many applications. Therefore, ease of use is important in order to keep the effort required to generate logs low.

The API is oriented to methods of `logrus` / `Log4J`. The abstraction from an actually used log framework allows a flexible choice between structured and traditional log formats.

### User circle

The user base here is not the general public, but rather developers at Cloudogu or open source committers.

## Applied Use-cases

1. Logging in CLI-Apps
   1. Logging für zusätzliche Information auf der CLI
   2. Log-Splitting (z. B. `cesapp`) 
      - Syslog-Stream: ALLE Logs (Print- und Level-Logs (Error...Debug))
      - StdOut:
        - ALLE Print-Logs (`logger.Printx()`) werden ausgegeben
          - ersetzt `fmt.Printx()` damit auch TTY-Ausgaben im Logfile erscheinen
        - Level-Logs werden je nach eingestelltem Log-Level (Filtering) ausgegeben
2. Logging in systemd-Services
   1. Syslog-Stream: Level-Logs je nach eingestelltem Log-Level (Filtering) 
3. Logging in Containern (Dogus, K8s-Komponenten)
   1. StdOut: Level-Logs je nach eingestelltem Log-Level (Filtering)

