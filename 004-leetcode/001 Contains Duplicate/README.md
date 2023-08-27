[217. Contains Duplicate](https://leetcode.com/problems/contains-duplicate)

<b>Solution 1</b>

```python
from typing import List


class Solution:
    def containsDuplicate(self, nums: List[int]) -> bool:
        len_nums = len(nums)
        len_nums_as_set = len(set(nums))
        if len_nums > len_nums_as_set:
            return True
        return False

```

<b>Solution 2</b>

```python
from typing import List


class Solution:
    def containsDuplicate(self, nums: List[int]) -> bool:
        return len(nums) != len(set(nums))

```