package connectapi

type Request[T Model, R Relationships] struct {
	Data RequestData[T, R] `json:"data"`
}

type RequestData[T Model, R Relationships] struct {
	ID            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"`
	Data          T      `json:"attributes"`
	Relationships R      `json:"relationships,omitempty"`
}
