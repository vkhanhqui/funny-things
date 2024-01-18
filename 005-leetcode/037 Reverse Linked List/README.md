[206. Reverse Linked List](https://leetcode.com/problems/reverse-linked-list)

```cs
public class Solution {
    public ListNode ReverseList(ListNode head) {
        ListNode prevNode = null;
        var currentNode = head;
        while (currentNode != null){
            var tempNode = currentNode.next;
            currentNode.next = prevNode;
            prevNode = currentNode;
            currentNode = tempNode;
        }
        return prevNode;
    }
}
```