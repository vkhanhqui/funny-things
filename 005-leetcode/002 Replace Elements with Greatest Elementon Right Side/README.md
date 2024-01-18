[1299. Replace Elements with Greatest Element on Right Side](https://leetcode.com/problems/replace-elements-with-greatest-element-on-right-side)

<b>Solution 1</b>

```python
class Solution:
    def replaceElements(self, arr: List[int]) -> List[int]:
        # start_at_last
        index = len(arr) - 1
        max_value = -1
        while index > -1:
            cur_value = arr[index]
            arr[index] = max_value
            max_value = max(cur_value, max_value)
            index -= 1
        return arr

```

<b>Solution 2</b>

```python
from typing import List


class Solution:
    def replaceElements(self, arr: List[int]) -> List[int]:
        for i in range(len(arr)-1):
            arr[i] = max(arr[i+1:])
        arr[-1] = -1
        return arr

```