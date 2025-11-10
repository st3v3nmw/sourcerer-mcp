package parser_test

import (
	"testing"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type JavaScriptParserTestSuite struct {
	ParserBaseTestSuite
}

func (s *JavaScriptParserTestSuite) SetupSuite() {
	s.ParserBaseTestSuite.SetupSuite()

	var err error
	s.parser, err = parser.NewJavaScriptParser(s.workspaceRoot)
	s.Require().NoError(err)
}

func (s *JavaScriptParserTestSuite) TestFunctionParsing() {
	chunks := s.getChunks("javascript/functions.js")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Function",
			path:    "simple_function",
			summary: "function simple_function() {",
			source: `// Test file for JavaScript function definitions

// Simple function declaration
function simple_function() {
    return "hello";
}`,
			startLine: 1,
			endLine:   6,
		},
		{
			name:    "Function With Params",
			path:    "function_with_params",
			summary: "function function_with_params(a, b) {",
			source: `// Function with parameters and return value
function function_with_params(a, b) {
    return a + b;
}`,
			startLine: 8,
			endLine:   11,
		},
		{
			name:    "Arrow Function",
			path:    "arrow_function",
			summary: "const arrow_function = () => {",
			source: `// Arrow function assigned to variable
const arrow_function = () => {
    return "arrow";
};`,
			startLine: 13,
			endLine:   16,
		},
		{
			name:    "Arrow With Params",
			path:    "arrow_with_params",
			summary: "const arrow_with_params = (x, y) => {",
			source: `// Arrow function with parameters
const arrow_with_params = (x, y) => {
    return x * y;
};`,
			startLine: 18,
			endLine:   21,
		},
		{
			name:    "Function Expression",
			path:    "function_expression",
			summary: "const function_expression = function() {",
			source: `// Function expression assigned to variable
const function_expression = function() {
    return "expression";
};`,
			startLine: 23,
			endLine:   26,
		},
		{
			name:    "Named Function Expression",
			path:    "named_function_expression",
			summary: "const named_function_expression = function namedFn() {",
			source: `// Named function expression
const named_function_expression = function namedFn() {
    return "named expression";
};`,
			startLine: 28,
			endLine:   31,
		},
		{
			name:    "Async Function",
			path:    "async_function",
			summary: "async function async_function() {",
			source: `// Async function
async function async_function() {
    return "async result";
}`,
			startLine: 33,
			endLine:   36,
		},
		{
			name:    "Async Arrow",
			path:    "async_arrow",
			summary: "const async_arrow = async () => {",
			source: `// Async arrow function
const async_arrow = async () => {
    return "async arrow";
};`,
			startLine: 38,
			endLine:   41,
		},
		{
			name:    "Generator Function",
			path:    "generator_function",
			summary: "function* generator_function() {",
			source: `// Generator function
function* generator_function() {
    yield 1;
    yield 2;
}`,
			startLine: 43,
			endLine:   47,
		},
		{
			name:    "Exported Function",
			path:    "exported_function",
			summary: "function exported_function(input) {",
			source: `// Exported functions
export function exported_function(input) {
    return input.length > 0;
}`,
			startLine: 49,
			endLine:   52,
		},
		{
			name:    "Async Exported Function",
			path:    "async_exported_function",
			summary: "async function async_exported_function(data) {",
			source: `export async function async_exported_function(data) {
    return await Promise.resolve(data);
}`,
			startLine: 54,
			endLine:   56,
		},
		{
			name:    "Default Exported Function",
			path:    "default_exported_function",
			summary: "function default_exported_function() {",
			source: `export default function default_exported_function() {
    return "default";
}`,
			startLine: 58,
			endLine:   60,
		},
		{
			name:    "Exported Generator",
			path:    "exported_generator",
			summary: "function* exported_generator() {",
			source: `export function* exported_generator() {
    yield "a";
    yield "b";
}`,
			startLine: 62,
			endLine:   65,
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
			s.Equal("javascript/functions.js::"+test.path, chunk.ID())
		})
	}
}

func (s *JavaScriptParserTestSuite) TestClassParsing() {
	chunks := s.getChunks("javascript/classes.js")

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
			summary: "class SimpleClass {",
			source: `// Test file for JavaScript class definitions

// Simple class with no methods
class SimpleClass {
}`,
			startLine: 1,
			endLine:   5,
		},
		{
			name:    "Class With Methods",
			path:    "ClassWithMethods",
			summary: "class ClassWithMethods {",
			source: `// Class with constructor and methods
class ClassWithMethods {
    constructor(value) {
        this.value = value;
    }

    // Instance method
    getValue() {
        return this.value;
    }

    // Setter method
    setValue(newValue) {
        this.value = newValue;
    }

    // Static method
    static createDefault() {
        return new ClassWithMethods(0);
    }

    // Getter method
    get displayValue() {
        return ` + "`Value: ${this.value}`" + `;
    }
}`,
			startLine: 7,
			endLine:   32,
		},
		{
			name:    "Extended Class",
			path:    "ExtendedClass",
			summary: "class ExtendedClass extends ClassWithMethods {",
			source: `// Class with inheritance
class ExtendedClass extends ClassWithMethods {
    constructor(value, name) {
        super(value);
        this.name = name;
    }

    // Override parent method
    getValue() {
        return ` + "`${this.name}: ${this.value}`" + `;
    }

    // New method
    getName() {
        return this.name;
    }
}`,
			startLine: 34,
			endLine:   50,
		},
		{
			name:    "Class With Privates",
			path:    "ClassWithPrivates",
			summary: "class ClassWithPrivates {",
			source: `// Class with private fields
class ClassWithPrivates {
    #privateField = 42;

    getPrivate() {
        return this.#privateField;
    }
}`,
			startLine: 52,
			endLine:   59,
		},
		{
			name:    "Exported Class",
			path:    "ExportedClass",
			summary: "class ExportedClass {",
			source: `// Exported classes
export class ExportedClass {
    constructor(value) {
        this.value = value;
    }

    getValue() {
        return this.value;
    }
}`,
			startLine: 61,
			endLine:   70,
		},
		{
			name:    "Exported APIProvider",
			path:    "APIProvider",
			summary: "class APIProvider {",
			source: `export class APIProvider {
    constructor(name) {
        this.name = name;
    }

    process() {
        return ` + "`Processing: ${this.name}`" + `;
    }
}`,
			startLine: 72,
			endLine:   80,
		},
		{
			name:    "Default Exported Plugin",
			path:    "DefaultPlugin",
			summary: "class DefaultPlugin {",
			source: `export default class DefaultPlugin {
    constructor(config) {
        this.config = config;
    }
}`,
			startLine: 82,
			endLine:   86,
		},
		// ClassWithMethods members
		{
			name:    "ClassWithMethods::constructor",
			path:    "ClassWithMethods::constructor",
			summary: "constructor(value) {",
			source: `constructor(value) {
        this.value = value;
    }`,
			startLine: 9,
			endLine:   11,
		},
		{
			name:    "ClassWithMethods::getValue",
			path:    "ClassWithMethods::getValue",
			summary: "getValue() {",
			source: `// Instance method
    getValue() {
        return this.value;
    }`,
			startLine: 13,
			endLine:   16,
		},
		{
			name:    "ClassWithMethods::setValue",
			path:    "ClassWithMethods::setValue",
			summary: "setValue(newValue) {",
			source: `// Setter method
    setValue(newValue) {
        this.value = newValue;
    }`,
			startLine: 18,
			endLine:   21,
		},
		{
			name:    "ClassWithMethods::createDefault",
			path:    "ClassWithMethods::createDefault",
			summary: "static createDefault() {",
			source: `// Static method
    static createDefault() {
        return new ClassWithMethods(0);
    }`,
			startLine: 23,
			endLine:   26,
		},
		{
			name:    "ClassWithMethods::displayValue",
			path:    "ClassWithMethods::displayValue",
			summary: "get displayValue() {",
			source: `// Getter method
    get displayValue() {
        return ` + "`Value: ${this.value}`" + `;
    }`,
			startLine: 28,
			endLine:   31,
		},
		// ExtendedClass members
		{
			name:    "ExtendedClass::constructor",
			path:    "ExtendedClass::constructor",
			summary: "constructor(value, name) {",
			source: `constructor(value, name) {
        super(value);
        this.name = name;
    }`,
			startLine: 36,
			endLine:   39,
		},
		{
			name:    "ExtendedClass::getValue",
			path:    "ExtendedClass::getValue",
			summary: "getValue() {",
			source: `// Override parent method
    getValue() {
        return ` + "`${this.name}: ${this.value}`" + `;
    }`,
			startLine: 41,
			endLine:   44,
		},
		{
			name:    "ExtendedClass::getName",
			path:    "ExtendedClass::getName",
			summary: "getName() {",
			source: `// New method
    getName() {
        return this.name;
    }`,
			startLine: 46,
			endLine:   49,
		},
		// ClassWithPrivates members
		{
			name:    "ClassWithPrivates::getPrivate",
			path:    "ClassWithPrivates::getPrivate",
			summary: "getPrivate() {",
			source: `getPrivate() {
        return this.#privateField;
    }`,
			startLine: 56,
			endLine:   58,
		},
		// ExportedClass members
		{
			name:    "ExportedClass::constructor",
			path:    "ExportedClass::constructor",
			summary: "constructor(value) {",
			source: `constructor(value) {
        this.value = value;
    }`,
			startLine: 63,
			endLine:   65,
		},
		{
			name:    "ExportedClass::getValue",
			path:    "ExportedClass::getValue",
			summary: "getValue() {",
			source: `getValue() {
        return this.value;
    }`,
			startLine: 67,
			endLine:   69,
		},
		// APIProvider members
		{
			name:    "APIProvider::constructor",
			path:    "APIProvider::constructor",
			summary: "constructor(name) {",
			source: `constructor(name) {
        this.name = name;
    }`,
			startLine: 73,
			endLine:   75,
		},
		{
			name:    "APIProvider::process",
			path:    "APIProvider::process",
			summary: "process() {",
			source: `process() {
        return ` + "`Processing: ${this.name}`" + `;
    }`,
			startLine: 77,
			endLine:   79,
		},
		// DefaultPlugin members
		{
			name:    "DefaultPlugin::constructor",
			path:    "DefaultPlugin::constructor",
			summary: "constructor(config) {",
			source: `constructor(config) {
        this.config = config;
    }`,
			startLine: 83,
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
			s.Equal("javascript/classes.js::"+test.path, chunk.ID())
		})
	}
}

func (s *JavaScriptParserTestSuite) TestTestFileParsing() {
	chunks := s.getChunks("javascript/tests/module.test.js")

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
			name:      "File Comment",
			path:      "85bcbce29203af48",
			summary:   "// Test file for JavaScript test patterns",
			source:    `// Test file for JavaScript test patterns`,
			startLine: 1,
			endLine:   1,
			fileType:  "tests",
		},
		{
			name:    "Test Simple Function",
			path:    "test_simple_function",
			summary: "function test_simple_function() {",
			source: `function test_simple_function() {
    expect(true).toBe(true);
}`,
			startLine: 5,
			endLine:   7,
			fileType:  "tests",
		},
		{
			name:    "Test Sample Class",
			path:    "TestSample",
			summary: "class TestSample {",
			source: `class TestSample {
    test_method() {
        expect(1 + 1).toBe(2);
    }

    test_another() {
        expect(true).toBeTruthy();
    }
}`,
			startLine: 9,
			endLine:   17,
			fileType:  "tests",
		},
		{
			name:    "Test Sample Test",
			path:    "25c23a1c567a7424",
			summary: "describe('Sample test suite', () => {",
			source: `describe('Sample test suite', () => {
    it('should pass basic test', () => {
        expect(2 + 2).toBe(4);
    });

    it('should handle async operations', async () => {
        const result = await Promise.resolve('test');
        expect(result).toBe('test');
    });
});`,
			startLine: 19,
			endLine:   28,
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
			s.Equal("javascript/tests/module.test.js::"+test.path, chunk.ID())
		})
	}
}

func (s *JavaScriptParserTestSuite) TestVariableParsing() {
	chunks := s.getChunks("javascript/variables.js")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Var",
			path:    "simple_var",
			summary: "var simple_var = \"hello\";",
			source: `// Test file for JavaScript variable declarations

// var declarations
var simple_var = "hello";`,
			startLine: 1,
			endLine:   4,
		},
		{
			name:      "Multiple Var",
			path:      "6e4ceec03e9d5435",
			summary:   "var multiple_var_1, multiple_var_2;",
			source:    `var multiple_var_1, multiple_var_2;`,
			startLine: 5,
			endLine:   5,
		},
		{
			name:      "Initialized Var",
			path:      "initialized_var",
			summary:   "var initialized_var = 42;",
			source:    `var initialized_var = 42;`,
			startLine: 6,
			endLine:   6,
		},
		{
			name:    "Simple Let",
			path:    "simple_let",
			summary: "let simple_let = \"world\";",
			source: `// let declarations
let simple_let = "world";`,
			startLine: 8,
			endLine:   9,
		},
		{
			name:      "Destructured Let",
			path:      "destructured_let",
			summary:   "let destructured_let = {a: 1, b: 2};",
			source:    `let destructured_let = {a: 1, b: 2};`,
			startLine: 10,
			endLine:   10,
		},
		{
			name:      "Array Destructured",
			path:      "array_destructured",
			summary:   "let array_destructured = [1, 2, 3];",
			source:    `let array_destructured = [1, 2, 3];`,
			startLine: 11,
			endLine:   11,
		},
		{
			name:    "Simple Const",
			path:    "simple_const",
			summary: "const simple_const = \"constant\";",
			source: `// const declarations
const simple_const = "constant";`,
			startLine: 13,
			endLine:   14,
		},
		{
			name:    "Object Const",
			path:    "object_const",
			summary: "const object_const = {",
			source: `const object_const = {
    name: "test",
    value: 123
};`,
			startLine: 15,
			endLine:   18,
		},
		{
			name:      "Arrow Function Const",
			path:      "arrow_function_const",
			summary:   "const arrow_function_const = (x) => x * 2;",
			source:    "const arrow_function_const = (x) => x * 2;",
			startLine: 19,
			endLine:   19,
		},
		{
			name:    "Function Expression",
			path:    "function_expression",
			summary: "const function_expression = function(a, b) {",
			source: `// Function expressions assigned to variables
const function_expression = function(a, b) {
    return a + b;
};`,
			startLine: 21,
			endLine:   24,
		},
		{
			name:    "Named Function Expression",
			path:    "named_function_expression",
			summary: "const named_function_expression = function calculator(x, y) {",
			source: `const named_function_expression = function calculator(x, y) {
    return x - y;
};`,
			startLine: 26,
			endLine:   28,
		},
		{
			name:    "Simple Arrow",
			path:    "simple_arrow",
			summary: "const simple_arrow = () => \"simple\";",
			source: `// Arrow functions with different syntaxes
const simple_arrow = () => "simple";`,
			startLine: 30,
			endLine:   31,
		},
		{
			name:      "Arrow With Params",
			path:      "arrow_with_params",
			summary:   "const arrow_with_params = (a, b) => a + b;",
			source:    "const arrow_with_params = (a, b) => a + b;",
			startLine: 32,
			endLine:   32,
		},
		{
			name:    "Arrow With Body",
			path:    "2370e5908d333e83",
			summary: "const arrow_with_body = (x) => {",
			source: `const arrow_with_body = (x) => {
    const result = x * 2;
    return result;
};`,
			startLine: 33,
			endLine:   36,
		},
		{
			name:    "Async Arrow Func",
			path:    "7d9a98184af91c01",
			summary: "const async_arrow_func = async (data) => {",
			source: `// Async arrow function
const async_arrow_func = async (data) => {
    const response = await fetch(data);
    return response;
};`,
			startLine: 38,
			endLine:   42,
		},
		{
			name:    "Destructuring Assignment",
			path:    "69d42be7a90295c0",
			summary: "let {name, age} = person;",
			source: `// Complex variable declarations
let {name, age} = person;`,
			startLine: 44,
			endLine:   45,
		},
		{
			name:      "Array Destructuring Assignment",
			path:      "d7688b14f2c5cde",
			summary:   "const [first, ...rest] = numbers;",
			source:    `const [first, ...rest] = numbers;`,
			startLine: 46,
			endLine:   46,
		},
		{
			name:      "Global Var",
			path:      "globalVar",
			summary:   "var globalVar = window.something || \"default\";",
			source:    `var globalVar = window.something || "default";`,
			startLine: 47,
			endLine:   47,
		},
		{
			name:    "Exported Const",
			path:    "EXPORTED_CONST",
			summary: "const EXPORTED_CONST = \"constant value\";",
			source: `// Exported variables
export const EXPORTED_CONST = "constant value";`,
			startLine: 49,
			endLine:   50,
		},
		{
			name:      "Exported API_KEY",
			path:      "API_KEY",
			summary:   "const API_KEY = \"api-key-12345\";",
			source:    `export const API_KEY = "api-key-12345";`,
			startLine: 51,
			endLine:   51,
		},
		{
			name:      "Exported Let",
			path:      "exported_let",
			summary:   "let exported_let = 42;",
			source:    `export let exported_let = 42;`,
			startLine: 52,
			endLine:   52,
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
			s.Equal("javascript/variables.js::"+test.path, chunk.ID())
		})
	}
}

func (s *JavaScriptParserTestSuite) TearDownSuite() {
	if s.parser != nil {
		s.parser.Close()
	}
}

func TestJavaScriptParserTestSuite(t *testing.T) {
	suite.Run(t, new(JavaScriptParserTestSuite))
}
