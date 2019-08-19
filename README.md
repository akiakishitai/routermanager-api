# RouterManager-API

WiFi強度の確認やWiFiサービスの再起動、デバイスの再起動・終了をするAPIプログラム。

## 作った理由

イーサネットコンバータ化したLinuxデバイスがWiFi切れやすく、そのたびにSSH接続してコマンド打つのが面倒。
これを解消するためAPIから確認・再起動させたかった。

APIからシステムの電源操作するという邪道なことしているので、外部公開しないよう注意してください。

## 実施環境

* Server: PINE A64
* OS: Armbian

## 設定

* APIポート

    デフォルトは `55555` ポート。

* 対象サーバの指定  
    本サンプルでは `scripts/ansihle/hosts/staging` で指定。  
    サーバのIPアドレスは `$HOME/.ssh/ssh_config` からの利用を想定。

*  `sudo` パスワードの指定  
    `root` 権限が必要なコマンドが含まれるため、 `sudo` パスワードを`ansible_become_password` で指定する。  
    ``

## インストール

1. コンパイル

    環境変数 `GOOS`, `GOARCH` は対象サーバによって変更。

    ```shell
    GOOS=linux GOARCH=arm64 go build -o build/routermanager-api cmd/main.go
    ```

1. ホストで `ansible playbook` 実行

    Ansible Vault のパスワードを `$HOME/.ansible/vault_password` ファイルで指定

    ```shell
    ansible-playbook \
        --inventory=scripts/ansihle/hosts/staging \
        --vault-id=$HOME/.ansible/vault_password \
        scripts/ansihle/main.yaml
    ```

リモートサーバ上で本APIプログラムが `systemd` サービスとして自動起動される。

## API

[api/openapi.yaml](api/openapi.yaml)

## Example

```shell
$ curl http://pine64.router.io:55555/v1/wifi
{"device":"wlx1","connected":true,"signal":54}

$ curl -X POST http://pine64.router.io:55555/v1/wifi

$ curl -X POST http://pine64.router.io:55555/v1/sys/poweroff -d "reboot"
```
