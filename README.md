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
Go言語製なので環境変数のHTTP_PROXYが利いてくれます。
リバースプロキシをプロキシごしに使いたい場合にnginxだと少し面倒なのと、
サブドメインへのリンクが貼ってあったら、
oriとdestでリンクを書き換えることができます(Content-Typeヘッダーがtext/htmlの場合のみ)