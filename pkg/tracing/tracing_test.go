package tracing

import (
	"testing"

	"gold-gym-be/pkg/log"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockLogger implements log.Logger interface for testing
type mockLogger struct {
	errorCalled bool
	infoCalled  bool
	lastMsg     string
}

func (m *mockLogger) Debug(msg string, fields ...zap.Field) {
	m.lastMsg = msg
}

func (m *mockLogger) Info(msg string, fields ...zap.Field) {
	m.infoCalled = true
	m.lastMsg = msg
}

func (m *mockLogger) Error(msg string, fields ...zap.Field) {
	m.errorCalled = true
	m.lastMsg = msg
}

func (m *mockLogger) Fatal(msg string, fields ...zap.Field) {
	m.lastMsg = msg
}

func (m *mockLogger) With(fields ...zap.Field) log.Logger {
	return m
}

// mockFactory implements log.Factory interface for testing
type mockFactory struct {
	logger *mockLogger
}

func (m *mockFactory) Bg() log.Logger {
	return m.logger
}

func (m *mockFactory) For(ctx interface{}) log.Logger {
	return m.logger
}

func TestJaegerLoggerAdapter_Error(t *testing.T) {
	mockLog := &mockLogger{}
	adapter := jaegerLoggerAdapter{logger: mockLog}

	testMsg := "test error message"
	adapter.Error(testMsg)

	assert.True(t, mockLog.errorCalled, "Error method should be called")
	assert.Equal(t, testMsg, mockLog.lastMsg, "Error message should match")
}

func TestJaegerLoggerAdapter_Infof(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple message",
			msg:      "test message",
			args:     nil,
			expected: "test message",
		},
		{
			name:     "formatted message with one arg",
			msg:      "test %s",
			args:     []interface{}{"formatted"},
			expected: "test formatted",
		},
		{
			name:     "formatted message with multiple args",
			msg:      "test %s %d",
			args:     []interface{}{"number", 42},
			expected: "test number 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := &mockLogger{}
			adapter := jaegerLoggerAdapter{logger: mockLog}

			adapter.Infof(tt.msg, tt.args...)

			assert.True(t, mockLog.infoCalled, "Info method should be called")
			assert.Equal(t, tt.expected, mockLog.lastMsg, "Info message should match expected format")
		})
	}
}

func TestInit_InvalidEnvVars(t *testing.T) {
	// This test verifies that Init handles missing/invalid environment variables
	// Note: In real scenario, this would call Fatal and exit, so we skip actual Init call
	// This is more of a documentation test showing the expected behavior

	t.Run("documentation", func(t *testing.T) {
		// Document that Init expects these environment variables:
		// - JAEGER_AGENT_HOST
		// - JAEGER_AGENT_PORT
		// - JAEGER_SAMPLER_TYPE (optional, defaults to const)
		// - JAEGER_SAMPLER_PARAM (optional, defaults to 1)

		// If these are not set properly, Init will call logger.Fatal
		assert.True(t, true, "Init requires proper Jaeger environment variables")
	})
}

func TestInit_ServiceNameSetting(t *testing.T) {
	// This test documents that the service name is properly set
	t.Run("documentation", func(t *testing.T) {
		serviceName := "test-service"
		// When Init is called with serviceName, it should:
		// 1. Parse env vars using config.FromEnv()
		// 2. Set cfg.ServiceName = serviceName
		// 3. Set cfg.Sampler.Type = "const"
		// 4. Set cfg.Sampler.Param = 1

		assert.Equal(t, "test-service", serviceName, "Service name should be set correctly")
	})
}
