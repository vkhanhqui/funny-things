[448. Find All Numbers Disappeared in an Array](https://leetcode.com/problems/find-all-numbers-disappeared-in-an-array)

<b>Solution 1</b>

```python
from typing import List


class Solution:
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        len_num = len(nums)
        one_to_n = [i for i in range(1, len_num+1)]
        rs = []
        for num in one_to_n:
            if num not in nums:
                rs.append(num)
        return rs

```

<b>Solution 2</b>


```python
from typing import List


class Solution:
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        set_nums = set(nums)
        return [
            index
            for index in range(1, len(nums) + 1)
            if index not in set_nums
        ]

```