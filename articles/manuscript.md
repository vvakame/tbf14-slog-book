# slog

https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md

そもそもしっかりした型が必要ならそれなりの対応するべき。slogが適切な候補ではないのでは。

vet

https://go.googlesource.com/go/+/refs/heads/master/src/log/slog/
https://go.googlesource.com/go/+/refs/heads/master/src/log/slog/internal/slogtest/
https://go-review.googlesource.com/c/exp/+/469075

> Handler authors and others may wish to use Value.Resolve
instead of calling LogValue directly.

HandleでCurrentUser読んだら字面に反して色々やってる処理だったのでメモリ使用量が一気にエグいことになってAPIがエラーを返す様になった話。

NewLogLogger 間違えがち

https://github.com/samber/slog-multi
