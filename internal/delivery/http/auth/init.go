package goldgym

import (
	"context"
	"gold-gym-be/internal/entity/auth/v2"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvc interface {
	LoginUser(ctx context.Context, _user, _password string, _host string, device string) (auth.Token, map[string]interface{}, string, error)
}

type IgoldgymauthSvc interface {
	RefreshToken(ctx context.Context, refreshToken string) (auth.Token, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Handler struct {
	goldgymSvc     IgoldgymSvc
	goldgymauthSvc IgoldgymauthSvc
	tracer         opentracing.Tracer
	logger         jaegerLog.Factory
}

// New for bridging product handler initialization
func New(is IgoldgymSvc, auth IgoldgymauthSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvc:     is,
		goldgymauthSvc: auth,
		tracer:         tracer,
		logger:         logger,
	}
}
