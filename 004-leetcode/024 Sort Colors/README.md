[75. Sort Colors](https://leetcode.com/problems/sort-colors)


```python
from typing import List
from collections import defaultdict


class Solution:
    def sortColors(self, nums: List[int]) -> None:
        """
        Do not return anything, modify nums in-place instead.
        """
        new_dict = defaultdict(list)
        for num in nums:
            new_dict[num].append(num)
        prev_slice = 0
        for i in range(0, 3, 1):
            nums[prev_slice:len(new_dict[i])+prev_slice] = new_dict[i]
            prev_slice += len(new_dict[i])

```
