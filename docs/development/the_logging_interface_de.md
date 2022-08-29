# Das Logging-Interface `core.Logger`

## Zusammenfassung

Diese Seite beschreibt Hintergründe und Use-cases zu dem Logging-Interface [`core.Logger`](../../core.logger.go).

Das hier beschriebene Interface bewältigt einige Nöte, die im Zusammenspiel von Logging, die sich in der Entwicklung von mehreren Softwareanwendungen ergeben:

- Bedarf an Level-Logging (`logger.Debug(msg)` versus `logger.V(42).Log(msg)`)
- Bedarf an unterschiedlichen Datenstromablagen
   - Syslog
   - TTY
   - Stdout für Container
- Leichte Benutzung ist wichtiger als ein stark vereinfachtes Log-Interface
   - `logger.Debug(msg)` vs `logger.Log(log.DebugLevel, "msg")`
- Client bestimmt, WIE geloggt werden soll, nicht das darüber liegende Library

Bei Cloudogu entwickeln wir viele Applikationen mit Golang. Darunter sind auch Golang-Bibliotheken (wie diese hier) aus Gründen von [DRY](https://clean-code-developer.com/grades/grade-1-red/#Don8217t_Repeat_Yourself_DRY). Eine Gemeinsamkeit in den jeweiligen Applikationen stellt das Logging dar.

Logging ist ein Querschnittsaspekt, d. h. ein Aspekt berührt nicht einen einzelnen (meist fachlichen) Bereich. Vielmehr berührt dieser Aspekt mehrere Bereiche auf die gleiche Weise. So soll bspw. in einem `storage`-Paket das Logging auf die gleiche Weise funktionieren wie in einem `api`-Paket.

"Log-Leveling" nennt man die Einteilung von Meldungen, die einem oder mehreren Log-Niveaus genügen. Eine geringe Anzahl von Log-Leveln vereinfacht die Benutzung und Zuordnung bei der Anwendungsentwicklung.

## Benutzung

### Den Logger dieser Bibliothek konfigurieren

Diese Bibliothek erzeugt Log- und Konsolenausgaben, die an einen eigenen, gültigen Logger delegiert werden können. 

Dieser Logger muss das [Logger-Interface](../../core/logger.go) dieser Bibliothek implementieren, damit diese alle Log-
Meldungen erfolgreich ausgeben kann. 

Um dies zu erreichen muss das Logger-Interface in der eigenen Anwendung implementiert werden. Dazu wird ein eigenes
Struct (im Beispiel `libraryLogger`) erzeugt, dass die Log-Lösung der eigenen Anwendung enthält (z. B. `logr`, 
`logrus`, `log`, usw.). Interface-Delegates sorgen dafür, dass alle Log-Aufrufe an den internen Logger weitergegeben werden.

Danach kann die Logger-Instanz in der Bibliothek mit der eigenen Instanz überschrieben werden.

### Beispiel

Die eigene Interface-Implementierung mit der Referenz auf den internen Anwendungslogger:

```go
package example

type libraryLogger struct {
	logger *yourInternalLogger
}

func (l *libraryLogger) Debug(args ...interface{})   { ... }
func (l *libraryLogger) Info(args ...interface{})    { ... }
func (l *libraryLogger) Warning(args ...interface{}) { ... }
func (l *libraryLogger) Error(args ...interface{})   { ... }
func (l *libraryLogger) Print(args ...interface{})   { ... }

func (l *libraryLogger) Debugf(format string, args ...interface{})   { ... }
func (l *libraryLogger) Infof(format string, args ...interface{})    { ... }
func (l *libraryLogger) Warningf(format string, args ...interface{}) { ... }
func (l *libraryLogger) Errorf(format string, args ...interface{})   { ... }
func (l *libraryLogger) Printf(format string, args ...interface{})   { ... }
```

Nun überschreibt eine Referenz der eigenen Interface-Implementierung den Logger in der Bibliothek:

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

## Interface-Architektur

Die folgenden Abschnitte gehen auf diese Aspekte ein.

### Tatsächliches Log-Verhalten

Die fachliche Applikation (also die Anwendung, die tatsächlich ausgeführt wird) legt fest, wie die Log-Meldungen aussehen und wohin diese Logs fließen sollen. Dies ist wichtig, weil die Applikation besser als eine abstrakte Bibliothek weiß, in welchem Szenario Log-Meldungen am meisten Sinn ergeben.

Das äußere Erscheinungsbild wird von zwei Punkten beeinflusst:
1. Das Format der Log-Meldungen
2. Die Ausgabemedium der Log-Meldungen 

#### Log-Verhalten: Format der Log-Meldungen

Das erste Log-Verhalten bezieht sich auf die Struktur der Log-Meldung an sich.
Log-Formate lassen sich grob in zwei Kategorien einteilen. Diese besitzen gegensätzliche Vor- und Nachteile:

| Format       | Lesbarkeit Menschen | Verarbeitbarkeit Maschinen | Frameworks                         |
|--------------|---------------------|----------------------------|------------------------------------|
| Traditionell | gut                 | schlecht                   | Golang `log`, `glog`               |
| Strukturiert | schlecht            | gut                        | `logrus`, `logr`, `log15`, `gokit` |

Dabei kommt es auf den jeweiligen Anwendungsfall der Applikation an, welches Format verwendet werden soll. 

##### Traditionelles Format

```
log.Infof("A group of %v walrus emerges from the ocean", 10)

I0522 19:59:41.842000   11996 main.go:18] A group of 10 walrus emerges from the ocean
```

##### Strukturiertes Format

```
log15.Info("A group of walrus emerges from the ocean", "animal", "walrus", "size", 10)

JSON  : {"animal":"walrus","lvl":3,"msg":"A group of walrus emerges from the ocean", "size":10,"t":"2015-05-21T20:51:53.594-04:00"}
Logfmt: t=2015-05-21T20:48:06-0400 lvl=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
```

#### Log-Verhalten: Ausgabemedium der Log-Meldungen

Neben dem genannten Log-Format spielt auch das Ausgabemedium eine Rolle in der Bestimmung des gewünschten Log-Verhaltens.

Beispiele von Ausgabemedien:
- StdOut / StdErr
- Syslog-Datenstrom
- Direkte Ablage in Dateien
- Kombinationen davon (z. B. Log-Meldung mit bestimmtem Log-Level auf TTY; gleichzeitig ALLE Log-Meldungen nach Syslog)

Dabei kann das gleiche Ausgabemedium in unterschiedlichen Kontexten auch unterschiedlich interpretiert werden.

Kontrastbeispiel für `StdOut`:
- eine CLI-Anwendung loggt Warnungen auf ein TTY - ein Mensch soll unmittelbar die Meldung lesen
- ein Kubernetes-Container loggt auf StdOut - erst sollen sämtliche Log-Meldungen des Clusters aggregiert werden, dann erst liest ein Mensch die Meldung (z. B. durch ELK-Stack)

### Log-Level

Um relevante Log-Meldungen (z. B. bei Debugging) schneller zu identifizieren, sollen hauptsächlich Log-Meldungen mit einem Level produziert werden. Das Interface definiert daher eine Reihe von Funktionen, die diese Levels auf einfache Weise zur Verfügung stellen. Hierbei wurde anstelle eines minimalen Interfaces abgewogen, dass Interface-Implementierungen leicht zu benutzen sein sollten.

Diese werden in fünf Gruppen eingeteilt:

- `Print`
  - wird stets ausgegeben, besonders interessant für CLI-Anwendungen
- `Error`
  - ein Anwendungsfehler ist aufgetreten, ein Stacktrace gibt evtl. Aufschluss über den Codepfad, in dem der Fehler aufgetreten ist
- `Warn`
  - ein möglicherweise inkorrekter Zustand wurde wahrgenommen, der evtl. die Anwendung bei weiterem Betrieb zu einem Fehler führt
  - enthält das `Error`-Log-Level
- `Info`
  - eine allgemeine Information
  - enthält das `Warn`-Log-Level sowie alle darüber liegenden Log-Level
- `Debug`
  - zur Fehlerbestimmung gibt die Anwendung aktuelle Ausführungspfade mit Kontextinformatoinen aus
  - enthält das `Info`-Log-Level sowie alle darüber liegenden Log-Level

### Über die API

Wenn man viel Code schreibt, dann gelangt man bei vielen Applikationen schnell an den Punkt, an vielen unterschiedlichen Codestellen Meldungen zu loggen. Daher ist eine einfache Bedienung wichtig, um den Aufwand bei der Erzeugung von Logs gering zu halten.

Die API orientiert sich an Methoden von `logrus` / `Log4J`. Die Abstraktion von einem tatsächlich verwendeten Log-Framework ermöglicht eine flexible Auswahl zwischen strukturierten und traditionellen Log-Formaten.

### Benutzerkreis

Der Benutzerkreis liegt hierbei nicht auf der Allgemeinheit, sondern eher auf Entwickler bei Cloudogu oder Open Source Committer.

## Use-cases für Logging

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

