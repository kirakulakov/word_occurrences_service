# Words occurrences from posts comments service.
<br>

## Before start
- ```
  DROP TABLE word; CREATE TABLE word(post_id INT, word TEXT, count INT);
  ```
- Also specify Postgres `url` in `config/config.yml` or from `env`.

<br>

#### Run service
```
go run cmd/app/main.go 
```

<br>

#### See docs
```
http://127.0.0.1:8080/swagger/index.html
```
