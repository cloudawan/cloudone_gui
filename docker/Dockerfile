FROM    centos:centos6

# Bundle app source
COPY . /src



# Priviledge
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/cloudone_gui/cloudone_gui

EXPOSE  8443
CMD ["/src/run.sh"]