[929. Unique Email Addresses](https://leetcode.com/problems/unique-email-addresses)


```python
from typing import List


class Solution:

    def parse_email(self, email: str) -> str:
        local, domain = email[:email.index('@')], email[email.index('@'):]
        local = local.replace('.', '')
        try:
            local = local[:local.index('+')]
        except:
            pass
        return local + domain

    def numUniqueEmails(self, emails: List[str]) -> int:
        rs = set(
            map(
                self.parse_email,
                emails
            )
        )
        return len(rs)

```
