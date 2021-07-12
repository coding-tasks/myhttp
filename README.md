# MyHTTP

`myhttp` is a small command line tool that makes HTTP requests and prints the address of the request along with the MD5
hash of the response.

```sh
$ myhttp -parallel 2 ankit.pl github.com amazon.com google.com

https://ankit.pl    dd0f4c92548e2c0f2fbda56fc723773f
https://github.com  510a202b4884a6baef1f0887faf43c4b
https://amazon.com  e226ce8c962ecc7acb28a10531625425
https://google.com  d877bc2867734491299b4baa352606d8
```

### Installation

Install the runnable binary to your `$GOPATH/bin`

```sh
go get github.com/coding-tasks/myhttp
```

Or, you can build and run it locally with

```sh
make build
./myhttp -h
```

### Assumptions/Limitations

- Only the `200` status code from the host is considered as a valid response.
- HTTP client timeout is set to 10 seconds and is not configurable via the command line.
- In case of error, the error message is displayed in the output instead of the md5 hash.
- The app will append `https://` scheme automatically if it is not provided. However, it doesn't validate if the
  provided URL is in the expected format beforehand.
- There is no upper limit on the number of workers, so you can technically spin up any number of worker routines as you
  want.

### Design

The tool is designed based on the concept of bounded parallelism where we spin up N number of worker threads that
communicate via a shared channel, `host` and `out` in our case. We continuously feed URL we want to fetch data from to
the `host` channel. The worker that is free will pick up the host, computes the result, and send it to the `out`
channel. The results are displayed from the `out` channel as soon as they are available.  
