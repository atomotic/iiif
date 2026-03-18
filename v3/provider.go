package v3

type Provider struct {
	ID       string     `json:"id,omitempty"`
	Type     string     `json:"type,omitempty"`
	Label    Label      `json:"label,omitempty"`
	Homepage []Homepage `json:"homepage,omitempty"`
	Logo     []Logo     `json:"logo,omitempty"`
	SeeAlso  []SeeAlso  `json:"seeAlso,omitempty"`
}

type Logo struct {
	ID      string    `json:"id,omitempty"`
	Type    string    `json:"type,omitempty"`
	Format  string    `json:"format,omitempty"`
	Service []Service `json:"service,omitempty"`
}
