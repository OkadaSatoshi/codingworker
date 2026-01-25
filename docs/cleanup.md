# クリーンアップ手順

検証環境からOllama、Aiderを削除する手順。

## 1. Ollamaの削除

```bash
# サービス停止
brew services stop ollama

# ダウンロードしたモデルの削除
ollama rm qwen2.5-coder:1.5b

# Ollamaアンインストール
brew uninstall ollama

# Ollamaデータディレクトリの削除（モデル、設定等）
rm -rf ~/.ollama
```

## 2. Aiderの削除

```bash
# Aiderアンインストール
pip3 uninstall aider-chat -y

# キャッシュの削除（オプション）
rm -rf ~/.aider
```

## 3. 削除確認

```bash
# コマンドが見つからなければOK
which ollama
which aider

# ディレクトリが存在しなければOK
ls ~/.ollama
ls ~/.aider
```

## 4. その他（オプション）

### Homebrewキャッシュのクリア

```bash
brew cleanup
```

### pipキャッシュのクリア

```bash
pip3 cache purge
```
