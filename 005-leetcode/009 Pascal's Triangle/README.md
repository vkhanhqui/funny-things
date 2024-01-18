[118. Pascal's Triangle](https://leetcode.com/problems/pascals-triangle)


```python
from typing import List


class Solution:
    def generate(self, numRows: int) -> List[List[int]]:
        result = [[1]]
        for _ in range(numRows - 1):
            temp_row = [0] + result[-1] + [0]
            current_row = []
            for j in range(len(temp_row) - 1):
                current_row.append(temp_row[j] + temp_row[j+1])
            result.append(current_row)
        return result

```
