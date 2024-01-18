[875. Koko Eating Bananas](https://leetcode.com/problems/koko-eating-bananas/)

```python
class Solution:
    def minEatingSpeed(self, piles: List[int], h: int) -> int:
        left, right = 1, max(piles)
        k = right
        while left <= right:
            mid = (left + right) // 2
            total_time = 0
            if mid != 0:
                for p in piles:
                    total_time += math.ceil(p / mid)
            if total_time > h:
                left = mid + 1
            else:
                right = mid - 1
                k = min(k, mid)
        return k

```