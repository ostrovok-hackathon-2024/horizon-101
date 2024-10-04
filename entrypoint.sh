#!/usr/bin/env bash

export SERVER_URI=:7777
/project/bin/server &
/project/bin/client "$@"
kill %
