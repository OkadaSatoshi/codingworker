# Phase 0: PoC テスト

## 概要

Aider + Ollama (qwen2.5-coder:7b) の実用性を検証するためのテストケース。

## テストケース

| # | ディレクトリ | タスク | 言語 |
|:--|:---|:---|:---|
| 1 | test-case-1-fizzbuzz/ | FizzBuzz スクリプト生成 | Go |
| 2 | test-case-2-csv/ | CSV ソート処理生成 | Go |
| 3 | test-case-3-bugfix/ | バグ修正 | Go |

## 実行手順

### 前提条件

- Ollama インストール済み
- qwen2.5-coder:7b ダウンロード済み
- Aider インストール済み

### テストケース1: FizzBuzz

```bash
cd test-case-1-fizzbuzz
git init
aider --model ollama_chat/qwen2.5-coder:7b

# Aider内で以下を入力:
# Create a Go program that prints FizzBuzz from 1 to 100
```

**成功基準**:
- コードが生成される
- `go run main.go` で正しく動作する

### テストケース2: CSV処理

```bash
cd test-case-2-csv
git init
aider --model ollama_chat/qwen2.5-coder:7b

# Aider内で以下を入力:
# Create a Go program that reads sample.csv and sorts it by the first column, then outputs to stdout
```

**成功基準**:
- コードが生成される
- sample.csv を正しくソートして出力する

### テストケース3: バグ修正

```bash
cd test-case-3-bugfix
git init
aider --model ollama_chat/qwen2.5-coder:7b buggy.go

# Aider内で以下を入力:
# Fix the bug in buggy.go
```

**成功基準**:
- バグが検出される
- 正しく修正される
- テストが通る

## 計測項目

| 項目 | 目標値 |
|:---|:---|
| 推論速度 | 60秒/応答以内 |
| メモリ使用量 | 14GB以下 |
| 成功率 | 2/3以上 |

## 結果記録

テスト結果は `results/performance.md` に記録する。
