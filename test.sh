#!/bin/bash -e

# clone the webthing-tester
if [ ! -d webthing-tester ]; then
    git clone https://github.com/WebThingsIO/webthing-tester
fi
pip3 install --user -r webthing-tester/requirements.txt

# build and test the single-thing example

# go run examples/single-thing/single-thing.go & EXAMPLE_PID=$!
# EXAMPLE_PID = go pid, not single-thing pid

go run examples/single-thing/single-thing.go >/dev/null 2>&1 &

sleep 5

function get_pid_by_listened_port() {
  pattern_str="*:8888"
  pid=$(ss -n -t -l -p | grep "$pattern_str" | column -t | awk -F ',' '{print $(NF-1)}')
  [[ $pid =~ "pid" ]] && pid=$(echo $pid | awk -F '=' '{print $NF}')
  EXAMPLE_PID=$pid
#   echo EXAMPLE_PID
}
get_pid_by_listened_port

# ./webthing-tester/test-client.py --skip-websocket --debug || ! echo 'Test failed' ; killall -9 single-thing
./webthing-tester/test-client.py --skip-websocket || ! echo 'Test failed' ; killall -9 single-thing

echo "single-thing test done!"
# kill -9 $EXAMPLE_PID


# build and test the multiple-things example
# ignore all print to std. >/dev/null 2>&1 &
go run examples/multiple-things/multiple-things.go > /dev/null 2>&1 &
sleep 5
get_pid_by_listened_port

# ignore test result and kill process
# ./webthing-tester/test-client.py --path-prefix "/0" --skip-websocket --debug || ! echo 'Test failed' ; killall -9 multiple-things
./webthing-tester/test-client.py --path-prefix "/0" --skip-websocket || ! echo 'Test failed' ; killall -9 multiple-things

killall -9 multiple-things
# kill -9 $EXAMPLE_PID