package redistool

import (
	"context"

	"go.uber.org/zap"
)

type ZapLogger zap.SugaredLogger

func (l *ZapLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	(*zap.SugaredLogger)(l).Infof(format, v...)
}
