FROM ubuntu 
COPY nrm . 
COPY config.yaml .
RUN  apt-get -y update && apt-get -y install apt-utils 
CMD ["./nrm"]