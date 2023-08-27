[125. Valid Palindrome](https://leetcode.com/problems/valid-palindrome)

```cs
using System.Text.RegularExpressions;

public class Solution {

        public bool IsPalindrome(string s)
        {
            var resp = false;
            var parsedS = s.ToLower().Replace(" ", "");
            Regex rgx = new Regex("[^a-zA-Z0-9]");
            parsedS = rgx.Replace(parsedS, "");
            if (parsedS.Length < 2)
            {
                resp = true;
            }
            else
            {
                var dividedLen = parsedS.Length / 2;
                var headS = parsedS.Substring(0, dividedLen);
                var tailS = parsedS.Substring(dividedLen, dividedLen);
                if (parsedS.Length %2 != 0)
                {
                    tailS = parsedS.Substring(dividedLen+1, dividedLen);
                }
                var reversedTailS = string.Join("", tailS.Reverse().ToArray());
                if (headS.Equals(reversedTailS))
                {
                    resp = true;
                }
            }
            return resp;
        }
}
```