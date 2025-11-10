package parser_test

import (
	"testing"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type PythonParserTestSuite struct {
	ParserBaseTestSuite
}

func (s *PythonParserTestSuite) SetupSuite() {
	s.ParserBaseTestSuite.SetupSuite()

	var err error
	s.parser, err = parser.NewPythonParser(s.workspaceRoot)
	s.Require().NoError(err)
}

func (s *PythonParserTestSuite) TestFunctionParsing() {
	chunks := s.getChunks("python/functions.py")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:      "Docstring Hashing",
			path:      "b4c6199b06bf4bfb",
			summary:   `"""Test file for Python function definitions."""`,
			source:    `"""Test file for Python function definitions."""`,
			startLine: 1,
			endLine:   1,
		},
		{
			name:    "Simple Function",
			path:    "simple_function",
			summary: "def simple_function():",
			source: `# Simple function with no parameters
def simple_function():
    pass`,
			startLine: 3,
			endLine:   5,
		},
		{
			name:    "Function With Params",
			path:    "function_with_params",
			summary: "def function_with_params(a, b):",
			source: `# Function with parameters and return value
def function_with_params(a, b):
    return a + b`,
			startLine: 7,
			endLine:   9,
		},
		{
			name:    "Decorated Function",
			path:    "decorated_function",
			summary: "def decorated_function():",
			source: `# Property decorator example
@property
def decorated_function():
    return "decorated"`,
			startLine: 11,
			endLine:   14,
		},
		{
			name:    "Static Method",
			path:    "static_method",
			summary: "def static_method():",
			source: `# Static method decorator
@staticmethod
def static_method():
    return "static"`,
			startLine: 16,
			endLine:   19,
		},
		{
			name:    "Class Method",
			path:    "class_method",
			summary: "def class_method(cls):",
			source: `# Class method decorator
@classmethod
def class_method(cls):
    return "class method"`,
			startLine: 21,
			endLine:   24,
		},
		{
			name:    "Async Function",
			path:    "async_function",
			summary: "async def async_function():",
			source: `# Async function example
async def async_function():
    return "async"`,
			startLine: 26,
			endLine:   28,
		},
		{
			name:    "Generator Function",
			path:    "generator_function",
			summary: "def generator_function():",
			source: `# Generator function example
#  foo bar baz

def generator_function():
    yield 1
    yield 2`,
			startLine: 30,
			endLine:   35,
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
			s.Equal("python/functions.py::"+test.path, chunk.ID())
		})
	}
}

func (s *PythonParserTestSuite) TestClassParsing() {
	chunks := s.getChunks("python/classes.py")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Class",
			path:    "SimpleClass",
			summary: "class SimpleClass:",
			source: `# Simple class with no methods
class SimpleClass:
    pass`,
			startLine: 3,
			endLine:   5,
		},
		{
			name:    "Class With Methods",
			path:    "ClassWithMethods",
			summary: "class ClassWithMethods:",
			source: `class ClassWithMethods:
    value = -1

    # Constructor method
    def __init__(self):
        self.value = 0

    # Instance method
    def method(self):
        return self.value

    # Property method with decorator
    @property
    def property_method(self):
        return self.value * 2`,
			startLine: 7,
			endLine:   21,
		},
		{
			name:    "Inherited Class",
			path:    "InheritedClass",
			summary: "class InheritedClass(ClassWithMethods):",
			source: `# Class with inheritance
class InheritedClass(ClassWithMethods):
    # Override parent method
    def method(self):
        return super().method() + 1`,
			startLine: 23,
			endLine:   27,
		},
		{
			name:    "Decorated Class",
			path:    "DecoratedClass",
			summary: "class DecoratedClass:",
			source: `# Decorated class using dataclass
@dataclass
class DecoratedClass:
    name: str
    age: int`,
			startLine: 29,
			endLine:   33,
		},
		// ClassWithMethods members
		{
			name:      "ClassWithMethods::value",
			path:      "ClassWithMethods::value",
			summary:   "value = -1",
			source:    `value = -1`,
			startLine: 8,
			endLine:   8,
		},
		{
			name:    "ClassWithMethods::__init__",
			path:    "ClassWithMethods::__init__",
			summary: "def __init__(self):",
			source: `# Constructor method
    def __init__(self):
        self.value = 0`,
			startLine: 10,
			endLine:   12,
		},
		{
			name:    "ClassWithMethods::method",
			path:    "ClassWithMethods::method",
			summary: "def method(self):",
			source: `# Instance method
    def method(self):
        return self.value`,
			startLine: 14,
			endLine:   16,
		},
		{
			name:    "ClassWithMethods::property_method",
			path:    "ClassWithMethods::property_method",
			summary: "def property_method(self):",
			source: `# Property method with decorator
    @property
    def property_method(self):
        return self.value * 2`,
			startLine: 18,
			endLine:   21,
		},
		// InheritedClass members
		{
			name:    "InheritedClass::method",
			path:    "InheritedClass::method",
			summary: "def method(self):",
			source: `# Override parent method
    def method(self):
        return super().method() + 1`,
			startLine: 25,
			endLine:   27,
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
			s.Equal("python/classes.py::"+test.path, chunk.ID())
		})
	}
}

func (s *PythonParserTestSuite) TestTestFileParsing() {
	chunks := s.getChunks("python/tests/test_module.py")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
		fileType  string
	}{
		{
			name:      "Test Simple Function",
			path:      "e36d1e5889771889",
			summary:   `"""Test file for Python test patterns."""`,
			source:    `"""Test file for Python test patterns."""`,
			startLine: 1,
			endLine:   1,
			fileType:  "tests",
		},
		{
			name:    "Test Simple Function",
			path:    "test_simple_function",
			summary: "def test_simple_function():",
			source: `def test_simple_function():
    assert True`,
			startLine: 5,
			endLine:   6,
			fileType:  "tests",
		},
		{
			name:    "Test Sample Class",
			path:    "TestSample",
			summary: "class TestSample(unittest.TestCase):",
			source: `class TestSample(unittest.TestCase):
    def test_method(self):
        self.assertEqual(1 + 1, 2)

    def test_another(self):
        self.assertTrue(True)`,
			startLine: 8,
			endLine:   13,
			fileType:  "tests",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			chunk, exists := chunks[test.path]
			s.Require().True(exists, "chunk %s not found", test.path)
			s.Require().NotNil(chunk)

			s.Equal(test.fileType, chunk.Type)
			s.Equal(test.path, chunk.Path)
			s.Equal(test.summary, chunk.Summary)
			s.Equal(test.source, chunk.Source)
			s.Equal(test.startLine, int(chunk.StartLine))
			s.Equal(test.endLine, int(chunk.EndLine))
			s.Equal("python/tests/test_module.py::"+test.path, chunk.ID())
		})
	}
}

func (s *PythonParserTestSuite) TearDownSuite() {
	if s.parser != nil {
		s.parser.Close()
	}
}

func TestPythonParserTestSuite(t *testing.T) {
	suite.Run(t, new(PythonParserTestSuite))
}
