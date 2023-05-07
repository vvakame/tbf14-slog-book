# slog解説

Goの標準ライブラリになる予定のslog^[https://pkg.go.dev/golang.org/x/exp/slog]について解説します。
APIの概要と使い方、落とし穴などを順に見ていきます。

## slogの特徴

slogは構造化ロギング（structured logging）を行うためのライブラリです。
s+logという名前がそのことを表していますね。

通常、ログというとテキストメッセージを想像すると思います。
対して、構造化ロギングではテキストメッセージ+Key-Valueペアで付帯データを表現することができます。
構造化ログのデータはたとえばJSONなどにエンコードされて出力されます。

構造化ログには、たとえば指定されたオプションの値を記録する用途に使えます。
同じログを出力するのにも、いくつかやり方があるので、ここでは3パターンを掲載しておきます。
`context.Context` を引数に取らない関数もありますが、宗教上の理由により説明しません。

<!-- maprange:../code/structured_logging_test.go,example1 -->
```go title=属性を明示的に組み立てる
slog.InfoCtx(
  ctx, "start processing",
  slog.Bool("verbose", verbose),
)
```
<!-- maprange.end -->

<!-- maprange:../code/structured_logging_test.go,example2 -->
```go title=属性をKey-Valueの2つの値をペアとして組み立てる
slog.InfoCtx(
  ctx, "start processing",
  "verbose", verbose,
)
```
<!-- maprange.end -->

<!-- maprange:../code/structured_logging_test.go,example3 -->
```go title=レベルを明示し属性も明示的に組み立てる
slog.LogAttrs(
  ctx,
  slog.LevelInfo,
  "start processing",
  slog.Bool("verbose", verbose),
)
```
<!-- maprange.end -->

筆者の好みは一番最初の例なので、本書では基本的にこの書き方で書いていきます。
世間的には2番目の例が好まれる気がします。

上記コードの出力例は次の通りです。
後述するHandlerによって出力結果のフォーマットはいくらでも変化します。

```text title=出力結果。実際は1行のJSON（LDJSON）として表示される
{
  "time":"2023-05-07T17:09:29.781259+09:00",
  "level":"INFO",
  "msg":"start processing",
  "verbose":true
}
```

構造化されたログの活用を考えたとき、プログラムコード中で属性を明示的に追加するだけではなく、コンテキストに依存した暗黙的な属性も含めることができるのもメリットです。
たとえばそのログが出力されたソースコード上の位置や、どのリクエストに紐づくログなのかというトレース情報や、現在の操作者が誰かといった情報などです。
先の例では、ログを出力した時刻とログレベルがデフォルトで出力されています。

構造化データは機械可読性（machine-readable）が高いのがメリットです。
たとえばBigQuery^[https://cloud.google.com/bigquery/]などのデータレイクに蓄積されたログは、個別の属性に対して検索が容易になるためログをデータとして活用しやすくなります。

## slogのAPI

slogのAPIの概要と使い方について簡単に見ていきます。
詳細はドキュメント^[https://pkg.go.dev/golang.org/x/exp/slog]を参照してください。

slogのコア部分はHandler^[https://pkg.go.dev/golang.org/x/exp/slog#Handler] interfaceとLogger^[https://pkg.go.dev/golang.org/x/exp/slog#Logger] struct、この2つによって成り立っています。
Handlerは自分で自由に実装ができて、渡されたRecord（後述）をどこかしらに出力する役目を持ちます。
LoggerはそんなHandlerにメッセージやデータを渡すため、プログラムコードからの値をRecord（後述）に加工します。
Loggerにはカスタマイズの余地はありません。

Handlerの仕様を次に示します。
自分でカスタマイズしたHandlerを作りたい場合、次の4つのメソッドを満たしてやればよいことになります。

```go
type Handler interface {
  // 指定されたレベルのログを出力するかしないか制御
  Enabled(context.Context, Level) bool
  // Recordをどこかに出力する処理を担当
  Handle(context.Context, Record) error
  // デフォルトで付属させる属性が指定されたHandlerを返す
  WithAttrs(attrs []Attr) Handler
  // 指定された名前のグループに属性を出力するHandlerを返す
  WithGroup(name string) Handler
}
```

Handlerを元にLoggerを作ります。
また、作ったLoggerが `slog.DebugCtx` etcで使われるよう、デフォルトのLoggerに指定します。

<!-- maprange:../code/structured_logging_test.go,newLogger -->
```go title=Loggerの作成とデフォルト指定
logger := slog.New(h)
slog.SetDefault(logger)
```
<!-- maprange.end -->

似た名前の関数として `slog.NewLogLogger` ^[https://pkg.go.dev/golang.org/x/exp/slog#NewLogLogger]も存在します。
こちらはLoggerを作成する関数ではなくて、 `log` パッケージのLoggerを作る、似て非なるものなので間違えないようにしましょう^[筆者は片手では収まらないくらいの回数間違えました]。

HandlerにはデフォルトでTextHandler^[https://pkg.go.dev/golang.org/x/exp/slog#TextHandler]とJSONHandler^[https://pkg.go.dev/golang.org/x/exp/slog#JSONHandler]の2つが準備されています。
`slog.NewTextHandler` や `slog.NewJSONHandler` でそれぞれのHandlerを作ることができます。

<!-- maprange:../code/default_handler_test.go,textHandler -->
```go title=TextHandlerを作成する
h = slog.NewTextHandler(os.Stdout)
logger = slog.New(h)
logger.InfoCtx(
  ctx, "start processing",
  slog.Bool("verbose", true),
)
```
<!-- maprange.end -->

また、 `slog.HandlerOptions` を通じて多少のカスタマイズを施すこともできます。
属性のKey, Valueそれぞれを書き換えたり、カスタマイズの自由度はかなりあります。

<!-- maprange:../code/default_handler_test.go,textHandlerWithHandlerOptions -->
```go title=TextHandlerを作成する
h = slog.HandlerOptions{
  // 呼び出し元コードの出力
  AddSource: true,
  // 出力するログレベル
  Level: slog.LevelDebug,
  // 属性の置き換え・削除など
  ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
    a.Key += "!"
    a.Value = slog.StringValue(a.Value.String() + "?")
    return a
  },
}.NewTextHandler(os.Stdout)
logger = slog.New(h)
logger.DebugCtx(
  ctx, "start processing",
  slog.Bool("verbose", true),
)
```
<!-- maprange.end -->

このコードの実行結果はこちら。
ReplaceAttrでKeyの末尾に `!` 、Valueの末尾に `?` を追加するようにしたので、実際にそのような出力になっています。

```text title=TextHandler+ReplaceAttrの出力サンプル。実際は1行で出力される
time!="2023-05-07 18:28:49.897963 +0900 JST?"
level!=DEBUG?
source!=(omitted)/default_handler_test.go:53?
msg!="start processing?"
verbose!=true?
```

自分でHandlerを実装する場合でも、通常は標準のJSONHandlerを内部で持ち、渡す値を事前に加工する実装になる場合が大半でしょう。

続いて、 `slog.Attr` について説明します。
ログのうち、属性を表す構造がこれです。
コード中では `slog.InfoCtx(ctx, "msg", "b", true)` のように実行されますが、Handlerにたどり着く時点では `slog.Attr{Key: "b", Value: slog.BoolValue(true)}` というような構造に変換されています。
先に見せた例のように、 `slog.InfoCtx(ctx, "msg", slog.Bool("b", true))` が等価な書き方で、 `slog.Bool` は `slog.Attr` の組み立て用の関数というわけです。

Handlerが直接処理するのは `slog.Record` です。
次に示す構造と、中に持つ `slog.Attr` を処理できる `func (r Record) Attrs(f func(Attr) bool)` を持っています。
Handlerはこれらの情報を駆使してログを出力します。

```go
type Record struct {
  // 出力が行われようとした時間
  Time time.Time
  // メッセージ
  Message string
  // ログレベル
  Level Level
  // Recordが組み立てられた時の呼び出し元のプログラムカウンタ
  // runtime.CallersFrames でフレームの情報が得られる
  PC uintptr
}
```

注意点として、`Handler.WithAttrs` で渡された `slog.Attr` は処理済であるとものとして、 `slog.Record` には含まれません。
Handlerがこれらをどう処理するべきかは、TextHandlerやJSONHandlerの内部実装である `slog.commonHandler` の実装を読んでみると面白いでしょう。
概要だけ述べると、 `Handler.WithAttrs` で渡された `slog.Attr` は早期に処理され内部的に `[]byte` で保持され、Handleが呼ばれるたびに処理しないように最適化されています。

次に解説するのは `slog.Value` です。
この値はstructなため直接組み立てることはできません。
`slog.AnyValue`、`slog.BoolValue`、`slog.DurationValue`、`slog.GroupValue` などなどさまざまな関数が用意されています。

面白い要素として、 `slog.LogValuer` があります。
これは渡された値の評価を実際にログ出力する直前（ `Handler.Handle` などの来る直前 ）まで遅延させることができます。
このインタフェースの定義は次の通りです。

```go
type LogValuer interface {
  LogValue() Value
}
```

公式のgodocでは、パスワードなどのセンシティブな値をログ出力する時にマスクさせる使い方や、シンプルに計算コストが重い値を計算するときに使う例がありました。
この値はLoggerが実際の `slog.Record` を生成する時に `LogValuer` な値ではなくなるまで再帰的に処理されます。
このため、Handlerの実装では `slog.Value` の値が `LogValuer` である場合は検討しなくてよくなっています。


TODO
Groupの話


## 実用の仕方

TODO
vet



GCPの話をします。
Cloud Runなどの利用時、ログ出力が一定のルール^[https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields]に従うとCloud Loggingがその形式に従って構造化ログとして解釈してくれます。

筆者は自作のgcpslogライブラリ^[https://pkg.go.dev/github.com/vvakame/sdlog/gcpslog]を利用してログ出力を行っています。
技術書典のWebサイトで実際に利用していて、元気に構造化ログを出力しています。
gcpslogで提供しているHandlerの実装は内部的には `slog.JSONHandler` を持っていて、 `Handler.Handle` でいくつかの追加の属性をRecordに追加し、JSONHandlerに流す、という処理を行っています。





## slogに不向きなこと

slogに不向きな用途についても解説しておきます。
構造化ロギングが行えるといえど、与えられる属性の自由度を制限することはできません。
いわば、スキーマレスなデータを出力することになります。

一貫したデータ構造がある場合は、しっかりとログ出力を設計し、安易にslog経由でログ出力を行うのは控えるべきでしょう。
たとえばGraphQLのリクエストログや、特定のデータの更新履歴など、特定のスキーマが存在する場合、特定の出力先にスキーマが一貫して存在している形で出力するべきでしょう。

スキーマを固定したログ出力にslogを使うべきか、別のものを使うべきかは設計次第ですが、特別なものが特別であるとわかるよう、slogを使わない選択をするのがわかりやすいのではないでしょうか。

## 落とし穴の紹介

TODO
HandleでCurrentUser読んだら字面に反して色々やってる処理だったのでメモリ使用量が一気にエグいことになってAPIがエラーを返す様になった話。
NewLogLogger 間違えがち

## コミュニティのライブラリ

TODO
https://pkg.go.dev/github.com/samber/slog-multi
https://pkg.go.dev/github.com/neilotoole/slogt


## 雑感

TODO

## 物置

https://go.googlesource.com/go/+/refs/heads/master/src/log/slog/
https://go.googlesource.com/go/+/refs/heads/master/src/log/slog/internal/slogtest/
https://go-review.googlesource.com/c/exp/+/469075

> Handler authors and others may wish to use Value.Resolve
instead of calling LogValue directly.

やっぱFromContextいるくない… *testing.T に紐付けようとするとdefault logger変えるのもあれだし
