# PacketListener · 网络流量监测与入侵检测

一个用 Go 实现的网络入侵检测系统（IDS），配有 Next.js 管理界面。选择网卡后，可以实时抓包，也可以按自定义规则匹配流量、记录事件、发送告警或封禁来源 IP。管理界面可用于查看日志、编辑规则和维护黑白名单。

## 快速开始

### 启动 Go 后端

后端的抓包和封禁流程以 Linux 为运行环境，需要 Go 1.22.1 或更高版本、C 编译器、libpcap 开发库、Bash、`iptables`、`ipset` 和可连接的 MySQL 实例。请在隔离的测试环境中运行。程序启动时会修改防火墙规则，收到退出信号后会执行 `ipset flush`。

在 Debian 或 Ubuntu 上，先安装编译和抓包依赖，再克隆仓库。

```bash
sudo apt install build-essential libpcap-dev iptables ipset
git clone https://github.com/dadawre/packet-listener.git
cd packet-listener
```

启动前，编辑 [conf/ids.conf](conf/ids.conf) 中的 `config { ... }`。

- 将 `sql_host`、`sql_user` 和 `sql_password` 改成你的 MySQL 连接信息。保持 `sql_schema = IDSdata`，后端启动时会检查数据库和日志表，并在需要时创建。
- `ids_mode = 1` 启用规则检测，`ids_mode = 2` 在终端查看抓到的数据包。只测试抓包和规则时，可设 `start_proxy = 0`。
- 将 `mail_*` 参数指向测试用的 SMTP 服务。后端启动时会把一封测试邮件加入发送队列。

在仓库根目录编译并启动，然后按终端提示输入要监听的网卡编号。

```bash
go build -o ids .
sudo ./ids
```

后端 API 监听 `8080` 端口。启用代理时，HTTP/HTTPS 代理监听 `8899` 端口。

### 启动管理界面

前端源码在同一仓库的 `master` 分支。`main` 中的 `fore-end` 指向该分支的前端提交。仓库缺少 `.gitmodules`，直接克隆 `master` 即可获取源码。

准备 Node.js 20 或更高版本和 npm，在另一个终端中运行以下命令。

```bash
git clone --branch master --single-branch https://github.com/dadawre/packet-listener.git packet-listener-ui
cd packet-listener-ui
npm ci
npm run dev -- --hostname 127.0.0.1
```

用浏览器打开登录页。前端从浏览器访问 `http://localhost:8080`，所以浏览器需要能通过该地址连接 Go 后端。[打开管理界面](http://127.0.0.1:3000)

登录账号和密码由后端的 `account`、`password` 配置项指定。登录后可查看事件记录和封禁列表，并编辑规则与黑白名单。

## 配置与规则

### 定义变量和加载规则

后端启动时读取 `conf/ids.conf`。用 `include` 加载规则文件，用 `$变量名` 引用已定义的变量。相对路径以启动目录为准，`#` 用于配置文件注释。`ipv4var` 和 `ipv6var` 接受单个地址或 CIDR 网段。IPv6 地址使用完整的八组写法。`portvar` 接受单个端口或端口列表。

```text
var ruleDir rulers
include $ruleDir/myruler
ipv4var labNetwork 192.168.56.0/24
ipv6var labNetworkV6 2001:0DB8:0000:0000:0000:0000:0000:0000/64
portvar webPorts [80,443]

config {
    max_coroutine = 256
    log_path = ./log
}
```

### 编写检测规则

在现有配置中加入 `labNetwork` 变量，再把下面的规则加入 `rulers/myruler` 并重启后端。来自 `192.168.56.0/24` 且发往 TCP 80 或 443 端口的流量命中时，日志会填入数据包的源地址和目的地址。[仓库自带规则示例](rulers/myruler)

```text
IPv4 src $labNetwork & ( TCP dst 80 | TCP dst 443 ) priority 5 action record name lab_web msg:"Web traffic from %s to %s" IPv4.src IPv4.dst
```

- 用 `&`、`|` 和括号组合条件，运算符和括号两侧留空格。同一组内从左到右求值，混用 `&` 和 `|` 时可用括号明确分组。
- TCP 和 UDP 端口条件使用数字。匹配多个端口时，用 `|` 连接条件。
- `priority` 使用正整数，数值越小越先匹配。每个数据包执行首条命中规则的动作。
- 用 `name` 给规则命名。`msg:"..."` 指定日志文本，消息后的 `协议.字段` 按顺序填入格式占位符。

| 动作 | 命中后的行为 |
| --- | --- |
| `pass` | 放行并记录命中事件到 MySQL。 |
| `record` | 写入普通日志并记录事件到 MySQL。 |
| `alert` | 写入告警日志并记录事件到 MySQL。 |
| `send` | 写入告警日志，将事件加入邮件队列，并记录到 MySQL。 |
| `drop` | 通过 `ipset` 封禁来源 IPv4 地址，并写入告警日志和 MySQL。可在动作后加 `banTime 15` 指定封禁 15 分钟。 |

### 规则字段

协议名和字段名区分大小写。抓包模式可以显示 ARP，当前规则解析器接受下表中的关键字。

| 关键字 | 字段 |
| --- | --- |
| `Eth` | `src`, `dst`, `length`, `type` |
| `IPv4` | `src`, `dst`, `checksum`, `protocol`, `TTL`, `offset`, `flags`, `id`, `length`, `TOS`, `IHL`, `version` |
| `IPv6` | `src`, `dst`, `hlimit`, `nexthead`, `length`, `label`, `class`, `version` |
| `TCP` | `src`, `dst`, `seq`, `seqack`, `offset`, `URG`, `ACK`, `PSH`, `RST`, `SYN`, `FIN`, `window`, `checksum`, `urgptr` |
| `UDP` | `src`, `dst`, `length`, `checksum` |
| `OS` | `cpu`, `gpu`, `memory` |

## 仓库结构

| 路径 | 内容 |
| --- | --- |
| [`main.go`](main.go) | 后端入口，初始化配置并启动抓包或规则检测。 |
| [`back_end/`](back_end/) | 抓包、协议解析、规则匹配、日志存储、HTTP API 和代理。 |
| [`conf/`](conf/) | 主配置、黑白名单和防火墙配置。 |
| [`rulers/`](rulers/) | 检测规则示例。 |
| [`bash/`](bash/) | 初始化和维护 `iptables`、`ipset` 的脚本。 |
| [`master` 分支](https://github.com/dadawre/packet-listener/tree/master) | Next.js 管理界面源码。 |
