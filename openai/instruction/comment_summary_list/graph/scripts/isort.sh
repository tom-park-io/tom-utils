#!/bin/bash

# dir
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

cd $SCRIPT_DIR/../

poetry run isort .
