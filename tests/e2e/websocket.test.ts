import WebSocket from 'ws';
import axios from 'axios';

jest.setTimeout(1000000); // タイムアウトを延長

describe('WebSocket E2E Tests with Go Server', () => {
  let ws: WebSocket; // WebSocketクライアントのインスタンスを保持する変数
  let authToken: string;

  // 全てのテストの前に実行されるセットアップ処理
  beforeAll(async () => {
    try {
      // 認証リクエストを送信してトークンを取得
      const uniqueEmail = `e2e_test_${Date.now()}@test.com`;
      const response = await axios.post('http://localhost:8083/api/user/create', {
        email: uniqueEmail,
        password: 'e2e_test_password'
      });
      authToken = response.headers['authorization'].split(' ')[1];

      // 取得したトークンを使ってWebSocket接続を確立
      const headers = {
        'Authorization': `Bearer ${authToken}`
      };
      ws = new WebSocket('ws://localhost:8083/ws', { headers });

      await new Promise<void>((resolve, reject) => {
        ws.on('open', () => {
          console.log('WebSocket connection established');
          resolve(); // WebSocket接続が確立されたらセットアップ完了
        });

        ws.on('error', (error) => {
          console.error('WebSocket connection error:', error);
          reject(error); // 接続エラーが発生した場合はテスト失敗
        });
      });
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Setup failed: ${error.message}`); // 認証リクエストに失敗した場合はテスト失敗
      } else {
        throw new Error('Setup failed: An unknown error occurred'); // 不明なエラーの場合
      }
    }
  });

  // 全てのテストの後に実行されるクリーンアップ処理
  afterAll(() => {
    ws.close(); // WebSocket接続を閉じる
  });

  // WebSocket接続が正しく確立されるかをテスト
  test('should establish a WebSocket connection', (done) => {
    // WebSocket接続が確立された時に発生するイベント
    ws.on('open', () => {
      console.log('WebSocket is open');
      expect(ws.readyState).toBe(WebSocket.OPEN); // WebSocketがOPEN状態であることを確認
      done(); // テスト完了
    });

    // すでにオープンしている場合の処理
    if (ws.readyState === WebSocket.OPEN) {
      console.log('WebSocket was already open');
      expect(ws.readyState).toBe(WebSocket.OPEN);
      done(); // テスト完了
    }

    ws.on('error', (error) => {
      console.error('Test WebSocket connection error:', error);
      done(error); // エラー発生時はテスト失敗
    });
  });
});
