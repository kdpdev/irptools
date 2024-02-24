package errs

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	Wrap   = WithFrame
	Error  = NewWithFrame
	Errorf = NewWithFramef
)

func NewWithStack(text string) error {
	return WithStackSkip(errors.New(text), 2)
}

func NewWithStackf(format string, args ...any) error {
	return WithStackSkip(fmt.Errorf(format, args...), 2)
}

func WithStack(err error) error {
	return WithStackSkip(err, 1)
}

func NewWithFrame(text string) error {
	return WithFrameSkip(errors.New(text), 1)
}

func NewWithFramef(format string, args ...any) error {
	return WithFrameSkip(fmt.Errorf(format, args...), 1)
}

func WithFrame(err error) error {
	return WithFrameSkip(err, 1)
}

func WithStackSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &errWithDetails{err: err, details: newStackDetails(skip)}
}

func WithFrameSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &errWithDetails{err: err, details: newFrameDetails(skip)}
}

func ensureWithoutW(format string) {
	if strings.Index(format, "%w") != -1 {
		panic("bad format:\n  The '%w' format option is not supported.\n  Use 'With...(fmt.Errorf(format, args))' instead.")
	}
}

type errWithDetails struct {
	err     error
	details errDetails
}

type errDetails interface {
	Print(w io.Writer) (int, error)
}

func (this *errWithDetails) String() string {
	return this.err.Error()
}

func (this *errWithDetails) Error() string {
	return this.err.Error()
}

func (this *errWithDetails) Unwrap() error {
	return this.err
}

func (this *errWithDetails) Format(s fmt.State, verb rune) {
	if formatter, ok := this.err.(fmt.Formatter); ok {
		formatter.Format(s, verb)
	} else {
		_, _ = fmt.Fprintf(s, string([]rune{'%', verb}), this.err)
	}

	if verb == 'v' {
		_, _ = this.details.Print(s)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newFrameDetails(skip int) *frameDetails {
	return &frameDetails{frame: getStack(skip)[0]}
}

func newStackDetails(skip int) *stackDetails {
	return &stackDetails{frames: getStack(skip)}
}

type frameDetails struct {
	frame
}

func (this *frameDetails) Print(w io.Writer) (int, error) {
	return printFrame(this.frame, "\n", w)
}

type stackDetails struct {
	frames
}

func (this *stackDetails) Print(w io.Writer) (int, error) {
	n, e := io.WriteString(w, "\n  stack:\n")
	if e != nil {
		return n, e
	}
	for _, f := range this.frames {
		nn, e := printFrame(f, "\n    ", w)
		n += nn
		if e != nil {
			return n, e
		}
	}
	return n, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func printFrame(f frame, prefix string, w io.Writer) (int, error) {
	return fmt.Fprintf(w, "%s%s:%d %s", prefix, frameFile(f), frameLine(f), frameShortFnName(f))
}

func getStack(skip int) frames {
	const mandatorySkip = 4
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+mandatorySkip, pcs[:])
	stack := make([]frame, n)
	for i, pc := range pcs[:n] {
		stack[i] = frame(pc)
	}
	return stack
}

type frame = uintptr
type frames = []frame

func framePc(f frame) uintptr {
	return uintptr(f) - 1
}

func frameFile(f frame) string {
	fn := runtime.FuncForPC(framePc(f))
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(framePc(f))
	return file
}

func frameLine(f frame) int {
	fn := runtime.FuncForPC(framePc(f))
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(framePc(f))
	return line
}

func frameFnName(f frame) string {
	fn := runtime.FuncForPC(framePc(f))
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func frameShortFnName(f frame) string {
	name := frameFnName(f)
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
