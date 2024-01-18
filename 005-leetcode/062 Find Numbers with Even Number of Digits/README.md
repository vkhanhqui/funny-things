[Find Numbers with Even Number of Digits](https://leetcode.com/explore/learn/card/fun-with-arrays/521/introduction/3237/)

```python
class Solution:
    def findNumbers(self, nums: List[int]) -> int:
        return sum([len(str(num)) % 2 == 0 for num in nums])

```
