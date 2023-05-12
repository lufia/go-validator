package validator

import (
	"bytes"
	"context"
	"io"
	"reflect"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

var (
	defaultLanguage = language.English
	DefaultCatalog  = catalog.NewBuilder(catalog.Fallback(defaultLanguage))
	defaultPrinter  = message.NewPrinter(defaultLanguage, message.Catalog(DefaultCatalog))
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
	return &namedArg{name: name}
}

type Arg interface {
	ValueOf(v any) any
}

type namedArg struct {
	name string
}

func (a *namedArg) ValueOf(v any) any {
	p := reflect.ValueOf(v)
	if p.Kind() == reflect.Pointer {
		p = p.Elem()
	}
	for _, f := range reflect.VisibleFields(p.Type()) {
		if f.Tag.Get("arg") != a.name {
			continue
		}
		v := p.FieldByIndex(f.Index)
		return v.Interface()
	}
	return nil
}
