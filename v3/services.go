package v3

type Services struct {
	ID      string    `json:"@id,omitempty"`
	Type    string    `json:"@type,omitempty"`
	Profile string    `json:"profile,omitempty"`
	Label   string    `json:"label,omitempty"`
	Service []Service `json:"service,omitempty"`
}
