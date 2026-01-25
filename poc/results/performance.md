# PoC パフォーマンス結果

## M4 MacBook Pro での検証結果（2025-01-25）

### テスト環境

- **マシン**: MacBook Pro M4
- **OS**: macOS
- **Ollama**: 0.14.3
- **Model**: qwen2.5-coder:7b
- **Aider**: 0.86.1

### テスト結果

#### テストケース1: FizzBuzz（Go）

| 項目 | 結果 |
|:---|:---|
| 実行時間 | 数十秒（詳細未計測） |
| 成功/失敗 | ✅ 成功 |
| 生成コード品質 | 良好（構文エラーなし、正常動作） |
| 備考 | `mise exec -- go run main.go` で動作確認済み |

#### テストケース2: CSV処理（Go）

| 項目 | 結果 |
|:---|:---|
| 実行時間 | 数十秒（詳細未計測） |
| 成功/失敗 | ✅ 成功 |
| 生成コード品質 | 良好（name列でソート、正常動作） |
| 備考 | sample.csv を読み込み、名前順にソートして出力 |

#### テストケース3: バグ修正（Go）

| 項目 | 結果 |
|:---|:---|
| 実行時間 | 数十秒（詳細未計測） |
| 成功/失敗 | ✅ 成功 |
| バグ検出できたか | ✅ off-by-one エラーを検出 |
| 正しく修正されたか | ✅ `i := 1; i <= len` → `i := 0; i < len` に修正 |
| テスト通過 | ✅ `go test -v` 成功 |
| 備考 | Aiderが自動でバグを特定し修正 |

### リソース使用量

| 項目 | 値 |
|:---|:---|
| Ollama メモリ使用量 | 約5GB |
| Aider メモリ使用量 | 約500MB |
| 合計メモリ使用量 | 約6GB |
| CPU使用率（平均） | M4では軽負荷 |

### Go/No-Go 判断

#### 判断基準チェック

- [x] 推論速度: 60秒/応答以内（M4では十分高速）
- [x] メモリ使用量: 14GB以下（約6GB）
- [x] テストケース成功率: 3/3（100%）

#### 判断結果

**結果**: ✅ Go

**理由**:
- 3つのテストケースすべて成功
- M4での推論速度は実用レベル
- メモリ使用量も問題なし

**次のアクション**:
- MBP 2018（Intel）での検証を実施
- Phase 2（Worker開発）を並行して進行

---

## MBP 2018 での検証結果（進行中）

### テスト環境

- **マシン**: MacBook Pro 2018 (Intel Core i7, 16GB RAM)
- **OS**: macOS
- **Ollama**: ___
- **Model**: qwen2.5-coder:7b
- **Aider**: ___

### 判明した課題

1. **repo-map処理が遅い**: CPU推論が遅く、repo-mapの処理でタイムアウト
2. **環境変数**: `OLLAMA_API_BASE` の設定が必要な場合あり

### 推奨設定（Intel Mac）

```bash
# 環境変数
export OLLAMA_API_BASE=http://127.0.0.1:11434

# repo-mapを無効化して軽量化
~/.local/bin/aider --model ollama_chat/qwen2.5-coder:7b --map-tokens 0

# またはより軽量なモデルを使用
ollama pull qwen2.5-coder:3b
~/.local/bin/aider --model ollama_chat/qwen2.5-coder:3b
```

### テスト結果

（検証中）

---

## 備考

- M4 MacBook Pro: Apple Silicon の高速推論により快適に動作
- MBP 2018 (Intel): CPU推論のため低速、設定調整が必要
- 7Bモデルは Intel Mac では厳しい可能性、3Bモデル推奨
