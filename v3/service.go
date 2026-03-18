package v3

type Service struct {
	ID      string    `json:"id,omitempty"`
	AtID    string    `json:"@id,omitempty"`
	Type    string    `json:"type,omitempty"`
	AtType  string    `json:"@type,omitempty"`
	Profile string    `json:"profile,omitempty"`
	Service []Service `json:"service,omitempty"`
}
