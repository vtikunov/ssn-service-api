package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type coreWithLevel struct {
	zapcore.Core
	level zapcore.Level
}

// Enabled returns true if the given level is at or above this level.
func (c *coreWithLevel) Enabled(l zapcore.Level) bool {
	return c.level.Enabled(l)
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Callers must use Check before calling Write.
func (c *coreWithLevel) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// With adds structured context to the Core.
func (c *coreWithLevel) With(fields []zapcore.Field) zapcore.Core {
	return &coreWithLevel{
		c.Core.With(fields),
		c.level,
	}
}

// WithLevel - возвращает опцию заменяющую логгер на логгер с необходимым уровнем.
func WithLevel(lvl zapcore.Level) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &coreWithLevel{core, lvl}
	})
}
