package micro

type Handler[T any] func() HandlerResponse[T]

type HandlerResponse[T any] struct {
	Service    Service[T]
	Response   interface{}
	Pagination bool
	Sort       bool
}
