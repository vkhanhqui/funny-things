[215. Kth Largest Element in an Array](https://leetcode.com/problems/kth-largest-element-in-an-array)

```python
class Solution:
    def findKthLargest(self, nums: List[int], k: int) -> int:
        nums.sort()
        return nums[len(nums)-k]
```