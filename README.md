
# Masterlab 实时通信服务器

由ueligo socket框架精简而来，可作为 游戏，聊天，异步服务器。
 
## golang 环境安装： 
   * Windows 安装示例 [windows安装golang环境](windows安装golang环境)  
   
  
   * Centos安装示例 [centos安装golang环境](centos安装golang环境)  
  
  

   * Ubuntu安装示例 [ubuntu安装golang环境](ubuntu安装golang环境)  

  

## 编译
下载 masterlab_socket 源码
```
git clone https://github.com/gopeak/masterlab_socket.git
cd masterlab_socket

go build
```

## 运行
masterlab_socket有两个配置文件, `config.toml` 是主配置文件，有端口和数据库连接的配置等信息。`cron.json`是定时执行任务配置。
直接运行编译后的exe文件，出现类似下图则安装成功
   ![](http://www.masterlab.vip/docs/images/masterlab_socket/masterlab_socket_win.png)  


 
