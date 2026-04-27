package connectapi

type ListResponse[T Model, R Relationships] struct {
	Data  []ResponseData[T, R] `json:"data"`
	Links ListLinks            `json:"links"`
}

type ListLinks struct {
	Next string `json:"next,omitempty"`
}
