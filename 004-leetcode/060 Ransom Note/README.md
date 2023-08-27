[383. Ransom Note](https://leetcode.com/problems/ransom-note)

```python
from collections import Counter

class Solution:
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        count_ransom, count_magazine = Counter(ransomNote), Counter(magazine)
        if count_ransom & count_magazine == count_ransom:
            return True
        return False

```
