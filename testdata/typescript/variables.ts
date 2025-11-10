// Test file for TypeScript variable declarations

// var declarations with types
var simple_var: string = "hello";
var typed_number: number = 42;
var inferred_var = "inferred";

// let declarations with types
let simple_let: string = "world";
let array_let: number[] = [1, 2, 3];
let tuple_let: [string, number] = ["test", 42];

// const declarations with types
const simple_const: string = "constant";
const object_const: { name: string; value: number } = {
    name: "test",
    value: 123
};

// Arrow function with types assigned to variables
const arrow_function_const = (x: number): number => x * 2;

// Function expressions with types assigned to variables
const function_expression = function(a: string, b: string): string {
    return a + b;
};

const named_function_expression = function calculator(x: number, y: number): number {
    return x - y;
};

// Arrow functions with different type syntaxes
const simple_arrow = (): string => "simple";
const arrow_with_params = (a: number, b: number): number => a + b;
const arrow_with_body = (x: number): number => {
    const result: number = x * 2;
    return result;
};

// Async arrow function with types
const async_arrow_func = async (data: string): Promise<Response> => {
    const response: Response = await fetch(data);
    return response;
};

// Complex variable declarations with types
let destructured_object: { name: string; age: number };
const destructured_array: [number, ...number[]] = [1, 2, 3];
let union_type: string | number = "could be either";
let optional_type: string | undefined;

// Generic variable
let generic_array: Array<string> = ["one", "two"];
let record_type: Record<string, number> = { a: 1, b: 2 };

// Class instance with type
const class_instance: Date = new Date();

// Type assertions
let assertion_var = "hello" as string;
let angle_bracket_assertion = <number>42;

// Readonly and const assertions
const readonly_array = [1, 2, 3] as const;
let readonly_object: Readonly<{ x: number }> = { x: 10 };

declare var jQuery: any;

// Exported variables
export const EXPORTED_CONST = "constant value";

export const VIEW_TYPE_REVIEW = "tutor-review";

export let exportedLet: number = 42;
