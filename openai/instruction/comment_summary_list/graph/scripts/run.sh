#!/bin/bash
# This script runs the main data plotting application and the unit tests.

# dir
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

cd $SCRIPT_DIR/../

poetry run python -m graph.data_plotter.main
