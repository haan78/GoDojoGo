#!/bin/bash
set -euo pipefail

set -a
source .env.server
set +a

HOST="root@$SERVERADDR"
REMOTE_DIR="/opt/GoDojoGo"
SERVICE="dojo"

SSH="ssh -o BatchMode=yes -o ConnectTimeout=10 -i \"$SERVERKEY\" $HOST"
SCP="scp -o BatchMode=yes -o ConnectTimeout=10 -i \"$SERVERKEY\""

echo "==> Cleaning output..."
rm -rf ./output/*
mkdir -p ./output

echo "==> Building..."
go build -C ./ -o ./output/GoDojoGo
go build -C ./cron  -o ../output/Cron

if [[ ! -e "./output/GoDojoGo" || ! -e "./output/Cron" ]]; then
  echo "Compile failed!"
  exit 1
fi

echo "==> Checking service state..."
if eval "$SSH \"systemctl is-active --quiet $SERVICE\""; then
  echo "==> Service is running -> stopping..."
  eval "$SSH \"systemctl stop $SERVICE\""
else
  echo "==> Service is not running -> continue..."
fi

echo "==> Ensuring remote folder exists: $REMOTE_DIR"
eval "$SSH \"mkdir -p '$REMOTE_DIR'\""

echo "==> Uploading binaries to $REMOTE_DIR ..."
# Upload both binaries
$SCP ./output/GoDojoGo ./output/Cron "$HOST:$REMOTE_DIR/"

echo "==> Setting permissions..."
eval "$SSH \"chmod +x '$REMOTE_DIR/GoDojoGo' '$REMOTE_DIR/Cron'\""

echo "==> Starting service..."
eval "$SSH \"systemctl start $SERVICE\""

echo "==> Verifying service..."
if eval "$SSH \"systemctl is-active --quiet $SERVICE\""; then
  echo "✅ Deploy complete: service '$SERVICE' is RUNNING."
  eval "$SSH \"systemctl --no-pager status $SERVICE | head -n 20\""
else
  echo "❌ Deploy failed: service '$SERVICE' is NOT running."
  echo "---- status ----"
  eval "$SSH \"systemctl --no-pager status $SERVICE || true\""
  echo "---- last logs ----"
  eval "$SSH \"journalctl -u $SERVICE -n 80 --no-pager || true\""
  exit 1
fi