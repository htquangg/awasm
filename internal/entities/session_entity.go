package entities

type AuthenticatorAssuranceLevel int

const (
	AAL0 AuthenticatorAssuranceLevel = iota
	AAL1
	AAL2
	AAL3
)

func (aal AuthenticatorAssuranceLevel) String() string {
	switch aal {
	case AAL0:
		return "aal0"
	case AAL1:
		return "aal1"
	case AAL2:
		return "aal2"
	case AAL3:
		return "aal3"
	default:
		return ""
	}
}
