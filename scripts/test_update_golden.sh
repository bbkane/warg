#!/bin/bash

# This is mainly so I can whitelist the VS Code agent to run this

WARG_TEST_UPDATE_GOLDEN=1 go test ./...
