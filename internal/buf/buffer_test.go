package buf_test

import (
	"testing"

	"github.com/bilus/experiments-chat/internal/buf"
)

func TestBufferWritesAndReturnsText(t *testing.T) {
	var buffer buf.Buffer

	if _, err := buffer.Write([]byte("hello")); err != nil {
		t.Fatalf("write first chunk: %v", err)
	}
	if _, err := buffer.Write([]byte(" world")); err != nil {
		t.Fatalf("write second chunk: %v", err)
	}

	if got := buffer.String(); got != "hello world" {
		t.Fatalf("expected accumulated text, got %q", got)
	}
}
