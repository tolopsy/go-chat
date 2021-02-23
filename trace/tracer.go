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
