# isucon13-tr
## APPの環境構築
DBの接続情報は事前にAPPで使われている環境変数に挿入して下さい。

## DBの環境構築
コンテナを起動する。
```bash
$ docker compose up -d
```

mysql-clientでアクセスする。
```bash
$ mysql -u root -p -h 127.0.0.1 -P 3306
```
