package middleware

import (
	"context"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type ImiddlewareSvc interface {
}

type IgoldgymSvc interface {
	GetGoldUserByEmail(ctx context.Context, email string) (string, error)
}

type IgoldgymSvcStock interface {
}

type (
	// Handler ...
	Handler struct {
		middlewareSvc   ImiddlewareSvc
		goldgymSvc      IgoldgymSvc
		goldgymSvcStock IgoldgymSvcStock

		tracer opentracing.Tracer
		logger jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(im ImiddlewareSvc, is IgoldgymSvc, isst IgoldgymSvcStock, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		middlewareSvc:   im,
		goldgymSvc:      is,
		goldgymSvcStock: isst,
		tracer:          tracer,
		logger:          logger,
	}
}
