[682. Baseball Game](https://leetcode.com/problems/baseball-game)

```cs
public class Solution {
    public int CalPoints(string[] operations) {
        var arr = new List<int>();
        foreach (var item in operations)
        {
            var lastIndex = arr.Count-1;
            switch (item)
            {
                case "C":
                    arr.RemoveAt(lastIndex);
                    break;
                case "D":
                    arr.Add(arr[lastIndex]*2);
                    break;
                case "+":
                    arr.Add(arr[lastIndex]+arr[lastIndex-1]);
                    break;
                default:
                    arr.Add(Int32.Parse(item));
                    break;
            }
        }
        return arr.Sum();
    }
}
```