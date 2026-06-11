package logger

import (
	"bytes"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := New(INFO)
	logger.SetOutput(&buf)

	logger.Info("test message %s", "arg")

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected output, got empty")
	}

	if !bytes.Contains(buf.Bytes(), []byte("INFO")) {
		t.Error("expected INFO level in output")
	}

	if !bytes.Contains(buf.Bytes(), []byte("test message arg")) {
		t.Error("expected message in output")
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := New(WARN)
	logger.SetOutput(&buf)

	logger.Debug("should not appear")
	logger.Info("should not appear")
	logger.Warn("should appear")

	if bytes.Contains(buf.Bytes(), []byte("should not appear")) {
		t.Error("DEBUG and INFO should be filtered")
	}

	if !bytes.Contains(buf.Bytes(), []byte("should appear")) {
		t.Error("WARN should not be filtered")
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := New(ERROR)
	logger.SetOutput(&buf)

	logger.Error("error occurred: %d", 42)

	if !bytes.Contains(buf.Bytes(), []byte("ERROR")) {
		t.Error("expected ERROR level in output")
	}

	if !bytes.Contains(buf.Bytes(), []byte("error occurred: 42")) {
		t.Error("expected formatted message in output")
	}
}

func TestLogger_DefaultLogger(t *testing.T) {
	SetLevel(INFO)

	if defaultLogger.level != INFO {
		t.Error("default logger level should be INFO")
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.level.String())
		}
	}
}