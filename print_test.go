package validator

import (
	"context"
	"fmt"
	"io"
	"testing"

	"golang.org/x/text/message"
)

type testPrinter struct{}

func (*testPrinter) Fprintf(w io.Writer, key message.Reference, a ...any) (int, error) {
	format := key.(string)
	format = "[" + format + "]"
	return fmt.Fprintf(w, format, a...)
}

func TestContextDefault(t *testing.T) {
	type Err struct {
		Body string `arg:"body"`
	}
	ctx := context.Background()
	e := &Err{Body: "hello"}
	format := "err: %v"
	s := ctxPrint(ctx, e, format, []Arg{ByName("body")})
	if want := "err: hello"; s != want {
		t.Errorf("ctxPrint(%q) = %q; want %q", format, s, want)
	}
}

func TestContextPrinter(t *testing.T) {
	type Err struct {
		Body string `arg:"body"`
	}
	ctx := WithPrinter(context.Background(), &testPrinter{})
	e := &Err{Body: "hello"}
	format := "err: %v"
	s := ctxPrint(ctx, e, format, []Arg{ByName("body")})
	if want := "[err: hello]"; s != want {
		t.Errorf("ctxPrint(%q) writes %q; want %q", format, s, want)
	}
}
