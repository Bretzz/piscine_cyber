# Source - https://stackoverflow.com/a/7069755
# Posted by Flexo - Save the data dump, modified by community. See post 'Timeline' for change history
# Retrieved 2026-11-11, License - CC BY-SA 4.0

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "$package - ssh trough Tor"
      echo " "
      echo "$package [options] application [arguments]"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "-p, --port=PORT           specify SSH port on the tor address"
      echo "-x, --proxy=PROXY         specify the address of the tor instance to use as proxy"
      exit 0
      ;;
    -p|--port)
      shift
      if test $# -gt 0; then
        export PROCESS=$1
      else
        echo "no process specified"
        exit 1
      fi
      shift
      ;;
    --action*)
      export PROCESS=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    -o)
      shift
      if test $# -gt 0; then
        export OUTPUT=$1
      else
        echo "no output dir specified"
        exit 1
      fi
      shift
      ;;
    --output-dir*)
      export OUTPUT=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    *)
      break
      ;;
  esac
done
