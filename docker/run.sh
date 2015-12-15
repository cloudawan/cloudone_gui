#!/bin/bash

# Use environment
sed -i "s/{{KUBEAPI_CLUSTER_HOST_AND_PORT}}/$KUBEAPI_CLUSTER_HOST_AND_PORT/g" /src/cloudone_gui/conf/app.conf

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


