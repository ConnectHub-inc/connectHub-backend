# connectHub-backend

## Development
### Format
goimports, gofmt
```makefile
make fmt
```
### Generate
go generate
```makefile
make generate
```
### Build Go
```makefile
make build
```
### Lint
```makefile
make lint
```
これは ./docker/back/... を対象として golangci-lint を実行するため、実用的ではありません。 実際には、次のように PKG 変数を指定し、package を限定した状態で実行することをお勧めします。
```makefile
make lint PKG=./docker/back/infra/...
```
後述の make lint-diff では差分のみを対象とするため、既存のコードには利用できません。 既存のコードの品質改善を行いたい場合には make lint が有用です。

### Lint (diff)
```makefile
make lint-diff
```
PKG 指定もできます。
```
make lint-diff PKG="./docker/back/infra/..."
```
make lint だと既存の指摘が多く、追加したコードに対する解析結果が判別しにくいため、 develop branch との差分の解析結果を表示する lint-diff を用意しています。 利用するには事前に reviewdog のインストールが必要です。
```shell
curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```
### Test
```makefile
make test
```

```makefile
make e2e-test
```
