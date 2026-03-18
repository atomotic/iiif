package v3

type Structure struct {
	ID            string             `json:"id,omitempty"`
	Type          string             `json:"type,omitempty"`
	Label         Label              `json:"label,omitempty"`
	Supplementary *Supplementary     `json:"supplementary,omitempty"`
	Items         []StructureItem    `json:"items,omitempty"`
}

type Supplementary struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type StructureItem struct {
	ID            string          `json:"id,omitempty"`
	Type          string          `json:"type,omitempty"`
	Label         Label           `json:"label,omitempty"`
	Supplementary *Supplementary  `json:"supplementary,omitempty"`
	Items         []StructureItem `json:"items,omitempty"`
	Source        string          `json:"source,omitempty"`
	Selector      *Selector       `json:"selector,omitempty"`
}

type Selector struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
