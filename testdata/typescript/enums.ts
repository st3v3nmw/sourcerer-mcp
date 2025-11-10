// Test file for TypeScript enum definitions

// Simple numeric enum
enum SimpleEnum {
    First,
    Second,
    Third
}

// String enum
enum StringEnum {
    Red = "red",
    Green = "green",
    Blue = "blue"
}

// Mixed enum
enum MixedEnum {
    None = 0,
    Read = 1,
    Write = 2,
    Execute = 4,
    Description = "permissions"
}

// Computed enum
enum ComputedEnum {
    Base = 1,
    Double = Base * 2,
    Triple = Base * 3
}

// Const enum
const enum ConstEnum {
    Tiny = 1,
    Small = 2,
    Medium = 4,
    Large = 8
}

// Enum with explicit values
enum StatusEnum {
    Pending = "PENDING",
    InProgress = "IN_PROGRESS",
    Completed = "COMPLETED",
    Failed = "FAILED"
}

// Reverse mapped enum
enum BiDirectionalEnum {
    Up = "UP",
    Down = "DOWN",
    Left = "LEFT",
    Right = "RIGHT"
}

// Enum in namespace
namespace EnumNamespace {
    export enum NestedEnum {
        Option1 = 1,
        Option2 = 2
    }
}

// Exported enums
export enum ExportedEnum {
    Alpha = "alpha",
    Beta = "beta",
    Gamma = "gamma"
}

export enum Status {
    Pending = "pending",
    Active = "active",
    Complete = "complete"
}
