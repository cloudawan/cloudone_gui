#!/bin/bash

# Copy configuration
mkdir -p /etc/cloudone_gui
cp /src/cloudone_gui/conf/app.conf /etc/cloudone_gui/app.conf
# Not overwrite
cp -n /src/cloudone_gui/conf/development_cert.pem /etc/cloudone_gui/cert.pem
cp -n /src/cloudone_gui/conf/development_key.pem /etc/cloudone_gui/key.pem

# Use environment
sed -i "s/{{CLOUDONE_HOST}}/$CLOUDONE_HOST/g" /etc/cloudone_gui/app.conf
sed -i "s/{{CLOUDONE_PORT}}/$CLOUDONE_PORT/g" /etc/cloudone_gui/app.conf
sed -i "s/{{CLOUDONE_ANALYSIS_HOST}}/$CLOUDONE_ANALYSIS_HOST/g" /etc/cloudone_gui/app.conf
sed -i "s/{{CLOUDONE_ANALYSIS_PORT}}/$CLOUDONE_ANALYSIS_PORT/g" /etc/cloudone_gui/app.conf

cd /src/cloudone_gui
./cloudone_gui &

while :
do
	sleep 1
done


