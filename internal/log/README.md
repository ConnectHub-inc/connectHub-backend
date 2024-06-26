# ログレベル定義

| ログレベル | 定義 | 本番環境でのロギング |
| ----------------- | -------- | -------- |
| DEBUG 　　　   | デバッグ情報。主に開発中に使用される詳細な情報。 | 通常は無効。トラブルシューティングが必要な場合にのみ有効。 |
| INFO　　　 | 一般的な情報メッセージ。アプリケーションの正常な動作を示す。 | 必要に応じて有効。通常は少量の情報として記録される。 |
| NOTICE　　　 | 重要だが、緊急ではない情報。注意を喚起するためのメッセージ。 | 一般的に有効。重要なイベントや状況を記録。 |
| WARNING　　　 | 潜在的な問題を示す警告メッセージ。 | 通常有効。潜在的な問題の早期発見に役立つ。 |
| ERROR　　　 | エラーメッセージ。アプリケーションの動作に支障をきたす可能性のある問題。頻発しなければ未対応でもいい | 必ず有効。エラーの発生とその原因を特定するために重要。 |
| CRITICAL 　　　| 致命的なエラー。アプリケーションの継続が困難な重大な問題。基本的には1回発生したら対応や調査が必要 | 必ず有効。緊急対応が必要な重大な問題を記録。 |
