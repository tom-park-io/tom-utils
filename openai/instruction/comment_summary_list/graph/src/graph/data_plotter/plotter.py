import os
from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns
from matplotlib.ticker import MultipleLocator


def plot_data(
    df: pd.DataFrame, title: str = "Sample Plot", output_file: str | Path = "output_plot.png"
):
    """
    Plot data using Seaborn & Matplotlib.
    Assumes df.index for x-axis and 'value' column for y-axis.
    Creates output directory if it doesn't exist.
    """
    print("📊 Plotting data...")
    
    output_path = Path(output_file)
    output_dir = output_path.parent
    os.makedirs(output_dir, exist_ok=True)
    
    plt.figure(figsize=(10, 6))
    sns.lineplot(x=df.index, y="value", data=df, marker="o")
    plt.title(title)
    plt.xlabel("Date")
    plt.ylabel("Value")
    plt.tight_layout()
    plt.savefig(output_path)
    plt.show()
    print(f"✅ Plot saved as {output_path}")

def plot_multiple_results(
    data_frames: list[pd.DataFrame],
    file_names: list[str],
    x_col: str,
    y_col: str,
    title: str = "Comparison of Results",
    output_file: str | Path = "comparison_plot.png"
):
    """
    Plots data from multiple DataFrames on a single line graph.
    Each DataFrame represents a different file/dataset.
    Creates output directory if it doesn't exist.
    Args:
        data_frames: A list of pandas DataFrames, each containing x and y data.
        file_names: A list of names corresponding to each DataFrame, used for legend.
        x_col: The name of the column to use for the x-axis.
        y_col: The name of the column to use for the y-axis.
        title: The title of the plot.
        output_file: The filename for saving the plot (can be str or Path).
    """
    print(f"📊 Plotting multiple results...")
    
    output_path = Path(output_file)
    output_dir = output_path.parent
    os.makedirs(output_dir, exist_ok=True)
    
    fig, ax = plt.subplots(figsize=(12, 7))
    
    min_x, max_x = float('inf'), float('-inf')
    min_y, max_y = float('inf'), float('-inf')

    for i, df in enumerate(data_frames):
        if x_col not in df.columns:
            print(f"⚠️ Warning: X column '{x_col}' not found in data from '{file_names[i]}'. Skipping.")
            continue
        if y_col not in df.columns:
            print(f"⚠️ Warning: Y column '{y_col}' not found in data from '{file_names[i]}'. Skipping.")
            continue
        
        if not df[x_col].empty:
            min_x = min(min_x, df[x_col].min())
            max_x = max(max_x, df[x_col].max())
        if not df[y_col].empty:
            min_y = min(min_y, df[y_col].min())
            max_y = max(max_y, df[y_col].max())

        sns.lineplot(x=x_col, y=y_col, data=df, label=file_names[i], marker="o", ax=ax)
        
    ax.set_title(title)
    ax.set_xlabel(x_col.replace("_", " ").capitalize())
    ax.set_ylabel(y_col.replace("_", " ").capitalize())
    ax.legend(title="File")
    ax.grid(True)
    
    if min_x != float('inf') and max_x != float('-inf'):
        ax.set_xticks(np.arange(int(np.floor(min_x)), int(np.ceil(max_x)) + 1, 1))
        ax.xaxis.set_major_locator(MultipleLocator(1))

    if min_y != float('inf') and max_y != float('-inf'):
        ax.yaxis.set_major_locator(MultipleLocator(5))
        ax.yaxis.set_minor_locator(MultipleLocator(1))

    plt.tight_layout()
    plt.savefig(output_path)
    plt.show()
    print(f"✅ Plot saved as {output_path}")
