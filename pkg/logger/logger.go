package pkg_logger

import (
	"context"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.cognotif/internal/repository/constant"
)

type Logger struct {
	*logrus.Logger
}

type LogFormatter struct{}

func NewLogger() *Logger {
	return &Logger{
		Logger: &logrus.Logger{
			Out:          io.MultiWriter(os.Stdout),
			ReportCaller: true,
			Level:        logrus.DebugLevel,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000",
				CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
					fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
					return "", fileName
				},
			},
		},
	}
}

func (x *Logger) Hashcode(c context.Context) *logrus.Entry {
	return x.Logger.WithField("hashcode", c.Value(constant.Hashcode{}).(string))
}
