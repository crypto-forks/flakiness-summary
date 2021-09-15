FROM golang:1.16-buster

COPY flakiness-summary.sh /flakiness-summary.sh
COPY process_results.py /process_results.py

# Install git
RUN apt-get update     
RUN apt-get install -y git

# Install cmake
RUN apt install -y cmake

ENTRYPOINT ["/flakiness-summary.sh"]