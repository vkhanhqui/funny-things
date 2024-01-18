[496. Next Greater Element I](https://leetcode.com/problems/next-greater-element-i)


<b> Solution 2 is applied the monotonic stack so that it is easy to understand and faster than the solution 1</b>

![alt text](https://github.com/vkhanhqui/leetcode-solution/blob/main/015%20Next%20Greater%20Element%20I/ahihi.png?raw=true)


<b>Solution 1</b>


```python
from typing import List


class Solution:
    def nextGreaterElement(
        self,
        nums1: List[int],
        nums2: List[int]
    ) -> List[int]:
        result = []
        len_nums1 = len(nums1)
        len_nums2 = len(nums2)
        index_1 = 0
        index_2 = 0
        is_eq = False
        while (
            index_1 < len_nums1
        ):
            if (
                not is_eq and
                nums1[index_1] == nums2[index_2]
            ):
                is_eq = True
            if (
                is_eq and
                nums2[index_2] > nums1[index_1]
            ):
                result.append(nums2[index_2])
                index_1 += 1
                index_2 = 0
                is_eq = False
                continue
            index_2 += 1
            if (
                index_2 > len_nums2 - 1
            ):
                result.append(-1)
                index_1 += 1
                index_2 = 0
                is_eq = False
        return result

```

<b>Solution 2</b>


```python
from typing import List


class Solution:
    def nextGreaterElement(
        self,
        nums1: List[int],
        nums2: List[int]
    ) -> List[int]:
        nums1_indices = {
            number: index
            for index, number in enumerate(nums1)
        }
        result = [-1] * len(nums1)
        stack = []
        for index in range(len(nums2)):
            current_value_2 = nums2[index]
            while stack and current_value_2 > stack[-1]:
                number_1 = stack.pop()
                index_1 = nums1_indices[number_1]
                result[index_1] = current_value_2
            if current_value_2 in nums1_indices:
                stack.append(current_value_2)
        return result

```
