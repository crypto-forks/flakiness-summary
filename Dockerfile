FROM golang:1.16-buster

RUN apt update 

# Install git   
RUN apt install -y git

# Install cmake
RUN apt install -y cmake

COPY flakiness-summary.sh /home/flakiness-summary.sh
COPY process_results.go /home/process_results.go

WORKDIR /home/

ENTRYPOINT ["/home/flakiness-summary.sh"]