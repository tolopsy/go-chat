package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(i ...interface{}) {
	fmt.Fprint(t.out, i...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// A tracer that does nothing
type nilTracer struct{}

func (n *nilTracer) Trace(i ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}
