package goldgym

import (
	"context"
	authV2 "gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	jaegerLog "gold-gym-be/pkg/log"
	pb "gold-gym-be/proto"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type IgoldgymSvc interface {
	GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	GetGoldUserDataByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	LoginUser(ctx context.Context, user, password, host, device string) (authV2.Token, map[string]interface{}, string, error)
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
	GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
}

type Handler struct {
	pb.UnimplementedGoldGymServiceServer
	goldgymSvc IgoldgymSvc
	tracer     opentracing.Tracer
	logger     jaegerLog.Factory
}

func NewHandler(goldgymSvc IgoldgymSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvc: goldgymSvc,
		tracer:     tracer,
		logger:     logger,
	}
}

func (h *Handler) GetGoldUser(ctx context.Context, req *pb.GetGoldUserRequest) (*pb.GetGoldUserResponse, error) {
	var spanCtx opentracing.SpanContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		spanCtx, _ = h.tracer.Extract(opentracing.TextMap, metadataTextMap(md))
	}

	span := h.tracer.StartSpan("GetGoldUser", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("gRPC request received", zap.String("method", "GetGoldUser"))

	users, err := h.goldgymSvc.GetGoldUser(ctx)
	if err != nil {
		h.logger.For(ctx).Error("Failed to get gold users", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get gold users: %v", err)
	}

	pbUsers := make([]*pb.GoldUser, 0, len(users))
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.GoldUser{
			GoldId:                int32(user.GoldId),
			GoldEmail:             user.GoldEmail,
			GoldPassword:          user.GoldPassword,
			GoldNama:              user.GoldNama,
			GoldNomorhp:           user.GoldNomorHp,
			GoldNomorkartu:        user.GoldNomorKartu,
			GoldCvv:               user.GoldCvv,
			GoldExpireddate:       user.GoldExpireddate,
			GoldNamapemegangkartu: user.GoldPemegangKartu,
		})
	}

	h.logger.For(ctx).Info("Successfully retrieved gold users", zap.Int("count", len(pbUsers)))

	return &pb.GetGoldUserResponse{
		Users: pbUsers,
	}, nil
}

func (h *Handler) GetGoldUserByEmail(ctx context.Context, req *pb.GetGoldUserByEmailRequest) (*pb.GetGoldUserByEmailResponse, error) {
	var spanCtx opentracing.SpanContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		spanCtx, _ = h.tracer.Extract(opentracing.TextMap, metadataTextMap(md))
	}

	span := h.tracer.StartSpan("GetGoldUserByEmail", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("gRPC request received",
		zap.String("method", "GetGoldUserByEmail"),
		zap.String("email", req.Email))

	if req.Email == "" {
		h.logger.For(ctx).Error("Email parameter is required")
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	user, err := h.goldgymSvc.GetGoldUserDataByEmail(ctx, req.Email)
	if err != nil {
		h.logger.For(ctx).Error("Failed to get gold user by email", zap.Error(err))

		if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			return nil, status.Errorf(codes.NotFound, "user with email %s not found", req.Email)
		}

		return nil, status.Errorf(codes.Internal, "failed to get gold user: %v", err)
	}

	pbUser := &pb.GoldUser{
		GoldId:                int32(user.GoldId),
		GoldEmail:             user.GoldEmail,
		GoldPassword:          "", // EMPTY for security - never return password
		GoldNama:              user.GoldNama,
		GoldNomorhp:           user.GoldNomorHp,
		GoldNomorkartu:        user.GoldNomorKartu,
		GoldCvv:               user.GoldCvv,
		GoldExpireddate:       user.GoldExpireddate,
		GoldNamapemegangkartu: user.GoldPemegangKartu,
	}

	h.logger.For(ctx).Info("Successfully retrieved gold user by email",
		zap.Int("user_id", user.GoldId))

	return &pb.GetGoldUserByEmailResponse{
		User: pbUser,
	}, nil
}

func (h *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var spanCtx opentracing.SpanContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		spanCtx, _ = h.tracer.Extract(opentracing.TextMap, metadataTextMap(md))
	}

	span := h.tracer.StartSpan("LoginUser", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("gRPC request received",
		zap.String("method", "LoginUser"),
		zap.String("email", req.Email))

	if req.Email == "" || req.Password == "" {
		h.logger.For(ctx).Error("Email and password are required")
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	token, userData, _, err := h.goldgymSvc.LoginUser(ctx, req.Email, req.Password, "grpc", "")
	if err != nil {
		h.logger.For(ctx).Error("Failed to login user", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	h.logger.For(ctx).Info("User logged in successfully")

	return &pb.LoginUserResponse{
		Token:     token.AccessToken,
		UserId:    userData["user_id"].(string),
		UserEmail: userData["user_email"].(string),
		UserName:  userData["user_name"].(string),
	}, nil
}

func (h *Handler) InsertGoldUser(ctx context.Context, req *pb.InsertGoldUserRequest) (*pb.InsertGoldUserResponse, error) {
	var spanCtx opentracing.SpanContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		spanCtx, _ = h.tracer.Extract(opentracing.TextMap, metadataTextMap(md))
	}

	span := h.tracer.StartSpan("InsertGoldUser", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("gRPC request received",
		zap.String("method", "InsertGoldUser"),
		zap.String("email", req.GoldEmail))

	if req.GoldEmail == "" || req.GoldPassword == "" || req.GoldNama == "" {
		h.logger.For(ctx).Error("Required fields are missing")
		return nil, status.Errorf(codes.InvalidArgument, "email, password, and name are required")
	}

	user := goldEntity.GetGoldUsers{
		GoldEmail:         req.GoldEmail,
		GoldPassword:      req.GoldPassword,
		GoldNama:          req.GoldNama,
		GoldNomorHp:       req.GoldNomorhp,
		GoldNomorKartu:    req.GoldNomorkartu,
		GoldCvv:           req.GoldCvv,
		GoldExpireddate:   req.GoldExpireddate,
		GoldPemegangKartu: req.GoldNamapemegangkartu,
	}

	result, err := h.goldgymSvc.InsertGoldUser(ctx, user)
	if err != nil {
		h.logger.For(ctx).Error("Failed to insert user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	h.logger.For(ctx).Info("User created successfully")

	return &pb.InsertGoldUserResponse{
		UserId:  result.(string),
		Message: "User created successfully",
	}, nil
}

func (h *Handler) GetAllSubscription(ctx context.Context, req *pb.GetAllSubscriptionRequest) (*pb.GetAllSubscriptionResponse, error) {
	var spanCtx opentracing.SpanContext
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		spanCtx, _ = h.tracer.Extract(opentracing.TextMap, metadataTextMap(md))
	}

	span := h.tracer.StartSpan("GetAllSubscription", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("gRPC request received", zap.String("method", "GetAllSubscription"))

	subscriptions, err := h.goldgymSvc.GetAllSubscription(ctx)
	if err != nil {
		h.logger.For(ctx).Error("Failed to get subscriptions", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get subscriptions: %v", err)
	}

	pbSubscriptions := make([]*pb.Subscription, 0, len(subscriptions))
	for _, sub := range subscriptions {
		pbSubscriptions = append(pbSubscriptions, &pb.Subscription{
			GoldNamapaket:       sub.GoldNamaPaket,
			GoldNamalayanan:     sub.GoldNamaLayanan,
			GoldHarga:           sub.GoldHarga,
			GoldJadwal:          sub.GoldJadwal,
			GoldListlatihan:     sub.GoldListLatihan,
			GoldJumlahpertemuan: int32(sub.GoldJumlahpertemuan),
			GoldDurasi:          int32(sub.GoldDurasi),
		})
	}

	h.logger.For(ctx).Info("Successfully retrieved subscriptions", zap.Int("count", len(pbSubscriptions)))

	return &pb.GetAllSubscriptionResponse{
		Subscriptions: pbSubscriptions,
	}, nil
}

type metadataTextMap metadata.MD

func (m metadataTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range m {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}
