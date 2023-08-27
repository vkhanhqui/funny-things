[535. Encode and Decode TinyURL](https://leetcode.com/problems/encode-and-decode-tinyurl)


```python
class Codec:
    def __init__(self):
        self.url = 'http://tinyurl.com/'
        self.hash_url = {}

    def encode(self, longUrl: str) -> str:
        """Encodes a URL to a shortened URL.
        """
        self.hash_url['en'] = f'{self.url}{len(longUrl)}'
        self.hash_url['de'] = longUrl
        return self.hash_url['en']


    def decode(self, shortUrl: str) -> str:
        """Decodes a shortened URL to its original URL.
        """
        return self.hash_url['de']

```
