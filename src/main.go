package main

import (
	"config"
	"flag"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"logic"
	"net/http"
)

func main() {
	var fileName string

	flag.StringVar(&fileName, "conf", "config.toml", "Configuration file to start game")
	flag.Parse()
	glog.Errorln("Configuration is", fileName)

	// load config
	err := config.ParseToml(fileName)
	if err != nil {
		glog.Errorln("配置文件.toml出错")
		glog.Fatal(err)
	}
	logic.InitAccount(config.GetWxAccount().Account)

	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/wechatcheck", check)

	e.Start(config.GetWxAccount().Port)
}

func check(c echo.Context) error {
	longUrl := c.QueryParam("url")
	return c.JSON(http.StatusOK, logic.Check(longUrl))
}
