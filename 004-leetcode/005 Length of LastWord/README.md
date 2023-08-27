[58. Length of Last Word](https://leetcode.com/problems/length-of-last-word)

```python
class Solution:
    def isSubsequence(self, s: str, t: str) -> bool:
        len_s, len_t = len(s), len(t)
        index_t = 0
        while (
            len_s > 0 and
            index_t < len_t
        ):
            if s[0] == t[index_t]:
                s = s[1:]
                len_s -= 1
            index_t += 1
        return len_s == 0

```
