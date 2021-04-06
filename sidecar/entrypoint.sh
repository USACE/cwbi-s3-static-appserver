#!/bin/bash

if [ -z "$S3_ENDPOINT_URL" ]
then
    CMD="aws s3 sync s3://${S3_BUCKET} /data/"
else
    CMD="aws s3 --endpoint-url ${S3_ENDPOINT_URL} sync s3://${S3_BUCKET} /data/"
fi

# Sleep for 20 seconds; Used in local development to wait for Minio to initialize
echo "Sleeping for 20 Seconds ..."
sleep 20

while true
do
    echo "Sync With Bucket ${S3_BUCKET}; $(date)"
    $CMD
    sleep 120
done
