// Test file for TypeScript test patterns

import { expect } from '@jest/globals';

// Simple test function
function test_simple_function(): void {
    expect(true).toBe(true);
}

// Test class with typed methods
class TestSample {
    test_method(): void {
        expect(1 + 1).toBe(2);
    }

    async test_async_method(): Promise<void> {
        const result: number = await Promise.resolve(42);
        expect(result).toBe(42);
    }

    test_generic_method<T>(value: T): T {
        expect(value).toBeDefined();
        return value;
    }
}

// Jest-style test suites with types
describe('Sample test suite', () => {
    it('should pass basic test', () => {
        expect(2 + 2).toBe(4);
    });

    it('should handle async operations', async (): Promise<void> => {
        const result: string = await Promise.resolve('test');
        expect(result).toBe('test');
    });

    it('should test with typed parameters', (): void => {
        const data: { name: string; value: number } = {
            name: 'test',
            value: 123
        };
        expect(data.name).toBe('test');
        expect(data.value).toBe(123);
    });
});

// Interface for test data
interface TestData {
    input: string;
    expected: number;
}

// Test with interface
describe('Interface test suite', () => {
    const testCases: TestData[] = [
        { input: 'test1', expected: 1 },
        { input: 'test2', expected: 2 }
    ];

    testCases.forEach((testCase: TestData) => {
        it(`should handle ${testCase.input}`, () => {
            expect(testCase.input.length).toBeGreaterThan(0);
            expect(testCase.expected).toBeGreaterThan(0);
        });
    });
});
