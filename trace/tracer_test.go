package trace

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	newTracer := New(&buf)

	if newTracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		newTracer.Trace("Trace Package on scene!")
		if buf.String() != "Trace Package on scene!\n" {
			t.Errorf("Tracer should not write '%s'.", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
