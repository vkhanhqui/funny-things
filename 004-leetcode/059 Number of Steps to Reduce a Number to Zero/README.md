[1342. Number of Steps to Reduce a Number to Zero](https://leetcode.com/problems/number-of-steps-to-reduce-a-number-to-zero)

```python
class Solution:
    def numberOfSteps(self, num: int) -> int:
        step = 0
        while num > 0:
            num = self.__process(num)
            step += 1
        return step

    def __process(self, current_num: int) -> int:
        if current_num % 2 == 0:
            return current_num / 2
        return current_num - 1

```
