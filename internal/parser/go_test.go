package parser_test

import (
	"testing"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type GoParserTestSuite struct {
	ParserBaseTestSuite
}

func (s *GoParserTestSuite) SetupSuite() {
	s.ParserBaseTestSuite.SetupSuite()

	var err error
	s.parser, err = parser.NewGoParser(s.workspaceRoot)
	s.Require().NoError(err)
}

func (s *GoParserTestSuite) TestFunctionParsing() {
	chunks := s.getChunks("go/functions.go")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:      "Package Comments Hashing",
			path:      "f2bcc925c6085e27",
			summary:   "// Function tests",
			source:    `// Function tests`,
			startLine: 1,
			endLine:   1,
		},
		{
			name:    "Imports Hashing",
			path:    "44983311c5db2e3",
			summary: "import (",
			source: `import (
	"context"
	"fmt"
)`,
			startLine: 5,
			endLine:   8,
		},
		{
			name:    "Simple Function",
			path:    "SimpleFunction",
			summary: "func SimpleFunction(x int) string {",
			source: `// SimpleFunction demonstrates basic function parsing
func SimpleFunction(x int) string {
	return fmt.Sprintf("%d", x)
}`,
			startLine: 19,
			endLine:   22,
		},
		{
			name:    "Multiple Params Function",
			path:    "MultipleParams",
			summary: "func MultipleParams(a string, b int, c bool) (string, error) {",
			source: `// MultipleParams shows function with multiple parameters and return values
func MultipleParams(a string, b int, c bool) (string, error) {
	if c {
		return fmt.Sprintf("%s-%d", a, b), nil
	}

	return "", fmt.Errorf("invalid")
}`,
			startLine: 24,
			endLine:   31,
		},
		{
			name:    "No Params Function",
			path:    "NoParams",
			summary: "func NoParams() {",
			source: `// NoParams function with no parameters
func NoParams() {
	fmt.Println("no params")
}`,
			startLine: 33,
			endLine:   36,
		},
		{
			name:    "No Return Function",
			path:    "NoReturn",
			summary: "func NoReturn(x int) {",
			source: `// NoReturn function with no return values
func NoReturn(x int) {
	fmt.Printf("got %d\n", x)
}`,
			startLine: 38,
			endLine:   41,
		},
		{
			name:    "Empty Function",
			path:    "EmptyFunction",
			summary: "func EmptyFunction() {}",
			source: `// EmptyFunction with empty body
func EmptyFunction() {}`,
			startLine: 43,
			endLine:   44,
		},
		{
			name:    "Complex Signature",
			path:    "ComplexSignature",
			summary: "func ComplexSignature(ctx context.Context, data map[string]interface{}, opts ...func(*Config))...",
			source: `// ComplexSignature with various parameter types
func ComplexSignature(ctx context.Context, data map[string]interface{}, opts ...func(*Config)) (*Result, error) {
	return &Result{Success: true}, nil
}`,
			startLine: 46,
			endLine:   49,
		},
		{
			name:    "Variadic Function",
			path:    "VariadicFunction",
			summary: "func VariadicFunction(first string, others ...int) int {",
			source: `// VariadicFunction with variadic parameters
func VariadicFunction(first string, others ...int) int {
	return len(others)
}`,
			startLine: 51,
			endLine:   54,
		},
		{
			name:    "Generic Function",
			path:    "GenericFunction",
			summary: "func GenericFunction[T any](items []T) T {",
			source: `// GenericFunction with type parameters
func GenericFunction[T any](items []T) T {
	var zero T
	if len(items) == 0 {
		return zero
	}

	return items[0]
}`,
			startLine: 56,
			endLine:   64,
		},
		{
			name:    "Duplicate Name Function",
			path:    "DuplicateNameFunction",
			summary: "func DuplicateNameFunction() string {",
			source: `// DuplicateNameFunction - testing duplicate function names

func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 66,
			endLine:   70,
		},
		{
			name:    "Duplicate Name Function - 2",
			path:    "DuplicateNameFunction-2",
			summary: "func DuplicateNameFunction() string {",
			source: `// DuplicateNameFunction (2)
// Testing duplicate function names
func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 72,
			endLine:   76,
		},
		{
			name:    "Duplicate Name Function - 3",
			path:    "DuplicateNameFunction-3",
			summary: "func DuplicateNameFunction() string {",
			source: `// DuplicateNameFunction (3)
//
// Testing duplicate function names
func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 78,
			endLine:   83,
		},
		{
			name:      "Standalone Comment Hashing",
			path:      "d5d69632b4d3fba5",
			summary:   "// A standalone comment",
			source:    `// A standalone comment`,
			startLine: 85,
			endLine:   85,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			chunk, exists := chunks[test.path]
			s.Require().True(exists, "chunk %s not found", test.path)
			s.Require().NotNil(chunk)

			s.Equal("src", chunk.Type)
			s.Equal(test.path, chunk.Path)
			s.Equal(test.summary, chunk.Summary)
			s.Equal(test.source, chunk.Source)
			s.Equal(test.startLine, int(chunk.StartLine))
			s.Equal(test.endLine, int(chunk.EndLine))
			s.Equal("go/functions.go::"+test.path, chunk.ID())
		})
	}
}

func (s *GoParserTestSuite) TestMethodParsing() {
	chunks := s.getChunks("go/methods.go")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Value Receiver Method",
			path:    "User::GetName",
			summary: "func (u User) GetName() string {",
			source: `// GetName is a value receiver method
func (u User) GetName() string {
	return u.Name
}`,
			startLine: 9,
			endLine:   12,
		},
		{
			name:    "Pointer Receiver Method",
			path:    "User::SetName",
			summary: "func (u *User) SetName(name string) {",
			source: `// SetName is a pointer receiver method
func (u *User) SetName(name string) {
	u.Name = name
}`,
			startLine: 14,
			endLine:   17,
		},
		{
			name:    "Service Add User Method",
			path:    "Service::AddUser",
			summary: "func (s *Service) AddUser(user User) {",
			source: `// AddUser adds a user to the service
func (s *Service) AddUser(user User) {
	s.users = append(s.users, user)
}`,
			startLine: 24,
			endLine:   27,
		},
		{
			name:    "Service Find User Method",
			path:    "Service::FindUser",
			summary: "func (s *Service) FindUser(id int) *User {",
			source: `// FindUser finds a user by ID
func (s *Service) FindUser(id int) *User {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i]
		}
	}

	return nil
}`,
			startLine: 29,
			endLine:   38,
		},
		{
			name:    "Service Count Method",
			path:    "Service::Count",
			summary: "func (s Service) Count() int {",
			source: `// Count returns the number of users
func (s Service) Count() int {
	return len(s.users)
}`,
			startLine: 40,
			endLine:   43,
		},
		{
			name:    "Generic Pointer Receiver Method",
			path:    "Repository::Add",
			summary: "func (r *Repository[T]) Add(item T) {",
			source: `// Add adds an item to the repository
func (r *Repository[T]) Add(item T) {
	r.items = append(r.items, item)
}`,
			startLine: 50,
			endLine:   53,
		},
		{
			name:    "Generic Value Receiver Method",
			path:    "Repository::Get",
			summary: "func (r Repository[T]) Get(index int) T {",
			source: `// Get retrieves an item by index
func (r Repository[T]) Get(index int) T {
	return r.items[index]
}`,
			startLine: 55,
			endLine:   58,
		},
		{
			name:    "Generic Size Method",
			path:    "Repository::Size",
			summary: "func (r Repository[T]) Size() int {",
			source: `// Size returns the number of items
func (r Repository[T]) Size() int {
	return len(r.items)
}`,
			startLine: 60,
			endLine:   63,
		},
		{
			name:    "Service A Helper Method",
			path:    "ServiceA::Helper",
			summary: "func (s ServiceA) Helper() string {",
			source: `func (s ServiceA) Helper() string {
	return "service A helper"
}`,
			startLine: 69,
			endLine:   71,
		},
		{
			name:    "Service B Helper Method",
			path:    "ServiceB::Helper",
			summary: "func (s ServiceB) Helper() string {",
			source: `func (s ServiceB) Helper() string {
	return "service B helper"
}`,
			startLine: 73,
			endLine:   75,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			chunk, exists := chunks[test.path]
			s.Require().True(exists, "chunk %s not found", test.path)
			s.Require().NotNil(chunk)

			s.Equal("src", chunk.Type)
			s.Equal(test.path, chunk.Path)
			s.Equal(test.summary, chunk.Summary)
			s.Equal(test.source, chunk.Source)
			s.Equal(test.startLine, int(chunk.StartLine))
			s.Equal(test.endLine, int(chunk.EndLine))
			s.Equal("go/methods.go::"+test.path, chunk.ID())
		})
	}
}

func (s *GoParserTestSuite) TestTypeParsing() {
	chunks := s.getChunks("go/types.go")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Basic Struct",
			path:    "BasicStruct",
			summary: "type BasicStruct struct {",
			source: `// BasicStruct demonstrates struct type parsing
type BasicStruct struct {
	Field1 string
	Field2 int
}`,
			startLine: 3,
			endLine:   7,
		},
		{
			name:    "Empty Struct",
			path:    "EmptyStruct",
			summary: "type EmptyStruct struct{}",
			source: `// EmptyStruct demonstrates empty struct
type EmptyStruct struct{}`,
			startLine: 9,
			endLine:   10,
		},
		{
			name:    "Embedded Struct",
			path:    "EmbeddedStruct",
			summary: "type EmbeddedStruct struct {",
			source: `// EmbeddedStruct demonstrates struct with embedded fields
type EmbeddedStruct struct {
	BasicStruct
	ExtraField bool
}`,
			startLine: 12,
			endLine:   16,
		},
		{
			name:    "Simple Interface",
			path:    "SimpleInterface",
			summary: "type SimpleInterface interface {",
			source: `// SimpleInterface demonstrates interface parsing
type SimpleInterface interface {
	Method1() string
	Method2(int) error
}`,
			startLine: 18,
			endLine:   22,
		},
		{
			name:    "Empty Interface",
			path:    "EmptyInterface",
			summary: "type EmptyInterface interface{}",
			source: `// EmptyInterface demonstrates empty interface
type EmptyInterface interface{}`,
			startLine: 24,
			endLine:   25,
		},
		{
			name:    "Embedded Interface",
			path:    "EmbeddedInterface",
			summary: "type EmbeddedInterface interface {",
			source: `// EmbeddedInterface demonstrates interface composition
type EmbeddedInterface interface {
	SimpleInterface
	Method3() bool
}`,
			startLine: 27,
			endLine:   31,
		},
		{
			name:    "Generic Type",
			path:    "GenericType",
			summary: "type GenericType[T any] struct {",
			source: `// GenericType demonstrates generic type declaration
type GenericType[T any] struct {
	Value T
}`,
			startLine: 33,
			endLine:   36,
		},
		{
			name:    "Constrained Generic",
			path:    "ConstrainedGeneric",
			summary: "type ConstrainedGeneric[T comparable] struct {",
			source: `// ConstrainedGeneric demonstrates generic with constraints
type ConstrainedGeneric[T comparable] struct {
	Key   T
	Value string
}`,
			startLine: 38,
			endLine:   42,
		},
		{
			name:    "Multiple Generics",
			path:    "MultipleGenerics",
			summary: "type MultipleGenerics[K comparable, V any] map[K]V",
			source: `// MultipleGenerics demonstrates multiple type parameters
type MultipleGenerics[K comparable, V any] map[K]V`,
			startLine: 44,
			endLine:   45,
		},
		{
			name:    "Type Alias",
			path:    "TypeAlias",
			summary: "type TypeAlias = string",
			source: `// TypeAlias demonstrates type alias
type TypeAlias = string`,
			startLine: 47,
			endLine:   48,
		},
		{
			name:    "Custom Type",
			path:    "CustomType",
			summary: "type CustomType string",
			source: `// CustomType demonstrates custom type based on existing type
type CustomType string`,
			startLine: 50,
			endLine:   51,
		},
		{
			name:    "Consts Block Hashing",
			path:    "796012d7ee1311f0",
			summary: "const (",
			source: `// Constants for testing const parsing
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	MaxRetries     = 5
)`,
			startLine: 53,
			endLine:   58,
		}, {
			name:    "Single Constant",
			path:    "DefaultTimeout",
			summary: "const DefaultTimeout = 30",
			source: `// Single constant
const DefaultTimeout = 30`,
			startLine: 60,
			endLine:   61,
		},
		{
			name:    "Vars Block Hashing",
			path:    "3a12064c73be460c",
			summary: "var (",
			source: `// Variables for testing var parsing
var (
	GlobalCounter int
	SystemReady   bool = true
	ConfigPath    string
)`,
			startLine: 63,
			endLine:   68,
		},
		{
			name:    "Single Variable",
			path:    "DefaultConfig",
			summary: "var DefaultConfig = BasicStruct{",
			source: `// Single variable
var DefaultConfig = BasicStruct{
	Field1: "default",
	Field2: 42,
}`,
			startLine: 70,
			endLine:   74,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			chunk, exists := chunks[test.path]
			s.Require().True(exists, "chunk %s not found", test.path)
			s.Require().NotNil(chunk)

			s.Equal("src", chunk.Type)
			s.Equal(test.path, chunk.Path)
			s.Equal(test.summary, chunk.Summary)
			s.Equal(test.source, chunk.Source)
			s.Equal(test.startLine, int(chunk.StartLine))
			s.Equal(test.endLine, int(chunk.EndLine))
			s.Equal("go/types.go::"+test.path, chunk.ID())
		})
	}
}

func (s *GoParserTestSuite) TestTestFileParsing() {
	chunks := s.getChunks("go/tests_test.go")

	chunk, exists := chunks["TestSimple"]
	s.Require().True(exists, "chunk %s not found", "TestSimple")
	s.Require().NotNil(chunk)

	s.Equal("tests", chunk.Type)
	s.Equal("TestSimple", chunk.Path)
	s.Equal("func TestSimple(t *testing.T) {", chunk.Summary)
	s.Equal(`// TestSimple is a basic test function
func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}`, chunk.Source)
	s.Equal(5, int(chunk.StartLine))
	s.Equal(10, int(chunk.EndLine))
	s.Equal("go/tests_test.go::TestSimple", chunk.ID())
}

func TestGoParserTestSuite(t *testing.T) {
	suite.Run(t, new(GoParserTestSuite))
}
