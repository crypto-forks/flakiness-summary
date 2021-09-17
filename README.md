# Flakiness summary docker action

This action runs the test suite multiple times, processes the results, and outputs a flakiness summary.

## Inputs

### `test-suite`

**Required** The name of the test suite to run. Valid options are `"unit"`, `"crypto-unit"`, `"integration-unit"`, and `"integration"`.

### `commit-sha`

**Required** The commit to run the test suite on.

## Example usage

uses: onflow/flakiness-summary@v1
with:
  test-suite: 'unit'
  commit-sha: ${{ github.sha }}