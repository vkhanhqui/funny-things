[155. Min Stack](https://leetcode.com/problems/min-stack)

```cs
public class MinStack {

    private Stack<int> Values { get; set; }
    private Stack<int> Mins { get; set; }

    public MinStack() {
        this.Values = new Stack<int>();
        this.Mins = new Stack<int>();
    }

    public void Push(int val) {
        var min = (this.Mins.Count == 0 || val < this.Mins.Peek()) ? val : this.Mins.Peek();
        this.Mins.Push(min);
        this.Values.Push(val);
    }

    public void Pop() {
        this.Mins.Pop();
        this.Values.Pop();
    }

    public int Top() {
        return this.Values.Peek();
    }

    public int GetMin() {
        return this.Mins.Peek();
    }
}

/**
 * Your MinStack object will be instantiated and called as such:
 * MinStack obj = new MinStack();
 * obj.Push(val);
 * obj.Pop();
 * int param_3 = obj.Top();
 * int param_4 = obj.GetMin();
 */
```