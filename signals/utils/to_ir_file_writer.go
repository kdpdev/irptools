package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"irptools/signals/irp"
	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fs"
)

func NewIrFileWriter(filePath string) (*IrFileWriter, error) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	dirPath, _ := filepath.Split(filePath)
	_, err = fs.EnsureDirExists(dirPath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	const jsonExt = ".ir"
	if strings.LastIndex(filePath, jsonExt) != len(filePath)-len(jsonExt) {
		if filePath[len(filePath)-1] != '.' {
			filePath += "."
		}
		filePath += "ir"
	}

	file, err := fs.CreateWriteOnlyFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	_, err = fmt.Fprintln(file, "Filetype: IR signals file")
	if err != nil {
		return nil, errs.Wrap(err)
	}

	_, err = fmt.Fprintln(file, "Version: 1")
	if err != nil {
		return nil, errs.Wrap(err)
	}

	encoder := NewIrEncoder(file)

	return &IrFileWriter{
		file:    file,
		encoder: encoder,
	}, nil
}

type IrFileWriter struct {
	file    *os.File
	encoder *IrEncoder
}

func (this *IrFileWriter) Consume(signal signal.Signal) error {
	return this.encoder.Encode(signal)
}

func (this *IrFileWriter) Close() error {
	return this.file.Close()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewIrEncoder(writer io.Writer) *IrEncoder {
	return &IrEncoder{w: writer}
}

type IrEncoder struct {
	w io.Writer
}

func (this *IrEncoder) Encode(s signal.Signal) error {

	lines := []string{}
	lines = append(lines, "#", " ")
	lines = append(lines, "name: ", s.Function)
	if s.Protocol != "" {
		lines = append(lines, "type: ", "parsed")
		lines = append(lines, "protocol: ", s.Protocol)
		lines = append(lines, "address: ", this.format4Bytes(s.Code.Address))
		lines = append(lines, "command: ", this.format4Bytes(s.Code.Command))
	} else {
		lines = append(lines, "type: ", "raw")
		lines = append(lines, "frequency: ", strconv.FormatUint(uint64(s.Frequency), 10))
		lines = append(lines, "duty_cycle: ", "0.330000")
		lines = append(lines, "data: ", this.formatSignalData(s.Data))
	}

	for i := 0; i < len(lines); i += 2 {
		_, err := fmt.Fprint(this.w, lines[i])
		if err != nil {
			return errs.Wrap(err)
		}
		_, err = fmt.Fprintln(this.w, lines[i+1])
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (this *IrEncoder) formatSignalData(data irp.SignalData) string {
	str := fmt.Sprintf("%v", data)
	str = strings.TrimLeft(str, "[")
	str = strings.TrimRight(str, "]")
	str = strings.ReplaceAll(str, ",", " ")
	return str
}

func (this *IrEncoder) format4Bytes(data [4]uint8) string {
	lines := [...]string{
		fmt.Sprintf("%02X", data[0]),
		fmt.Sprintf("%02X", data[1]),
		fmt.Sprintf("%02X", data[2]),
		fmt.Sprintf("%02X", data[3]),
	}
	return strings.Join(lines[:], " ")
}
