# Nested Comments

### Intro
A root comment can have up to 10 nested layers.
```
1. Comment A (root)
   └── 2. Reply to A
       └── 3. Reply to 2
            ...
            └── 10. Reply to 9
```

Instead of using recursion, I will build the feature  using `Stack` and `Map`.

### Table Design

```sql
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,

    content TEXT NOT NULL,
    path VARCHAR(255), -- stack of its parents

    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

### Setting Up

1. Start all services:
```
make compose-all
```

2. Shutdown all services:
```
make compose-down
```

### curl Sample

To keep it simple, I hardcoded the following ids in the script. Feel free to replace them:
- Post id: `019826d1-1933-7650-94f2-c04c076fbba6`
- User id: `019826d2-c7fd-7862-a1cf-3ac0cced6110`

1. Create a comment layer 1
    ```bash
    curl -X POST http://localhost:3000/comments \
    -H "Content-Type: application/json" \
    -d '{
        "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
        "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
        "content": "This is a new comment"
    }'
    ```

2. Create a comment layer 2
    ```bash
    curl -X POST http://localhost:3000/comments \
    -H "Content-Type: application/json" \
    -d '{
        "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
        "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
        "content": "This is a comment layer 2",
        "path": "{comment layer 1 id}"
    }'
    ```

3. Create a comment layer 3
    ```bash
    curl -X POST http://localhost:3000/comments \
    -H "Content-Type: application/json" \
    -d '{
        "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
        "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
        "content": "This is a comment layer 3",
        "path": "{comment layer 2 path}/{comment layer 2 id}"
    }'
    ```

4. List comments by post
    ```bash
    curl -x GET http://localhost:3000/comments/post/019826d1-1933-7650-94f2-c04c076fbba6
    ```

    Response:
    ```json
    [
        {
            "id": "019826d5-13b7-7788-935b-408bac22c8d2",
            "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
            "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
            "content": "This is a new comment",
            "createdAt": "2025-07-20T07:55:55.192Z",
            "updatedAt": "2025-07-20T07:55:55.192Z",
            "path": null,
            "replies": [
                {
                    "id": "019826d5-6a9b-7515-8025-e86b316f72fd",
                    "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
                    "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
                    "content": "This is a comment layer 2",
                    "createdAt": "2025-07-20T07:56:17.435Z",
                    "updatedAt": "2025-07-20T07:56:17.435Z",
                    "path": "019826d5-13b7-7788-935b-408bac22c8d2",
                    "replies": [
                        {
                            "id": "019826d6-f127-7268-ae69-cb56f047ded3",
                            "postId": "019826d1-1933-7650-94f2-c04c076fbba6",
                            "userId": "019826d2-c7fd-7862-a1cf-3ac0cced6110",
                            "content": "This is a comment layer 3",
                            "createdAt": "2025-07-20T07:57:57.415Z",
                            "updatedAt": "2025-07-20T07:57:57.415Z",
                            "path": "019826d5-13b7-7788-935b-408bac22c8d2/019826d5-6a9b-7515-8025-e86b316f72fd",
                            "replies": []
                        }
                    ]
                }
            ]
        }
    ]
    ```
