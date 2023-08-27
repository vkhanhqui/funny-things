[707. Design Linked List](https://leetcode.com/problems/design-linked-list)

```cs
public class ListNode
{
    public int val;
    public ListNode next;
    public ListNode prev;
    public ListNode(int val = 0)
    {
        this.val = val;
    }
}

public class MyLinkedList
{
    public ListNode left;
    public ListNode right;

    public MyLinkedList()
    {
        left = new ListNode();
        right = new ListNode();
        left.next = right;
        right.prev = left;
    }

    public int Get(int index)
    {
        var cur = left.next;
        int i = 0;
        while (cur != null && i < index)
        {
            i += 1;
            cur = cur.next;
        }
        if (cur != null && cur != this.right && i == index)
        {
            return cur.val;
        }
        return -1;
    }

    public void AddAtHead(int val)
    {
        var newNode = new ListNode(val);
        var prev = this.left;
        var next = this.left.next;

        newNode.prev = prev;
        newNode.next = next;

        next.prev = newNode;
        prev.next = newNode;
    }

    public void AddAtTail(int val)
    {
        var newNode = new ListNode(val);
        var prev = this.right.prev;
        var next = this.right;

        newNode.prev = prev;
        newNode.next = next;

        next.prev = newNode;
        prev.next = newNode;
    }

    public void AddAtIndex(int index, int val)
    {
        var cur = this.left.next;
        int i = 0;
        while (cur != null && i < index)
        {
            i += 1;
            cur = cur.next;
        }
        if (cur != null && i == index)
        {
            var newNode = new ListNode(val);
            var prev = cur.prev;

            newNode.next = cur;
            newNode.prev = prev;

            cur.prev = newNode;
            prev.next = newNode;
        }
    }

    public void DeleteAtIndex(int index)
    {
        var cur = this.left.next;
        int i = 0;
        while (cur != null && i < index)
        {
            i += 1;
            cur = cur.next;
        }
        if (cur != null && cur != this.right && i == index)
        {
            cur.next.prev = cur.prev;
            cur.prev.next = cur.next;
        }
    }
}
```