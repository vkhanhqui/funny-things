[122. Best Time to Buy and Sell Stock II](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-ii)


```python
from typing import List


class Solution:
    def maxProfit(self, prices: List[int]) -> int:
        total = 0
        for index in range(1, len(prices)):
            if prices[index] > prices[index-1]:
                total += (prices[index] - prices[index-1])
        return total

```
