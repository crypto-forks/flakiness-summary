# Flakiness summary docker action

This action runs the test suite multiple times, processes the results, and outputs a flakiness summary.

## Inputs

### `test-category`

**Required** The category of tests to run. Valid options are `"unit"`, `"crypto-unit"`, `"integration-unit"`, and `"integration"`.

### `commit-sha`

**Required** The commit to run the tests on.

## Example usage

uses: onflow/flakiness-summary@v1
with:
  test-category: 'unit'
  commit-sha: ${{ github.sha }}