package misc

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"irptools/utils/alg"
)

func SPrettyPrint(obj interface{}) string {
	b := strings.Builder{}
	_, _ = FPrettyPrint(&b, obj)
	return b.String()
}

func FPrettyPrint(writer io.Writer, obj interface{}) (int, error) {

	var impl func(w io.Writer, obj interface{}, indent int) (int, error)

	impl = func(w io.Writer, obj interface{}, indent int) (int, error) {
		if obj == nil {
			return fmt.Fprint(w, obj)
		}

		typeOf := reflect.TypeOf(obj)
		if typeOf.Kind() != reflect.Struct {
			return fmt.Fprint(w, obj)
		}

		val := reflect.ValueOf(obj)
		numField := val.NumField()
		n, err := fmt.Fprint(w, val.Type().Name(), " {")
		if err != nil {
			return n, err
		}

		if numField == 0 {
			nn, err := fmt.Fprint(w, "}")
			n += nn
			return n, err
		}

		nn, err := fmt.Fprint(w, "\n")
		n += nn
		if err != nil {
			return n, err
		}

		names := make([]string, 0, numField)
		values := make([]interface{}, 0, numField)
		for i := 0; i < numField; i++ {
			names = append(names, val.Type().Field(i).Name)
			values = append(values, val.Field(i).Interface())
		}

		maxNameIdx := alg.MaxElemIdxIf(names, func(lhs, rhs int) bool {
			return len(names[lhs]) < len(names[rhs])
		})

		maxNameLen := len(names[maxNameIdx])
		lineFmt := fmt.Sprintf("  %%-%vv%%-%vv = %%v\n", indent, maxNameLen)
		for i := 0; i < len(names); i++ {
			b := strings.Builder{}
			_, err := impl(&b, values[i], indent+maxNameLen+len(" = { "))
			if err != nil {
				return n, err
			}
			nn, err = fmt.Fprintf(w, lineFmt, "", names[i], b.String())
			n += nn
			if err != nil {
				return n, err
			}
		}

		nn, err = fmt.Fprintf(w, fmt.Sprintf("%%+%vv", indent+1), "}")
		n += nn

		return n, err
	}

	return impl(writer, obj, 0)
}
