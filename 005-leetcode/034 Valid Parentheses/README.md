[20. Valid Parentheses](https://leetcode.com/problems/valid-parentheses)

```cs
public class Solution {
    public bool IsValid(string s) {
        var st = new Stack<char>();
        var di = new Dictionary<char, char>() {
            [')'] = '(',
            [']'] = '[',
            ['}'] = '{',
        };

        foreach (var c in s){
            if(!di.ContainsKey(c)){
                st.Push(c);
            }else if(st.Count == 0 || st.Pop() != di[c]){
                return false;
            }
        }

        return st.Count == 0;
    }
}
```