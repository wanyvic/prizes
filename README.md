环境 mongodb 创建用户 并启动密码验证   doc/mongo_install.sh文件
编译 ./build.sh
启动示例
```
nohup sudo ./prizesd --db-server mongodb://massgrid:password@localhost:27017/docker --rpc-server localhost:9442 --rpc-username 5de7b1f3 --rpc-password password --time-Scale-Statement 30 -D -l debug >>debug.log &
```

命令参数
```
./prizesd -h
Flag shorthand -h has been deprecated, please use --help

Usage:  prizesd [OPTIONS]

A monitor for docker swarm

Options:
      --config-file string         Daemon configuration file (default "/etc/prizes/daemon.json")
      --db-server list             database host (default mongodb://localhost:27017)
  -D, --debug                      Enable debug mode    //debug模式
      --help                       Print usage
  -H, --host list                  Daemon socket(s) to connect to   //api address : -H tcp://localhost:2000
  -l, --log-level string           Set the logging level ("debug"|"info"|"warn"|"error"|"fatal") (default "info")
      --rpc-password string        Set MassGrid rpc password (default "password")
      --rpc-server list            MassGrid rpc host (default tcp://localhost:9442)
      --rpc-username string        Set MassGrid rpc username (default "user")
      --testnet                    Set massgrid testnet
  -t, --time-Scale int             Set record Millisecond time scale to database (default 3000) //docker 请求间隔时间
      --time-Scale-Statement int   Set time cycle for statement Minute (default 5) //结算间隔时间
  -v, --version                    Print version information and quit
```


模块
db  
    node    存储所有node信息
    service     存储所有server信息
    servicetimeaxis     存储每个server的租用时间轴信息
    task        //存储所有活动过的task信息
server 
    restAPI标准
    通过tcp、unix_sock 与massgridd 通讯
refresh
    优先队列 对server进行计时
prizeservice
    service 创建退款结算更新操作
