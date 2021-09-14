FROM golang

ENV NUM_RUNS=10

COPY flakiness-summary.sh /flakiness-summary.sh
COPY process_results.py /process_results.py

# Install git
RUN apt-get update     
RUN apt-get install -y git

ENTRYPOINT ["/flakiness-summary.sh"]