[1523. Count Odd Numbers in an Interval Range](https://leetcode.com/problems/count-odd-numbers-in-an-interval-range)

```cs
public class Solution {
    public int CountOdds(int low, int high) {
        int countOdd = 0;
        for(int i=low; i<= high; i++)
        {
            if(i%2 != 0)
            {
                countOdd += 1;
            }
        }
        return countOdd;
    }
}

```
