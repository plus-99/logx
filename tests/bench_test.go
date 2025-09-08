package logx_test

import (
	"os"
	"testing"

	"github.com/plus99/logx"
	logrus "github.com/sirupsen/logrus"
	"github.com/rs/zerolog"
)

func BenchmarkLogxInfo(b *testing.B) {
	l := logx.New()
	l.SetLevel(logx.InfoLevel)
	l.SetEncoder(logx.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z07:00"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.WithFields(logx.Fields{"n": i, "user": "u"}).Info("bench")
	}
}

func BenchmarkLogrusInfo(b *testing.B) {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logrus.WithFields(logrus.Fields{"n": i, "user": "u"}).Info("bench")
	}
}

func BenchmarkZerologInfo(b *testing.B) {
	w := os.Stdout
	logger := zerolog.New(w).With().Timestamp().Logger()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info().Int("n", i).Str("user", "u").Msg("bench")
	}
}