[Merge Sorted Array](https://leetcode.com/explore/learn/card/fun-with-arrays/525/inserting-items-into-an-array/3253/)

```python
class Solution:
    def merge(self, nums1: List[int], m: int, nums2: List[int], n: int) -> None:
        """
        Do not return anything, modify nums1 in-place instead.
        """
        if m + n == 0:
            return

        for i in range(n):
            nums1[i+m] = nums2[i]

        nums1.sort()

```
