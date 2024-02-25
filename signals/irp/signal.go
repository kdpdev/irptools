package irp

type Frequency = uint

type Address = [4]uint8
type Command = [4]uint8

type SignalCode struct {
	Address Address `json:"address"`
	Command Command `json:"command"`
}

type Micros = uint

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewSignalData() SignalData {
	return SignalData{}
}

type SignalData []Micros

func (this *SignalData) Add(durs ...Micros) {
	*this = append(*this, durs...)
}

func (this *SignalData) Pop() {
	if len(*this) > 0 {
		*this = (*this)[:len(*this)-1]
	}
}

func (this *SignalData) AddBits8(bits [8]bool, bit1Mark, bit1Space, bit0Mark, bit0Space Micros) {
	this.AddBits(bits[:], bit1Mark, bit1Space, bit0Mark, bit0Space)
}

func (this *SignalData) AddBits16(bits [16]bool, bit1Mark, bit1Space, bit0Mark, bit0Space Micros) {
	this.AddBits(bits[:], bit1Mark, bit1Space, bit0Mark, bit0Space)
}

func (this *SignalData) AddBits32(bits [32]bool, bit1Mark, bit1Space, bit0Mark, bit0Space Micros) {
	this.AddBits(bits[:], bit1Mark, bit1Space, bit0Mark, bit0Space)
}

func (this *SignalData) AddBits(bits []bool, bit1Mark, bit1Space, bit0Mark, bit0Space Micros) {
	for _, b := range bits {
		if b {
			this.Add(bit1Mark, bit1Space)
		} else {
			this.Add(bit0Mark, bit0Space)
		}
	}
}

func (this *SignalData) Duration() Micros {
	result := Micros(0)
	for _, d := range *this {
		result += d
	}
	return result
}

func (this *SignalData) Clone() SignalData {
	if this == nil {
		return nil
	}

	clone := make(SignalData, len(*this))
	copy(clone, *this)

	return clone
}

func (this *SignalData) WithPauseTie(minPause Micros) SignalData {
	length := len(*this)

	if length == 0 {
		*this = append(*this, 0, minPause)
		return *this
	}

	if length%2 == 1 {
		*this = append(*this, minPause)
		return *this
	}

	if (*this)[length-1] < minPause {
		(*this)[length-1] = minPause
	}

	return *this
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
