package irp

func NewIrpUnsupported(protocol string, frequency Frequency) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: frequency,
		decode: func(code SignalCode) (SignalData, error) {
			return SignalData{}, NewUnsupportedProtocolError(protocol)
		},
	}
}
