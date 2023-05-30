#!/bin/sh

# Health check to make sure the sql-server is ready completely
sleep 15
while ! nc -z sql-server 1433 ; do
    echo "Waiting for the SQL Server"
    sleep 5
done

echo "Execute sql script"
/opt/mssql-tools/bin/sqlcmd -U sa -P MyP@ssw0rd123 -i chatapp.sql

echo "Done sql script"
/opt/mssql-tools/bin/sqlcmd -U sa -P MyP@ssw0rd123 -S localhost -Q "SELECT name FROM sys.databases"
