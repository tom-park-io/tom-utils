#!/bin/bash

# dir
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

cd $SCRIPT_DIR/../

# choose new or init(manual)

# # auto
# poetry new myproject

# # manual (auto 로 생성되는 폴더 구조를 왠만해서 컨벤션으로 따름)
poetry init -n

# touch README.md
# poetry install
poetry install --no-root

poetry add --group dev pytest
poetry add pandas matplotlib seaborn
