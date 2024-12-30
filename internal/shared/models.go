package shared

type CommonResponse struct {
	Message string `json:"message"`
}

type DataResponse[T any] struct {
	Data T `json:"data"`
}
