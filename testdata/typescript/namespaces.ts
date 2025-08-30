// Test file for TypeScript namespace definitions

// Simple namespace
namespace SimpleNamespace {
    export const value = "namespace value";
    export function helper(): string {
        return "namespace function";
    }
}

// Nested namespaces
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
}

// Namespace with classes and interfaces
namespace UtilityNamespace {
    export interface Logger {
        log(message: string): void;
    }

    export class ConsoleLogger implements Logger {
        log(message: string): void {
            console.log(`[LOG]: ${message}`);
        }
    }

    export function createLogger(): Logger {
        return new ConsoleLogger();
    }
}

// Module declaration (ambient namespace)
declare namespace ExternalLibrary {
    interface Options {
        timeout: number;
    }

    function initialize(options: Options): void;
}

// Global augmentation
declare global {
    namespace NodeJS {
        interface ProcessEnv {
            CUSTOM_VAR: string;
        }
    }
}

// Namespace merging
namespace MergedNamespace {
    export const first = "first";
}

namespace MergedNamespace {
    export const second = "second";
}

// Namespace with generics
namespace GenericNamespace {
    export interface Container<T> {
        value: T;
    }

    export function wrap<T>(value: T): Container<T> {
        return { value };
    }
}
