package api

const (
	StringValue = iota
	IntValue
	FloatValue
)

type Config[T interface{}] struct {
	value T
}

type Service interface {
	GetConfig() Config
	SetConfig(config Config)
}
