[1491. Average Salary Excluding the Minimum and Maximum Salary](https://leetcode.com/problems/average-salary-excluding-the-minimum-and-maximum-salary)

```cs
public class Solution {
    public double Average(int[] salary) {
        int max = 0, min = salary[0];
        double avg = 0;
        foreach(int num in salary)
        {
            max = Math.Max(num, max);
            min = Math.Min(num, min);
            avg += num;
        }
        return (avg - max - min)/ (salary.Length - 2);
    }
}

```
