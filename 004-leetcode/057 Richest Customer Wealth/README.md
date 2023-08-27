[1672. Richest Customer Wealth](https://leetcode.com/problems/richest-customer-wealth)

```python
class Solution:
    def maximumWealth(self, accounts: List[List[int]]) -> int:
        return max([sum(acc) for acc in accounts])

```
