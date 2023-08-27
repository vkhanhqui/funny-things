[Squares of a Sorted Array](https://leetcode.com/explore/learn/card/fun-with-arrays/521/introduction/3240/)

```python
class Solution:
    def sortedSquares(self, nums: List[int]) -> List[int]:
        squares = [num*num for num in nums]
        squares.sort()
        return squares

```
