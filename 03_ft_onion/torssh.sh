#!/bin/sh

PORT=4242
PROXY=127.0.0.1:9050

if [ $# != 1 ]; then
    echo "Missing torssh target"
    exit 1
fi


ssh -v -p $PORT -o ProxyCommand="nc -X 5 -x $PROXY %h %p" $1