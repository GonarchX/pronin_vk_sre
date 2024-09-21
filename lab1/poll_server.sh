#!/bin/bash

URL="http://localhost:8888/status"
STATUS_FILE="/etc/opt/webapp/status.txt"
LOG_FILE="/etc/opt/webapp/error.txt"

response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" -X GET $URL)
body=$(echo $response | sed -e 's/HTTPSTATUS.*//')
http_status=$(echo $response | sed -e 's/.*HTTPSTATUS://')

if [[ "$body" == "OK" ]]; then
    echo "SUCCESS - $(date) - HTTP Status: $http_status - Response Body: '$body'" > $STATUS_FILE
else
    echo "ERROR - $(date) - HTTP Status: $http_status - Response Body: '$body'" > $STATUS_FILE
    echo "Error: server returns response with invalid status code (HTTP Status: '$http_status', Response Body: '$body')" >> $LOG_FILE
fi