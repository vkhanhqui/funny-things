[912. Sort an Array](https://leetcode.com/problems/sort-an-array/)

### Insertion sort: Time Limit Exceeded
```python
class Solution:
    def sortArray(self, nums: List[int]) -> List[int]:
        len_nums = len(nums)
        for i in range(len_nums):
            j = i - 1
            while j > -1 and nums[j+1] < nums[j]:
                temp = nums[j+1]
                nums[j+1] = nums[j]
                nums[j] = temp
                j -= 1
        return nums

```

### Merge Sort: Passed
```python
class Solution:
    def sortArray(self, nums: List[int]) -> List[int]:
        len_nums = len(nums)
        if len_nums <= 1:
            return nums

        # Divide the list into halves & sort within the merge
        mid = len_nums // 2
        left_nums = self.sortArray(nums[:mid])
        right_nums = self.sortArray(nums[mid:])

        return self.merge(left_nums, right_nums)

    def merge(
        self,
        left_nums: List[int],
        right_nums: List[int]
    ) -> List[int]:
        merged_nums = []
        l, r = 0, 0
        # Compare elements from both halves and merge them in sorted order
        while (l < len(left_nums) and
            r < len(right_nums)
        ):
            if left_nums[l] < right_nums[r]:
                merged_nums.append(left_nums[l])
                l += 1
            else:
                merged_nums.append(right_nums[r])
                r += 1
        # Add remaining elements from the left and right half
        merged_nums += left_nums[l:] + right_nums[r:]
        return merged_nums

```

### Quick Sort: Memory Limit Exceeded
```python
class Solution:
    def sortArray(self, nums: List[int]) -> List[int]:
        # If the list has 1 or 0 elements, it is already sorted
        if len(nums) <= 1:
            return nums

        # Choose the first element as the pivot
        pivot = nums[0]

        # Divide the list into two partitions:
        # - elements less than or equal to the pivot
        # - elements greater than the pivot
        less = [x for x in nums[1:] if x <= pivot]
        greater = [x for x in nums[1:] if x > pivot]

        # Recursively sort the "less" and "greater" partitions
        sorted_less = self.sortArray(less)
        sorted_greater = self.sortArray(greater)

        # Combine the sorted partitions and the pivot
        return sorted_less + [pivot] + sorted_greater

```