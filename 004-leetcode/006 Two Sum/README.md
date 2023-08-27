[1. Two Sum](https://leetcode.com/problems/two-sum)

```python
from typing import List


class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        seen_values = {}
        for index, value in enumerate(nums):
            remaining_value = target - value
            if remaining_value in seen_values:
                return [seen_values[remaining_value], index]
            seen_values[value] = index

```
