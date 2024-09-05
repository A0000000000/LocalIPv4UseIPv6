# 已过时，新项目详见[IPv4ByIPv6](https://github.com/A0000000000/IPv4ByIPv6)

# 使用公网IPv6进行内网IPv4组网
* 仅支持Linux系统（借助于Linux下的tun/tap实现）
* 仅支持arp与ip协议

# 配置文件
* Local下为本机配置
  * ip：本机的一个内网ipv4地址，必填
  * mask：组网后的子网掩码，必填
* Remote为远程设备的内网ipv4与公网ipv6地址，选填

# 端口监听
* []:1113端口用于转发不同设备子网中的数据包，使用TCP6协议
* :8080使用go-zero为本机提供的http服务，用于对配置文件中的Remote子项进行增删改查。注意：这里的更新不会持久化到文件，仅本次运行生效

# API
* GET /config：配置文件以JSON形式返回
* GET /config/:ipv4：返回一组Remote的配置
* Post /config：增加一条Remote的记录，json data传参
  * Params：
    * ipv4 string
    * ipv6 string
* PUT /config：更新一条Remote的记录
  * Params：
    * ipv4 string
    * ipv6 string
* DELETE /config/:ipv4：删除一条Remote的配置
