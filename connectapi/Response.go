package connectapi

type Response[T Model, R Relationships] struct {
	Data ResponseData[T, R] `json:"data"`
}

type ResponseData[T Model, R Relationships] struct {
	ID            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"`
	Data          T      `json:"attributes"`
	Relationships R      `json:"relationships,omitempty"`
}
