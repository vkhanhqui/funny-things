[26. Remove Duplicates from Sorted Array](https://leetcode.com/problems/remove-duplicates-from-sorted-array)

### C#
```cs
using System.Collections;

public class Solution {
    public int RemoveDuplicates(int[] nums) {
        var filteredNums = new ArrayList() {-999};
        for (int i = 0; i < nums.Length; i++)
        {
            if (!filteredNums[filteredNums.Count-1].Equals(nums[i])){
                nums[filteredNums.Count-1] = nums[i];
                filteredNums.Add(nums[i]);
            }
        }
        return filteredNums.Count-1;
    }
}
```

### Python
```python
class Solution:
    def removeDuplicates(self, nums: List[int]) -> int:
        filtered_nums = [nums[0]]
        for num in nums:
            if filtered_nums[-1] != num:
                nums[len(filtered_nums)] = num
                filtered_nums.append(num)
        return len(filtered_nums)
```