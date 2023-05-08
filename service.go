package micro

type Service[T any] func(ctx *Context[T]) (interface{}, Error)
