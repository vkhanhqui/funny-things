[Third Maximum Number](https://leetcode.com/explore/learn/card/fun-with-arrays/523/conclusion/3231/)

```python
class Solution:
    def thirdMax(self, nums: List[int]) -> int:
        new_set = sorted(set(nums))

        listed_eles = list(new_set)

        if len(listed_eles) > 2:
            return listed_eles[-3]
        return listed_eles[-1]

```
