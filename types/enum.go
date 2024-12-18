package types

type EnumMapper interface {
	InMap() error
	ToMap() any
}
