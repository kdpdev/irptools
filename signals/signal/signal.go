package signal

import (
	"fmt"

	"irptools/signals/irp"
	jsonutils "irptools/utils/json"
)

type Signal struct {
	Id        string         `json:"id"`
	Source    string         `json:"source"`
	Brand     string         `json:"brand"`
	Device    string         `json:"device"`
	Function  string         `json:"function"`
	Protocol  string         `json:"protocol"`
	Frequency irp.Frequency  `json:"frequency"`
	Data      irp.SignalData `json:"data"`
	Code      irp.SignalCode `json:"code"`
}

func (this *Signal) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		m := make(map[string]interface{})
		err := jsonutils.Cast(this, &m)
		if err == nil {
			_, _ = s.Write([]byte("signal.Signal:"))
			for k, v := range m {
				_, _ = fmt.Fprintf(s, "\n  %s: %v", k, v)
			}
		} else {
			_, _ = fmt.Fprintf(s, "%v", err)
		}
	default:
		_, _ = s.Write([]byte("signal.Signal"))
	}
}
