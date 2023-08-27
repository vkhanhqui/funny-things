[Duplicate Zeros](https://leetcode.com/explore/learn/card/fun-with-arrays/525/inserting-items-into-an-array/3245/)

```python
class Solution:
    def duplicateZeros(self, arr: List[int]) -> None:
        """
        Do not return anything, modify arr in-place instead.
        """
        new_arr = []
        len_arr = len(arr)
        i = 0
        while i < len_arr:
            cur_value = arr[i]
            if cur_value == 0:
                new_arr.append(cur_value)
                len_arr -= 1
            new_arr.append(cur_value)
            i += 1

        for i in range(len(arr)):
            arr[i] = new_arr[i]

```
