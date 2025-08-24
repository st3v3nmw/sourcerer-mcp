"""Test file for Python test patterns."""

import unittest

def test_simple_function():
    assert True

class TestSample(unittest.TestCase):
    def test_method(self):
        self.assertEqual(1 + 1, 2)

    def test_another(self):
        self.assertTrue(True)
