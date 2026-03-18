package v3

type Canvas struct {
	ID          string           `json:"id,omitempty"`
	Type        string           `json:"type,omitempty"`
	Label       Label            `json:"label,omitempty"`
	Height      int              `json:"height,omitempty"`
	Width       int              `json:"width,omitempty"`
	Items       []AnnotationPage `json:"items,omitempty"`
	Annotations []AnnotationPage `json:"annotations,omitempty"`
}
