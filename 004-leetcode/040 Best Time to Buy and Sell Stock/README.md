[121. Best Time to Buy and Sell Stock](https://leetcode.com/problems/best-time-to-buy-and-sell-stock)

```cs
public class Solution {
    public int MaxProfit(int[] prices) {
        int maxProfit = 0;
        int minPrice = prices[0];
        foreach(var curPrice in prices){
            if(minPrice > curPrice){
                minPrice = curPrice;
            }
            var curProfit = curPrice - minPrice;
            if(maxProfit < curProfit){
                maxProfit = curProfit;
            }
        }
        return maxProfit;
    }
}
```