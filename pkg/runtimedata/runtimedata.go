package runtimedata

type LayoutMap interface {
	Layout() RuntimeData
}

type RuntimeData interface {
	Data() ([]byte, error)
}
