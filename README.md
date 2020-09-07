# chat-room

# 环境
go 1.13版本，Mac机器可以直接执行二进制程序（src/roomsvr/roomsvr，src/websocketcli/websocketcli），其他机器环境先执行go build。

# 运行

## roomsvr

聊天房间进程，可以创建房间或者加入某个已存在的房间，执行下面命令启动房间进程，监听8080端口。使用了websocket库github.com/gorilla/websocket。

```
./rommsvr
```

## websocketcli

测试程序，启动websocket客户端连接roomsvr。可以启动多个websocketcli进程进行测试。

### 创建房间
```
./websocketcli -addr 127.0.0.1:8080 -cmd create_room -user roland

```
创建后会返回房间的ID，例如：
```
SYSTEM: roland create room 1 success
SYSTEM: roland join room
```

### 加入房间
需要指定房间ID。
```
./websocketcli -addr 127.0.0.1:8080 -cmd join_room -user stone -roomid 1
```
加入房间后，会立即收到该房间最近50条消息，并提示加入房间成功。
```
roland: dsdf
roland: sdfdsf
roland: sdf
roland: sd
roland: fd
roland: sf
roland: ds
roland: fd
SYSTEM: stone join room
```

### 发送聊天消息
直接在终端输入，roomsvr会对聊天消息进行关键词屏蔽。
```
you are fecker
roland: you are ****er
```


### 获取popular单词
在终端输入：/popular 10，回车会返回最近10s的popular单词：
```
SYSTEM: popular word is wo
```

### 获取用户信息
在终端输入：/stats [user]，回车会返回用户登录时长：
```
SYSTEM: 00d 00h 00m 08s
```

# 扩展性

roomsvr可以水平扩展，但是还需要一个router服务（未开发），客户端创建房间时需要向router询问负载最小的roomsvr地址，然后再进行连接。同样，在加入房间时也询问房间对应的roomsvr地址。router是无状态的，也可以进程水平扩展，各个roomsvr定期向redis集群上报负载情况。
