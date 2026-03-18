package v3

type Metadata struct {
	Label Label `json:"label,omitempty"`
	Value Label `json:"value,omitempty"`
}

func NewMetadata(label string, value string) Metadata {
	return Metadata{
		Label: Label{"en": {label}},
		Value: Label{"none": {value}},
	}
}
