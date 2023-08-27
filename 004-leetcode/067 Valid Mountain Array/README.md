[Valid Mountain Array](https://leetcode.com/explore/learn/card/fun-with-arrays/527/searching-for-items-in-an-array/3251/)

```python
class Solution:
    def validMountainArray(self, arr: List[int]) -> bool:
        len_arr = len(arr)
        if len(arr) < 3:
            return False

        # start_at_first
        i = 0
        while i < len_arr - 1 and arr[i] < arr[i+1]:
            i += 1

        # start_at_last
        j = len(arr) - 1
        while j > 0 and arr[j] < arr[j-1]:
            j -= 1

        # if i == j => The top of the mountain
        return (
            i == j and
            i != 0 and
            j != len(arr) - 1
        )

```
