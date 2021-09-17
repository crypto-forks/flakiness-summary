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

# TODO: output a flakiness summary
# see https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-output-parameter

print('Finished.')