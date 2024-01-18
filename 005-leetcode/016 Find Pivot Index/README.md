[724. Find Pivot Index](https://leetcode.com/problems/find-pivot-index)


```python
from typing import List


class Solution:
    def pivotIndex(self, nums: List[int]) -> int:
        total = sum(nums)
        len_nums = len(nums)
        left_side = 0
        for i in range(len_nums):
            right_side = total - nums[i] - left_side
            if left_side == right_side:
                return i
            left_side += nums[i]
        return -1

```
