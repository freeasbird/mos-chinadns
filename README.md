# mos-chinadns

一个简单的DNS解析器，实现基于IP的分流，支持DoH，支持IPv6，支持[EDNS Client Subnet](https://tools.ietf.org/html/rfc7871)。

---

## 命令帮助

    -c string
            [路径]json配置文件路径
    -gen string
            [路径]生成一个json配置文件模板至该路径
    -dir string
            [路径]变更程序的工作目录

    -v    调试模式，更多的log输出

使用方式请参考[配置示例](#配置示例)

## 配置示例

如果不清楚如何配置，以下是一种常用的防污染与分流配置：

*下文的`必需`表示该参数这个分流方案的必要参数，不表示是程序的必要参数*

* 监听地址(bind_addr)(必需)：按需设置。
* 本地服务器(local-server)(必需)：一个低延时但会被污染大陆服务器，用于解析大陆域名。比如阿里DNS 223.5.5.5。
* 远程服务器(remote-server)(必需)：一个无污染的服务器。用于解析非大陆域名。比如OpenDNS的非常规端口 208.67.222.222:443。
* 本地服务器IP白名单(local-allowed-ip-list)(必需)：中国大陆IP列表，用于区别大陆与非大陆结果。最新的列表可以从[这里](https://github.com/LisonFan/china_ip_list)获得。
* 本地服务器域名黑名单(local-blocked-domain-list)(非必需)：强制这些域名用远程服务器解析。用于强制打开这些域名的国外版而非中国版。
* 本地服务器IP黑名单(local-blocked-ip-list)(非必需)：希望被屏蔽的IP列表，比如运营商的广告服务器IP。

**json配置文件模板**

*直接复制使用需先去掉注释*

    {
        "bind_addr": "127.0.0.1:53",            // IP:端口
        "local_server": "223.5.5.5:53",         // IP:端口
        "remote_server": "208.67.222.222:443",  // IP:端口
        "remote_server_url": "",                // URL
        "remote_server_skip_verify": false,     // ture 或 false
        "remote_server_delay_start": 0,         // 整数 单位毫秒
        "local_allowed_ip_list": "/path/to/your/chn/ip/list",
        "local_blocked_ip_list": "",            // 路径
        "local_blocked_domain_list": "",        // 路径
        "remote_ecs_subnet": ""                 // CIDR 如：1.2.3.0/24
    }

远程服务器(remote-server)支持DoH，只需将服务器IP地址填入`remote_server`，URL填入`remote_server_url`。

如使用[Google 的 DNS over Https](https://developers.google.com/speed/public-dns/docs/doh)，只需填入：

    ...
    "remote_server": "8.8.8.8:443",
    "remote_server_url": "https://dns.google/dns-query",
    ...

**其他非必要选项说明**

* remote_server_skip_verify: 跳过DoH服务器身份验证。**高危选项，会破坏DoH的安全性，仅在知道自己在干什么的情况下启用**
* remote_server_delay_start: 如果在设定时间(单位毫秒)后local_server无响应，则开始请求remote_server。将该值设定为local_server的延时可显著降低请求remote_server的次数。0表示将同时发送请求。
* remote_ecs_subnet: 关于ECS作用请参考[EDNS Client Subnet](https://tools.ietf.org/html/rfc7871)。**仅少数DNS服务商提供ECS支持**

**其他使用方式**

* 仅指定`remote_server`。当作普通带ECS的DoH客户端使用

## 本地服务器黑白名单

**工作流程**

1. 如果指定了域名黑名单->匹配域名->黑名单中的域名将被远程服务器解析
2. 如果指定了IP黑名单->匹配本地服务器返回的IP->丢弃黑名单中的结果。
3. 如果指定了IP白名单->匹配本地服务器返回的IP->不在白名单的结果将被丢弃。


**域名黑名单格式**

由正则表达式构成，每个表达式一行

**IP黑/白名单格式**

由单个IP或CIDR构成，每个表达式一行，支持IPv6，比如：

    1.0.1.0/24
    1.0.2.0/23
    1.0.8.0/21
    2001:dd8:1a::/48

    2.2.2.2
    3.3.3.3
    2001:ccd:1a

## Open Source Components / Libraries

部分设计参考

* [ChinaDNS](https://github.com/shadowsocks/ChinaDNS): [GPLv3](https://github.com/shadowsocks/ChinaDNS/blob/master/COPYING)

依赖

* [sirupsen/logrus](https://github.com/sirupsen/logrus): [MIT](https://github.com/sirupsen/logrus/blob/master/LICENSE)
* [miekg/dns](https://github.com/miekg/dns): [LICENSE](https://github.com/miekg/dns/blob/master/LICENSE)
* [valyala/fasthttp](https://github.com/valyala/fasthttp):[MIT](https://github.com/valyala/fasthttp/blob/master/LICENSE)
