#!/bin/bash -e

# clone the webthing-tester
if [ ! -d webthing-tester ]; then
    git clone https://github.com/WebThingsIO/webthing-tester
fi
pip3 install --user -r webthing-tester/requirements.txt

# build and test the single-thing example
go run examples/single-thing/single-thing.go &
EXAMPLE_PID=$!
sleep 5
./webthing-tester/test-client.py --debug
kill -15 $EXAMPLE_PID

# build and test the multiple-things example
# go run examples/multiple-things/multiple-things.go &
# EXAMPLE_PID=$!
# sleep 5
# ./webthing-tester/test-client.py --path-prefix "/0"
# kill -15 $EXAMPLE_PID