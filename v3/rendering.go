package v3

type Rendering struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Label  Label  `json:"label,omitempty"`
	Format string `json:"format,omitempty"`
}
