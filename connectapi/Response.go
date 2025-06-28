package connectapi

type Response[T Model] struct {
	Data ResponseData[T] `json:"data"`
}

type ResponseData[T Model] struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Data T      `json:"attributes"`
}
