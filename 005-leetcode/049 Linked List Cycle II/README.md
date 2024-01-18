[142. Linked List Cycle II](https://leetcode.com/problems/linked-list-cycle-ii)

```cs
/**
 * Definition for singly-linked list.
 * public class ListNode {
 *     public int val;
 *     public ListNode next;
 *     public ListNode(int x) {
 *         val = x;
 *         next = null;
 *     }
 * }
 */
public class Solution {
    public ListNode DetectCycle(ListNode head) {
        ListNode fast = head, slow = head;
        bool isNotEq = true;
        while(fast != null && fast.next != null && isNotEq)
        {
            slow = slow.next;
            fast = fast.next.next;
            if(slow == fast){
                isNotEq = false;
            }
        }

        if (isNotEq)
        {
            return null;
        }

        while (head != slow) {
            head = head.next;
            slow = slow.next;
        }
        return slow;
    }
}

```
