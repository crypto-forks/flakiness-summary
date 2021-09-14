import sys
import json

for line in sys.stdin:
    print(line)
    obj = json.loads(line)
    # TODO: process object

# TODO: write results to DB

print("Finished.")