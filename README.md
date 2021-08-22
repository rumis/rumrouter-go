# RumRouter-Go

RumRouter是一个基于注解自动生成路由代码的工具。它可以避免人工书写路由代码带来的大量重复劳动，也可以避免人为引入的错误。

## 快速使用

### 安装

    go get github.com/rumis/rumrouter-go/cmd/rumrouter

安装完成后会在`$GOPATH/bin`目录下生成可执行文件`rumrouter`，可以通过`rumrouter -help`验证是否安装成功并查看支持的命令行参数：

    $ rumrouter -help

    Usage of rumrouter:
    -i string
            input source code dirctory,default is current work dirctory (default "/home/ubuntu/workspace/annotion/echodemo")
    -o string
            router code dirctory,default is rumrouter (default "/home/ubuntu/workspace/annotion/echodemo/rumrouter")
    -p string
            package name,default is rumrouter (default "rumrouter")
    -t string
            framework type, now echo,echov4,gin is available (default "echov4")


`-i` 我们项目的根目录，必须包含`go.mod`文件，也就是要求我们的项目必须是使用`go mod`做包管理的。

`-o` 最终生成的路由代码存放路径，需要使用绝对路径，默认为项目目录下的`rumrouter`文件夹下。

`-p` 生成的路由代码的包名称，需要和参数`-i`配合使用，这样最终生成的代码才方便直接使用。

`-t` 框架类型，目前支持`gin`,`echo`,`echov4`三种，默认为`echov4`

### 使用

#### 注解路由定义

目前本工具可以识别三种注解定义：中间件、路由组、路由，分别对应一般go框架里的middleware、routergroup、router三种结构。

##### 中间件

    // @Middleware(name="auth")

中间件的定义比较简单，仅包含`name`一个参数，它相当于一个索引，需要在项目中全局唯一。在生成代码时，路由组以及路由中配置的中间名名称会被替换为实际的实现。

##### 路由组

    // @RouterGroup(middleware="auth",prefix="/user")

路由组包含两个属性，`prefix`表示本路由组中所有路由的统一前缀。`middleware`表示应用于本路由组的中间件名称，名称在定义中间件注解时定义。支持一次定义多个中间件名称，通过逗号隔开，按照定义顺序中间件会依次执行。

##### 路由

    // @Router(method="options",path="/getconfig",middleware="auth")

路由包含三个属性，`middleware`和路由组中的使用方式完全一致。 `method`表示该路由支持HTTP METHOD，需要注意的是mehtod在使用时一定需要和我们使用的框架匹配，一定是需要框架支持的，否则会报错。`path`表示该路由的部分请求路径，此处和路由组中的`prefix`组成完成的请求路径。


#### 示例代码

    package test

    import (
        "net/http"

        "github.com/labstack/echo/v4"
    )

    // Person the person controller
    // @RouterGroup(prefix="/person",middleware="auth")
    type Person struct{}

    // GetAge get the person age
    // @Router(path="/age",method="GET")
    func (p *Person) GetAge(c echo.Context) error {
        return c.String(http.StatusOK, "1")
    }

    // GetName get the person name
    // @Router(path="/name",method="GET,POST")
    func (p Person) GetName(c echo.Context) error {
        return c.String(http.StatusOK, "test liu")
    }

#### 代码生成

项目根目录执行如下命令：

    rumrouter -t echov4

生成如下代码，`rumrouter/echo.v4.gen.go`

    package rumrouter

    import (
        "github.com/labstack/echo/v4"
        "demo.com/echodemo/api/test"
    )

    func InitRouter(app *echo.Echo) {

        g0 := app.Group("/person")
        gPersonInst := test.Person{}
        g0.GET("/age", gPersonInst.GetAge)
        g0.GET("/name", gPersonInst.GetName)
        g0.POST("/name", gPersonInst.GetName)
    }

项目`main`函数中添加路由的引用

    package main

    import (
        "github.com/labstack/echo/v4"
        "demo.com/echodemo/rumrouter"
    )

    func main() {

        e := echo.New()

        rumrouter.InitRouter(e)

        e.Logger.Fatal(e.Start(":7070"))
    }
