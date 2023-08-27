[509. Fibonacci Number](https://leetcode.com/problems/fibonacci-number)

```cs
public class Solution {
    public int Fib(int n) {
        if(n == 0){
            return 0;
        }else if(n == 1){
            return 1;
        }
        return Fib(n-1) + Fib(n-2);
    }
}

```
