#!/usr/bin/env bash

# hyperfine "docker run -i $(docker build -q .) < ./input.csv"
export SERVER_URI=:7777
./bin/server &
hyperfine "./bin/client < ./input.csv"
kill %
