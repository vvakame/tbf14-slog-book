# slog解説

Goの標準ライブラリになる予定のslog^[https://pkg.go.dev/golang.org/x/exp/slog]について解説します。
slogはGo 1.21で `log/slog` として標準ライブラリの一員になる予定です。

本章ではAPIの概要と使い方、落とし穴などを順に見ていきます。

使っている `golang.org/x/exp` のバージョンは次の通りです。

<!-- mapfile:./exp.version.txt -->
```
v0.0.0-20230510235704-dd950f8aeaea
```
<!-- mapfile.end -->

## slogの特徴と構造

slogは構造化ロギング（structured logging）を行うためのライブラリです。
s+logという名前がそのことを表していますね。

通常、ログというとテキストメッセージを想像すると思います。
対して、構造化ロギングではテキストメッセージ+Key-Valueペアで付帯データを表現することができます。
構造化ログのデータはたとえばJSONなどにエンコードされて出力されます。
たとえば、指定されたオプションの値をメッセージと共に記録する用途などに使えます。

まずは、slogの概要を簡単に説明しておおまかな作りを把握した後、それぞれの詳しい説明をしていきます。

大まかな構造として、 `Logger`、 `Handler`、 `Attr`、 `Record` が重要な構造です。
また、 `Level`、 `Value`、他にも細々としたものがありますが、だいたい名前どおりのものが多いですのでここでは割愛します。

重要な4要素がプログラムに対してどのように関わってくるかを説明します。

1. ログをどこにどのように出力するか決定する `Handler` を作成
1. `Handler` から `Logger` を作成
1. `Logger` を使ってログ出力を行う。メッセージ+ `Attr` を引数にわたす
1. `Logger` は各種状態+引数を組み合わせてメッセージ+ `Attr` の集合体である `Record` を組み立て、 `Handler` に渡す
1. `Handler` がどこかになにかを出力してログ出力完了

上記の流れからわかるように、slogのカスタマイズのキモは `Handler` になります。
`Handler` をカスタマイズすることで、どういったフォーマットにするか、出力先はどこなのかを制御できます。

`Logger` はstructなので、カスタマイズできません。
メッセージや `Attr` を `Record` に組み立てる役割を果たします。

ログは出力のレベルとメッセージと属性[^attr-translate]（`Attr`, attributes）によって構成されます。
ここでいう属性はKey-Valueのペアで、メッセージとは独立しているので、機械可読な構造を持つ、ということになります。

[^attr-translate]: attributesを属性と訳すのが正しいかちょっと迷いましたが、この本では属性と訳すことにします

`Record` はログ作成時刻、ログ出力レベル、メッセージ、PC（プログラムカウンタ、ログ出力したソースコード上の位置を特定するのに使う）、そして複数の `Attr` を持つデータ構造です。
`Logger` でログ出力が行われた時、いくらかの前処理などを経て `Record` が組み立てられ、 `Handler` に最終的な出力処理が委ねられます。

厳密にいうと、 `With` や `WithGroup` が事前に呼ばれている場合、いくつか `Handler` にて `Record` の処理よりさらに前に処理済の要素がある可能性はあります。

## API解説

slogのAPIの概要と使い方について簡単に見ていきます。
詳細はドキュメント^[https://pkg.go.dev/golang.org/x/exp/slog]を参照してください。

slogのコア部分はHandler^[https://pkg.go.dev/golang.org/x/exp/slog#Handler] interfaceとLogger^[https://pkg.go.dev/golang.org/x/exp/slog#Logger] struct、この2つによって成り立っています。
Handlerは自分で自由に実装ができて、渡されたRecordをどこかしらに出力する役目を持ちます。
LoggerはそんなHandlerにメッセージやデータを渡すため、プログラムコードからの値をRecordに加工します。
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
また、作ったLoggerがパッケージの関数である `slog.DebugCtx` などで使われるよう、デフォルトのLoggerに指定します。

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
h = slog.NewTextHandler(os.Stdout, nil)
logger = slog.New(h)
logger.InfoCtx(
  ctx, "start processing",
  slog.Bool("verbose", true),
)
```
<!-- maprange.end -->

また、 `slog.HandlerOptions` を通じて多少のカスタマイズを施すこともできます。
属性のKey, Valueそれぞれを書き換えたり、カスタマイズの自由度があります。

<!-- maprange:../code/default_handler_test.go,textHandlerWithHandlerOptions -->
```go title=TextHandlerを作成する
h = slog.NewTextHandler(
  os.Stdout,
  &slog.HandlerOptions{
    // 呼び出し元コードの出力
    AddSource: true,
    // 出力するログレベル
    Level: slog.LevelDebug,
    // 属性の置き換え・削除など
    ReplaceAttr: func(
      groups []string, a slog.Attr,
    ) slog.Attr {
      a.Key += "!"
      a.Value =
        slog.StringValue(a.Value.String() + "?")
      return a
    },
  },
)
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

すでに `slog.InfoCtx` などの例が出てきましたが、ログ出力を行うAPIについて説明します。
同じログを出力するのにも、いくつかやり方があるので、ここでは3パターンを掲載しておきます。
`context.Context` を引数に取らない関数もありますが、 `ctx` を引数に取る書き方が普及してほしいという筆者の私情により説明しません。

<!-- maprange:../code/structured_logging_test.go,example1 -->
```go title=属性を明示的に組み立てる
slog.InfoCtx(
  ctx, "start processing",
  slog.String("userID", userID),
)
```
<!-- maprange.end -->

<!-- maprange:../code/structured_logging_test.go,example2 -->
```go title=属性をKey-Valueの2つの値をペアとして組み立てる
slog.InfoCtx(
  ctx, "start processing",
  "userID", userID,
)
```
<!-- maprange.end -->

<!-- maprange:../code/structured_logging_test.go,example3 -->
```go title=レベルを明示し属性も明示的に組み立てる
slog.LogAttrs(
  ctx,
  slog.LevelInfo,
  "start processing",
  slog.String("userID", userID),
)
```
<!-- maprange.end -->

筆者の好みは一番最初に出した、明示的に `Attr` を組み立てるやり方が好きなので、本書では基本的にこの書き方でいきます。
世間的には2番目の例が好まれる気はします。

上記コードの出力例は次の通りです。
これは `slog.JSONHandler` の出力例です。
Handlerによって結果のフォーマットはいくらでも変化します。

```text title=出力結果。実際は1行のJSON（LDJSON）として表示される
{
  "time":"2023-05-07T17:09:29.781259+09:00",
  "level":"INFO",
  "msg":"start processing",
  "userID":"a1b2c3"
}
```

構造化されたログの活用を考えたとき、プログラムコード中で属性を明示的に追加するだけではなく、実行時の文脈に依存した（暗黙的な）属性も含めることができるのもメリットです。
たとえばそのログが出力されたソースコード上の位置や、どのリクエストに紐づくログなのかというトレース情報や、現在の操作者が誰かといった情報などです。
先の例では、ログを出力した時刻とログレベルが明示的に指定してないにも関わらず出力されています。

構造化データは機械可読性（machine-readable）が高いのがメリットだと書きました。
この出力例を見ると、ログ出力をBigQuery^[https://cloud.google.com/bigquery/]などのデータレイクに蓄積し、活用できるイメージが湧くと思います。
疑似的なクエリですが `SELECT * FROM Logs WHERE userID = "a1b2c3"` のように特定のユーザに対するログ一覧を見る、といった用途がイメージできますね。
個別の属性に対して検索などの加工ができるようになると、ログをデータとして活用しやすくなります。

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

注意点として、`Handler.WithAttrs` （もしくは `Logger.With` ）で渡された `slog.Attr` はHandlerによって事前に処理済であるとものとして、 `slog.Record` には含まれません。
Handlerがこれらをどう処理するべきかは、TextHandlerやJSONHandlerの内部実装である `slog.commonHandler` の実装を読んでみると面白いでしょう。
概要だけ述べると、 `Handler.WithAttrs` で渡された `slog.Attr` は早期に処理され内部的に `[]byte` で保持され、Handleが呼ばれるたびにシリアライズされないように最適化されています。

次に解説するのは `slog.Value` です。
この値はstructなため直接組み立てることはできません。
`slog.AnyValue`、`slog.BoolValue`、`slog.DurationValue`、`slog.GroupValue` などなどさまざまな関数が用意されています。

`slog.Value` は `slog.Kind` という、値が種類のものかを示す値を持ちます。
`slog.Value` の値をなんらかの理由で自力で処理したい場合、お世話になることもあるかもしれません。
すこし `reflect` パッケージの構造に似ていますね。

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

サンプルを見てみましょう。

<!-- maprange:../code/log_valuer_test.go,sensitiveDef -->
```go title=値を加工し、固定のマスクされた値を返す実装
var _ slog.LogValuer = SensitiveString("")

type SensitiveString string

func (SensitiveString) LogValue() slog.Value {
  return slog.StringValue("**censored**")
}

```
<!-- maprange.end -->

<!-- maprange:../code/log_valuer_test.go,groupDef -->
```go title=複数の値をグループ化し、1つの値として返す実装（Group化）
var _ slog.LogValuer = (*User)(nil)

type User struct {
  ID       string
  Password SensitiveString
}

func (user *User) LogValue() slog.Value {
  return slog.GroupValue(
    slog.String("id", user.ID),
    slog.Any("password", user.Password),
  )
}

```
<!-- maprange.end -->

<!-- maprange:../code/log_valuer_test.go,emit -->
```go title=前述の2つのLogValuer実装を出力してみる
slog.InfoCtx(
  ctx, "print LogValuer value",
  slog.Any("user", &User{
    ID:       "123",
    Password: "53c237",
  }),
)
```
<!-- maprange.end -->

このコードを実際に出力した結果は次の通りです。
例によって実際には1行で出力されているものを筆者が手でインデントをつけています。

```text title=出力結果。userの値はグループ化され、パスワードはマスクされている。
{
  "time":"2023-05-13T18:34:55.929369+09:00",
  "level":"INFO",
  "msg":"print LogValuer value",
  "user":{
    "id":"123",
    "password":"**censored**"
  }
}
```

便利ですね。

もし、手動で `slog.Value` の値を解決したい場合は `value.Resolve()` を呼びます。
値が `LogValuer` である場合、再帰的に値を解決してくれます。
`LogValue()` を直接呼ぶより簡単かつ確実です。

<!-- maprange:../code/log_valuer_test.go,resolve -->
```go title=Resolveは再帰的に値を解決してくれる
v := slog.AnyValue(&User{
  ID:       "123",
  Password: "53c237",
})
// Kind = LogValuer
t.Log(v.Kind())

v = v.Resolve()
// Kind = Group
t.Log(v.Kind())
```
<!-- maprange.end -->


Groupの話が出てきたのでGroupについて解説します。
Groupはその名の通り、1つの属性の値として、グループ化した値（JSONでいうオブジェクト的な構造）を持つことができます。
slogのHandlerにおいて、最終的な出力がJSONであるとは限りません。
様々な出力形態を取りうることから、グループ化された値をJSONで突っ込めばいいというわけにはいきません。
JSONHandlerとTextHandler、それぞれでの出力例を確認してみます。

まずはJSONHandlerでの出力を行うコードです。

<!-- maprange:../code/group_test.go,json -->
```go title=グループ化された値をJSONHandlerで出力
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.InfoCtx(
  ctx, "print group value",
  slog.Group(
    "group",
    slog.String("foo", "bar"),
    slog.Int("fizz", 4),
  ),
)
```
<!-- maprange.end -->

これに対する出力は次の通り。
例によってインデントなどは筆者によるものです。

```text title=JSONでの出力結果
{
  "time":"2023-05-13T18:49:48.969808+09:00",
  "level":"INFO",
  "msg":"print group value",
  "group":{
    "foo":"bar",
    "fizz":4
  }
}
```

次にTextHandlerで出力を行うコードです。
1行目のHandlerの生成以外はJSONHandlerの例と完全に同一です。

<!-- maprange:../code/group_test.go,text -->
```go title=グループ化された値をTextHandlerで出力
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
logger.InfoCtx(
  ctx, "print group value",
  slog.Group(
    "group",
    slog.String("foo", "bar"),
    slog.Int("fizz", 4),
  ),
)
```
<!-- maprange.end -->

これに対する出力は次の通り。
例によってインデントなどは筆者によるものです。

```text title=Textでの出力結果
time=2023-05-13T18:49:48.969+09:00
level=INFO
msg="print group value"
group.foo=bar
group.fizz=4
```

JSONの例ではJSONのオブジェクトになっていましたが、Textの例では属性は個別に出力され、Key名がGroupの名前と連結されていることがわかります。
よって、slogで値のグループ化を行う場合はGroupを使って値を組み立てる必要があります。

Group化の実装を手で書くのは面倒ですので、そのうち `json` タグを見てLogValuerの実装を自動生成するツールが登場する気がします。

最後に、ログ出力レベルである `Level` と、それに関連する `Leveler`、 `LevelVar` についても解説しておきます。

slogでは、ログ出力のレベルは4段階あります。
でかい数字＝より重要 という図式です。

```go title=Levelと各値の定義
type Level int

const (
  LevelDebug Level = -4
  LevelInfo  Level = 0
  LevelWarn  Level = 4
  LevelError Level = 8
)
```

それぞれの値がなぜこれになったのかというのがドキュメントに色々書いてあります。
0をInfoにしたかった。
柔軟性をもたせるために各レベルの間にはすこしの隙間があるほうが都合がよい。
たとえばGoogle Cloud LoggingはInfoとWarnの間にNoticeがある。
あとはOpenTelemetry^[https://opentelemetry.io/]の仕様などの兼ね合いもみてこの定義にしたらしいです。

LevelはLevelerインタフェースを実装しています。
JSONHandlerやTextHandlerで使うHandlerOptionsに設定できるLevelも、実際はLevelerの値を設定できます。

Levelerの定義を見てみましょう。

```go title=Levelerの定義
type Leveler interface {
  Level() Level
}
```

呼び出すとLevelを返すインタフェースですね。
つまり、動的に適用するLevelを変更できる構造になっています。

これを利用した `LevelVar` というstructもあります。
このstructはLevelerを実装しており、かつ動的に値が変更可能になっています。

<!-- maprange:../code/level_leveler_test.go,levelVar -->
```go title=LevelVar利用例
var levelVar slog.LevelVar

h := slog.NewTextHandler(
  os.Stdout,
  &slog.HandlerOptions{
    Level: &levelVar,
    // 出力の見やすさのため time を削除する関数（実装略）
    ReplaceAttr: removeTimeKey,
  },
)
logger := slog.New(h)

ls := []slog.Level{
  slog.LevelDebug,
  slog.LevelInfo,
  slog.LevelWarn,
  slog.LevelError,
}
for _, l := range ls {
  levelVar.Set(l)
  logger.WarnCtx(ctx, "warning!", "l", l)
}
```
<!-- maprange.end -->

```text title=Debug、Info、Warnで出力される
level=WARN msg=warning! l=DEBUG
level=WARN msg=warning! l=INFO
level=WARN msg=warning! l=WARN
```

LevelVarを使う以外にも、HandlerのEnabledに手を加えるなど、出力レベルを制御する方法はいくつかあります。
必要に応じて使い分けましょう。

次に説明する要素は `slog.Source` です。
`slog.Record` のPCを `runtime.CallersFrames` を使って、Function、File、Lineに変換したものです。
つまり、実行時の関数名やソースコードの場所などの情報が詰まっています。
`slog.HandlerOptions` の `AddSource` をtrueにした場合、 このSource structが値をつめた状態で提供されます。

<!-- maprange:../code/add_source_test.go,source -->
```go title=AddSourceをtrueに設定してログ出力
h := slog.NewJSONHandler(
  os.Stdout,
  &slog.HandlerOptions{
    AddSource: true,
  },
)

logger := slog.New(h)
logger.InfoCtx(ctx, "emit source")
```
<!-- maprange.end -->

このコードを出力した時のログは次の通りです。
`source` が追加されていますね。

```text title="source"が追加されている
{
  "time":"2023-05-15T03:34:30.720356+09:00",
  "level":"INFO",
  "source":{
    "function":"example/code.Test_addSource",
    "file":"/work/code/add_source_test.go",
    "line":22
  },
  "msg":"emit source"
}
```

さて、ここまで様々な例を見てきましたが、明示的に指定していない属性がいくつか出力されていました。
`"time"`、`"level"`、`"source"`、`"message"` の4つです。
これらは定数が用意されているため、置き換え処理などに利用することができます。
例を見てみましょう。

<!-- maprange:../code/default_keys_test.go,keys -->
```go title=4つの出力のKeyをリネームしてみる
h := slog.NewJSONHandler(
  os.Stdout,
  &slog.HandlerOptions{
    AddSource: true,
    ReplaceAttr: func(
      groups []string, a slog.Attr,
    ) slog.Attr {
      switch a.Key {
      case slog.TimeKey:
        a.Key = "t"
      case slog.LevelKey:
        a.Key = "l"
      case slog.SourceKey:
        a.Key = "s"
      case slog.MessageKey:
        a.Key = "m"
      }
      return a
    },
  },
)

logger := slog.New(h)
logger.InfoCtx(ctx, "rename keys")
```
<!-- maprange.end -->

出力例も確認します。
Keyの名前が置き換えられていますね。

```text title=Keyが置き換わっている
{
  "t":"2023-05-15T03:39:44.978996+09:00",
  "l":"INFO",
  "s":{
    "function":"example/code.Test_defaultKeys",
    "file":"/work/code/default_keys_test.go",
    "line":37
  },
  "m":"rename keys"
}
```

最後に、 `Logger.With` と `Logger.WithGroup` を説明します。
Handlerの定義に `WithAttrs` と `WithGroup` があったのを覚えていますか？
LoggerのWithとWithGroupは、内部的にはLogger自身をcloneして内部のHandlerのWithAttrs、WithGroupを呼び、それを新しいHandlerにしたLoggerを返しています。

実際のコード例と出力結果をそれぞれ見てみましょう。

<!-- maprange:../code/with_test.go,with -->
```go title=Logger.Withを使ってみる
logger := slog.New(h)
logger = logger.With(
  slog.String("key1", "value1"),
)
logger.InfoCtx(
  ctx, "emit with With",
  slog.String("key2", "value2"),
)
```
<!-- maprange.end -->

```text title=Withで指定したAttrが常に追加される
{
  "time":"2023-05-15T03:49:53.689317+09:00",
  "level":"INFO",
  "msg":"emit with With",
  "key1":"value1",
  "key2":"value2"
}
```

<!-- maprange:../code/with_test.go,withGroup -->
```go title=Logger.WithGroupを使ってみる
logger := slog.New(h)
logger = logger.WithGroup("child")
logger.InfoCtx(
  ctx, "emit with WithGroup",
  slog.String("key3", "value3"),
)
```
<!-- maprange.end -->

```text title=WithGroupで指定したグループの配下に指定したAttrが出力される
{
  "time":"2023-05-15T03:49:53.68944+09:00",
  "level":"INFO",
  "msg":"emit with WithGroup",
  "child":{
    "key3":"value3"
  }
}
```

面白いですね。
特に `Logger.With` は実行している環境（production or staging とか）などの常に出力しておきたい固定値を出力するのに適しています。
これらの機能はHandlerに事前に渡されるため、あらかじめシリアライズされバイト列で保持するなどの工夫が可能になるため、より効率的なログ出力を行える場合があります。

これは `Logger.With` などに限った話ではありませんが、slogでは属性のKeyが重複した場合でも、なにかしらのチェックは行われないため、容赦なく重複したフィールドが出力されます。
こういった事象が発生すると困る場合、自前のHandlerで何らかのチェック処理を行ってやる必要があるでしょう。

<!-- maprange:../code/with_test.go,duplicatedAttrs -->
```go title=Withや明示的なAttrの指定で重複が発生してしまった
logger := slog.New(h)
logger = logger.With(
  slog.String("key", "value1"),
)
logger.InfoCtx(
  ctx, "emit... but duplicated keys",
  slog.String("key", "value2"),
  slog.String("key", "value3"),
)
```
<!-- maprange.end -->

```text title="key"というフィールドが3つも出力されてしまった
{
  "time":"2023-05-15T03:49:53.689452+09:00",
  "level":"INFO",
  "msg":"emit... but duplicated keys",
  "key":"value1",
  "key":"value2",
  "key":"value3"
}
```

## 実用例の紹介

筆者がどのようにslogを実用しているかを紹介していきます。

まずは、新規にプログラムを作る場合です。
筆者は本書を執筆するためにptproc^[https://github.com/vvakame/ptproc]というCLIツールを作成しました。
CLIツールであっても、なんらかのログ出力ライブラリは必要です。
標準に搭載されている `log` パッケージは出力をカスタマイズできず、使いづらいため敬遠。
今までだと `logr` ^[https://github.com/go-logr/logr]を採用していたのですが、その代替としてslogを利用できます。

新規でプログラムを作る場合、特になにも考えずにslogを使い始めることができます。
プログラム中では `slog.DebugCtx` などのパッケージ関数を使ってデフォルトのロガーに処理を任せてます。
main関数やテスト用の関数などでそれぞれ適切なHandlerを作成し、それを利用したLoggerをデフォルトに設定します。

テスト用の関数で `*testing.T` などをHandler内部で利用してログ出力を行いたい場合、いくつか注意点があります。
`t.Parallel()` でテストを並列実行した時、デフォルトのロガーをテストごとに切り替えてしまうと、テストの実行コンテキストとHandlerの後ろにいる `*testing.T` が食い違ってしまう可能性があります。
これを解決するためには、 `context.Context` に絡めた何かしらの仕掛けが必要になるでしょう。
今のところ、これを具合良く行ってくれるライブラリは存在していません。

次に、既存のプログラムを更新する場合です。
筆者の場合、技術書典WebのAPIサーバがそれに該当しました。
このAPIサーバは大昔はGoogle App Engineで稼働していたため、歴史的経緯で google.golang.org/appengine/log ^[https://pkg.go.dev/google.golang.org/appengine/log] と同じAPIを持つ自作の互換ライブラリを利用していました。
これは旧APIを正規表現や置換を駆使してslog経由での出力に無理やり置き換えました。
このため、 `slog.DebugCtx(ctx, fmt.Sprintf("BookShelfItemID: %s is already exists", list[0].ID))` のような、本来はメッセージ+属性に分解できるコードがいたるところに残っています。
とはいえ、slogへの移行自体はそれなりにサクッとできたので、使いたい場所があればslog本来の使い方をできるようになりました。

現在、このAPIサーバはGoogle CloudのCloud Run上で稼働しています。
Cloud Runなどの利用時、ログ出力が一定のルール^[https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields]に従うとCloud Loggingがその形式に従って構造化ログとして解釈してくれます。

筆者は自作のgcpslogライブラリ^[https://pkg.go.dev/github.com/vvakame/sdlog/gcpslog]を利用してログ出力を行っています。
gcpslogで提供しているHandlerの実装は内部的には `slog.JSONHandler` を持っていて、 `Handler.Handle` で使用上追加可能ないくつかの属性をRecordに追加し、JSONHandlerに流す、という処理を行っています。

さて、slogへの移行はそれなりにサクッとできた、と書きましたが、実はこの移行過程で痛い目を見たのでケーススタディとして紹介します。
次のコードは、APIサーバ内で実際に使っているカスタムHandlerの実装の一部です。

```go title=ログ出力時に現在のログインユーザのIDも記録する
func (h *handler) Handle(
  ctx context.Context, record slog.Record,
) error {
  user := domains.CurrentUser(ctx)
  if user != nil {
    record.AddAttrs(
      slog.String("currentUserID", string(user.ID)),
    )
  }

  return h.base.Handle(ctx, record)
}
```

特定のユーザのログが一発で一覧化できると、問い合わせ対応のときに便利じゃん！
…って、思ったんですよね。

本来であれば、ユーザが明らかになった時に `Logger.With` などを使って効率よく出力するべきなのですが、それをやるのは実際めんどくさい…。
コード中ではデフォルトロガー経由で出力します。
ですので、 `context.Context` を使ってリクエストごとに別々のLoggerを生成してHandler内部で振り分けをしてやる構造が必要になり、実装が面倒です。
slogはデフォルトでは `WithContext` や `FromContext` ライクなユーティリティ関数を提供しません。
これがあるとデファクトがあるので考えやすかった気もするんですが…。

で、ログ出力というのはめちゃめちゃ頻繁に呼ばれるため、先のコード中の `domains.CurrentUser` の処理が重たいと大変なことになってしまうわけです。
もちろん、筆者もそれはわきまえていますので、ログインユーザのデータはオンメモリで一瞬で返ってくる構造にしてあります。
ところが、ステージング環境ではセッションからユーザを特定できなかった場合ヘッダを見てごまかす、テスト用の仕組みがありました。
その処理が結構重くて、未ログインユーザがステージング環境にアクセスするとメモリをガンガン消費していって大変なことになる…といった面白アクシデントが発生しました。
未ログイン時の処理が重たいのは想定外だった…普段はログインしっぱなしで開発してるので…。

みなさんは可能な限り `Logger.With` などを活用できるよう、設計段階で工夫を凝らしましょう。
という教訓でした。

## 間違った使い方

slogを間違った使い方をしてしまった時の話をします。
本書ではAttrの組み立ては明示的に `slog.String("key", "value")` 形式で行っていましたが、引数を2つ別々にわたす、 `"key", "value"` 形式でAttrを指定することもできたのでした。
では、keyを指定し忘れて、 `"value"` のように1つだけ値を投げ込むとどうなるのか？確認してみましょう。

<!-- maprange:../code/badkey_test.go,badkey -->
```go title=Attrに当たる部分がKey-Valueではない
logger.InfoCtx(ctx, "test", "value(not-key)")
```
<!-- maprange.end -->

これの実行結果は次の通りです。
Key-Valueのペアがなかったので、Attrの作成に失敗して `BADKEY!` というKeyが表れています。

```test title=Key-Valueのペアじゃなかったので BADKEY! が登場
time=2023-05-14T17:39:08.989+09:00
level=INFO
msg=test
!BADKEY=value(not-key)
```

今のところ、この問題は静的には検出されません。
ですが、 `go vet` でこの課題を検出できるようにしようというProposal^[https://github.com/golang/go/issues/59407]が提出され、すでにAcceptedになっています。
そのためのパッチ^[https://go-review.googlesource.com/c/tools/+/487475]もレビュー中なので、おそらくGo 1.21リリース時には `go vet` をちゃんと利用していればこの問題は実行より前に検出できるはずです。
今のうちから、ちゃんと `go vet` を利用するようにしておきましょう。

## slogに不向きなこと

slogに不向きな用途についても解説しておきます。
構造化ロギングが行えるといえど、与えられる属性の自由度を制限することはできません。
出力するデータ構造は、いわばスキーマレスになるわけです。

一貫したデータ構造がある場合は、しっかりとログ出力を設計し、slog経由以外の、定型データを出力するようにするべきです。
この用途にslogを使うのは、便利ではなくなると思います。
たとえばGraphQLのリクエストログや、特定のデータの更新履歴など、特定のスキーマが存在する場合、一貫したスキーマのもと、特定の出力先に出力するべきでしょう。

スキーマを固定したログ出力にslogを使うべきか、別のものを使うべきかは設計次第ですが、特別なものが特別であるとわかるよう、slogを使わない選択をするのがわかりやすいのではないでしょうか。

## Handlerのテスト

slogにはHandlerのテストを行ってくれるパッケージがあります。
`golang.org/x/exp/slog/slogtest` です。
標準ライブラリになった時には `testing/slogtest` になる予定です。

存在する関数は1つだけで、 `slogtest.TestHandler` がそれです。
引数として `Handler` と、生成された出力のパーサーを渡します。
`Record` からバイト列に変換するのが `Handler` の仕事なので、出力されたバイト列を元のデータ構造に戻せたらHandlerは壊れてないね！と調べてくれる感じです。
若干の脳筋さを感じます。

標準の `slog.JSONHandler` のテストを行うコードを書いてみます。

<!-- maprange:../code/slogtest_test.go,slogtest -->
```go title=slogtestでslog.JSONHandlerをテストする
var buf bytes.Buffer
h := slog.NewJSONHandler(&buf, nil)
err := slogtest.TestHandler(
  h, func() []map[string]any {
    ss := strings.Split(buf.String(), "\n")
    ret := make([]map[string]any, 0, len(ss))

    for _, s := range ss {
      if s == "" {
        continue
      }

      result := map[string]any{}
      err := json.Unmarshal([]byte(s), &result)
      if err != nil {
        t.Fatal(err)
      }

      ret = append(ret, result)
    }
    return ret
  },
)
if err != nil {
  t.Error(err)
}
```
<!-- maprange.end -->

## コミュニティのライブラリ

コミュニティには既にいくつか、slog用ライブラリが存在しています。
5月11日にslog自体に破壊的変更が入ってしまったので、本項目執筆時点（5月14日）ではだいたいのライブラリは壊れています。
利用時は各ライブラリのライセンスを確認のうえ適切に使ってください。

また、slog自体がまだexperimentalな存在なので、市民権を得ているライブラリはありません。
悪意のある実装が存在しないかなどの検討は本書に掲載されている、されていないに関わらず個別によくよく吟味するべきです。
かくいう筆者も、ここで紹介はするものの実際に利用しているライブラリは今のところありません。

`github.com/samber/slog-multi` ^[https://pkg.go.dev/github.com/samber/slog-multi]はslogのHandlerを複数組み合わせたり、条件によって振り分けたり、ロードバランシングを行ったりと多機能なライブラリです。
筆者は複数のログ出力先を利用していないのですが、出力のレベルがErrorであればSlackに通知する、などのユースケースはあるかもしれません。
技術書典Webでもデータ不整合の検知などはSlackに通知していますが、slog経由でできると嬉しいのかも…。

このライブラリの作者のsamberさんは、他にも `slog-slack` や `slog-fluentd` など、色々なライブラリを作っているので他のライブラリもチェックしてみると自分にマッチしたものがあるかもしれません。

`github.com/galecore/xslog` ^[https://pkg.go.dev/github.com/galecore/xslog]は様々なサービス、ライブラリをバックエンドとするHandlerのコレクションライブラリです。
OpenTelemetryやSentry、zap、testingパッケージ向けの実装などが揃っています。
一番最初にslogのHandlerの実装例を勉強するのに丁度いいかもしれません。

`github.com/neilotoole/slogt` ^[https://pkg.go.dev/github.com/neilotoole/slogt]は `testing.T` を引数にとるHandlerです。
テストコードから呼び出されたslogでのログ出力を `t.Log` などに出力してくれるライブラリです。
テスト失敗時、slog経由でのログ出力を確認したい場合はかなりあるので、便利そうなライブラリです。

## おわりに

slogパッケージを説明しました。
既存ライブラリを踏まえ、実用的で使いやすいライブラリになっていると思います。
`context.Context` にLoggerを持たせたり取り出したりする標準的な手法がほしい気がしますが、そこはコミュニティの今後のユースケースにかかってくるでしょう。
今後はslogは幅広い局面で使われることになると思いますので今後のコミュニティのアクションが非常に楽しみです。
みなさんも面白い使い方、便利な使い方を編み出したら、ぜひSNSやブログなどでシェアしてみてください。
