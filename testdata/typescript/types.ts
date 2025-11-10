// Test file for TypeScript type definitions

// Simple type alias
type SimpleType = string;

// Union type
type UnionType = string | number | boolean;

// Intersection type
type IntersectionType = { name: string } & { age: number };

// Generic type alias
type GenericType<T> = {
    value: T;
    optional?: T;
};

// Function type alias
type FunctionType = (input: string) => boolean;

// Complex function type with overloads
type OverloadedFunction = {
    (input: string): string;
    (input: number): number;
    (input: boolean): boolean;
};

// Object type alias
type ObjectType = {
    readonly id: string;
    name: string;
    age?: number;
    [key: string]: any;
};

// Array type aliases
type StringArray = string[];
type NumberTuple = [number, number, string?];

// Conditional type
type ConditionalType<T> = T extends string ? string[] : T[];

// Mapped type
type MappedType<T> = {
    [K in keyof T]: T[K] | null;
};

// Template literal type
type TemplateType = `prefix_${string}_suffix`;

// Utility type usage
type PartialUser = Partial<{ name: string; age: number }>;
type RequiredUser = Required<{ name?: string; age?: number }>;

// Enum-like type
type StatusType = 'pending' | 'completed' | 'failed';

// Recursive type
type RecursiveType<T> = T | RecursiveType<T>[];

// Type with generic constraints
type ConstrainedType<T extends Record<string, any>> = {
    data: T;
    keys: keyof T;
};

// Exported types
export type ExportedType = string | number;

export type ConfigOptions = {
    enabled: boolean;
    timeout: number;
};
