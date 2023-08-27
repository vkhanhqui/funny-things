[303. Range Sum Query - Immutable](https://leetcode.com/problems/range-sum-query-immutable)

```cs
public class NumArray {
    private int[] _nums;

    public NumArray(int[] nums) {
        this._nums = nums;
    }

    public int SumRange(int left, int right) {
        return this._nums.Skip(left).Take(right - left + 1).Sum();
    }
}

/**
 * Your NumArray object will be instantiated and called as such:
 * NumArray obj = new NumArray(nums);
 * int param_1 = obj.SumRange(left,right);
 */
```