# はじめに

この本はGoの次期標準ライブラリであるslogについて、実際に利用して、それを解説するものです。
技術書典Web<span class="footnote">https://techbookfest.org/</span>や本書を執筆するために作ったptproc<span class="footnote">https://github.com/vvakame/ptproc</span>などでslogを実用しています。
本書を執筆している2023年5月の段階では、slog搭載のGo 1.21はまだまだリリースされていないため、`golang.org/x/exp/slog` <span class="footnote">https://pkg.go.dev/golang.org/x/exp/slog</span> を利用し、これについて解説していきます。

使っている `golang.org/x/exp` のバージョンは次の通りです。

<!-- mapfile:./exp.version.txt -->
```
v0.0.0-20230510235704-dd950f8aeaea
```
<!-- mapfile.end -->

本パッケージはexperimentalなだけあって、稀に破壊的変更が入ります。
5月11日時点で `slog.NewJSONHandler` や `slog.NewTextHandler` 、その他のシグニチャが変わったりしてドヒャりました<span class="footnote">https://github.com/golang/exp/commit/dd950f8aeaea</span>。
パッケージ利用時は本書記述のものから変更が入っている可能性を念頭においてください。

本書にかかれていることの多くは公式のプロポーザル<span class="footnote">https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md</span>にも書かれている内容です。
自分で実用を開始してみてからプロポーザルを読み返したら、「困ったこと全部書いてあるじゃん！」と落ち込んだりもしました。
そんな筆者の体験と実用上の注意点をまとめたものが本書になります。

Goには長らくの間、構造化ロギングライブラリの標準が存在していませんでした。
Go標準のlogパッケージ<span class="footnote">https://pkg.go.dev/log</span>は存在はすれど、構造化されたログを出力する柔軟性を持っていません。
そのため、筆者はプロセス起動時などの限られたシチュエーションでのみ利用していました。

この隙間を埋めるためにコミュニティによって様々なロギングライブラリが作られました。
Uberのzap<span class="footnote">https://pkg.go.dev/go.uber.org/zap</span>、zerolog<span class="footnote">https://pkg.go.dev/github.com/rs/zerolog</span>、Kubernetesでも使われるlogr<span class="footnote">https://pkg.go.dev/github.com/go-logr/logr</span>などが登場してきました。

これらの知見を活かし、遂に登場した標準ライブラリがslogです。
APIの設計的にはlogrに近い構造で、使いやすく、かつカスタマイズしやすい構造になっていると思います。
今後Goコミュニティがどう発展していくかはまだわからないですが、標準であるslogが広く使われていくことは間違いないでしょう。
