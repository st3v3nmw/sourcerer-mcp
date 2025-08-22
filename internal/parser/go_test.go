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
		startLine uint
		endLine   uint
	}{
		{
			name:    "Imports Hashing",
			path:    "44983311c5db2e3",
			summary: "import (",
			source: `import (
	"context"
	"fmt"
)`,
			startLine: 3,
			endLine:   6,
		},
		{
			name:      "Comments Hashing",
			path:      "399e014a89b03d0a",
			summary:   "// SimpleFunction demonstrates basic function parsing",
			source:    `// SimpleFunction demonstrates basic function parsing`,
			startLine: 17,
			endLine:   17,
		},
		{
			name:    "Simple Function",
			path:    "SimpleFunction",
			summary: "func SimpleFunction(x int) string {",
			source: `func SimpleFunction(x int) string {
	return fmt.Sprintf("%d", x)
}`,
			startLine: 18,
			endLine:   20,
		},
		{
			name:    "Multiple Params Function",
			path:    "MultipleParams",
			summary: "func MultipleParams(a string, b int, c bool) (string, error) {",
			source: `func MultipleParams(a string, b int, c bool) (string, error) {
	if c {
		return fmt.Sprintf("%s-%d", a, b), nil
	}

	return "", fmt.Errorf("invalid")
}`,
			startLine: 23,
			endLine:   29,
		},
		{
			name:    "No Params Function",
			path:    "NoParams",
			summary: "func NoParams() {",
			source: `func NoParams() {
	fmt.Println("no params")
}`,
			startLine: 32,
			endLine:   34,
		},
		{
			name:    "No Return Function",
			path:    "NoReturn",
			summary: "func NoReturn(x int) {",
			source: `func NoReturn(x int) {
	fmt.Printf("got %d\n", x)
}`,
			startLine: 37,
			endLine:   39,
		},
		{
			name:      "Empty Function",
			path:      "EmptyFunction",
			summary:   "func EmptyFunction() {}",
			source:    `func EmptyFunction() {}`,
			startLine: 42,
			endLine:   42,
		},
		{
			name:    "Complex Signature",
			path:    "ComplexSignature",
			summary: "func ComplexSignature(ctx context.Context, data map[string]interface{}, opts ...func(*Config))...",
			source: `func ComplexSignature(ctx context.Context, data map[string]interface{}, opts ...func(*Config)) (*Result, error) {
	return &Result{Success: true}, nil
}`,
			startLine: 45,
			endLine:   47,
		},
		{
			name:    "Variadic Function",
			path:    "VariadicFunction",
			summary: "func VariadicFunction(first string, others ...int) int {",
			source: `func VariadicFunction(first string, others ...int) int {
	return len(others)
}`,
			startLine: 50,
			endLine:   52,
		},
		{
			name:    "Generic Function",
			path:    "GenericFunction",
			summary: "func GenericFunction[T any](items []T) T {",
			source: `func GenericFunction[T any](items []T) T {
	var zero T
	if len(items) == 0 {
		return zero
	}

	return items[0]
}`,
			startLine: 55,
			endLine:   62,
		},
		{
			name:    "Duplicate Name Function",
			path:    "DuplicateNameFunction",
			summary: "func DuplicateNameFunction() string {",
			source: `func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 65,
			endLine:   67,
		},
		{
			name:    "Duplicate Name Function - 2",
			path:    "DuplicateNameFunction-2",
			summary: "func DuplicateNameFunction() string {",
			source: `func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 70,
			endLine:   72,
		},
		{
			name:    "Duplicate Name Function - 3",
			path:    "DuplicateNameFunction-3",
			summary: "func DuplicateNameFunction() string {",
			source: `func DuplicateNameFunction() string {
	return "duplicate name"
}`,
			startLine: 75,
			endLine:   77,
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
			s.Equal(test.startLine, chunk.StartLine)
			s.Equal(test.endLine, chunk.EndLine)
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
		startLine uint
		endLine   uint
	}{
		{
			name:    "Value Receiver Method",
			path:    "User::GetName",
			summary: "func (u User) GetName() string {",
			source: `func (u User) GetName() string {
	return u.Name
}`,
			startLine: 10,
			endLine:   12,
		},
		{
			name:    "Pointer Receiver Method",
			path:    "User::SetName",
			summary: "func (u *User) SetName(name string) {",
			source: `func (u *User) SetName(name string) {
	u.Name = name
}`,
			startLine: 15,
			endLine:   17,
		},
		{
			name:    "Service Add User Method",
			path:    "Service::AddUser",
			summary: "func (s *Service) AddUser(user User) {",
			source: `func (s *Service) AddUser(user User) {
	s.users = append(s.users, user)
}`,
			startLine: 25,
			endLine:   27,
		},
		{
			name:    "Service Find User Method",
			path:    "Service::FindUser",
			summary: "func (s *Service) FindUser(id int) *User {",
			source: `func (s *Service) FindUser(id int) *User {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i]
		}
	}

	return nil
}`,
			startLine: 30,
			endLine:   38,
		},
		{
			name:    "Service Count Method",
			path:    "Service::Count",
			summary: "func (s Service) Count() int {",
			source: `func (s Service) Count() int {
	return len(s.users)
}`,
			startLine: 41,
			endLine:   43,
		},
		{
			name:    "Generic Pointer Receiver Method",
			path:    "Repository::Add",
			summary: "func (r *Repository[T]) Add(item T) {",
			source: `func (r *Repository[T]) Add(item T) {
	r.items = append(r.items, item)
}`,
			startLine: 51,
			endLine:   53,
		},
		{
			name:    "Generic Value Receiver Method",
			path:    "Repository::Get",
			summary: "func (r Repository[T]) Get(index int) T {",
			source: `func (r Repository[T]) Get(index int) T {
	return r.items[index]
}`,
			startLine: 56,
			endLine:   58,
		},
		{
			name:    "Generic Size Method",
			path:    "Repository::Size",
			summary: "func (r Repository[T]) Size() int {",
			source: `func (r Repository[T]) Size() int {
	return len(r.items)
}`,
			startLine: 61,
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
			s.Equal(test.startLine, chunk.StartLine)
			s.Equal(test.endLine, chunk.EndLine)
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
		startLine uint
		endLine   uint
	}{
		{
			name:    "Basic Struct",
			path:    "BasicStruct",
			summary: "type BasicStruct struct {",
			source: `type BasicStruct struct {
	Field1 string
	Field2 int
}`,
			startLine: 4,
			endLine:   7,
		},
		{
			name:      "Empty Struct",
			path:      "EmptyStruct",
			summary:   "type EmptyStruct struct{}",
			source:    `type EmptyStruct struct{}`,
			startLine: 10,
			endLine:   10,
		},
		{
			name:    "Embedded Struct",
			path:    "EmbeddedStruct",
			summary: "type EmbeddedStruct struct {",
			source: `type EmbeddedStruct struct {
	BasicStruct
	ExtraField bool
}`,
			startLine: 13,
			endLine:   16,
		},
		{
			name:    "Simple Interface",
			path:    "SimpleInterface",
			summary: "type SimpleInterface interface {",
			source: `type SimpleInterface interface {
	Method1() string
	Method2(int) error
}`,
			startLine: 19,
			endLine:   22,
		},
		{
			name:      "Empty Interface",
			path:      "EmptyInterface",
			summary:   "type EmptyInterface interface{}",
			source:    `type EmptyInterface interface{}`,
			startLine: 25,
			endLine:   25,
		},
		{
			name:    "Embedded Interface",
			path:    "EmbeddedInterface",
			summary: "type EmbeddedInterface interface {",
			source: `type EmbeddedInterface interface {
	SimpleInterface
	Method3() bool
}`,
			startLine: 28,
			endLine:   31,
		},
		{
			name:    "Generic Type",
			path:    "GenericType",
			summary: "type GenericType[T any] struct {",
			source: `type GenericType[T any] struct {
	Value T
}`,
			startLine: 34,
			endLine:   36,
		},
		{
			name:    "Constrained Generic",
			path:    "ConstrainedGeneric",
			summary: "type ConstrainedGeneric[T comparable] struct {",
			source: `type ConstrainedGeneric[T comparable] struct {
	Key   T
	Value string
}`,
			startLine: 39,
			endLine:   42,
		},
		{
			name:      "Multiple Generics",
			path:      "MultipleGenerics",
			summary:   "type MultipleGenerics[K comparable, V any] map[K]V",
			source:    `type MultipleGenerics[K comparable, V any] map[K]V`,
			startLine: 45,
			endLine:   45,
		},
		{
			name:      "Type Alias",
			path:      "TypeAlias",
			summary:   "type TypeAlias = string",
			source:    `type TypeAlias = string`,
			startLine: 48,
			endLine:   48,
		},
		{
			name:      "Custom Type",
			path:      "CustomType",
			summary:   "type CustomType string",
			source:    `type CustomType string`,
			startLine: 51,
			endLine:   51,
		},
		{
			name:    "Consts Block Hashing",
			path:    "796012d7ee1311f0",
			summary: "const (",
			source: `const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	MaxRetries     = 5
)`,
			startLine: 54,
			endLine:   58,
		}, {
			name:      "Single Constant",
			path:      "DefaultTimeout",
			summary:   "const DefaultTimeout = 30",
			source:    `const DefaultTimeout = 30`,
			startLine: 61,
			endLine:   61,
		},
		{
			name:    "Vars Block Hashing",
			path:    "3a12064c73be460c",
			summary: "var (",
			source: `var (
	GlobalCounter int
	SystemReady   bool = true
	ConfigPath    string
)`,
			startLine: 64,
			endLine:   68,
		},
		{
			name:    "Single Variable",
			path:    "DefaultConfig",
			summary: "var DefaultConfig = BasicStruct{",
			source: `var DefaultConfig = BasicStruct{
	Field1: "default",
	Field2: 42,
}`,
			startLine: 71,
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
			s.Equal(test.startLine, chunk.StartLine)
			s.Equal(test.endLine, chunk.EndLine)
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
	s.Equal(`func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}`, chunk.Source)
	s.Equal(uint(6), chunk.StartLine)
	s.Equal(uint(10), chunk.EndLine)
	s.Equal("go/tests_test.go::TestSimple", chunk.ID())
}

func TestGoParserTestSuite(t *testing.T) {
	suite.Run(t, new(GoParserTestSuite))
}
