#!/bin/sh
# wait-for-db.sh

set -e

host="$1"
shift
cmd="$@"

until pg_isready -h "$(echo $host | cut -d':' -f1)" -p "$(echo $host | cut -d':' -f2)"; do
  >&2 echo "Postgres não está disponível ainda - aguardando..."
  sleep 1
done

>&2 echo "Postgres está disponível - continuando"
exec $cmd 