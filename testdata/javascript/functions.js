// Test file for JavaScript function definitions

// Simple function declaration
function simple_function() {
    return "hello";
}

// Function with parameters and return value
function function_with_params(a, b) {
    return a + b;
}

// Arrow function assigned to variable
const arrow_function = () => {
    return "arrow";
};

// Arrow function with parameters
const arrow_with_params = (x, y) => {
    return x * y;
};

// Function expression assigned to variable
const function_expression = function() {
    return "expression";
};

// Named function expression
const named_function_expression = function namedFn() {
    return "named expression";
};

// Async function
async function async_function() {
    return "async result";
}

// Async arrow function
const async_arrow = async () => {
    return "async arrow";
};

// Generator function
function* generator_function() {
    yield 1;
    yield 2;
}

// Exported functions
export function exported_function(input) {
    return input.length > 0;
}

export async function async_exported_function(data) {
    return await Promise.resolve(data);
}

export default function default_exported_function() {
    return "default";
}

export function* exported_generator() {
    yield "a";
    yield "b";
}
