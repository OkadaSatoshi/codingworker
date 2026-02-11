# クリーンアップ手順

検証環境からOllama、Aiderを削除する手順。

## 1. Ollamaの削除

```zsh
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

```zsh
# Aiderアンインストール
uv tool uninstall aider-chat

# キャッシュの削除（オプション）
rm -rf ~/.aider
```

## 3. 削除確認

```zsh
# コマンドが見つからなければOK
which ollama
which aider

# ディレクトリが存在しなければOK
ls ~/.ollama
ls ~/.aider
```

## 4. その他（オプション）

### Homebrewキャッシュのクリア

```zsh
brew cleanup
```

### uvキャッシュのクリア

```zsh
uv cache clean
```
