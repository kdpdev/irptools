package irp

// http://www.hifi-remote.com/johnsfine/DecodeIR.html

var supportedIrpCreators = map[string]func(protocol string) Irp{
	"kaseikyo":  NewIrpKaseikyo,
	"nec":       NewIrpNec,
	"necext":    NewIrpNecExt,
	"nec42":     NewIrpNec42,
	"rc5":       NewIrpRc5,
	"rc6":       NewIrpRc6,
	"samsung32": NewIrpSamsung32,
	"sirc":      NewIrpSirc12,
	"sirc12":    NewIrpSirc12,
	"sirc15":    NewIrpSirc15,
}

type Irp interface {
	Protocol() string
	Frequency() Frequency
	Decode(code SignalCode) (SignalData, error)
}

func GetIrp(protocol string) (Irp, error) {
	irp := supportedIrps[protocol]
	if irp == nil {
		return nil, Wrap(NewUnsupportedProtocolError(protocol))
	}
	return irp, nil
}

var supportedIrps = createSupportedIrps()

func createSupportedIrps() map[string]Irp {
	result := map[string]Irp{}
	for protocol, create := range supportedIrpCreators {
		result[protocol] = create(protocol)
	}
	return result
}

type irpImpl struct {
	protocol  string
	frequency Frequency
	decode    func(code SignalCode) (SignalData, error)
}

func (this *irpImpl) Protocol() string {
	return this.protocol
}

func (this *irpImpl) Frequency() Frequency {
	return this.frequency
}

func (this *irpImpl) Decode(code SignalCode) (SignalData, error) {
	return this.decode(code)
}
