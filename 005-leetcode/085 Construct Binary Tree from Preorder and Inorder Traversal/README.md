[105. Construct Binary Tree from Preorder and Inorder Traversal](https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
class Solution:
    def buildTree(self, preorder: List[int], inorder: List[int]) -> Optional[TreeNode]:
        if not preorder or not inorder:
            return None
        # Create the root node using the first element of the preorder list
        root = TreeNode(preorder[0])

        # Find the index of the root value in the inorder list
        mid = inorder.index(root.val)

        # Recursively build the left and right subtrees

        # Left subtree:
        # preorder list contains elements from index 1 to mid+1
        # inorder list contains elements from index 0 to mid-1
        root.left = self.buildTree(preorder[1:mid + 1], inorder[:mid])

        # Right subtree:
        # preorder list contains elements from index mid+1 to the end
        # inorder list contains elements from index mid+1 to the end
        root.right = self.buildTree(preorder[mid + 1:], inorder[mid + 1:])
        return root

```