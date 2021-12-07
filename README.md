# logrep
用于处理这种形式的多行日志：
```
2021/12/07 03:45:23 [error] 21009#21009: *26474 broken header: "GET / HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip

" while reading PROXY protocol, client: 127.0.0.1, server: 127.0.0.1:80
2021/12/07 03:45:23 [error] 21009#21009: *26475 broken header: "GET / HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip

" while reading PROXY protocol, client: 127.0.0.1, server: 127.0.0.1:80
2021/12/07 03:45:23 [error] 21009#21009: *26476 broken header: "GET / HTTP/1.1
Host: localhost
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip

" while reading PROXY protocol, client: 127.0.0.1, server: 127.0.0.1:80
```

使用方式类似grep

```
$ ./logrep -b '^[0-9]{4}' keep test.log 

match:        "keep"
file:         "test.log"
blockheading: "^[0-9]{4}"
delim:        false
nohighlight:  false
ignorecase:   false

2021/12/07 03:45:24 [error] 21010#21010: *26485 broken header: "GET / HTTP/1.1
Host: localhost
Accept-Encoding: identity
connection: keep-alive

" while reading PROXY protocol, client: 127.0.0.1, server: 127.0.0.1:80
2021/12/07 03:45:24 [error] 21010#21010: *26486 broken header: "GET / HTTP/1.1
Host: localhost
Accept-Encoding: identity
connection: keep-alive

" while reading PROXY protocol, client: 127.0.0.1, server: 127.0.0.1:80
2021/12/07 03:45:24 [error] 21010#21010: *26487 broken header: "GET / HTTP/1.1
Host: 127.0.0.1
Accept-Encoding: identity
connection: keep-alive

==== total 3 matches ====

```
