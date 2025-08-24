from dataclasses import dataclass

# Simple class with no methods
class SimpleClass:
    pass

class ClassWithMethods:
    value = -1

    # Constructor method
    def __init__(self):
        self.value = 0

    # Instance method
    def method(self):
        return self.value

    # Property method with decorator
    @property
    def property_method(self):
        return self.value * 2

# Class with inheritance
class InheritedClass(ClassWithMethods):
    # Override parent method
    def method(self):
        return super().method() + 1

# Decorated class using dataclass
@dataclass
class DecoratedClass:
    name: str
    age: int
