package agentkapp

import (
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/pluralsh/kuberentes-agent/pkg/agentcfg"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func (a *App) logger(levelEnum agentcfg.LogLevelEnum, sync zapcore.WriteSyncer) (*zap.Logger, zap.AtomicLevel, error) {
	level, err := logz.LevelFromString(levelEnum.String())
	if err != nil {
		return nil, zap.NewAtomicLevel(), err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	return zap.New(
		zapcore.NewCore(
			&agentIdEncoder{
				Encoder: zapcore.NewJSONEncoder(logz.NewProductionEncoderConfig()),
				agentId: a.AgentId,
			},
			sync,
			atomicLevel,
		),
		zap.ErrorOutput(sync),
	), atomicLevel, nil
}

// agentIdEncoder wraps a zapcore.Encoder to add agent id field if agent id is available.
type agentIdEncoder struct {
	zapcore.Encoder
	agentId *ValueHolder[int64]
}

func (e *agentIdEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	id, ok := e.agentId.tryGet()
	if ok {
		l := len(fields)
		f := make([]zapcore.Field, l+1)
		copy(f, fields)
		f[l] = logz.AgentId(id)
		fields = f
	}
	return e.Encoder.EncodeEntry(entry, fields)
}

func (e *agentIdEncoder) Clone() zapcore.Encoder {
	return &agentIdEncoder{
		Encoder: e.Encoder.Clone(),
		agentId: e.agentId,
	}
}
