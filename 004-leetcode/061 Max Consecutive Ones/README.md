[Max Consecutive Ones](https://leetcode.com/explore/learn/card/fun-with-arrays/521/introduction/3238/)

```python
class Solution:
    def findMaxConsecutiveOnes(self, nums: List[int]) -> int:
        max_cons = 0
        cur_max_cons = 0
        for num in nums:
            if num == 1:
                cur_max_cons += 1
                max_cons = max(max_cons, cur_max_cons)
                continue
            cur_max_cons = 0
        return max_cons

```
