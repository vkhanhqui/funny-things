[706. Design HashMap](https://leetcode.com/problems/design-hashmap)

```cs
public class MyHashMap {
    private int[] arr = Enumerable.Repeat(-1, 10000000).ToArray();
    public MyHashMap() {

    }

    public void Put(int key, int value) {
        this.arr[key] = value;
    }

    public int Get(int key) {
        return this.arr[key];
    }

    public void Remove(int key) {
        this.arr[key] = -1;
    }
}

/**
 * Your MyHashMap object will be instantiated and called as such:
 * MyHashMap obj = new MyHashMap();
 * obj.Put(key,value);
 * int param_2 = obj.Get(key);
 * obj.Remove(key);
 */
```