[700. Search in a Binary Search Tree](https://leetcode.com/problems/search-in-a-binary-search-tree/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def searchBST(self, root: Optional[TreeNode], val: int) -> Optional[TreeNode]:
        if not root:
            return None

        root_val = root.val
        if root_val > val:
            return self.searchBST(root.left, val)
        elif root_val < val:
            return self.searchBST(root.right, val)
        return root

```