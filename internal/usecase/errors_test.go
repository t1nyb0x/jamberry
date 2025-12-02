package usecase

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple message",
			message: "validation failed",
		},
		{
			name:    "japanese message",
			message: "❌ URLを入力してください。",
		},
		{
			name:    "empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &ValidationError{Message: tt.message}
			if err.Error() != tt.message {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.message)
			}
		})
	}
}

func TestNotFoundError(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple message",
			message: "not found",
		},
		{
			name:    "japanese message",
			message: "🔍 該当する結果は見つかりませんでした。",
		},
		{
			name:    "empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &NotFoundError{Message: tt.message}
			if err.Error() != tt.message {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.message)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ValidationError",
			err:  &ValidationError{Message: "test"},
			want: true,
		},
		{
			name: "NotFoundError",
			err:  &NotFoundError{Message: "test"},
			want: false,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "NotFoundError",
			err:  &NotFoundError{Message: "test"},
			want: true,
		},
		{
			name: "ValidationError",
			err:  &ValidationError{Message: "test"},
			want: false,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.err); got != tt.want {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}
