### Bulk Document

- **Request Body:** (example)
  ```json
    {"id": "1", "name": "John Doe", "age": 30, "title":"some title1", "content":"some content1"}
    {"id": "2", "name": "Jane Doe", "age": 25, "title":"some title2", "content":"some content2"}
    {"name": "Jake Smith", "age": 40, "title":"some title3", "content":"some content3"}
  ```

*Note: If the request does not include an `id` field, a universally unique identifier (UUID) will be automatically generated for identification purposes.*
