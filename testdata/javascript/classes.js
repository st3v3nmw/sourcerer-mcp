// Test file for JavaScript class definitions

// Simple class with no methods
class SimpleClass {
}

// Class with constructor and methods
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
        return `Value: ${this.value}`;
    }
}

// Class with inheritance
class ExtendedClass extends ClassWithMethods {
    constructor(value, name) {
        super(value);
        this.name = name;
    }

    // Override parent method
    getValue() {
        return `${this.name}: ${this.value}`;
    }

    // New method
    getName() {
        return this.name;
    }
}

// Class with private fields
class ClassWithPrivates {
    #privateField = 42;

    getPrivate() {
        return this.#privateField;
    }
}
