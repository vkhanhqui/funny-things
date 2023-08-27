[242. Valid Anagram](https://leetcode.com/problems/valid-anagram)

```python
class Solution:
    def isAnagram(self, s: str, t: str) -> bool:
        return sorted(s) == sorted(t)

```
