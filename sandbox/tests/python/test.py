import unittest
from main import sum

class TestAddMethods(unittest.TestCase):
    def test_add(self):
        from_test = 1 + 4
        from_main = sum(1, 4)
        self.assertEqual(from_test, from_main, "sum is broken")

if __name__ == '__main__':
    unittest.main()
