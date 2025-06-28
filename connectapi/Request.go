package connectapi

type Request[T Model] struct {
	Data RequestData[T] `json:"data"`
}

type RequestData[T Model] struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Data Model  `json:"attributes"`
}
