#!/usr/bin/env bash

/project/bin/server &
/project/bin/client "$@"
kill %
