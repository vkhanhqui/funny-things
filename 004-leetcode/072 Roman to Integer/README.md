[13. Roman to Integer](https://leetcode.com/problems/roman-to-integer)

```python
class Solution:
    def romanToInt(self, s: str) -> int:
        converter = {
            "I": 1,
            "V": 5,
            "X": 10,
            "L": 50,
            "C": 100,
            "D": 500,
            "M": 1000,
        }
        rs = 0
        for i in range(len(s)-1):
            cur_c = converter[s[i]]
            next_c = converter[s[i+1]]
            if cur_c < next_c:
                rs -= cur_c
                continue
            rs += cur_c
        return rs+converter[s[-1]]

```
