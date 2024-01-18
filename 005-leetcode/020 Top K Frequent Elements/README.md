[347. Top K Frequent Elements](https://leetcode.com/problems/top-k-frequent-elements)

<b>Solution 1</b>

```python
from collections import Counter
from typing import List


class Solution:
    def topKFrequent(
        self,
        nums: List[int],
        k: int
    ) -> List[int]:
        nums_counter = Counter(nums).most_common(k)
        return list(
            map(lambda x: x[0], nums_counter)
        )

```

<b>Solution 2</b>

```python
from collections import Counter
from typing import List


class Solution:
    def topKFrequent(
        self,
        nums: List[int],
        k: int
    ) -> List[int]:
        result = []
        nums_counter = Counter(nums).items()
        len_nums = len(nums)
        frequent_nums = [[] for _ in range((len_nums + 1))]
        for key, v in nums_counter:
            frequent_nums[v].append(key)
        for i in range(len_nums, 0, -1):
            for n in frequent_nums[i]:
                result.append(n)
                k -= 1
            if k == 0:
                return result

```