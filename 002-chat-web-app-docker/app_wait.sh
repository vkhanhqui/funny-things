#!/bin/sh

# Sleep 60s in order to wait for the sql-server is ready
sleep 60

# Start Tomcat
catalina.sh run
