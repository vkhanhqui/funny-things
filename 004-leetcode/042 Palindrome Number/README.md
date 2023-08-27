[9. Palindrome Number](https://leetcode.com/problems/palindrome-number)

```cs
public class Solution {
    public bool IsPalindrome(int x) {
        if (x < 0){
            return false;
        }
        int reversedX = 0;
        int temp = x;
        while(temp != 0){
            reversedX = (reversedX * 10) + (temp % 10);
            temp /= 10;
        }
        return reversedX == x;
    }
}

```