// Test file for TypeScript interface definitions

// Simple interface
interface SimpleInterface {
    name: string;
    value: number;
}

// Interface with optional properties
interface OptionalInterface {
    required: string;
    optional?: number;
    readonly id: string;
}

// Interface with method signatures
interface MethodInterface {
    getName(): string;
    setValue(value: number): void;
    process(data: string): Promise<boolean>;
}

// Generic interface
interface GenericInterface<T, U = string> {
    data: T;
    metadata: U;
    transform<V>(input: T): V;
}

// Extending interfaces
interface ExtendedInterface extends SimpleInterface, MethodInterface {
    category: string;
    tags: string[];
}

// Interface with index signature
interface IndexedInterface {
    [key: string]: any;
    [key: number]: string;
}

// Interface with call signature
interface CallableInterface {
    (input: string): boolean;
    property: number;
}

// Interface with construct signature
interface ConstructableInterface {
    new (value: string): { result: string };
}

// Namespace interface
namespace Interfaces {
    export interface NamespacedInterface {
        value: boolean;
    }
}

// Merged interface declarations
interface MergedInterface {
    first: string;
}

interface MergedInterface {
    second: number;
}

// Exported interfaces
export interface ExportedInterface {
    id: string;
    process(): void;
}

export interface LLMProvider {
    name: string;
    generate(prompt: string): Promise<string>;
}
