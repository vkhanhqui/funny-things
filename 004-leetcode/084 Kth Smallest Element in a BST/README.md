[230. Kth Smallest Element in a BST](https://leetcode.com/problems/kth-smallest-element-in-a-bst/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def kthSmallest(self, root: Optional[TreeNode], k: int) -> int:
        ordered_list = []
        curr = root
        while curr or ordered_list:
            # Catch the left ones of the current node
            while curr:
                ordered_list.append(curr)
                curr = curr.left
            # Get the smallest
            curr = ordered_list.pop()
            k -= 1
            if k == 0:
                return curr.val
            # Go to the right side of the smallest
            curr = curr.right

```