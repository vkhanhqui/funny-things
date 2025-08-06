#!/bin/sh

# Install netcat to execute the "nc" command
apt-get update && apt-get install -y netcat

# Enable job control
set -m

# Start sql-server
./opt/mssql/bin/sqlservr & ./sql_wait.sh

# Display jobs output in the terminal
fg
