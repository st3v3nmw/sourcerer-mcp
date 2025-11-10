// Test file for TypeScript class definitions

// Simple class with no methods
class SimpleClass {
}

// Class with typed constructor and methods
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
        return `Value: ${this.value}`;
    }

    // Setter property
    set displayValue(val: string) {
        const match = val.match(/Value: (\d+)/);
        if (match) {
            this.value = parseInt(match[1]);
        }
    }
}

// Class with inheritance and generics
class ExtendedClass<T> extends ClassWithMethods {
    private name: T;

    constructor(value: number, name: T) {
        super(value);
        this.name = name;
    }

    // Override parent method with different return type
    getValue(): string {
        return `${this.name}: ${super.getValue()}`;
    }

    // New method with generic return
    getName(): T {
        return this.name;
    }
}

// Abstract class
abstract class AbstractClass {
    protected abstract process(): void;

    public execute(): void {
        this.process();
    }
}

// Interface implementation
interface Processable {
    process(): void;
}

class ImplementedClass extends AbstractClass implements Processable {
    process(): void {
        console.log("processing");
    }
}

// Class with readonly and optional properties
class ConfigClass {
    readonly id: string;
    name?: string;
    private _config: Record<string, any> = {};

    constructor(id: string, name?: string) {
        this.id = id;
        this.name = name;
    }
}

// Decorated

@sealed
class BugReport {
  type = "report";
  title: string;

  constructor(t: string) {
    this.title = t;
  }
}

// Exported classes
export class ExportedClass {
    private value: number;

    constructor(value: number) {
        this.value = value;
    }

    getValue(): number {
        return this.value;
    }
}

export class OpenRouterProvider implements Processable {
    name = "OpenRouter";

    process(): void {
        console.log("processing");
    }
}

export default class TutorPlugin {
    private config: any;

    constructor(config: any) {
        this.config = config;
    }
}
