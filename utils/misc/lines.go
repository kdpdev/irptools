package misc

import (
	"bufio"
	"io"

	"irptools/utils/errs"
)

func NewStreamLinesIter(stream io.Reader) func() (string, bool, error) {
	scanner := bufio.NewScanner(stream)
	return func() (string, bool, error) {
		ok := scanner.Scan()
		if !ok {
			return "", false, errs.Wrap(scanner.Err())
		}
		line := scanner.Text()
		return line, true, nil
	}
}

func EnumStreamLines(stream io.ReadCloser, callback func(string) (bool, error)) (err error) {
	next := NewStreamLinesIter(stream)

	defer func() {
		closeErr := stream.Close()
		if err == nil {
			err = errs.Wrap(closeErr)
		}
	}()

	var line string
	var ok bool

	for {
		line, ok, err = next()
		if !ok || err != nil {
			return errs.Wrap(err)
		}

		ok, err = callback(line)
		if !ok || err != nil {
			return errs.Wrap(err)
		}
	}
}
