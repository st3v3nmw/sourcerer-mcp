// Test file for JavaScript variable declarations

// var declarations
var simple_var = "hello";
var multiple_var_1, multiple_var_2;
var initialized_var = 42;

// let declarations
let simple_let = "world";
let destructured_let = {a: 1, b: 2};
let array_destructured = [1, 2, 3];

// const declarations
const simple_const = "constant";
const object_const = {
    name: "test",
    value: 123
};
const arrow_function_const = (x) => x * 2;

// Function expressions assigned to variables
const function_expression = function(a, b) {
    return a + b;
};

const named_function_expression = function calculator(x, y) {
    return x - y;
};

// Arrow functions with different syntaxes
const simple_arrow = () => "simple";
const arrow_with_params = (a, b) => a + b;
const arrow_with_body = (x) => {
    const result = x * 2;
    return result;
};

// Async arrow function
const async_arrow_func = async (data) => {
    const response = await fetch(data);
    return response;
};

// Complex variable declarations
let {name, age} = person;
const [first, ...rest] = numbers;
var globalVar = window.something || "default";

// Exported variables
export const EXPORTED_CONST = "constant value";
export const API_KEY = "api-key-12345";
export let exported_let = 42;
