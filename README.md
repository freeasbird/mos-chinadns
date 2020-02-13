# mos-chinadns

一个适用于中国的DNS解析器，实现DNS分流与防污染，支持IPv6，支持[EDNS Client Subnet](https://tools.ietf.org/html/rfc7871)。

---

## 命令帮助

    -bind-addr string
            [地址:端口]监听地址 e.g. '127.0.0.1:53'(必需)
    -c string
            [路径]json配置文件路径
    -dir string
            [路径]变更程序的工作目录
    -gen string
            [路径]生成一个json配置文件模板至该路径
    -local-server string
            [地址:端口]本地服务器地址 (必需)
    -local-blocked-domain-list string
            [路径]域名黑名单位置
    -local-allowed-ip-list string
            [路径]IP白名单位置
    -local-blocked-ip-list string
            [路径]IP黑名单位置
    -remote-server string
            [地址:端口]远程服务器地址 (必需)
    -use-tcp string
            [l|r|l_r] 是否使用TCP协议。'l'表示仅本地服务器使用TCP, 'r'表示仅远程服务器使用TCP, 'l_r'表示都使用TCP
    -remote-ecs-subnet
            [CIDR格式IP段] 向远程服务器的请求将使用EDNS0并附带包含有此地址的ESC信息 e.g. '1.2.3.0/24'
    -v    调试模式，更多的log输出

使用方式请参考[配置示例](#配置示例)

## 配置示例

如果不清楚如何配置，以下是一种常用的防污染与分流配置：

* 本地服务器(local-server)：一个低延时但会被污染大陆服务器，用于解析大陆域名。(必需)
* 远程服务器(remote-server)：一个[不会被污染的服务器](#不会被污染的服务器)。用于解析非大陆域名。(必需)
* 本地服务器域名黑名单(local-blocked-domain-list)：强制这些域名用远程服务器解析。用于强制打开这些域名的国外版而非中国版。(非必需)
* 本地服务器IP白名单(local-allowed-ip-list)：中国大陆IP列表，用于区别大陆与非大陆结果。最新的列表可以从[这里](https://github.com/LisonFan/china_ip_list)获得。(必需)
* 本地服务器IP黑名单(local-blocked-ip-list)：希望被屏蔽的IP列表，比如运营商的广告服务器IP。(非必需)

## 本地服务器黑白名单工作流程

* 如果指定了域名黑名单->匹配域名->黑名单中的域名将被远程服务器解析
* 如果指定了IP黑名单->匹配本地服务器返回的IP->丢弃黑名单中的结果。
* 如果指定了IP白名单->匹配本地服务器返回的IP->不在白名单的结果将被丢弃。

## json配置文件模板

        {
                "bind_addr": "",
                "local_server": "",
                "remote_server": "",
                "use_tcp": "",
                "local_allowed_ip_list": "",
                "local_blocked_ip_list": "",
                "local_blocked_domain_list": "",
                "remote_ecs_subnet": ""
        }

## 域名黑名单

由正则表达式构成，每个表达式一行

## IP黑/白名单

由单个IP或CIDR构成，每个表达式一行，支持IPv6，比如：

    1.0.1.0/24
    1.0.2.0/23
    1.0.8.0/21
    2001:dd8:1a::/48

    2.2.2.2
    3.3.3.3
    2001:ccd:1a

## 无污染的服务器

远程服务器(remote_server)必须为无污染服务器，如`8.8.8.8`, `1.1.1.1`。因为某些不可抗拒原因，在大陆使用UDP或TCP直连是不稳定的。所以建议：

* 用代理将请求转发
* 在本地运行DoH客户端，比如[mos-doh-client](https://github.com/IrineSistiana/mos-doh-client)

## Open Source Components / Libraries

部分设计参考

* [ChinaDNS](https://github.com/shadowsocks/ChinaDNS): [GPLv3](https://github.com/shadowsocks/ChinaDNS/blob/master/COPYING)

依赖

* [sirupsen/logrus](https://github.com/sirupsen/logrus): [MIT](https://github.com/sirupsen/logrus/blob/master/LICENSE)
* [miekg/dns](https://github.com/miekg/dns): [LICENSE](https://github.com/miekg/dns/blob/master/LICENSE)