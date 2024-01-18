[412. Fizz Buzz](https://leetcode.com/problems/fizz-buzz)

```python
class Solution:
    def fizzBuzz(self, n: int) -> List[str]:
        return [self.__process(index) for index in range(1, n+1)]

    def __process(self, index: int) -> str:
        if index % 3 == 0 and index % 5 == 0:
            return "FizzBuzz"
        if index % 3 == 0:
            return "Fizz"
        if index % 5 == 0:
            return "Buzz"
        return str(index)

```
