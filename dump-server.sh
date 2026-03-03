#!/bin/bash

##run those commands below before run script mk67te4232jqtg4bvyg
# eval "$(ssh-agent -s)"
# ssh-add "/home/ali/.ssh/ankarakendo2"

set -euo pipefail

set -a
source .env.server
set +a

HOST="root@$SERVERADDR"
{ sleep 0.5; echo "$KEYPASS"; } | script -q /dev/null -c "ssh-add $SERVERKEY"

sshpass -p "$KEYPASS" ssh -o BatchMode=yes -o ConnectTimeout=10 -i "$SERVERKEY" "$HOST" \
  "mysqldump -u$SERVERMYSQLUSER -p$SERVERMYSQLPASS \
  --databases godojogo \
  --routines --triggers --events \
  --single-transaction \
  --set-gtid-purged=OFF" > mysql-init/01-schema.sql
