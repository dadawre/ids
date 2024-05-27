conf文件文件规则如下:

    预定义关键字:any ,max_coroutine,log_path,include等
    包含rules文件规则:include "规则名"
    IPv4变量定义规则:ipv4var IPv4变量名 IPv4或者IPv4网段/子网掩码位数 (IPv4采用.分10进制表示方法)
    IPv6变量定义规则:ipv6var IPv6变量名 IPv6或者IPv6网段/子网掩码位数 (IPv6采用冒号:加16进制表示)
    端口变量定义规则:portvar port变量名 port 或者 port数组 如[20,30]
    普通变量定义规则:var 变量名 值
    变量调用规则:$变量名
    注释规则:# 注释字符串
    参数配置规则:
    config
    {参数名 = 参数值}
    max_coroutine = 1024
    log_path = ./log
    }

rules文件规则如下: 
    
    首先协议执行操作有:放行(pass，不会有日志)，记录(record，将数据包记录到普通日志中)，
    告警(alert 将数据包记录到紧急日志中,并产生告警)，发送给总部(send 并记录到发送日志中)，
    丢弃(drop 丢弃数据包并记录到紧急日志中)
    协议名 {协议字段名 协议字段值} {&| | |~ 协议名 协议字段名 协议字段值} priority 优先级(优先级越高，越后执行,优先级值范围为0~65535) action 执行操作 [msg:"日志字符串"] 
    
    note:&&的优先级高于||,先运算完&&再再运算||,每种运算操作从左到右依次执行。
rule支持协议:
    
    ARP协议支持字段名 type,protocol,hsize,psize,option,shaddress,spaddress,dhaddress,dpaddress
    MAC协议支持字段名 src,dst,length,type
    IPv4协议支持字段名 src,dst,checksum,protocol,TTL,offset,flags,id,length,TOS,IHL,version
    TCP协议支持字段名 src,dst,seq,seqack,offset,URG,ACK,PSH,RST,SYN,FIN,window,checksum,urgptr
    UDP协议支持字段名 src,dst,length,checksum
    IPv6协议支持字段名 dst,src,hlimit,nexthead,length,label,class,version
    
