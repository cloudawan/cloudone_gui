#!/bin/bash

# Use environment
sed -i "s/{{KUBEAPI_HOST}}/$KUBEAPI_HOST/g" /src/kubernetes_management_gui/conf/app.conf
sed -i "s/{{KUBEAPI_PORT}}/$KUBEAPI_PORT/g" /src/kubernetes_management_gui/conf/app.conf

# Notice the IP and port are not the one inside Kuberntes cluster but how outside world use it since it is for client side java script
sed -i "s/{{NODE_HOST}}/$NODE_HOST/g" /src/kubernetes_management_gui/conf/app.conf
sed -i "s/{{NODE_PORT}}/$NODE_PORT/g" /src/kubernetes_management_gui/conf/app.conf


cd /src/kubernetes_management_gui
./kubernetes_management_gui &

while :
do
	sleep 1
done


