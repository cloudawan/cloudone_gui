FROM ubuntu:14.04

# Bundle app source
COPY . /src

# Privilege
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/cloudone_gui/cloudone_gui

EXPOSE  8443
CMD ["/src/run.sh"]