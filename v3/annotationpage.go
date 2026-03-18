package v3

type AnnotationPage struct {
	ID    string       `json:"id,omitempty"`
	Type  string       `json:"type,omitempty"`
	Items []Annotation `json:"items,omitempty"`
}
