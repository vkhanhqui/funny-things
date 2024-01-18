[1480. Running Sum of 1d Array](https://leetcode.com/problems/running-sum-of-1d-array)

```cs
public class Solution {
    public int[] RunningSum(int[] nums) {
        int numsLen = nums.Length;
        var rs = new int[numsLen];
        var sum = 0;
        for(int i = 0; i < numsLen; i++){
            sum += nums[i];
            rs[i] = sum;
        }
        return rs;
    }
}

```