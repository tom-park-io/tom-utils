# Graph Utility

This project is a Python utility for loading data and plotting graphs.

## Features

- Load data from CSV and JSON files.
- Generate plots from the loaded data.

## Installation

```bash
poetry install
```

## Usage

Run the main script:

```bash
sh scripts/run.sh
```

Or execute tests:

```bash
poetry run pytest
```

## Project Structure

- `src/graph/`: Main application code.
  - `data_loading/`: Modules for loading data.
  - `data_plotter/`: Modules for plotting data.
- `tests/`: Unit tests.
- `scripts/`: Helper scripts.
- `data/`: Sample data files (user needs to create, e.g., `data/json/test_result_1.json`).
- `plots/`: Directory where generated plots are saved.
