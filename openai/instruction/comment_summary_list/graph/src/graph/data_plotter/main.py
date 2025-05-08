import glob
import os
from pathlib import Path

import pandas as pd

from graph.data_loading.data_loader import load_raw_json
from graph.data_plotter.plotter import plot_multiple_results


def main():
    # 스크립트 파일의 위치를 기준으로 프로젝트 루트를 계산합니다.
    current_file_path = Path(__file__).resolve()
    project_root = current_file_path.parent.parent.parent.parent

    # JSON 파일들이 있는 디렉토리 경로
    json_dir_path = project_root / "data" / "json"
    plots_dir_path = project_root / "plots"  # 그래프 저장 디렉토리

    print(f"INFO: Searching for JSON files in: {json_dir_path}")

    # 디렉토리 내의 모든 .json 파일 경로 가져오기
    json_file_paths = glob.glob(str(json_dir_path / "*.json"))

    if not json_file_paths:
        print(f"WARNING: No JSON files found in {json_dir_path}. Exiting.")
        return

    all_data_frames = []
    all_file_names = []

    for file_path in sorted(json_file_paths):  # 파일 순서를 일관되게 하기 위해 정렬
        print(f"INFO: Processing file: {file_path}")
        try:
            # raw JSON 데이터 로드
            data = load_raw_json(Path(file_path))

            # 'articles_list' 추출
            articles_list = data.get("articles_list")
            if articles_list is None or not isinstance(articles_list, list):
                print(
                    f"WARNING: 'articles_list' not found or not a list in {file_path}. Skipping."
                )
                continue

            # result 점수와 인덱스 추출
            results = []
            for i, article in enumerate(articles_list):
                if isinstance(article, dict) and "result" in article:
                    results.append(
                        {"article_index": i, "result_score": article["result"]}
                    )
                else:
                    print(
                        f"WARNING: Item at index {i} in 'articles_list' from {file_path} is not a dict or missing 'result'. Skipping item."
                    )

            if not results:
                print(
                    f"WARNING: No valid result data extracted from {file_path}. Skipping file."
                )
                continue

            # DataFrame 생성
            df = pd.DataFrame(results)
            all_data_frames.append(df)
            all_file_names.append(
                os.path.basename(file_path)
            )  # 파일 이름만 추출 (예: test_result_1.json)

        except Exception as e:
            print(f"ERROR: Failed to process file {file_path}: {e}")

    if not all_data_frames:
        print("INFO: No data processed to plot. Exiting.")
        return

    # 그래프 출력 파일 경로 설정 (plots/ 디렉토리 아래)
    output_plot_file = plots_dir_path / "article_scores_comparison.png"

    # 모든 결과를 하나의 그래프로 플롯
    plot_multiple_results(
        data_frames=all_data_frames,
        file_names=all_file_names,
        x_col="article_index",
        y_col="result_score",
        title="Comparison of Article Result Scores by File",
        output_file=output_plot_file,  # 수정된 출력 파일 경로 사용
    )


if __name__ == "__main__":
    main()
