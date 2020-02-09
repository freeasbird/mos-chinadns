# mos-chinadns

一个适用于中国的DNS解析器，实现DNS分流与防污染。参考[ChinaDNS](https://github.com/shadowsocks/ChinaDNS)。

---

## 使用方式

    -b string
            [地址:端口]监听地址 (必需)
    -c string
            [路径]json配置文件路径
    -dir string
            [路径]变更程序的工作目录
    -gen string
            [路径]生成一个json配置文件模板至该路径
    -local string
            [地址:端口]本地服务器地址 (必需)
    -local-blocked-domain string
            [路径]域名黑名单位置
    -local-allowed-ip string
            [路径]IP白名单位置
    -local-blocked-ip string
            [路径]IP黑名单位置
    -remote string
            [地址:端口]远程服务器地址 (必需)
    -use-tcp string
            [l|r|l_r] 使用TCP协议，'l'表示仅本地服务器使用TCP, 'r'表示仅远程服务器使用TCP, 'l_r'表示都使用TCP
    -v    调试模式，更多的log输出


## 本地服务器黑白名单工作流程

* 如果指定了域名黑名单->匹配域名->黑名单中的域名将被远程服务器解析
* 如果指定了IP黑名单->匹配本地服务器返回的IP->丢弃黑名单中的结果。
* 如果指定了白名单->匹配本地服务器返回的IP->不在白名单的结果将被丢弃。

## json配置文件

    {
        "bind_addr": "",
        "local_server": "",
        "remote_server": "",
        "use_tcp": "",
        "local_allowed_ip_list": "",
        "local_blocked_ip_list": "",
        "local_blocked_domain_list": ""
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

## 配置示例

* 本地服务器：会被污染大陆服务器
* 远程服务器：不会被污染的服务器
* 本地服务器域名黑名单：会被污染的域名
* 本地服务器IP白名单：中国大陆IP
* 本地服务器IP黑名单：希望被屏蔽的IP，比如运营商的广告服务器IP

## Open Source Components / Libraries

部分设计参考

* [ChinaDNS](https://github.com/shadowsocks/ChinaDNS): [GPLv3](https://github.com/shadowsocks/ChinaDNS/blob/master/COPYING)

依赖

* [sirupsen/logrus](https://github.com/sirupsen/logrus): [MIT](https://github.com/sirupsen/logrus/blob/master/LICENSE)
* [miekg/dns](https://github.com/miekg/dns): [LICENSE](https://github.com/miekg/dns/blob/master/LICENSE)