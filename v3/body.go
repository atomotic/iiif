package v3

type Body struct {
	// Common fields
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`

	// Image/Video/Audio body fields
	Format  string    `json:"format,omitempty"`
	Service []Service `json:"service,omitempty"`
	Height  int       `json:"height,omitempty"`
	Width   int       `json:"width,omitempty"`

	// TextualBody fields
	Value    string `json:"value,omitempty"`
	Language string `json:"language,omitempty"`
}
