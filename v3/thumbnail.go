package v3

type Thumbnail struct {
	ID      string    `json:"id,omitempty"`
	Type    string    `json:"type,omitempty"`
	Format  string    `json:"format,omitempty"`
	Service []Service `json:"service,omitempty"`
}
