#!/bin/sh

PORT=4242
PROXY=127.0.0.1:9050

# Source - https://stackoverflow.com/a/7069755
# Posted by Flexo - Save the data dump, modified by community. See post 'Timeline' for change history
# Retrieved 2026-11-11, License - CC BY-SA 4.0

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "torssh - SSH through Tor"
      echo " "
      echo "Usage: torssh [options] [user@]host"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "-p, --port=PORT           specify SSH port on the tor address (default: 4242)"
      echo "-x, --proxy=PROXY         specify the address of the tor instance to use as proxy (default: 127.0.0.1:9050)"
      exit 0
      ;;
    -p|--port)
      shift
      if test $# -gt 0; then
        PORT=$1
      else
        echo "no port specified"
        exit 1
      fi
      shift
      ;;
    --port=*)
      PORT=$(echo "$1" | sed -e 's/^[^=]*=//g')
      shift
      ;;
    -x|--proxy)
      shift
      if test $# -gt 0; then
        PROXY=$1
      else
        echo "no proxy specified"
        exit 1
      fi
      shift
      ;;
    --proxy=*)
      PROXY=$(echo "$1" | sed -e 's/^[^=]*=//g')
      shift
      ;;
    -*)
      echo "Unknown option: $1"
      exit 1
      ;;
    *)
      break
      ;;
  esac
done

if [ $# != 1 ]; then
    echo "Missing torssh target"
    exit 1
fi

TARGET=$1

ssh -v -p "$PORT" -o ProxyCommand="nc -X 5 -x $PROXY %h %p" "$TARGET"