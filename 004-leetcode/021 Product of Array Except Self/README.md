[238. Product of Array Except Self](https://leetcode.com/problems/product-of-array-except-self)


```python
from typing import List


class Solution:
    def productExceptSelf(self, nums: List[int]) -> List[int]:
        len_nums = len(nums)
        rs = [1] * len_nums
        pre = 1
        for i in range(len_nums):
            rs[i] = pre
            pre *= nums[i]
        post = 1
        for i in range(len_nums-1, -1, -1):
            rs[i] *= post
            post *= nums[i]
        return rs

```
