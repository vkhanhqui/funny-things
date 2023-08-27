[1189. Maximum Number of Balloons](https://leetcode.com/problems/maximum-number-of-balloons)


```python
from collections import Counter


class Solution:
    def maxNumberOfBalloons(self, text: str) -> int:
        balloon_counter = Counter('balloon')
        text_counter = Counter(text)
        result = len(text)
        for k in balloon_counter:
            result = min(
                result,
                text_counter[k] // balloon_counter[k]
            )
        return result

```
