FROM golang:1.16-buster

RUN apt update 

# Install git   
RUN apt install -y git

# Install cmake
RUN apt install -y cmake

# Install python
RUN apt install -y python3

COPY flakiness-summary.sh /home/flakiness-summary.sh
COPY process_results.py /home/process_results.py

WORKDIR /home/

ENTRYPOINT ["/home/flakiness-summary.sh"]