# howto

1. `make`
2. set `$SERVER_URI` to `:$port` or `$host:$port` to the same value for both client and server
3.
```sh
./bin/client < ./input.csv
```

## docker

paste this into posix shell

```sh
img=$(docker build -q .)
ctnr=$(docker run -p 7777:7777 -d --rm -q -e SERVER_URI=:7777 $img)
docker exec -i $ctnr /project/bin/client < ./input.csv
docker stop $ctnr
```
