This application implements http server which proceeds each request
And responds with the total number of requests.

example of running:

```
$ go build
$ ./coding-challenge
```

or

```
$ go run main.go
```
by defaults server run with next parameters:
```
addr =  ":8080" // "http server address
file  = "storage.txt" // "file for persisting"
count = "60s" // "time for counting"
```
if you want change default values, you must give your values like this:

```
$ go run main.go -file=yourfile -addr=youraddr -count=yourcounttime
```