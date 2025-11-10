// Test file for TypeScript function definitions

// Simple function declaration with type annotations
function simple_function(): string {
    return "hello";
}

// Function with typed parameters and return value
function function_with_params(a: number, b: number): number {
    return a + b;
}

// Arrow function assigned to variable with types
const arrow_function = (): string => {
    return "arrow";
};

// Arrow function with typed parameters
const arrow_with_params = (x: number, y: number): number => {
    return x * y;
};

// Function expression assigned to variable with types
const function_expression = function(a: string, b: string): string {
    return a + b;
};

// Named function expression with types
const named_function_expression = function namedFn(value: any): string {
    return "named expression";
};

// Async function with Promise return type
async function async_function(): Promise<string> {
    return "async result";
}

// Async arrow function with typed parameter
const async_arrow = async (data: string): Promise<string> => {
    return "async arrow";
};

// Generator function with typed yield
function* generator_function(): Generator<number, void, unknown> {
    yield 1;
    yield 2;
}

// Generic function
function generic_function<T>(value: T): T {
    return value;
}

// Function with optional parameters
function optional_params(required: string, optional?: number): string {
    return required + (optional || 0);
}

// Function with rest parameters
function rest_params(first: string, ...rest: number[]): string {
    return first + rest.join(',');
}

// Function signature
function myFunc(x: number): string;

// Exported functions
export function exportedFunction(input: string): boolean {
    return input.length > 0;
}

export async function asyncExportedFunction(data: any): Promise<void> {
    await Promise.resolve(data);
}

export default function defaultExportedFunction(): string {
    return "default";
}
