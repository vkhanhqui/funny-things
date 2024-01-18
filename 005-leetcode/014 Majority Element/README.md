[169. Majority Element](https://leetcode.com/problems/majority-element)


<b>Solution 1</b>


```python
from typing import List
from collections import defaultdict


class Solution:
    def majorityElement(self, nums: List[int]) -> int:
        count_elements = defaultdict(int)
        for num in nums:
            count_elements[num] += 1
        return max(count_elements, key=count_elements.get)

```

<b>Solution 2</b>


```python
from typing import List
from collections import Counter


class Solution:
    def majorityElement(self, nums: List[int]) -> int:
        count_nums = Counter(nums).items()
        thredsold = len(nums)/2
        for k, v in count_nums:
            if v > thredsold:
                return k

```
