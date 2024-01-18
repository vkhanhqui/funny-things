[Check If N and Its Double Exist](https://leetcode.com/explore/learn/card/fun-with-arrays/527/searching-for-items-in-an-array/3250/)

```python
class Solution:
    def checkIfExist(self, arr: List[int]) -> bool:
        new_hashmap = {}
        for num in arr:
            if new_hashmap.get(num * 2):
                return True
            if new_hashmap.get(num / 2):
                return True
            new_hashmap[num] = True
```
