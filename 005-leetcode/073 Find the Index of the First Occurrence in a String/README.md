[28. Find the Index of the First Occurrence in a String](https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string)

```python
class Solution:
    def strStr(self, haystack: str, needle: str) -> int:
        if needle not in haystack:
            return -1
        first = haystack.split(needle)[0]
        if first == haystack[0:len(first)]:
            return len(first)
        return 0

```
