package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		tracer.Trace("Trace Package on scene!")
		if buf.String() != "Trace Package on scene!\n" {
			t.Errorf("Tracer should not write '%s'.", buf.String())
		}
	}
}

func New() {}
