"""Test file for Python function definitions."""

# Simple function with no parameters
def simple_function():
    pass

# Function with parameters and return value
def function_with_params(a, b):
    return a + b

# Property decorator example
@property
def decorated_function():
    return "decorated"

# Static method decorator
@staticmethod
def static_method():
    return "static"

# Class method decorator
@classmethod
def class_method(cls):
    return "class method"

# Async function example
async def async_function():
    return "async"

# Generator function example
#  foo bar baz

def generator_function():
    yield 1
    yield 2
