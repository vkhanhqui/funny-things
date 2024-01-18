[36. Valid Sudoku](https://leetcode.com/problems/valid-sudoku)


```python
from typing import List
from collections import defaultdict


class Solution:
    def isValidSudoku(self, board: List[List[str]]) -> bool:
        rows = defaultdict(set)
        cols = defaultdict(set)
        boxes = defaultdict(set)
        for i in range(9):
            for j in range(9):
                value = board[i][j]
                box = (i//3, j//3)
                if value == '.':
                    continue
                if (
                    value in rows[i] or
                    value in cols[j] or
                    value in boxes[box]
                ):
                    return False
                rows[i].add(value)
                cols[j].add(value)
                boxes[box].add(value)
        return True

```
