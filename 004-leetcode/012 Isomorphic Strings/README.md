[205. Isomorphic Strings](https://leetcode.com/problems/isomorphic-strings)

<b>Solution 1</b>


```python
class Solution:
    def isIsomorphic(self, s: str, t: str) -> bool:
        s_dict, t_dict = {}, {}
        for i in range(len(s)):
            char_s, char_t = s[i], t[i]
            if (
                (char_s in s_dict and s_dict[char_s] != char_t) or
                (char_t in t_dict and t_dict[char_t] != char_s)
            ):
                return False
            s_dict[char_s] = char_t
            t_dict[char_t] = char_s
        return True

```

<b>Solution 2</b>


```python
class Solution:
    def isIsomorphic(self, s: str, t: str) -> bool:
        if len(s) != len(t):
            return False

        s_dict = {}
        for i in range(len(s)):
            char_s, char_t = s[i], t[i]
            if (
                char_s in s_dict
            ):
                if s_dict[char_s] != char_t:
                    return False
            elif char_t in s_dict.values():
                return False
            s_dict[char_s] = char_t
        return True

```