package logger

import (
	"context"
	"fmt"
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

const (
	path     = "./logs/"
	debugLog = "debug.log"
	infoLog  = "info.log"
	warnLog  = "warn.log"
	errorLog = "error.log"

	TraceID = "TraceID"
)

func init() {
	debugWriter, err := rotatelogs.New(path+debugLog+"_%Y%m%d%H%M",
		rotatelogs.WithLinkName(path+debugLog),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(4)*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	warnWriter, err := rotatelogs.New(path+warnLog+"_%Y%m%d%H%M",
		rotatelogs.WithLinkName(path+warnLog),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(4)*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	infoWriter, err := rotatelogs.New(path+infoLog+"_%Y%m%d%H%M",
		rotatelogs.WithLinkName(path+infoLog),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(4)*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	errorWriter, err := rotatelogs.New(path+errorLog+"_%Y%m%d%H%M",
		rotatelogs.WithLinkName(path+errorLog),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(4)*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	//Log.SetReportCaller(true)

	log.SetOutput(&fakeWriter{})
	log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	log.SetLevel(logrus.DebugLevel)
	//log.SetReportCaller(true)

	wh := NewWriterHook(map[logrus.Level]io.Writer{
		logrus.DebugLevel: debugWriter,
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  warnWriter,
		logrus.ErrorLevel: errorWriter})
	log.AddHook(wh)
}

func LogDebug(ctx context.Context, msg string) {
	v := ctx.Value(TraceID)
	log.WithFields(logrus.Fields{"traceID": v}).Debug(msg)
}

func LogDebugf(ctx context.Context, msg string, v ...interface{}) {
	LogDebug(ctx, fmt.Sprintf(msg, v))
}

func LogInfo(ctx context.Context, msg string) {
	v := ctx.Value(TraceID)
	log.WithFields(logrus.Fields{"traceID": v}).Info(msg)
}

func LogInfof(ctx context.Context, msg string, v ...interface{}) {
	LogInfo(ctx, fmt.Sprintf(msg, v))
}

func LogWarn(ctx context.Context, msg string) {
	v := ctx.Value(TraceID)
	log.WithFields(logrus.Fields{"traceID": v}).Warn(msg)
}

func LogWarnf(ctx context.Context, msg string, v ...interface{}) {
	LogWarn(ctx, fmt.Sprintf(msg, v))
}

func LogError(ctx context.Context, msg string, err error) {
	v := ctx.Value(TraceID)
	log.WithFields(logrus.Fields{"traceID": v, "error": err}).Errorln(msg)
}

func LogErrorf(ctx context.Context, msg string, err error, v ...interface{}) {
	LogError(ctx, fmt.Sprintf(msg, v), err)
}

type fakeWriter struct{}

func (f *fakeWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type writerHook struct {
	writers map[logrus.Level]io.Writer
}

func NewWriterHook(writers map[logrus.Level]io.Writer) *writerHook {
	return &writerHook{writers: writers}
}

func (w *writerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w *writerHook) Fire(entry *logrus.Entry) error {
	bytes, err := entry.Bytes()
	if err != nil {
		fmt.Println("writerHook#Fire entry.Bytes err:", err)
		return err
	}

	if writer, ok := w.writers[entry.Level]; ok {
		if _, err := writer.Write(bytes); err != nil {
			fmt.Println("writerHook#Fire write err:", err)
		}
		return err
	}

	return nil
}
