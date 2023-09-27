# Words occurrences from posts comments service.
<br>

### Before start
```
DROP TABLE word; CREATE TABLE word(post_id INT, word TEXT, count INT);
```

<br>

#### Run service
```
go run cmd/app/main.go 
```
