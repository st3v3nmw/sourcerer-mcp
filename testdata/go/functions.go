package testdata

import (
	"context"
	"fmt"
)

type Config struct {
	Name string
}

type Result struct {
	Success bool
	Data    interface{}
}

// SimpleFunction demonstrates basic function parsing
func SimpleFunction(x int) string {
	return fmt.Sprintf("%d", x)
}

// MultipleParams shows function with multiple parameters and return values
func MultipleParams(a string, b int, c bool) (string, error) {
	if c {
		return fmt.Sprintf("%s-%d", a, b), nil
	}

	return "", fmt.Errorf("invalid")
}

// NoParams function with no parameters
func NoParams() {
	fmt.Println("no params")
}

// NoReturn function with no return values
func NoReturn(x int) {
	fmt.Printf("got %d\n", x)
}

// EmptyFunction with empty body
func EmptyFunction() {}

// ComplexSignature with various parameter types
func ComplexSignature(ctx context.Context, data map[string]interface{}, opts ...func(*Config)) (*Result, error) {
	return &Result{Success: true}, nil
}

// VariadicFunction with variadic parameters
func VariadicFunction(first string, others ...int) int {
	return len(others)
}

// GenericFunction with type parameters
func GenericFunction[T any](items []T) T {
	var zero T
	if len(items) == 0 {
		return zero
	}

	return items[0]
}

// DuplicateNameFunction - testing duplicate function names
func DuplicateNameFunction() string {
	return "duplicate name"
}

// DuplicateNameFunction (2) - testing duplicate function names
func DuplicateNameFunction() string {
	return "duplicate name"
}

// DuplicateNameFunction (3) - testing duplicate function names
func DuplicateNameFunction() string {
	return "duplicate name"
}
