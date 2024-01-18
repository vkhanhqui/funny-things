[70. Climbing Stairs](https://leetcode.com/problems/climbing-stairs)

```cs
public class Solution {
    public int ClimbStairs(int n) {
        int far_one = 1;
        int far_two = 1;
        for(int i = 0; i < n - 1; i++){
            int temp = far_one;
            far_one += far_two;
            far_two = temp;
        }
        return far_one;
    }
}

```
