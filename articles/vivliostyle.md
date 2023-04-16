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

[技術書典NFT](https://techbookfest.org/product/9GsUCh18KNL40JTqr4savb)のPDFを眺めてみる。
こいつはいつもどおりの[TechBoosterのReVIEW-Template](https://github.com/TechBooster/ReVIEW-Template)製なので、こいつの見た目に近くなるのが期待値。
いつもにやつに見た目似せたいならいつものでいいじゃん、という内なる声に屈していてはイノベーションなどないのだ、と克己していく。
うーん、やっぱりノドと小口がほぼ同じだなぁ。

というわけで、Vivliostyleでもなにか先達がいないかと思って「Vivliostyle ノド 小口」で検索すると次の記事が見つかった。
https://zenn.dev/macneko/articles/06aec138a357b9
ええやんこれで！！記事の最後とかでテーマ配布します！とか言い始めたら最高！と思いつつ読み勧めたけどそういうものはなかった。
残念無念。

