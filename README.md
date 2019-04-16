# rproxy
simple reverse proxy

## Usage:
```sh
git clone https://github.com/kokardy/rproxy
cd rproxy
go build .
./proxy -scheme http \
         -rhost www.yourtarget.com \
         -addr :8080 \
         -ori regexp_before \
         -dest regexp_after
```

自分のポートと別のドメインを1対1で対応したリバースプロキシを立てます。

簡単に既存のサイトのミラーを作れます。
(キャッシュはしないので処理の分散には使えない)

Go言語製なので環境変数のHTTP_PROXYが利いてくれます。

oriとdestでレスポンスを書き換えて返すことができます(Content-Typeヘッダーがtext/htmlの場合のみ)。

サブドメインへのリンクが貼ってある場合、

サブドメイン用にもう一つ実行して 
oriとdestでリンクを相互に書き換えればサブドメインのリンクも使えます。