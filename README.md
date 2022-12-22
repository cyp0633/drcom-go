# drcom-go

使用 Go 语言实现的 Drcom 客户端

兴趣项目，不保证可用性

可用 Golang 的交叉编译特性编译不同 ABI 和 OS 的二进制文件，便于在不同设备上运行

AGPLv3 协议授权

## 使用方法

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

生成配置文件的方法请参考 [drcom-generic Wiki](https://github.com/drcoms/drcom-generic/wiki/d%E7%89%88%E7%AE%80%E7%95%A5%E4%BD%BF%E7%94%A8%E5%92%8C%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E)

## 灵感来自

- [drcom-generic](https://github.com/drcoms/drcom-generic)
- [dogcom](https://github.com/mchome/dogcom)
- 无数分析 Drcom 协议的热心网友
