[1929. Concatenation of Array](https://leetcode.com/problems/concatenation-of-array)


```python
from typing import List


class Solution:
    def getConcatenation(self, nums: List[int]) -> List[int]:
        nums += nums
        return nums

```
