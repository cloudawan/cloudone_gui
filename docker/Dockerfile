FROM    centos:centos6

# Bundle app source
COPY . /src



# Priviledge
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/kubernetes_management_gui/kubernetes_management_gui

EXPOSE  8443
CMD ["/src/run.sh"]