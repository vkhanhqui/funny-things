[225. Implement Stack using Queues](https://leetcode.com/problems/implement-stack-using-queues)

```cs
public class MyStack {
    public Queue MyQueue;

    public MyStack() {
        MyQueue = new Queue();
    }

    public void Push(int x) {
        MyQueue.Enqueue(x);
    }

    public int Pop() {
        int length = MyQueue.Count;
        while(length != 1){
            MyQueue.Enqueue(MyQueue.Dequeue());
            length -= 1;
        }
        return (int)MyQueue.Dequeue();
    }

    public int Top() {
        int length = MyQueue.Count;
        while(length != 1){
            MyQueue.Enqueue(MyQueue.Dequeue());
            length -= 1;
        }
        int top =  (int)MyQueue.Dequeue();
        MyQueue.Enqueue(top);
        return top;
    }

    public bool Empty() {
        return MyQueue.Count == 0;
    }
}

/**
 * Your MyStack object will be instantiated and called as such:
 * MyStack obj = new MyStack();
 * obj.Push(x);
 * int param_2 = obj.Pop();
 * int param_3 = obj.Top();
 * bool param_4 = obj.Empty();
 */

```
