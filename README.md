# Coffee

A Content Manger System based on Violet

[Demo](https://coffee.zhenly.cn/)

[中文简介](https://github.com/XMatrixStudio/Coffee/blob/master/Doc/README.md)

## Using

部署MongoDB

填写配置文件

```yaml
Mongo:
  Host: 0.0.0.0
  Port: 27890
  User: username
  Password: password
  Name: coffee
Server:
  Host: 127.0.0.1
  Port: 30070
  Dev: true
  ThumbDir: ./Thumb
  UserDir: ./UserData
  TempDir: ./Temp
Violet:
  ClientID: xxxxxxxx
  ClientKey: xxxxxxxxxx
  ServerHost: https://oauth.xmatrix.studio/api/v2
  LoginURL: https://oauth.xmatrix.studio/Verify/Authorize
```

使用Docker打包镜像

```bash
docker build -t coffee .
```

运行

```bash
docker run -p 30070:30070 --name coffee_server -d coffee:latest
```





