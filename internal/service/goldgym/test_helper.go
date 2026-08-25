package goldgym

import (
	jaegerLog "gold-gym-be/pkg/log"

	"go.uber.org/zap"
)

func newTestLogger() jaegerLog.Factory {
	logger, _ := zap.NewDevelopment()
	return jaegerLog.NewFactory(logger)
}

func newTestService(repo RepoData) *Service {
	return New(repo, nil, nil, newTestLogger())
}
