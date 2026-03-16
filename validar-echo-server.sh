#!/bin/bash
MESSAGE="hello"
SERVER="server"
SERVER_PORT="12345"
NETWORK="tp0_testing_net"

RESPONSE=$(docker run --rm --network $NETWORK busybox\
 sh -c "echo '$MESSAGE' | nc -w 1 $SERVER $SERVER_PORT")

if [ "$RESPONSE" = "$MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi