# mos-chinadns

一个简单的DNS解析器，实现基于域名与IP的分流，支持DoH，IPv6，[EDNS Client Subnet](https://tools.ietf.org/html/rfc7871)。

---

- [mos-chinadns](#mos-chinadns)
  - [命令帮助](#命令帮助)
  - [配置文件说明](#配置文件说明)
  - [三分钟快速上手 & 预设配置](#三分钟快速上手--预设配置)
  - [分流效果](#分流效果)
  - [实现细节](#实现细节)
    - [黑白名单](#黑白名单)
    - [关于EDNS Client Subnet (ECS)](#关于edns-client-subnet-ecs)
    - [关于DNS-over-HTTPS (DoH)](#关于dns-over-https-doh)
    - [关于文件路径](#关于文件路径)
  - [Open Source Components / Libraries](#open-source-components--libraries)

## 命令帮助

    -c string
            [路径]json配置文件路径
    -gen string
            [路径]生成一个json配置文件模板至该路径
    -dir string
            [路径]变更程序的工作目录
    -dir2exe
            变更程序的工作目录至可执行文件的目录

    -v    调试模式，更多的log输出

## 配置文件说明

    {
        // [IP:端口][必需] 监听地址
        "bind_addr": "127.0.0.1:53", 

        // [IP:端口] 本地服务器地址 建议:一个低延时但会被污染大陆服务器，用于解析大陆域名。
        "local_server": "223.5.5.5:53",     

        // [bool] 本地服务器是否屏蔽非A或AAAA请求。
        "local_server_block_unusual_type": false,

        // [IP:端口] 远程服务器地址 建议:一个无污染的服务器。用于解析非大陆域名。   
        "remote_server": "8.8.8.8:443", 

        // [URL] 远程DoH服务器的url，如果填入，远程服务器将使用DoH协议
        "remote_server_url": "https://dns.google/dns-query",  

        // [bool] 是否跳过验证DoH服务器身份 高危选项，会破坏DoH的安全性
        "remote_server_skip_verify": false, 

        // [int] 单位毫秒 远程服务器延时启动时间
        // 如果在设定时间(单位毫秒)后local_server无响应或失败，则开始请求remote_server。
        // 如果local_server延时较低，将该值设定为120%的local_server的延时可显著降低请求remote_server的次数。
        // 0表示禁用延时，请求将同时发送。
        "remote_server_delay_start": 0, 

        // [路径] 本地服务器IP白名单 建议:中国大陆IP列表，用于区别大陆与非大陆结果。
        "local_allowed_ip_list": "/path/to/your/chn/ip/list", 

        // [路径] 本地服务器IP黑名单 建议:希望被屏蔽的IP列表，比如运营商的广告服务器IP。
        "local_blocked_ip_list": "/path/to/your/black/ip/list",
        
        // [路径] 强制使用本地服务器解析的域名名单 建议:中国的域名。
        "local_forced_domain_list": "/path/to/your/domain/list",

        // [路径] 本地服务器域名黑名单 建议:希望强制打开国外版而非中国版的域名。
        "local_blocked_domain_list": "/path/to/your/domain/list",

        // [CIDR] EDNS Client Subnet 
        "remote_ecs_subnet": "1.2.3.0/24"
    }

## 三分钟快速上手 & 预设配置

在这里下载最新版本：[release](https://github.com/IrineSistiana/mos-chinadns/releases)

一份最新的中国大陆IPv4与IPv6的地址表`chn.list`和域名表`chn_domain.list`已包含在release的zip包中。

`chn.list`数据来自[APNIC](https://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest)。`chn_domain.list`数据来自[felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) [LICENSE](https://github.com/felixonmars/dnsmasq-china-list/blob/master/LICENSE)。

将预设配置复制并保存至`config.json`，确保`chn.list`，`chn_domain.list`，`config.json`和`mos-chinadns`在同一目录。

用以下命令启动

    mos-chinadns -c config.json -dir2exe

<details><summary><code>预设配置1 通用 按大陆IP与域名分流</code></summary><br>

使用中国大陆IP表`chn.list`和域名表`chn_domain.list`分流。国内域名使用`阿里云DNS`解析，国际域名使用`OpenDNS`解析。

    {
        "bind_addr": "127.0.0.1:53",
        "local_server": "223.5.5.5:53",
        "remote_server": "208.67.222.222:443",
        "local_allowed_ip_list": "./chn.list",
        "local_forced_domain_list": "./chn_domain.list"
    }

</details>

<details><summary><code>预设配置2 通用 按大陆IP与域名分流 远程服务器使用DoH</code></summary><br>

使用中国大陆IP表`chn.list`和域名表`chn_domain.list`分流。国内域名使用`阿里云DNS`解析，国际域名使用[Google DoH](https://developers.google.com/speed/public-dns/docs/doh)解析。

    {
        "bind_addr": "127.0.0.1:53",
        "local_server": "223.5.5.5:53",
        "remote_server": "8.8.8.8:443",
        "remote_server_url": "https://dns.google/dns-query",
        "local_allowed_ip_list": "./chn.list",
        "local_forced_domain_list": "./chn_domain.list"
    }

</details>

<details><summary><code>预设配置3 单DoH客户端模式</code></summary><br>

使用[Google DoH](https://developers.google.com/speed/public-dns/docs/doh)作为上游服务器，解析所有域名。

建议启用ECS使解析更精确。[如何启用?](#关于edns-client-subnet-ecs)

    {
        "bind_addr": "127.0.0.1:53",
        "remote_server": "8.8.8.8:443",
        "remote_server_url": "https://dns.google/dns-query",
        "remote_server_skip_verify": false,
        "remote_server_delay_start": 0
    }

</details>

## 分流效果

国内域名交由`local_server`解析，无格外延时。国外域名将会由`remote_server`解析，确保无污染。

<details><summary><code>dig www.baidu.com 演示</code></summary><br>

    ubuntu@ubuntu:~$ dig www.baidu.com @192.168.1.1 -p5455

    ; <<>> DiG 9.11.3-1ubuntu1.11-Ubuntu <<>> www.baidu.com @192.168.1.1 -p5455
    ;; global options: +cmd
    ;; Got answer:
    ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 57335
    ;; flags: qr rd ra; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 1

    ;; OPT PSEUDOSECTION:
    ; EDNS: version: 0, flags:; udp: 4096
    ;; QUESTION SECTION:
    ;www.baidu.com.			IN	A

    ;; ANSWER SECTION:
    www.baidu.com.		561	IN	CNAME	www.a.shifen.com.
    www.a.shifen.com.	250	IN	A	36.152.44.96
    www.a.shifen.com.	250	IN	A	36.152.44.95

    ;; Query time: 4 msec
    ;; SERVER: 192.168.1.1#5455(192.168.1.1)
    ;; WHEN: Sun Mar 15 18:17:55 PDT 2020
    ;; MSG SIZE  rcvd: 149

</details>

<details><summary><code>dig www.google.com 演示</code></summary><br>

    ubuntu@ubuntu:~$ dig www.google.com @192.168.1.1 -p5455

    ; <<>> DiG 9.11.3-1ubuntu1.11-Ubuntu <<>> www.google.com @192.168.1.1 -p5455
    ;; global options: +cmd
    ;; Got answer:
    ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 2719
    ;; flags: qr rd ra; QUERY: 1, ANSWER: 6, AUTHORITY: 0, ADDITIONAL: 1

    ;; OPT PSEUDOSECTION:
    ; EDNS: version: 0, flags:; udp: 512
    ;; QUESTION SECTION:
    ;www.google.com.			IN	A

    ;; ANSWER SECTION:
    www.google.com.		280	IN	A	74.125.68.99
    www.google.com.		280	IN	A	74.125.68.105
    www.google.com.		280	IN	A	74.125.68.104
    www.google.com.		280	IN	A	74.125.68.103
    www.google.com.		280	IN	A	74.125.68.106
    www.google.com.		280	IN	A	74.125.68.147

    ;; Query time: 72 msec
    ;; SERVER: 192.168.1.1#5455(192.168.1.1)
    ;; WHEN: Sun Mar 15 18:19:20 PDT 2020
    ;; MSG SIZE  rcvd: 223

</details>

## 实现细节

### 黑白名单

**流程**

1. 如果指定了域名白名单->匹配域名->白名单中的域名将被本地服务器解析
2. 如果指定了域名黑名单->匹配域名->黑名单中的域名将被远程服务器解析
3. 如果指定了IP黑名单->匹配本地服务器返回的IP->丢弃黑名单中的结果。
4. 如果指定了IP白名单->匹配本地服务器返回的IP->不在白名单的结果将被丢弃。

远程服务器的结果一定会被接受

**域名黑/白名单格式**

采用按域向前匹配的方式，与dnsmasq匹配方式类似。每个表达式一行。

规则：

* `cn`相当于`*.cn`。会匹配所有以cn结尾的域名，`example.cn`，`www.google.cn`
* `google.com`相当于`*.google.com`。会匹配`www.google.com`, `www.l.google.com`，但不会匹配`www.google.cn`。

比如：

    cn
    google.com
    google.com.hk
    www.google.com.sg

**IP黑/白名单格式**

由单个IP或CIDR构成，每个表达式一行，支持IPv6，比如：

    1.0.1.0/24
    1.0.2.0/23
    1.0.8.0/21
    2001:dd8:1a::/48

    2.2.2.2
    3.3.3.3
    2001:ccd:1a

**性能**

IP列表与域名列表均已做性能优化。IP列表采用二分搜索，数据仅存储在一个对象上。域名列表采用Hash，数据仅存储在三个对象上。无需担心长列表的匹配时间与GC的压力。

### 关于EDNS Client Subnet (ECS)

`remote_ecs_subnet` 填入自己的IP段即可启用ECS。如不详请务必留空。

启用ECS最简单的方法:

- 百度搜索`IP`，得到自己的IP地址，如`1.2.3.4`
- 将最后一位变`0`，并加上`/24`。如`1.2.3.4`变`1.2.3.0/24`
- 将`1.2.3.0/24`填入`remote_ecs_subnet`

更多ECS资料请参考rfc文档：[EDNS Client Subnet](https://tools.ietf.org/html/rfc7871)

### 关于DNS-over-HTTPS (DoH)

填入同时填入`remote_server`和`remote_server_url`即可启用DoH模式。请求方式为[RFC 8484](https://tools.ietf.org/html/rfc8484) GET。

想了解有那些服务器支持DoH，请参阅[维基百科公共域名解析服务列表](https://en.wikipedia.org/wiki/Public_recursive_name_server)。

### 关于文件路径

建议使用`-dir2exe`选项将工作目录设置为程序所在目录，这样的话配置文件`-c`路径和配置文件中的路径可以是相对于程序的相对路径。

如过附加`-dir2exe`后程序启动报错那就只能启动程序前手动`cd`或者使用绝对路径。

## Open Source Components / Libraries

部分设计参考

* [ChinaDNS](https://github.com/shadowsocks/ChinaDNS): [GPLv3](https://github.com/shadowsocks/ChinaDNS/blob/master/COPYING)

依赖

* [sirupsen/logrus](https://github.com/sirupsen/logrus): [MIT](https://github.com/sirupsen/logrus/blob/master/LICENSE)
* [miekg/dns](https://github.com/miekg/dns): [LICENSE](https://github.com/miekg/dns/blob/master/LICENSE)
* [valyala/fasthttp](https://github.com/valyala/fasthttp):[MIT](https://github.com/valyala/fasthttp/blob/master/LICENSE)
