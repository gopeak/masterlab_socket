
# Masterlab 实时通信服务器

由ueligo socket框架精简而来，可作为 游戏，聊天，异步服务器。
 
## golang 环境安装： 
   * Windows 安装示例 [windows安装golang环境](https://github.com/gopeak/masterlab_socket/wiki/windows%E5%AE%89%E8%A3%85golang%E7%8E%AF%E5%A2%83)  
   
  
   * Centos安装示例 [centos安装golang环境](https://github.com/gopeak/masterlab_socket/wiki/centos%E5%AE%89%E8%A3%85golang%E7%8E%AF%E5%A2%83)  
  
  

   * Ubuntu安装示例 [ubuntu安装golang环境](https://github.com/gopeak/masterlab_socket/wiki/ubuntu%E5%AE%89%E8%A3%85golang%E7%8E%AF%E5%A2%83)  

  

## 编译
下载 masterlab_socket 源码
```
git clone https://github.com/gopeak/masterlab_socket.git
cd masterlab_socket

go build
```

## 运行
masterlab_socket有两个配置文件, `config.toml` 是主配置文件，有端口和数据库连接的配置等信息。`cron.json`是定时执行任务配置。  
相关命令  
```
启动命令： ./masterlab_scoket start
后台运行： ./masterlab_scoket start -d
停止后台进程：./masterlab_scoket stop
指定配置文件：./masterlab_scoket start start -c /xxx/config.toml
```

   ![](http://www.masterlab.vip/docs/images/masterlab_socket/masterlab_socket_win.png)  


 
