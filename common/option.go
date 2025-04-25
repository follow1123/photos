package common

type Option[T any] interface {
	Apply(*T)
}

type OptionFunc[T any] func(*T)

func (f OptionFunc[T]) Apply(opts *T) {
	f(opts)
}
