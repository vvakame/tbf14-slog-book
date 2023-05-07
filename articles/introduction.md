# はじめに

この本はGoの次期標準ライブラリであるslogについて、実際に利用して、それを解説するものです。
本書を執筆している2023年5月の段階では、slog搭載のGo 1.21はまだまだリリースされていないため、`golang.org/x/exp/slog`[^expslog] を代わりに利用しています。

[^expslog]: https://pkg.go.dev/golang.org/x/exp/slog

本書にかかれていることの多くは公式のプロポーザル[^proposal]にも書かれている内容です。
自分で実用を開始してみてからプロポーザルを読み返したら、「全部書いてあるじゃん！」と落ち込んだりもしました。
そんな筆者の体験と実用上の注意点をまとめたものが本書になります。

[^proposal]: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md

Goには長らくの間、構造化ロギングライブラリの標準が存在していませんでした。
Go標準のlogパッケージ[^stdlog]は存在はすれど、構造化されたログを出力する柔軟性を持っていないため、プロセス起動時などの限られたシチュエーションでのみ使われていました。

[^stdlog]: https://pkg.go.dev/log

この隙間を埋めるためにコミュニティによって様々なロギングライブラリが作られました。
Uberのzap^[https://pkg.go.dev/go.uber.org/zap]、zerolog^[https://pkg.go.dev/github.com/rs/zerolog]、Kubernetesでも使われるlogr^[https://pkg.go.dev/github.com/go-logr/logr]などが登場してきました。

これらの知見を活かし、遂に登場した標準ライブラリがslogです。
APIの設計的にはlogrに近い構造で、使いやすく、かつカスタマイズしやすい構造になっていると思います。
今後Goコミュニティがどう発展していくかはまだわからないですが、標準であるslogが広く使われていくことは間違いないでしょう。
