[876. Middle of the Linked List](https://leetcode.com/problems/middle-of-the-linked-list)

```cs
public class Solution {
    public int GetLength(ListNode head){
        ListNode temp = head;
        int length = 0;
        while(temp != null)
        {
            temp = temp.next;
            length += 1;
        }
        return length;
    }

    public ListNode MiddleNode(ListNode head) {
        int length = GetLength(head);
        for(int i = 0; i < length/2; i++)
        {
            head = head.next;
        }
        return head;
    }
}

```
