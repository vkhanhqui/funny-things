[128. Longest Consecutive Sequence](https://leetcode.com/problems/longest-consecutive-sequence)


```python
from typing import List
from collections import defaultdict


class Solution:
    def longestConsecutive(self, nums: List[int]) -> int:
        if nums:
            nums.sort()
            rs = defaultdict(set)
            still_consecutive = 1
            for index, num in enumerate(nums):
                if (
                    num == nums[index-1] + 1 or
                    num == nums[index-1]
                ):
                    rs[still_consecutive].add(num)
                    rs[still_consecutive].add(nums[index-1])
                else:
                    still_consecutive += 1
            max_index = max(rs, key=lambda k: len(rs[k]))
            return len(rs[max_index])
        return 0

```
