package validator

import (
	"bytes"
	"context"
	"io"
	"reflect"

	"golang.org/x/text/message"
)

type printerKey struct{}

// Printer is the interface that wraps Fprintf method.
type Printer interface {
	Fprintf(w io.Writer, key message.Reference, a ...any) (int, error)
}

func WithPrinter(ctx context.Context, p Printer) context.Context {
	return context.WithValue(ctx, printerKey{}, p)
}

func ctxPrint(ctx context.Context, v any, key message.Reference, args []Arg) string {
	var (
		w bytes.Buffer
		p Printer = defaultPrinter
	)
	if pp := ctx.Value(printerKey{}); pp != nil {
		p = pp.(Printer)
	}

	a := make([]any, len(args))
	for i, arg := range args {
		a[i] = arg.ValueOf(v)
	}
	p.Fprintf(&w, key, a...)
	return w.String()
}

func ByName(name string) Arg {
	return &namedArg{names: []string{name}}
}

// ByPath refers to a object by names.
func ByPath(names ...string) Arg {
	return &namedArg{names: names}
}

type Arg interface {
	ValueOf(v any) any
}

type namedArg struct {
	names []string
}

func (a *namedArg) ValueOf(v any) any {
	if len(a.names) == 0 {
		return nil
	}
	for _, name := range a.names {
		n, ok := valueOf(v, name)
		if !ok {
			return nil
		}
		v = n
	}
	return v
}

func valueOf(v any, name string) (any, bool) {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	for _, f := range reflect.VisibleFields(p.Type()) {
		if f.Tag.Get("arg") != name {
			continue
		}
		v := p.FieldByIndex(f.Index)
		return v.Interface(), true
	}
	return nil, false
}
