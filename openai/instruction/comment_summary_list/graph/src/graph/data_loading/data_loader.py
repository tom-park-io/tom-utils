import json
from pathlib import Path

import pandas as pd


def load_csv_data(filepath: str | Path) -> pd.DataFrame:
    """
    Load CSV data into a pandas DataFrame.
    Assumes 'date' column for parsing and setting as index.
    Can accept a string path or a pathlib.Path object.
    """
    if not isinstance(filepath, (str, Path)):
        raise TypeError(f"filepath must be a string or Path, not {type(filepath)}")
    print(f"📂 Loading CSV data from {str(filepath)}")
    try:
        df = pd.read_csv(filepath, parse_dates=["date"])
        df.set_index("date", inplace=True)
        print(f"✅ Loaded {len(df)} rows from CSV.")
        return df
    except Exception as e:
        print(f"❌ Error loading CSV data from {str(filepath)}: {e}")
        raise


def load_json_data(
    filepath: str | Path,
    orient: str = "records",
    lines: bool = False,
    date_columns: list[str] | None = None,
    index_col: str | None = None
) -> pd.DataFrame:
    """
    Load JSON data into a pandas DataFrame using pandas.read_json.
    Can accept a string path or a pathlib.Path object.
    Args:
        filepath: Path to the JSON file.
        orient: Pandas read_json orient parameter.
        lines: Pandas read_json lines parameter.
        date_columns: A list of column names to parse as dates.
        index_col: Column to set as index. If a date_column, it will be a DateTimeIndex.
    """
    if not isinstance(filepath, (str, Path)):
        raise TypeError(f"filepath must be a string or Path, not {type(filepath)}")
    print(f"📂 Loading JSON data into DataFrame from {str(filepath)}")
    try:
        df = pd.read_json(filepath, orient=orient, lines=lines)
        if date_columns:
            for col in date_columns:
                if col in df.columns:
                    try:
                        df[col] = pd.to_datetime(df[col])
                        print(f"ℹ️  Parsed column '{col}' as datetime.")
                    except Exception as e:
                        print(f"⚠️ Could not convert column '{col}' to datetime: {e}")
                else:
                    print(f"⚠️ Specified date column '{col}' not found in DataFrame.")
        
        if index_col and index_col in df.columns:
            try:
                df.set_index(index_col, inplace=True)
                print(f"ℹ️  Set column '{index_col}' as index.")
            except Exception as e:
                print(f"⚠️ Could not set column '{index_col}' as index: {e}")
        elif index_col:
            print(f"⚠️ Specified index column '{index_col}' not found in DataFrame.")
            
        print(f"✅ Loaded {len(df)} rows from JSON into DataFrame.")
        return df
    except Exception as e:
        print(f"❌ Error loading JSON data into DataFrame from {str(filepath)}: {e}")
        raise


def load_raw_json(filepath: str | Path) -> dict | list:
    """
    Load JSON data from a file into a Python dictionary or list.
    Can accept a string path or a pathlib.Path object.
    """
    if not isinstance(filepath, (str, Path)):
        raise TypeError(f"filepath must be a string or Path, not {type(filepath)}")
    print(f"📂 Loading raw JSON data from {str(filepath)}")
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            data = json.load(f)
        print(f"✅ Successfully loaded raw JSON data.")
        return data
    except FileNotFoundError:
        print(f"❌ Error: File not found at {str(filepath)}")
        raise
    except json.JSONDecodeError as e:
        print(f"❌ Error decoding JSON from {str(filepath)}: {e}")
        raise
    except Exception as e:
        print(f"❌ An unexpected error occurred while loading raw JSON from {str(filepath)}: {e}")
        raise 