[605. Can Place Flowers](https://leetcode.com/problems/can-place-flowers)

<b> LOL, I beat almost the players by this script</b>


![alt text](https://github.com/vkhanhqui/leetcode-solution/blob/main/013%20Can%20Place%20Flowers/ahihi.png?raw=true)



```python
from typing import List


class Solution:

    def canPlaceFlowers(
        self,
        flowerbed: List[int],
        n: int
    ) -> bool:
        places = 0
        flowerbed = [0] + flowerbed + [0]
        len_flowerbed = len(flowerbed)
        for index in range(1, len_flowerbed - 1):
            if (
                flowerbed[index] != 1 and
                flowerbed[index-1] != 1 and
                flowerbed[index+1] != 1
            ):
                flowerbed[index] = 1
                places += 1
        if places < n:
            return False
        return True

```
