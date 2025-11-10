package parser_test

import (
	"testing"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type TypeScriptParserTestSuite struct {
	ParserBaseTestSuite
}

func (s *TypeScriptParserTestSuite) SetupSuite() {
	s.ParserBaseTestSuite.SetupSuite()

	var err error
	s.parser, err = parser.NewTypeScriptParser(s.workspaceRoot)
	s.Require().NoError(err)
}

func (s *TypeScriptParserTestSuite) TestFunctionParsing() {
	chunks := s.getChunks("typescript/functions.ts")

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
			summary: "function simple_function(): string {",
			source: `// Test file for TypeScript function definitions

// Simple function declaration with type annotations
function simple_function(): string {
    return "hello";
}`,
			startLine: 1,
			endLine:   6,
		},
		{
			name:    "Function With Params",
			path:    "function_with_params",
			summary: "function function_with_params(a: number, b: number): number {",
			source: `// Function with typed parameters and return value
function function_with_params(a: number, b: number): number {
    return a + b;
}`,
			startLine: 8,
			endLine:   11,
		},
		{
			name:    "Arrow Function",
			path:    "arrow_function",
			summary: "const arrow_function = (): string => {",
			source: `// Arrow function assigned to variable with types
const arrow_function = (): string => {
    return "arrow";
};`,
			startLine: 13,
			endLine:   16,
		},
		{
			name:    "Arrow With Params",
			path:    "arrow_with_params",
			summary: "const arrow_with_params = (x: number, y: number): number => {",
			source: `// Arrow function with typed parameters
const arrow_with_params = (x: number, y: number): number => {
    return x * y;
};`,
			startLine: 18,
			endLine:   21,
		},
		{
			name:    "Function Expression",
			path:    "function_expression",
			summary: "const function_expression = function(a: string, b: string): string {",
			source: `// Function expression assigned to variable with types
const function_expression = function(a: string, b: string): string {
    return a + b;
};`,
			startLine: 23,
			endLine:   26,
		},
		{
			name:    "Named Function Expression",
			path:    "named_function_expression",
			summary: "const named_function_expression = function namedFn(value: any): string {",
			source: `// Named function expression with types
const named_function_expression = function namedFn(value: any): string {
    return "named expression";
};`,
			startLine: 28,
			endLine:   31,
		},
		{
			name:    "Async Function",
			path:    "async_function",
			summary: "async function async_function(): Promise<string> {",
			source: `// Async function with Promise return type
async function async_function(): Promise<string> {
    return "async result";
}`,
			startLine: 33,
			endLine:   36,
		},
		{
			name:    "Async Arrow",
			path:    "async_arrow",
			summary: "const async_arrow = async (data: string): Promise<string> => {",
			source: `// Async arrow function with typed parameter
const async_arrow = async (data: string): Promise<string> => {
    return "async arrow";
};`,
			startLine: 38,
			endLine:   41,
		},
		{
			name:    "Generator Function",
			path:    "generator_function",
			summary: "function* generator_function(): Generator<number, void, unknown> {",
			source: `// Generator function with typed yield
function* generator_function(): Generator<number, void, unknown> {
    yield 1;
    yield 2;
}`,
			startLine: 43,
			endLine:   47,
		},
		{
			name:    "Generic Function",
			path:    "generic_function",
			summary: "function generic_function<T>(value: T): T {",
			source: `// Generic function
function generic_function<T>(value: T): T {
    return value;
}`,
			startLine: 49,
			endLine:   52,
		},
		{
			name:    "Optional Params",
			path:    "optional_params",
			summary: "function optional_params(required: string, optional?: number): string {",
			source: `// Function with optional parameters
function optional_params(required: string, optional?: number): string {
    return required + (optional || 0);
}`,
			startLine: 54,
			endLine:   57,
		},
		{
			name:    "Rest Params",
			path:    "rest_params",
			summary: "function rest_params(first: string, ...rest: number[]): string {",
			source: `// Function with rest parameters
function rest_params(first: string, ...rest: number[]): string {
    return first + rest.join(',');
}`,
			startLine: 59,
			endLine:   62,
		},
		{
			name:    "Function Signature",
			path:    "myFunc",
			summary: "function myFunc(x: number): string;",
			source: `// Function signature
function myFunc(x: number): string;`,
			startLine: 64,
			endLine:   65,
		},
		{
			name:    "Exported Function",
			path:    "exportedFunction",
			summary: "function exportedFunction(input: string): boolean {",
			source: `// Exported functions
export function exportedFunction(input: string): boolean {
    return input.length > 0;
}`,
			startLine: 67,
			endLine:   70,
		},
		{
			name:    "Async Exported Function",
			path:    "asyncExportedFunction",
			summary: "async function asyncExportedFunction(data: any): Promise<void> {",
			source: `export async function asyncExportedFunction(data: any): Promise<void> {
    await Promise.resolve(data);
}`,
			startLine: 72,
			endLine:   74,
		},
		{
			name:    "Exported Default Function",
			path:    "defaultExportedFunction",
			summary: "function defaultExportedFunction(): string {",
			source: `export default function defaultExportedFunction(): string {
    return "default";
}`,
			startLine: 76,
			endLine:   78,
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
			s.Equal("typescript/functions.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestClassParsing() {
	chunks := s.getChunks("typescript/classes.ts")

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
			source: `// Test file for TypeScript class definitions

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
			source: `// Class with typed constructor and methods
class ClassWithMethods {
    private value: number;

    constructor(value: number) {
        this.value = value;
    }

    // Instance method with return type
    getValue(): number {
        return this.value;
    }

    // Setter method with typed parameter
    setValue(newValue: number): void {
        this.value = newValue;
    }

    // Static method with return type
    static createDefault(): ClassWithMethods {
        return new ClassWithMethods(0);
    }

    // Getter method with return type
    get displayValue(): string {
        return ` + "`Value: ${this.value}`" + `;
    }

    // Setter property
    set displayValue(val: string) {
        const match = val.match(/Value: (\d+)/);
        if (match) {
            this.value = parseInt(match[1]);
        }
    }
}`,
			startLine: 7,
			endLine:   42,
		},
		{
			name:    "Extended Class",
			path:    "ExtendedClass",
			summary: "class ExtendedClass<T> extends ClassWithMethods {",
			source: `// Class with inheritance and generics
class ExtendedClass<T> extends ClassWithMethods {
    private name: T;

    constructor(value: number, name: T) {
        super(value);
        this.name = name;
    }

    // Override parent method with different return type
    getValue(): string {
        return ` + "`${this.name}: ${super.getValue()}`" + `;
    }

    // New method with generic return
    getName(): T {
        return this.name;
    }
}`,
			startLine: 44,
			endLine:   62,
		},
		{
			name:    "Abstract Class",
			path:    "AbstractClass",
			summary: "abstract class AbstractClass {",
			source: `// Abstract class
abstract class AbstractClass {
    protected abstract process(): void;

    public execute(): void {
        this.process();
    }
}`,
			startLine: 64,
			endLine:   71,
		},
		{
			name:    "Processable Interface",
			path:    "Processable",
			summary: "interface Processable {",
			source: `// Interface implementation
interface Processable {
    process(): void;
}`,
			startLine: 73,
			endLine:   76,
		},
		{
			name:    "Implemented Class",
			path:    "ImplementedClass",
			summary: "class ImplementedClass extends AbstractClass implements Processable {",
			source: `class ImplementedClass extends AbstractClass implements Processable {
    process(): void {
        console.log("processing");
    }
}`,
			startLine: 78,
			endLine:   82,
		},
		{
			name:    "Config Class",
			path:    "ConfigClass",
			summary: "class ConfigClass {",
			source: `// Class with readonly and optional properties
class ConfigClass {
    readonly id: string;
    name?: string;
    private _config: Record<string, any> = {};

    constructor(id: string, name?: string) {
        this.id = id;
        this.name = name;
    }
}`,
			startLine: 84,
			endLine:   94,
		},
		{
			name:    "Decorated Class",
			path:    "BugReport",
			summary: "@sealed",
			source: `// Decorated

@sealed
class BugReport {
  type = "report";
  title: string;

  constructor(t: string) {
    this.title = t;
  }
}`,
			startLine: 96,
			endLine:   106,
		},
		{
			name:    "Exported Class",
			path:    "ExportedClass",
			summary: "class ExportedClass {",
			source: `// Exported classes
export class ExportedClass {
    private value: number;

    constructor(value: number) {
        this.value = value;
    }

    getValue(): number {
        return this.value;
    }
}`,
			startLine: 108,
			endLine:   119,
		},
		{
			name:    "Exported OpenRouterProvider",
			path:    "OpenRouterProvider",
			summary: "class OpenRouterProvider implements Processable {",
			source: `export class OpenRouterProvider implements Processable {
    name = "OpenRouter";

    process(): void {
        console.log("processing");
    }
}`,
			startLine: 121,
			endLine:   127,
		},
		{
			name:    "Exported Default Class",
			path:    "TutorPlugin",
			summary: "class TutorPlugin {",
			source: `export default class TutorPlugin {
    private config: any;

    constructor(config: any) {
        this.config = config;
    }
}`,
			startLine: 129,
			endLine:   135,
		},
		// ClassWithMethods members
		{
			name:      "ClassWithMethods::value",
			path:      "ClassWithMethods::value",
			summary:   "private value: number",
			source:    "private value: number",
			startLine: 9,
			endLine:   9,
		},
		{
			name:    "ClassWithMethods::constructor",
			path:    "ClassWithMethods::constructor",
			summary: "constructor(value: number) {",
			source: `constructor(value: number) {
        this.value = value;
    }`,
			startLine: 11,
			endLine:   13,
		},
		{
			name:    "ClassWithMethods::getValue",
			path:    "ClassWithMethods::getValue",
			summary: "getValue(): number {",
			source: `// Instance method with return type
    getValue(): number {
        return this.value;
    }`,
			startLine: 15,
			endLine:   18,
		},
		{
			name:    "ClassWithMethods::setValue",
			path:    "ClassWithMethods::setValue",
			summary: "setValue(newValue: number): void {",
			source: `// Setter method with typed parameter
    setValue(newValue: number): void {
        this.value = newValue;
    }`,
			startLine: 20,
			endLine:   23,
		},
		{
			name:    "ClassWithMethods::createDefault",
			path:    "ClassWithMethods::createDefault",
			summary: "static createDefault(): ClassWithMethods {",
			source: `// Static method with return type
    static createDefault(): ClassWithMethods {
        return new ClassWithMethods(0);
    }`,
			startLine: 25,
			endLine:   28,
		},
		{
			name:    "ClassWithMethods::displayValue getter",
			path:    "ClassWithMethods::displayValue",
			summary: "get displayValue(): string {",
			source: `// Getter method with return type
    get displayValue(): string {
        return ` + "`Value: ${this.value}`" + `;
    }`,
			startLine: 30,
			endLine:   33,
		},
		{
			name:    "ClassWithMethods::displayValue setter",
			path:    "ClassWithMethods::displayValue-2",
			summary: "set displayValue(val: string) {",
			source: `// Setter property
    set displayValue(val: string) {
        const match = val.match(/Value: (\d+)/);
        if (match) {
            this.value = parseInt(match[1]);
        }
    }`,
			startLine: 35,
			endLine:   41,
		},
		// ExtendedClass members
		{
			name:      "ExtendedClass::name",
			path:      "ExtendedClass::name",
			summary:   "private name: T",
			source:    "private name: T",
			startLine: 46,
			endLine:   46,
		},
		{
			name:    "ExtendedClass::constructor",
			path:    "ExtendedClass::constructor",
			summary: "constructor(value: number, name: T) {",
			source: `constructor(value: number, name: T) {
        super(value);
        this.name = name;
    }`,
			startLine: 48,
			endLine:   51,
		},
		{
			name:    "ExtendedClass::getValue",
			path:    "ExtendedClass::getValue",
			summary: "getValue(): string {",
			source: `// Override parent method with different return type
    getValue(): string {
        return ` + "`${this.name}: ${super.getValue()}`" + `;
    }`,
			startLine: 53,
			endLine:   56,
		},
		{
			name:    "ExtendedClass::getName",
			path:    "ExtendedClass::getName",
			summary: "getName(): T {",
			source: `// New method with generic return
    getName(): T {
        return this.name;
    }`,
			startLine: 58,
			endLine:   61,
		},
		// AbstractClass members
		{
			name:      "AbstractClass::process",
			path:      "AbstractClass::process",
			summary:   "protected abstract process(): void",
			source:    "protected abstract process(): void",
			startLine: 66,
			endLine:   66,
		},
		{
			name:    "AbstractClass::execute",
			path:    "AbstractClass::execute",
			summary: "public execute(): void {",
			source: `public execute(): void {
        this.process();
    }`,
			startLine: 68,
			endLine:   70,
		},
		// ImplementedClass members
		{
			name:    "ImplementedClass::process",
			path:    "ImplementedClass::process",
			summary: "process(): void {",
			source: `process(): void {
        console.log("processing");
    }`,
			startLine: 79,
			endLine:   81,
		},
		// ConfigClass members
		{
			name:      "ConfigClass::id",
			path:      "ConfigClass::id",
			summary:   "readonly id: string",
			source:    "readonly id: string",
			startLine: 86,
			endLine:   86,
		},
		{
			name:      "ConfigClass::name",
			path:      "ConfigClass::name",
			summary:   "name?: string",
			source:    "name?: string",
			startLine: 87,
			endLine:   87,
		},
		{
			name:      "ConfigClass::_config",
			path:      "ConfigClass::_config",
			summary:   "private _config: Record<string, any> = {}",
			source:    "private _config: Record<string, any> = {}",
			startLine: 88,
			endLine:   88,
		},
		{
			name:    "ConfigClass::constructor",
			path:    "ConfigClass::constructor",
			summary: "constructor(id: string, name?: string) {",
			source: `constructor(id: string, name?: string) {
        this.id = id;
        this.name = name;
    }`,
			startLine: 90,
			endLine:   93,
		},
		// BugReport members
		{
			name:      "BugReport::type",
			path:      "BugReport::type",
			summary:   "type = \"report\"",
			source:    `type = "report"`,
			startLine: 100,
			endLine:   100,
		},
		{
			name:      "BugReport::title",
			path:      "BugReport::title",
			summary:   "title: string",
			source:    "title: string",
			startLine: 101,
			endLine:   101,
		},
		{
			name:    "BugReport::constructor",
			path:    "BugReport::constructor",
			summary: "constructor(t: string) {",
			source: `constructor(t: string) {
    this.title = t;
  }`,
			startLine: 103,
			endLine:   105,
		},
		// ExportedClass members
		{
			name:      "ExportedClass::value",
			path:      "ExportedClass::value",
			summary:   "private value: number",
			source:    "private value: number",
			startLine: 110,
			endLine:   110,
		},
		{
			name:    "ExportedClass::constructor",
			path:    "ExportedClass::constructor",
			summary: "constructor(value: number) {",
			source: `constructor(value: number) {
        this.value = value;
    }`,
			startLine: 112,
			endLine:   114,
		},
		{
			name:    "ExportedClass::getValue",
			path:    "ExportedClass::getValue",
			summary: "getValue(): number {",
			source: `getValue(): number {
        return this.value;
    }`,
			startLine: 116,
			endLine:   118,
		},
		// OpenRouterProvider members
		{
			name:      "OpenRouterProvider::name",
			path:      "OpenRouterProvider::name",
			summary:   "name = \"OpenRouter\"",
			source:    `name = "OpenRouter"`,
			startLine: 122,
			endLine:   122,
		},
		{
			name:    "OpenRouterProvider::process",
			path:    "OpenRouterProvider::process",
			summary: "process(): void {",
			source: `process(): void {
        console.log("processing");
    }`,
			startLine: 124,
			endLine:   126,
		},
		// TutorPlugin members
		{
			name:      "TutorPlugin::config",
			path:      "TutorPlugin::config",
			summary:   "private config: any",
			source:    "private config: any",
			startLine: 130,
			endLine:   130,
		},
		{
			name:    "TutorPlugin::constructor",
			path:    "TutorPlugin::constructor",
			summary: "constructor(config: any) {",
			source: `constructor(config: any) {
        this.config = config;
    }`,
			startLine: 132,
			endLine:   134,
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
			s.Equal("typescript/classes.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestInterfaceParsing() {
	chunks := s.getChunks("typescript/interfaces.ts")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Interface",
			path:    "SimpleInterface",
			summary: "interface SimpleInterface {",
			source: `// Test file for TypeScript interface definitions

// Simple interface
interface SimpleInterface {
    name: string;
    value: number;
}`,
			startLine: 1,
			endLine:   7,
		},
		{
			name:    "Optional Interface",
			path:    "OptionalInterface",
			summary: "interface OptionalInterface {",
			source: `// Interface with optional properties
interface OptionalInterface {
    required: string;
    optional?: number;
    readonly id: string;
}`,
			startLine: 9,
			endLine:   14,
		},
		{
			name:    "Method Interface",
			path:    "MethodInterface",
			summary: "interface MethodInterface {",
			source: `// Interface with method signatures
interface MethodInterface {
    getName(): string;
    setValue(value: number): void;
    process(data: string): Promise<boolean>;
}`,
			startLine: 16,
			endLine:   21,
		},
		{
			name:    "Generic Interface",
			path:    "GenericInterface",
			summary: "interface GenericInterface<T, U = string> {",
			source: `// Generic interface
interface GenericInterface<T, U = string> {
    data: T;
    metadata: U;
    transform<V>(input: T): V;
}`,
			startLine: 23,
			endLine:   28,
		},
		{
			name:    "Extended Interface",
			path:    "ExtendedInterface",
			summary: "interface ExtendedInterface extends SimpleInterface, MethodInterface {",
			source: `// Extending interfaces
interface ExtendedInterface extends SimpleInterface, MethodInterface {
    category: string;
    tags: string[];
}`,
			startLine: 30,
			endLine:   34,
		},
		{
			name:    "Indexed Interface",
			path:    "IndexedInterface",
			summary: "interface IndexedInterface {",
			source: `// Interface with index signature
interface IndexedInterface {
    [key: string]: any;
    [key: number]: string;
}`,
			startLine: 36,
			endLine:   40,
		},
		{
			name:    "Callable Interface",
			path:    "CallableInterface",
			summary: "interface CallableInterface {",
			source: `// Interface with call signature
interface CallableInterface {
    (input: string): boolean;
    property: number;
}`,
			startLine: 42,
			endLine:   46,
		},
		{
			name:    "Constructable Interface",
			path:    "ConstructableInterface",
			summary: "interface ConstructableInterface {",
			source: `// Interface with construct signature
interface ConstructableInterface {
    new (value: string): { result: string };
}`,
			startLine: 48,
			endLine:   51,
		},
		{
			name:    "Namespace Interface",
			path:    "5a8ecd1335d8d887",
			summary: "namespace Interfaces {",
			source: `// Namespace interface
namespace Interfaces {
    export interface NamespacedInterface {
        value: boolean;
    }
}`,
			startLine: 53,
			endLine:   58,
		},
		{
			name:    "Merged Interface I",
			path:    "MergedInterface",
			summary: "interface MergedInterface {",
			source: `// Merged interface declarations
interface MergedInterface {
    first: string;
}`,
			startLine: 60,
			endLine:   63,
		},
		{
			name:    "Merged Interface II",
			path:    "MergedInterface-2",
			summary: "interface MergedInterface {",
			source: `interface MergedInterface {
    second: number;
}`,
			startLine: 65,
			endLine:   67,
		},
		{
			name:    "Exported Interface",
			path:    "ExportedInterface",
			summary: "interface ExportedInterface {",
			source: `// Exported interfaces
export interface ExportedInterface {
    id: string;
    process(): void;
}`,
			startLine: 69,
			endLine:   73,
		},
		{
			name:    "Exported LLMProvider",
			path:    "LLMProvider",
			summary: "interface LLMProvider {",
			source: `export interface LLMProvider {
    name: string;
    generate(prompt: string): Promise<string>;
}`,
			startLine: 75,
			endLine:   78,
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
			s.Equal("typescript/interfaces.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestTypeParsing() {
	chunks := s.getChunks("typescript/types.ts")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Type",
			path:    "SimpleType",
			summary: "type SimpleType = string;",
			source: `// Test file for TypeScript type definitions

// Simple type alias
type SimpleType = string;`,
			startLine: 1,
			endLine:   4,
		},
		{
			name:    "Union Type",
			path:    "UnionType",
			summary: "type UnionType = string | number | boolean;",
			source: `// Union type
type UnionType = string | number | boolean;`,
			startLine: 6,
			endLine:   7,
		},
		{
			name:    "Intersection Type",
			path:    "IntersectionType",
			summary: "type IntersectionType = { name: string } & { age: number };",
			source: `// Intersection type
type IntersectionType = { name: string } & { age: number };`,
			startLine: 9,
			endLine:   10,
		},
		{
			name:    "Generic Type",
			path:    "GenericType",
			summary: "type GenericType<T> = {",
			source: `// Generic type alias
type GenericType<T> = {
    value: T;
    optional?: T;
};`,
			startLine: 12,
			endLine:   16,
		},
		{
			name:    "Function Type",
			path:    "FunctionType",
			summary: "type FunctionType = (input: string) => boolean;",
			source: `// Function type alias
type FunctionType = (input: string) => boolean;`,
			startLine: 18,
			endLine:   19,
		},
		{
			name:    "Overloaded Function",
			path:    "OverloadedFunction",
			summary: "type OverloadedFunction = {",
			source: `// Complex function type with overloads
type OverloadedFunction = {
    (input: string): string;
    (input: number): number;
    (input: boolean): boolean;
};`,
			startLine: 21,
			endLine:   26,
		},
		{
			name:    "Object Type",
			path:    "ObjectType",
			summary: "type ObjectType = {",
			source: `// Object type alias
type ObjectType = {
    readonly id: string;
    name: string;
    age?: number;
    [key: string]: any;
};`,
			startLine: 28,
			endLine:   34,
		},
		{
			name:    "String Array",
			path:    "StringArray",
			summary: "type StringArray = string[];",
			source: `// Array type aliases
type StringArray = string[];`,
			startLine: 36,
			endLine:   37,
		},
		{
			name:      "Number Tuple",
			path:      "NumberTuple",
			summary:   "type NumberTuple = [number, number, string?];",
			source:    `type NumberTuple = [number, number, string?];`,
			startLine: 38,
			endLine:   38,
		},
		{
			name:    "Conditional Type",
			path:    "ConditionalType",
			summary: "type ConditionalType<T> = T extends string ? string[] : T[];",
			source: `// Conditional type
type ConditionalType<T> = T extends string ? string[] : T[];`,
			startLine: 40,
			endLine:   41,
		},
		{
			name:    "Mapped Type",
			path:    "MappedType",
			summary: "type MappedType<T> = {",
			source: `// Mapped type
type MappedType<T> = {
    [K in keyof T]: T[K] | null;
};`,
			startLine: 43,
			endLine:   46,
		},
		{
			name:    "Template Type",
			path:    "TemplateType",
			summary: "type TemplateType = `prefix_${string}_suffix`;",
			source: `// Template literal type
type TemplateType = ` + "`prefix_${string}_suffix`;",
			startLine: 48,
			endLine:   49,
		},
		{
			name:    "Partial User",
			path:    "PartialUser",
			summary: "type PartialUser = Partial<{ name: string; age: number }>;",
			source: `// Utility type usage
type PartialUser = Partial<{ name: string; age: number }>;`,
			startLine: 51,
			endLine:   52,
		},
		{
			name:      "Required User",
			path:      "RequiredUser",
			summary:   "type RequiredUser = Required<{ name?: string; age?: number }>;",
			source:    `type RequiredUser = Required<{ name?: string; age?: number }>;`,
			startLine: 53,
			endLine:   53,
		},
		{
			name:    "Status Type",
			path:    "StatusType",
			summary: "type StatusType = 'pending' | 'completed' | 'failed';",
			source: `// Enum-like type
type StatusType = 'pending' | 'completed' | 'failed';`,
			startLine: 55,
			endLine:   56,
		},
		{
			name:    "Recursive Type",
			path:    "RecursiveType",
			summary: "type RecursiveType<T> = T | RecursiveType<T>[];",
			source: `// Recursive type
type RecursiveType<T> = T | RecursiveType<T>[];`,
			startLine: 58,
			endLine:   59,
		},
		{
			name:    "Constrained Type",
			path:    "ConstrainedType",
			summary: "type ConstrainedType<T extends Record<string, any>> = {",
			source: `// Type with generic constraints
type ConstrainedType<T extends Record<string, any>> = {
    data: T;
    keys: keyof T;
};`,
			startLine: 61,
			endLine:   65,
		},
		{
			name:    "Exported Type",
			path:    "ExportedType",
			summary: "type ExportedType = string | number;",
			source: `// Exported types
export type ExportedType = string | number;`,
			startLine: 67,
			endLine:   68,
		},
		{
			name:    "Exported ConfigOptions",
			path:    "ConfigOptions",
			summary: "type ConfigOptions = {",
			source: `export type ConfigOptions = {
    enabled: boolean;
    timeout: number;
};`,
			startLine: 70,
			endLine:   73,
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
			s.Equal("typescript/types.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestEnumParsing() {
	chunks := s.getChunks("typescript/enums.ts")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Enum",
			path:    "SimpleEnum",
			summary: "enum SimpleEnum {",
			source: `// Test file for TypeScript enum definitions

// Simple numeric enum
enum SimpleEnum {
    First,
    Second,
    Third
}`,
			startLine: 1,
			endLine:   8,
		},
		{
			name:    "String Enum",
			path:    "StringEnum",
			summary: "enum StringEnum {",
			source: `// String enum
enum StringEnum {
    Red = "red",
    Green = "green",
    Blue = "blue"
}`,
			startLine: 10,
			endLine:   15,
		},
		{
			name:    "Mixed Enum",
			path:    "MixedEnum",
			summary: "enum MixedEnum {",
			source: `// Mixed enum
enum MixedEnum {
    None = 0,
    Read = 1,
    Write = 2,
    Execute = 4,
    Description = "permissions"
}`,
			startLine: 17,
			endLine:   24,
		},
		{
			name:    "Computed Enum",
			path:    "ComputedEnum",
			summary: "enum ComputedEnum {",
			source: `// Computed enum
enum ComputedEnum {
    Base = 1,
    Double = Base * 2,
    Triple = Base * 3
}`,
			startLine: 26,
			endLine:   31,
		},
		{
			name:    "Const Enum",
			path:    "ConstEnum",
			summary: "const enum ConstEnum {",
			source: `// Const enum
const enum ConstEnum {
    Tiny = 1,
    Small = 2,
    Medium = 4,
    Large = 8
}`,
			startLine: 33,
			endLine:   39,
		},
		{
			name:    "Status Enum",
			path:    "StatusEnum",
			summary: "enum StatusEnum {",
			source: `// Enum with explicit values
enum StatusEnum {
    Pending = "PENDING",
    InProgress = "IN_PROGRESS",
    Completed = "COMPLETED",
    Failed = "FAILED"
}`,
			startLine: 41,
			endLine:   47,
		},
		{
			name:    "Bidirectional Enum",
			path:    "BiDirectionalEnum",
			summary: "enum BiDirectionalEnum {",
			source: `// Reverse mapped enum
enum BiDirectionalEnum {
    Up = "UP",
    Down = "DOWN",
    Left = "LEFT",
    Right = "RIGHT"
}`,
			startLine: 49,
			endLine:   55,
		},
		{
			name:    "Enum Namespace",
			path:    "643245017f52983b",
			summary: "namespace EnumNamespace {",
			source: `// Enum in namespace
namespace EnumNamespace {
    export enum NestedEnum {
        Option1 = 1,
        Option2 = 2
    }
}`,
			startLine: 57,
			endLine:   63,
		},
		{
			name:    "Exported Enum",
			path:    "ExportedEnum",
			summary: "enum ExportedEnum {",
			source: `// Exported enums
export enum ExportedEnum {
    Alpha = "alpha",
    Beta = "beta",
    Gamma = "gamma"
}`,
			startLine: 65,
			endLine:   70,
		},
		{
			name:    "Exported Status",
			path:    "Status",
			summary: "enum Status {",
			source: `export enum Status {
    Pending = "pending",
    Active = "active",
    Complete = "complete"
}`,
			startLine: 72,
			endLine:   76,
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
			s.Equal("typescript/enums.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestNamespaceParsing() {
	chunks := s.getChunks("typescript/namespaces.ts")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine int
		endLine   int
	}{
		{
			name:    "Simple Namespace",
			path:    "843d1d59a469ba2e",
			summary: "namespace SimpleNamespace {",
			source: `// Test file for TypeScript namespace definitions

// Simple namespace
namespace SimpleNamespace {
    export const value = "namespace value";
    export function helper(): string {
        return "namespace function";
    }
}`,
			startLine: 1,
			endLine:   9,
		},
		{
			name:    "Outer Namespace",
			path:    "2888478eb6290af",
			summary: "namespace OuterNamespace {",
			source: `// Nested namespaces
namespace OuterNamespace {
    export namespace InnerNamespace {
        export const data = 42;
        export interface Config {
            enabled: boolean;
        }
    }

    export const config: InnerNamespace.Config = {
        enabled: true
    };
}`,
			startLine: 11,
			endLine:   23,
		},
		{
			name:    "Utility Namespace",
			path:    "ab6ebc6c6bfe2530",
			summary: "namespace UtilityNamespace {",
			source: `// Namespace with classes and interfaces
namespace UtilityNamespace {
    export interface Logger {
        log(message: string): void;
    }

    export class ConsoleLogger implements Logger {
        log(message: string): void {
            console.log(` + "`[LOG]: ${message}`" + `);
        }
    }

    export function createLogger(): Logger {
        return new ConsoleLogger();
    }
}`,
			startLine: 25,
			endLine:   40,
		},
		{
			name:    "External Library",
			path:    "e55a503d7227a55b",
			summary: "declare namespace ExternalLibrary {",
			source: `// Module declaration (ambient namespace)
declare namespace ExternalLibrary {
    interface Options {
        timeout: number;
    }

    function initialize(options: Options): void;
}`,
			startLine: 42,
			endLine:   49,
		},
		{
			name:    "Global Namespace",
			path:    "25f16e2c0e5708b4",
			summary: "declare global {",
			source: `// Global augmentation
declare global {
    namespace NodeJS {
        interface ProcessEnv {
            CUSTOM_VAR: string;
        }
    }
}`,
			startLine: 51,
			endLine:   58,
		},
		{
			name:    "Merged Namespace 1",
			path:    "68f0b782ae122c7",
			summary: "namespace MergedNamespace {",
			source: `// Namespace merging
namespace MergedNamespace {
    export const first = "first";
}`,
			startLine: 60,
			endLine:   63,
		},
		{
			name:    "Merged Namespace 2",
			path:    "ba7f8fd919fbf0a",
			summary: "namespace MergedNamespace {",
			source: `namespace MergedNamespace {
    export const second = "second";
}`,
			startLine: 65,
			endLine:   67,
		},
		{
			name:    "Generic Namespace",
			path:    "ac16ab1273e542c9",
			summary: "namespace GenericNamespace {",
			source: `// Namespace with generics
namespace GenericNamespace {
    export interface Container<T> {
        value: T;
    }

    export function wrap<T>(value: T): Container<T> {
        return { value };
    }
}`,
			startLine: 69,
			endLine:   78,
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
			s.Equal("typescript/namespaces.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TestTestFileParsing() {
	chunks := s.getChunks("typescript/tests/module.test.ts")

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
			path:      "5b931df60944a870",
			summary:   "// Test file for TypeScript test patterns",
			source:    `// Test file for TypeScript test patterns`,
			startLine: 1,
			endLine:   1,
			fileType:  "tests",
		},
		{
			name:    "Test Simple Function",
			path:    "test_simple_function",
			summary: "function test_simple_function(): void {",
			source: `// Simple test function
function test_simple_function(): void {
    expect(true).toBe(true);
}`,
			startLine: 5,
			endLine:   8,
			fileType:  "tests",
		},
		{
			name:    "Test Sample Class",
			path:    "TestSample",
			summary: "class TestSample {",
			source: `// Test class with typed methods
class TestSample {
    test_method(): void {
        expect(1 + 1).toBe(2);
    }

    async test_async_method(): Promise<void> {
        const result: number = await Promise.resolve(42);
        expect(result).toBe(42);
    }

    test_generic_method<T>(value: T): T {
        expect(value).toBeDefined();
        return value;
    }
}`,
			startLine: 10,
			endLine:   25,
			fileType:  "tests",
		},
		{
			name:    "Sample Test Suite",
			path:    "39749806ee9100df",
			summary: "describe('Sample test suite', () => {",
			source: `// Jest-style test suites with types
describe('Sample test suite', () => {
    it('should pass basic test', () => {
        expect(2 + 2).toBe(4);
    });

    it('should handle async operations', async (): Promise<void> => {
        const result: string = await Promise.resolve('test');
        expect(result).toBe('test');
    });

    it('should test with typed parameters', (): void => {
        const data: { name: string; value: number } = {
            name: 'test',
            value: 123
        };
        expect(data.name).toBe('test');
        expect(data.value).toBe(123);
    });
});`,
			startLine: 27,
			endLine:   46,
			fileType:  "tests",
		},
		{
			name:    "Test Data Interface",
			path:    "TestData",
			summary: "interface TestData {",
			source: `// Interface for test data
interface TestData {
    input: string;
    expected: number;
}`,
			startLine: 48,
			endLine:   52,
			fileType:  "tests",
		},
		{
			name:    "Interface Test Suite",
			path:    "609eeb753c9224f8",
			summary: "describe('Interface test suite', () => {",
			source: `// Test with interface
describe('Interface test suite', () => {
    const testCases: TestData[] = [
        { input: 'test1', expected: 1 },
        { input: 'test2', expected: 2 }
    ];

    testCases.forEach((testCase: TestData) => {
        it(` + "`should handle ${testCase.input}`" + `, () => {
            expect(testCase.input.length).toBeGreaterThan(0);
            expect(testCase.expected).toBeGreaterThan(0);
        });
    });
});`,
			startLine: 54,
			endLine:   67,
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
			s.Equal("typescript/tests/module.test.ts::"+test.path, chunk.ID())
		})
	}
}

func (s *TypeScriptParserTestSuite) TearDownSuite() {
	if s.parser != nil {
		s.parser.Close()
	}
}

func (s *TypeScriptParserTestSuite) TestVariableParsing() {
	chunks := s.getChunks("typescript/variables.ts")

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
			summary: `var simple_var: string = "hello";`,
			source: `// Test file for TypeScript variable declarations

// var declarations with types
var simple_var: string = "hello";`,
			startLine: 1,
			endLine:   4,
		},
		{
			name:      "Typed Number",
			path:      "typed_number",
			summary:   `var typed_number: number = 42;`,
			source:    `var typed_number: number = 42;`,
			startLine: 5,
			endLine:   5,
		},
		{
			name:      "Inferred Var",
			path:      "inferred_var",
			summary:   `var inferred_var = "inferred";`,
			source:    `var inferred_var = "inferred";`,
			startLine: 6,
			endLine:   6,
		},
		{
			name:    "Simple Let",
			path:    "simple_let",
			summary: `let simple_let: string = "world";`,
			source: `// let declarations with types
let simple_let: string = "world";`,
			startLine: 8,
			endLine:   9,
		},
		{
			name:      "Array Let",
			path:      "array_let",
			summary:   `let array_let: number[] = [1, 2, 3];`,
			source:    `let array_let: number[] = [1, 2, 3];`,
			startLine: 10,
			endLine:   10,
		},
		{
			name:      "Tuple Let",
			path:      "tuple_let",
			summary:   `let tuple_let: [string, number] = ["test", 42];`,
			source:    `let tuple_let: [string, number] = ["test", 42];`,
			startLine: 11,
			endLine:   11,
		},
		{
			name:    "Simple Const",
			path:    "simple_const",
			summary: `const simple_const: string = "constant";`,
			source: `// const declarations with types
const simple_const: string = "constant";`,
			startLine: 13,
			endLine:   14,
		},
		{
			name:    "Object Const",
			path:    "object_const",
			summary: `const object_const: { name: string; value: number } = {`,
			source: `const object_const: { name: string; value: number } = {
    name: "test",
    value: 123
};`,
			startLine: 15,
			endLine:   18,
		},
		{
			name:    "Arrow Function Const",
			path:    "arrow_function_const",
			summary: `const arrow_function_const = (x: number): number => x * 2;`,
			source: `// Arrow function with types assigned to variables
const arrow_function_const = (x: number): number => x * 2;`,
			startLine: 20,
			endLine:   21,
		},
		{
			name:    "Function Expression",
			path:    "function_expression",
			summary: `const function_expression = function(a: string, b: string): string {`,
			source: `// Function expressions with types assigned to variables
const function_expression = function(a: string, b: string): string {
    return a + b;
};`,
			startLine: 23,
			endLine:   26,
		},
		{
			name:    "Named Function Expression",
			path:    "named_function_expression",
			summary: `const named_function_expression = function calculator(x: number, y: number): number...`,
			source: `const named_function_expression = function calculator(x: number, y: number): number {
    return x - y;
};`,
			startLine: 28,
			endLine:   30,
		},
		{
			name:    "Simple Arrow",
			path:    "simple_arrow",
			summary: `const simple_arrow = (): string => "simple";`,
			source: `// Arrow functions with different type syntaxes
const simple_arrow = (): string => "simple";`,
			startLine: 32,
			endLine:   33,
		},
		{
			name:      "Arrow With Params",
			path:      "arrow_with_params",
			summary:   `const arrow_with_params = (a: number, b: number): number => a + b;`,
			source:    `const arrow_with_params = (a: number, b: number): number => a + b;`,
			startLine: 34,
			endLine:   34,
		},
		{
			name:    "Arrow With Body",
			path:    "767114ac2a5c8a3f",
			summary: `const arrow_with_body = (x: number): number => {`,
			source: `const arrow_with_body = (x: number): number => {
    const result: number = x * 2;
    return result;
};`,
			startLine: 35,
			endLine:   38,
		},
		{
			name:    "Async Arrow Func",
			path:    "b17a6b2599fdca36",
			summary: `const async_arrow_func = async (data: string): Promise<Response> => {`,
			source: `// Async arrow function with types
const async_arrow_func = async (data: string): Promise<Response> => {
    const response: Response = await fetch(data);
    return response;
};`,
			startLine: 40,
			endLine:   44,
		},
		{
			name:    "Destructured Object",
			path:    "destructured_object",
			summary: `let destructured_object: { name: string; age: number };`,
			source: `// Complex variable declarations with types
let destructured_object: { name: string; age: number };`,
			startLine: 46,
			endLine:   47,
		},
		{
			name:      "Destructured Array",
			path:      "destructured_array",
			summary:   `const destructured_array: [number, ...number[]] = [1, 2, 3];`,
			source:    `const destructured_array: [number, ...number[]] = [1, 2, 3];`,
			startLine: 48,
			endLine:   48,
		},
		{
			name:      "Union Type",
			path:      "union_type",
			summary:   `let union_type: string | number = "could be either";`,
			source:    `let union_type: string | number = "could be either";`,
			startLine: 49,
			endLine:   49,
		},
		{
			name:      "Optional Type",
			path:      "optional_type",
			summary:   `let optional_type: string | undefined;`,
			source:    `let optional_type: string | undefined;`,
			startLine: 50,
			endLine:   50,
		},
		{
			name:    "Generic Array",
			path:    "generic_array",
			summary: `let generic_array: Array<string> = ["one", "two"];`,
			source: `// Generic variable
let generic_array: Array<string> = ["one", "two"];`,
			startLine: 52,
			endLine:   53,
		},
		{
			name:      "Record Type",
			path:      "record_type",
			summary:   `let record_type: Record<string, number> = { a: 1, b: 2 };`,
			source:    `let record_type: Record<string, number> = { a: 1, b: 2 };`,
			startLine: 54,
			endLine:   54,
		},
		{
			name:    "Class Instance",
			path:    "class_instance",
			summary: `const class_instance: Date = new Date();`,
			source: `// Class instance with type
const class_instance: Date = new Date();`,
			startLine: 56,
			endLine:   57,
		},
		{
			name:    "Assertion Var",
			path:    "assertion_var",
			summary: `let assertion_var = "hello" as string;`,
			source: `// Type assertions
let assertion_var = "hello" as string;`,
			startLine: 59,
			endLine:   60,
		},
		{
			name:      "Angle Bracket Assertion",
			path:      "angle_bracket_assertion",
			summary:   `let angle_bracket_assertion = <number>42;`,
			source:    `let angle_bracket_assertion = <number>42;`,
			startLine: 61,
			endLine:   61,
		},
		{
			name:    "Readonly Array",
			path:    "readonly_array",
			summary: `const readonly_array = [1, 2, 3] as const;`,
			source: `// Readonly and const assertions
const readonly_array = [1, 2, 3] as const;`,
			startLine: 63,
			endLine:   64,
		},
		{
			name:      "Readonly Object",
			path:      "readonly_object",
			summary:   `let readonly_object: Readonly<{ x: number }> = { x: 10 };`,
			source:    `let readonly_object: Readonly<{ x: number }> = { x: 10 };`,
			startLine: 65,
			endLine:   65,
		},
		{
			name:      "Ambient Declaration",
			path:      "jQuery",
			summary:   `declare var jQuery: any;`,
			source:    `declare var jQuery: any;`,
			startLine: 67,
			endLine:   67,
		},
		{
			name:    "Exported Const",
			path:    "EXPORTED_CONST",
			summary: `const EXPORTED_CONST = "constant value";`,
			source: `// Exported variables
export const EXPORTED_CONST = "constant value";`,
			startLine: 69,
			endLine:   70,
		},
		{
			name:      "Exported VIEW_TYPE_REVIEW",
			path:      "VIEW_TYPE_REVIEW",
			summary:   `const VIEW_TYPE_REVIEW = "tutor-review";`,
			source:    `export const VIEW_TYPE_REVIEW = "tutor-review";`,
			startLine: 72,
			endLine:   72,
		},
		{
			name:      "Exported Let",
			path:      "exportedLet",
			summary:   `let exportedLet: number = 42;`,
			source:    `export let exportedLet: number = 42;`,
			startLine: 74,
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
			s.Equal("typescript/variables.ts::"+test.path, chunk.ID())
		})
	}
}

func TestTypeScriptParserTestSuite(t *testing.T) {
	suite.Run(t, new(TypeScriptParserTestSuite))
}
