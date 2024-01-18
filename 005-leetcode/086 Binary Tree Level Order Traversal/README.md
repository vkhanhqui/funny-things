[102. Binary Tree Level Order Traversal](https://leetcode.com/problems/binary-tree-level-order-traversal/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
from collections import deque
class Solution:
    def levelOrder(self, root: Optional[TreeNode]) -> List[List[int]]:
        resp = []
        new_queue = deque()
        if root:
            new_queue.append(root)
        while new_queue:
            level_elements = []
            for _ in range(len(new_queue)):
                cur = new_queue.popleft()
                if cur:
                    level_elements.append(cur.val)
                    if cur.left:
                        new_queue.append(cur.left)
                    if cur.right:
                        new_queue.append(cur.right)
            resp.append(level_elements)
        return resp

```