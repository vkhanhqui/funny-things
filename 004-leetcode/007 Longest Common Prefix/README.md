### This problem didn't explain clearly !!

[14. Longest Common Prefix](https://leetcode.com/problems/longest-common-prefix)

```python
from typing import List


class Solution:
    def longestCommonPrefix(self, strs: List[str]) -> str:
        if not strs:
            return ""
        if len(strs) == 1:
            return strs[0]
        strs.sort()
        result = ""
        for x, y in zip(strs[0], strs[-1]):
            if x == y:
                result += x
            else:
                break
        return result

```
