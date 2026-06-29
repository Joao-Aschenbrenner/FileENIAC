// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/observability/metrics"
	"github.com/ENIACSystems/FileENIAC/backend/internal/observability/tracing"
	"go.uber.org/zap"
)

type commandContext struct {
	context.Context
	correlationID string
	tracer        *tracing.Tracer
}

func newCommandContext() *commandContext {
	corrID := log.NewID()
	ctx := log.WithCorrelationID(context.Background(), corrID)
	return &commandContext{
		Context:       ctx,
		correlationID: corrID,
		tracer:        tracing.New(),
	}
}

func (cc *commandContext) log() *zap.Logger {
	return log.WithContext(cc.Context)
}

func (cc *commandContext) startStep(name string) *tracing.Step {
	return cc.tracer.Start(name)
}

func (cc *commandContext) traceOperation(name string, fn func() error) error {
	step := cc.startStep(name)
	start := time.Now()
	err := fn()
	step.EndWithError(err)
	elapsed := time.Since(start)
	fields := []zap.Field{
		zap.String("step", name),
		zap.Duration("duration", elapsed),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		metrics.Get().Counter("cmd.errors", 1)
	}
	metrics.Get().Counter("cmd.ops", 1)
	metrics.Get().Timer(name)
	cc.log().Info("operation", fields...)
	return err
}

func (cc *commandContext) close() {
	fields := []zap.Field{
		zap.String("correlation_id", cc.correlationID),
		zap.Int("steps", len(cc.tracer.Steps())),
	}
	for _, s := range cc.tracer.Steps() {
		fields = append(fields, zap.String("trace."+s.Name, s.Duration.String()))
	}
	cc.log().Info("command completed", fields...)
}
