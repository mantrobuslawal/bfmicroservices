package catalog

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNormalisePageSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   int
		want    int
		wantErr error
	}{
		{
			name:  "zero uses default",
			input: 0,
			want:  defaultPageSize,
		},
		{
			name:  "negative uses default",
			input: -10,
			want:  defaultPageSize,
		},
		{
			name:  "one is accepted",
			input: 1,
			want:  1,
		},
		{
			name:  "max page size is accepted",
			input: maxPageSize,
			want:  maxPageSize,
		},
		{
			name:    "above max returns invalid page size",
			input:   maxPageSize + 1,
			wantErr: ErrInvalidPageSize,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalisePageSize(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("normalisePageSize() error = %v, want %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Fatalf("normalisePageSize() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 6, 8, 10, 30, 0, 0, time.UTC)

	payload := catalogCursor{
		CreatedAt: createdAt,
		ID:        "prod_test_001",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal cursor payload: %v", err)
	}

	token := base64.RawURLEncoding.EncodeToString(data)

	got, err := decodeCursor(token)
	if err != nil {
		t.Fatalf("decodeCursor() error = %v, want nil", err)
	}

	if got.ID != payload.ID {
		t.Fatalf("decoded ID = %q, want %q", got.ID, payload.ID)
	}

	if !got.CreatedAt.Equal(payload.CreatedAt) {
		t.Fatalf("decoded CreatedAt = %s, want %s", got.CreatedAt, payload.CreatedAt)
	}
}

func TestDecodeCursorRejectsMalformedToken(t *testing.T) {
	t.Parallel()

	_, err := decodeCursor("not-valid-base64-%")

	if err == nil {
		t.Fatal("decodeCursor() error = nil, want error")
	}
}

func TestDecodeCursorRejectsMissingID(t *testing.T) {
	t.Parallel()

	payload := catalogCursor{
		CreatedAt: time.Date(2026, 6, 8, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal cursor payload: %v", err)
	}

	token := base64.RawURLEncoding.EncodeToString(data)

	_, err = decodeCursor(token)
	if err == nil {
		t.Fatal("decodeCursor() error = nil, want error")
	}
}

func TestDecodeCursorRejectsMissingCreatedAt(t *testing.T) {
	t.Parallel()

	payload := catalogCursor{
		ID: "prod_test_001",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal cursor payload: %v", err)
	}

	token := base64.RawURLEncoding.EncodeToString(data)

	_, err = decodeCursor(token)
	if err == nil {
		t.Fatal("decodeCursor() error = nil, want error")
	}
}
