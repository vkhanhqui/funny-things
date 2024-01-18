[Height Checker](https://leetcode.com/explore/learn/card/fun-with-arrays/523/conclusion/3228/)

```python
import copy
class Solution:
    def heightChecker(self, heights: List[int]) -> int:
        expected = copy.deepcopy(heights)
        expected.sort()

        count = 0
        for index, value in enumerate(heights):
            if value != expected[index]:
                count += 1
        return count

```
