## Limit
Query결과에 대해 n개의 갯수 제한을 줍니다.
````go
users, err := client.User.
    Query().
    Limit(n).
    All(ctx)
````

## Offset
Query 결과에서 결과를 볼 시작 부분을 지정해줍니다.
```go
users, err := client.User.
    Query().
    Offset(10).
    All(ctx)
```

## 기타 예시
```go
users, err := client.User.
    Query().
    Offset(1).
	Limit(2).
    All(ctx)
```
