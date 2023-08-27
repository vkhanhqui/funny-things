[1700. Number of Students Unable to Eat Lunch](https://leetcode.com/problems/number-of-students-unable-to-eat-lunch)

```cs
public class Solution {
    public int CountStudents(int[] students, int[] sandwiches) {
        int squares = 0; // students with numbers 0
        int circulars = 0; // students with numbers 1
        foreach(int student in students){
            if(student == 0)
            {
                squares += 1;
            }else
            {
                circulars += 1;
            }
        }

        foreach(int sandwiche in sandwiches){
            if(sandwiche == 0)
            {
                if(squares == 0) // no one want a square
                {
                    return circulars;
                }
                squares -= 1;
            }else
            {
                if(circulars == 0) // no one want a circular
                {
                    return squares;
                }
                circulars -= 1;
            }
        }

        return 0;
    }
}

```
