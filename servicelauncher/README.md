# service launcher

## 说明

用于在IDE里调试服务，实现方式是绕开了用plugin加载服务，直接通过对应服务的NewService()得到一个新服务，后面的加载仍然调用serviceloader/loader里的方法进行。

## 服务的配置位置

```servicelauncher/launcher/services.go```

如果要新增一个服务，修改里面的LoadService方法，在switch里新增对应服务名和对应包路径的调用。

## 启动参数

```servicelauncher match --config=../match/config.yml```

与serviceloader一样，但**实际使用并不需要go install**，而是在IDE里配置作为启动包，详见下面的IDE启动配置


## IDE启动配置

### VS Code

1. 先保存工作区（workspace）
2. 编辑steve文件夹下的.vscode的launch.json
3. 在json里的configuration里配置下面的json对象（以room为例），configuration是一个数组，可以配置多个
4. 配置之后在左侧第四个图标（Ctrl+Shift+D）里就有了

```
        {
            "name": "Launch room",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "remotePath": "",
            "port": 2345,
            "host": "127.0.0.1",
            "program": "${fileDirname}",
            "env": {},
            "args": ["room", "--config=${workspaceFolder}/room/config.yml"],
            "showLog": true,
            "trace": "verbose"
        }
```

### GoLand

目前没有找到GoLand的Run/Debug设置可以使用变量，所以只能使用绝对路径，就不能让配置通用了。

这是一个可用的设置的步骤（可能还有其他配置方式，可以补充）：

1. 新增一个Go Build
2. Run kind选择Directory
3. Directory要指到steve\servicelauncher的绝对路径
4. Working directory指到steve的绝对路径
5. Program arguments填入（以启动match为例）：match --config=match/config.yml

如果要启动多个服务，重复1-5步，第5步修改match为对应服务名即可