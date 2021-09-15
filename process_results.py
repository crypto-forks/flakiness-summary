#!/usr/bin/env python3

import sys
import json
import os

print(os.environ['COMMIT_SHA'])
print(os.environ['COMMIT_TIME'])

for line in sys.stdin:
    print(line)
    obj = json.loads(line)
    # TODO: process object

# TODO: write results to DB

print('Finished.')