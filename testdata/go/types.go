package testdata

// BasicStruct demonstrates struct type parsing
type BasicStruct struct {
	Field1 string
	Field2 int
}

// EmptyStruct demonstrates empty struct
type EmptyStruct struct{}

// EmbeddedStruct demonstrates struct with embedded fields
type EmbeddedStruct struct {
	BasicStruct
	ExtraField bool
}

// SimpleInterface demonstrates interface parsing
type SimpleInterface interface {
	Method1() string
	Method2(int) error
}

// EmptyInterface demonstrates empty interface
type EmptyInterface interface{}

// EmbeddedInterface demonstrates interface composition
type EmbeddedInterface interface {
	SimpleInterface
	Method3() bool
}

// GenericType demonstrates generic type declaration
type GenericType[T any] struct {
	Value T
}

// ConstrainedGeneric demonstrates generic with constraints
type ConstrainedGeneric[T comparable] struct {
	Key   T
	Value string
}

// MultipleGenerics demonstrates multiple type parameters
type MultipleGenerics[K comparable, V any] map[K]V

// TypeAlias demonstrates type alias
type TypeAlias = string

// CustomType demonstrates custom type based on existing type
type CustomType string

// Constants for testing const parsing
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	MaxRetries     = 5
)

// Single constant
const DefaultTimeout = 30

// Variables for testing var parsing
var (
	GlobalCounter int
	SystemReady   bool = true
	ConfigPath    string
)

// Single variable
var DefaultConfig = BasicStruct{
	Field1: "default",
	Field2: 42,
}
