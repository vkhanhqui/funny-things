[199. Binary Tree Right Side View](https://leetcode.com/problems/binary-tree-right-side-view/)

```python
# Definition for a binary tree node.
# class TreeNode:
#     def __init__(self, val=0, left=None, right=None):
#         self.val = val
#         self.left = left
#         self.right = right
from collections import deque
class Solution:
    def rightSideView(self, root: Optional[TreeNode]) -> List[int]:
        if not root:
            return []
        new_q = deque([root])
        resp = []
        while new_q:
            right_side = None
            for _ in range(len(new_q)):
                cur = new_q.popleft()
                if cur:
                    right_side = cur
                    new_q.append(cur.left)
                    new_q.append(cur.right)
            if right_side:
                resp.append(right_side.val)
        return resp


```