package conf

import (
	"encoding/json"
	_ "github.com/b3log/wide/i18n"
	"github.com/b3log/wide/util"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type User struct {
	Name     string
	Password string
}

type conf struct {
	Server                string
	StaticServer          string
	EditorChannel         string
	OutputChannel         string
	ShellChannel          string
	StaticResourceVersion string
	ContextPath           string
	StaticPath            string
	RuntimeMode           string
	Repos                 string
	UserRepos             string
	Users                 []User
}

var Wide conf
var rawWide conf

func Save() bool {
	// 可变部分
	rawWide.Users = Wide.Users

	// 原始配置文件内容
	bytes, err := json.MarshalIndent(rawWide, "", "    ")

	if nil != err {
		glog.Error(err)

		return false
	}

	if err = ioutil.WriteFile("conf/wide.json", bytes, 0644); nil != err {
		glog.Error(err)

		return false
	}

	return true
}

func Load() {
	bytes, _ := ioutil.ReadFile("conf/wide.json")

	err := json.Unmarshal(bytes, &Wide)
	if err != nil {
		glog.Error(err)

		os.Exit(-1)
	}

	// 保存未经变量替换处理的原始配置文件，用于写回时
	json.Unmarshal(bytes, &rawWide)

	ip, err := util.Net.LocalIP()
	if err != nil {
		glog.Error(err)

		os.Exit(-1)
	}

	glog.Infof("IP [%s]", ip)
	Wide.Server = strings.Replace(Wide.Server, "{IP}", ip, 1)
	Wide.StaticServer = strings.Replace(Wide.StaticServer, "{IP}", ip, 1)
	Wide.EditorChannel = strings.Replace(Wide.EditorChannel, "{IP}", ip, 1)
	Wide.OutputChannel = strings.Replace(Wide.OutputChannel, "{IP}", ip, 1)
	Wide.ShellChannel = strings.Replace(Wide.ShellChannel, "{IP}", ip, 1)

	// 获取当前执行路径
	file, _ := exec.LookPath(os.Args[0])
	pwd, _ := filepath.Abs(file)
	pwd = pwd[:strings.LastIndex(pwd, string(os.PathSeparator))]
	glog.Infof("pwd [%s]", pwd)

	Wide.Repos = strings.Replace(Wide.Repos, "{pwd}", pwd, 1)
	Wide.UserRepos = strings.Replace(Wide.UserRepos, "{pwd}", pwd, 1)

	glog.Info("Conf: \n" + string(bytes))
}
