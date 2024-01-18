[450. Delete Node in a BST](https://leetcode.com/problems/delete-node-in-a-bst/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def deleteNode(self, root: Optional[TreeNode], key: int) -> Optional[TreeNode]:
        if not root:
            return None
        if key < root.val:
            root.left = self.deleteNode(root.left, key)
        elif key > root.val:
            root.right = self.deleteNode(root.right, key)
        # Found the key node
        else:
            if not root.left:
                return root.right
            elif not root.right:
                return root.left
            # Find min val in the right
            min_right_node = root.right
            while min_right_node.left:
                min_right_node = min_right_node.left
            # Replace key node by min node
            root.val = min_right_node.val
            root.right = self.deleteNode(root.right, root.val)
        return root

```