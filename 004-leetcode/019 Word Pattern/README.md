[290. Word Pattern](https://leetcode.com/problems/word-pattern)


```python
class Solution:
    def wordPattern(self, pattern: str, s: str) -> bool:
        s_split = s.split()
        if len(pattern) != len(s_split):
            return False
        char_dict = {}
        word_dict = {}
        for p, w in zip(pattern, s_split):
            if (
                (p in char_dict and char_dict[p] != w) or
                (w in word_dict and word_dict[w] != p)
            ):
                return False
            char_dict[p] = w
            word_dict[w] = p
        return True

```
