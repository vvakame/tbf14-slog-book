# Vivliostyle使ってみた

諸般の事情で本書の内容自体は手堅くハズれがなくライブラリ規模もさほど大きくないものにしたので、反動で別のところで冒険する気持ちになりました。
今回は普段使っているRe:VIEWを離れ、Vivliostyleを試してみることにしました。

というわけで、この章ではVivliostyleを使い始めるにあたり辿った道筋をここに記しておきます。

とりあえずプロジェクト作成。

```shell
$ yarn create book tbf14-slog-book
...
? Description description
? Author name vvakame
? Author email vvakame@gmail.com
? License MIT
? choose theme @vivliostyle/theme-techbook - Techbook (技術同人誌) theme
...
```

作成されたファイルの内容は、うーん、だいぶ質素…。
GitHub Actionsでのビルドとかがないのはもちろん、ファイルも最低限の構成で、このまま本を書き始められる感じじゃないなぁ。
テーマを自分で作れるし、実際に配布している人もいるようだけど、公式からの導線が弱いのが辛いポイント。

PDFをビルドしてみたけど、物理本をターゲットにしたスタイリングになっている。
小口とノドの幅にかなり差があり、電子機器でPDFを読む時の体験が悪そう。
この辺はどこで設定変えられるんだろうな？


ついでに、epubを生成する方法も今のところわかってない。


https://vivliostyle.org/ja/make-books-with-create-book/

https://docs.vivliostyle.org/#/ja/create-book
