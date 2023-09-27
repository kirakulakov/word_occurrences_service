# Words occurrences from posts comments service.

<br>
<br>

### Before start
```
DROP TABLE word; CREATE TABLE word(post_id int, word text, count int);
```

<br>

##### Run service
```
go run cmd/app/main.go 
```
