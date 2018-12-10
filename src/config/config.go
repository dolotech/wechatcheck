package config

import (
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"os"
)


type WxAccount struct {
	Account [][]string
	Port string
}



// Config 配置类型
type PublicConfig struct {
	WxAccount WxAccount

}

// Opts Config 默认配置
var opts *PublicConfig

// ParseToml 解析配置文件
func ParseToml(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {

		glog.Errorln("没有找到配置文件 ...")

		return nil
	}
	opts = &PublicConfig{}
	_, err := toml.DecodeFile(file, opts)
	if err != nil {
		glog.Errorln("配置文件解析错误：", err)
		return err
	}
	glog.Infof("config is %v", opts)
	return nil
}

// Opts 获取配置
func Opts() *PublicConfig {
	return opts
}


func GetWxAccount() WxAccount {
	return opts.WxAccount
}
