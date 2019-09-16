/*
Package log implements a simple structured logging API designed with few assumptions. Designed for
centralized logging solutions such as Kinesis which require encoding and decoding before fanning-out
to handlers.

You may use this package with inline handlers, much like Logrus, however a centralized solution
is recommended so that apps do not need to be re-deployed to add or remove logging service
providers.
*/
package log
