#!/bin/bash

# Use environment
sed -i "s/{{KUBEAPI_HOST}}/$KUBEAPI_HOST/g" /src/cloudone_gui/conf/app.conf
sed -i "s/{{KUBEAPI_PORT}}/$KUBEAPI_PORT/g" /src/cloudone_gui/conf/app.conf

# Notice the IP and port are not the one inside Kuberntes cluster but how outside world use it since it is for client side java script
sed -i "s/{{NODE_HOST}}/$NODE_HOST/g" /src/cloudone_gui/conf/app.conf
sed -i "s/{{NODE_PORT}}/$NODE_PORT/g" /src/cloudone_gui/conf/app.conf

sed -i "s/{{CLOUDONE_HOST}}/$CLOUDONE_HOST/g" /src/cloudone_gui/conf/app.conf
sed -i "s/{{CLOUDONE_PORT}}/$CLOUDONE_PORT/g" /src/cloudone_gui/conf/app.conf
sed -i "s/{{CLOUDONE_ANALYSIS_HOST}}/$CLOUDONE_ANALYSIS_HOST/g" /src/cloudone_gui/conf/app.conf
sed -i "s/{{CLOUDONE_ANALYSIS_PORT}}/$CLOUDONE_ANALYSIS_PORT/g" /src/cloudone_gui/conf/app.conf

cd /src/cloudone_gui
./cloudone_gui &

while :
do
	sleep 1
done


