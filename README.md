# drcom-go

使用 Go 语言实现的 Drcom 客户端。

兴趣项目，本人实测可用，但不保证在其他环境下也可用。作者不对使用本项目带来的任何后果负责。

AGPLv3 协议授权（毕竟也算是用了其他人 AGPLv3 drcom 客户端的代码）。

## 下载

在 [Releases](https://github.com/cyp0633/drcom-go/releases) 页面或 [Actions](https://github.com/cyp0633/drcom-go/actions) 页面下载。

前者（可能）更稳定，后者则包含最新的特性。

## 使用方法

要使用本工具，需要 `drcom.conf` 配置文件和命令行参数，其中命令行参数部分见下。

```bash
drcom-go <command> [flags]
```

**command**: `dhcp` 或 `pppoe`，分别对应 DHCP 和 PPPoE 客户端

**flags**:

- `-h` / `--help` 查看帮助
- `-c` / `--conf` 指定配置文件路径，兼容 drcom-generic 的配置文件格式
- `-b` / `--bind-ip` 指定绑定的 IP 地址，用于绑定网卡
- `-l` / `--log` 指定日志文件路径，指定后将不输出到标准输出
- `-d` / `--daemon` 后台运行
- `-e` / `--eternal` 无限重连
- `-D` / `--debug` 输出调试信息

如 `drcom-go dhcp -c ./drcom.conf -e -d` 是一种常见用法，它将使用 DHCP 模式，在后台运行，登录失败也将尝试无限重连，并使用当前目录下的 `drcom.conf` 配置文件。

配置文件编写方法请见 [Wiki](https://github.com/cyp0633/drcom-go/wiki/%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)。

## 灵感来自

- [drcom-generic](https://github.com/drcoms/drcom-generic)
- [dogcom](https://github.com/mchome/dogcom)
- 无数分析 Drcom 协议的热心网友
