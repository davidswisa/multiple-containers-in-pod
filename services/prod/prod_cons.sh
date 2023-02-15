#!/bin/bash

# start prod
cd /prod && go run . &

# start cons
cd /cons && go run . &

wait -n

exit $?