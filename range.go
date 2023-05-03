package requiring

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lufia/go-pointer"
	"golang.org/x/exp/constraints"
)

type ordered = constraints.Ordered

func InRange[T ordered](min, max *T, opts ...any) Validator {
	var r inRangeValidator[T]
	r.min = min
	r.max = max
	for _, opt := range opts {
		switch v := opt.(type) {
		case InRangeViolationPrinter[T]:
			r.p = v
		}
	}
	return &r
}

type inRangeValidator[T ordered] struct {
	name     string
	min, max *T
	p        InRangeViolationPrinter[T]
}

func (r *inRangeValidator[T]) SetName(name string) {
	r.name = name
}

func (r *inRangeValidator[T]) Validate(v any) error {
	n := v.(T)
	ok := true
	if r.min != nil && n < *r.min {
		ok = false
	}
	if r.max != nil && n > *r.max {
		ok = false
	}
	if !ok {
		return &InRangeViolationError[T]{
			Name:  r.name,
			Value: n,
			Min:   r.min,
			Max:   r.max,
			rule:  r,
		}
	}
	return nil
}

type InRangeViolationError[T ordered] struct {
	Name     string
	Value    T
	Min, Max *T
	rule     *inRangeValidator[T]
}

func (e InRangeViolationError[T]) Error() string {
	p := e.rule.p
	if p == nil {
		p = &inRangePrinter[T]{}
	}
	var w bytes.Buffer
	p.Print(&w, e)
	return w.String()
}

type inRangePrinter[T ordered] struct{}

func (inRangePrinter[T]) Print(w io.Writer, e InRangeViolationError[T]) {
	if e.Name != "" {
		fmt.Fprintf(w, "'%s' is not ", e.Name)
	}
	fmt.Fprintf(w, "in range(%v ... %v)",
		pointer.NewFormatter[T](e.Min),
		pointer.NewFormatter[T](e.Max))
}

type InRangeViolationPrinter[T ordered] interface {
	Printer[InRangeViolationError[T]]
}

var (
	_ Validator                    = (*inRangeValidator[int])(nil)
	_ ViolationError               = (*InRangeViolationError[int])(nil)
	_ InRangeViolationPrinter[int] = (*inRangePrinter[int])(nil)
)
