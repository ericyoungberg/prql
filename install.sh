#!/bin/bash
#
# This script installs the system files for PrQL
# on *nix systems.
#

set -e;

# Install binaries

cp build/prql /usr/bin/prql
cp build/prqld /usr/bin/prqld

# Setup the working directory

PRQL_DIR=/var/lib/prql

mkdir -p $PRQL_DIR
touch "$PRQL_DIR/databases" "$PRQL_DIR/tokens"
cp config/prql.toml "$PRQL_DIR/prql.toml"

if which groupadd > /dev/null; then
  if ! grep -qE "^prql:" /etc/group; then
    groupadd prql
  fi

  chown -R root:prql $PRQL_DIR
else
  echo 'Could not create prql group'
fi



# Setup prqld service if possible

if which systemctl > /dev/null; then
  cp config/prqld.service /lib/systemd/system/prqld.service
  systemctl enable prqld
else
  echo 'Could not setup systemd service'
fi

echo 'PrQL Installed'
