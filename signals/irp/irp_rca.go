package irp

// https://www.sbprojects.net/knowledge/ir/rca.php

const (
	FrequencyRca = 56000
)

func NewIrpRca(protocol string) Irp {
	return NewIrpUnsupported(protocol, FrequencyRca)
}
