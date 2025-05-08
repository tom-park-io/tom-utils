import pandas as pd
import pytest

from graph.data_loading.data_loader import load_csv_data, load_json_data


# CSV 테스트 케이스
def test_load_csv_data(tmp_path):
    # 임시 CSV 파일 생성
    csv_content = "date,value\n2023-01-01,10\n2023-01-02,20"
    csv_file = tmp_path / "sample.csv"
    csv_file.write_text(csv_content)

    df = load_csv_data(str(csv_file))

    assert isinstance(df, pd.DataFrame)
    assert not df.empty
    assert len(df) == 2
    assert "value" in df.columns
    assert pd.api.types.is_datetime64_any_dtype(df.index)
    assert df.index.name == "date"
    assert df.loc[pd.Timestamp("2023-01-01"), "value"] == 10


def test_load_csv_data_file_not_found(tmp_path):
    with pytest.raises(FileNotFoundError):
        load_csv_data(str(tmp_path / "non_existent.csv"))


# JSON Lines (JSONL) 테스트 케이스
def test_load_json_data_lines(tmp_path):
    # 임시 JSONL 파일 생성
    jsonl_content = (
        '{"date": "2023-01-01", "value": 30, "category": "A"}\n'
        '{"date": "2023-01-02", "value": 40, "category": "B"}'
    )
    jsonl_file = tmp_path / "sample.jsonl"
    jsonl_file.write_text(jsonl_content)

    df = load_json_data(str(jsonl_file), lines=True)

    assert isinstance(df, pd.DataFrame)
    assert not df.empty
    assert len(df) == 2
    assert "value" in df.columns
    assert "category" in df.columns
    assert pd.api.types.is_datetime64_any_dtype(df.index)
    assert df.index.name == "date"
    assert df.loc[pd.Timestamp("2023-01-02"), "value"] == 40


# JSON Records 테스트 케이스
def test_load_json_data_records(tmp_path):
    # 임시 JSON 파일 (레코드 배열) 생성
    json_records_content = (
        '[{"timestamp": "2023-03-10T10:00:00", "measurement": 1.5},'
        ' {"timestamp": "2023-03-10T10:05:00", "measurement": 1.8}]'
    )
    json_file = tmp_path / "sample.json"
    json_file.write_text(json_records_content)

    df = load_json_data(str(json_file), orient="records")

    assert isinstance(df, pd.DataFrame)
    assert not df.empty
    assert len(df) == 2
    assert "measurement" in df.columns
    assert pd.api.types.is_datetime64_any_dtype(df.index)
    assert df.index.name == "timestamp"
    assert df.loc[pd.Timestamp("2023-03-10T10:05:00"), "measurement"] == 1.8


def test_load_json_data_no_date_column(tmp_path):
    json_content = '[{"value": 100}, {"value": 200}]'
    json_file = tmp_path / "no_date.json"
    json_file.write_text(json_content)

    df = load_json_data(str(json_file))
    assert df.index.name is None  # 날짜 열이 없으므로 기본 인덱스


def test_load_json_data_file_not_found(tmp_path):
    with pytest.raises(FileNotFoundError):
        load_json_data(str(tmp_path / "non_existent.json"))
