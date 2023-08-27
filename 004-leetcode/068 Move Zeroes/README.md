[Move Zeroes](https://leetcode.com/explore/learn/card/fun-with-arrays/511/in-place-operations/3157/)

```python
class Solution:
    def moveZeroes(self, nums: List[int]) -> None:
        """
        Do not return anything, modify nums in-place instead.
        """
        filtered_nums = [num for num in nums if num != 0]
        zero_nums = len(nums) - len(filtered_nums)
        for i in range(len(nums)):
            if i > len(filtered_nums) - 1:
                nums[i] = 0
            else:
                nums[i] = filtered_nums[i]

```
