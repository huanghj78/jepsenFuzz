# JepsenFuzz

## 简介

本框架灵感来自[jepsen io/jepsen](https://github.com/jepsen-io/jepsen)和[pingcap/tipocket](https://github.com/pingcap/tipocket)，**JepsenFuzz**是一个模糊测试驱动的故障注入框架，旨在测试分布式数据库。
JepsenFuzz保留了Jepsen的原始功能，同时结合了基于模糊算法的故障调度，提高了故障注入的智能性和有效性。
JepsenFuzz目前仍在开发中……


## 依赖

* [chaosblade](https://github.com/chaosblade-io/chaosblade) >= v1.7.4
（需要在被测机器上安装，并进行vim ~/.bash_profile，添加export PATH="$PATH:/usr/local/bin/chaosblade-1.7.4"（对应的位置）并source ~/.bashrc使其生效）


## 故障场景

* 节点宕机：random_kill, all_kill, minor_kill, major_kill
* 时钟偏移：small_skews, subcritical_skews, critical_skews, big_skews, huge_skews, strobe_skews
* 网络分区：partition, two_partition, multi_partition, all_partition
* 网络数据包故障：loss, delay, duplicate, corrupt
* 资源不足故障：random_cpufl, all_cpufl, major_cpufl, minor_cpufl, disk_burn, disk_fill, mem_fullload

### 创建新的故障注入场景
主要需要实现“Invoke”、“Recover”等函数，需要保证Recover操作可以消除之前所注入的故障，否则可能导致故障持续作用。随后还需要在/jepsenFuzz/cmd/util/suit.go中添加对应的故障触发parse；还需要在/jepsenFuzz/pkg/nemesis/nemesis.go中对您所添加的nemesis进行注册。


## 快速开始

### 运行

进入testcase中的某个文件夹后，修改“client.db”文件中的setup函数，更改其中的dsn语句，将“数据库用户名”和“数据库用户密码”修改为相应的信息
```go
dsn := fmt.Sprintf("postgres://数据库用户名:数据库用户密码@%s/test?sslmode=disable&target_session_attrs=read-write", addressesStr)
```
除此之外，还需要进入/jepsenFuzz/util/util.go中的ExecuteRemoteCommand函数，修改username和password，设置为您被测数据库所在机器的用户名和密码（以此使所有nemesis都使用该用户名和密码来进行故障注入）

随后即可回到testcase中的某个文件夹中，执行`make build`，生成二进制可执行文件，执行的一个例子如下：

```bash
./bin/vbank -node-addr 10.10.3.0:26000  -node-addr 10.10.4.26:26000 -node-addr 10.10.3.76:26000 -node-addr 10.10.1.9:26000 -node-addr 10.10.1.174:26000 -nemesis proc-kill
```


## 工作负载

* gauss-bank-standard/schedule：模拟银行交易的工作负载
* gauss-block-writer：模拟分块写的压力测试工作负载

### 创建新的工作负载
主要需要实现数据库的“SetUp”、“TearDown”、“Start”、“Execute”、“verify”、“initDB”等函数，并需要实现一些工作负载的具体执行的函数