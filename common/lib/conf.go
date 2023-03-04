package lib

import (
	"bytes"
	"database/sql"
	dlog "gateway-micro/common/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type BaseConf struct {
	DebugMode    string    `mapstructure:"debug_mode"`
	TimeLocation string    `mapstructure:"time_location"`
	Log          LogConfig `mapstructure:"log"`
	Base         struct {
		DebugMode    string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
}

type LogConfFileWriter struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type LogConfConsoleWriter struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

type LogConfig struct {
	Level string               `mapstructure:"log_level"`
	FW    LogConfFileWriter    `mapstructure:"file_writer"`
	CW    LogConfConsoleWriter `mapstructure:"console_writer"`
}

type MysqlMapConf struct {
	List map[string]*MySQLConf `mapstructure:"list"`
}

type MySQLConf struct {
	DriverName      string `mapstructure:"driver_name"`
	DataSourceName  string `mapstructure:"data_source_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

type RedisMapConf struct {
	List map[string]*RedisConf `mapstructure:"list"`
}

type RedisConf struct {
	ProxyList    []string `mapstructure:"proxy_list"`
	Password     string   `mapstructure:"password"`
	Db           int      `mapstructure:"db"`
	ConnTimeout  int      `mapstructure:"conn_timeout"`
	ReadTimeout  int      `mapstructure:"read_timeout"`
	WriteTimeout int      `mapstructure:"write_timeout"`
}

var (
	ConfBase        *BaseConf
	DBMapPool       map[string]*sql.DB
	GORMMapPool     map[string]*gorm.DB
	DBDefaultPool   *sql.DB
	GORMDefaultPool *gorm.DB
	ConfRedis       *RedisConf
	ConfRedisMap    *RedisMapConf
	ViperConfMap    map[string]*viper.Viper
)

// GetBaseConf 获取基本配置信息
func GetBaseConf() *BaseConf {
	return ConfBase
}

func InitBaseConf(path string) error {
	ConfBase = &BaseConf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
		return err
	}

	if ConfBase.DebugMode == "" {
		if ConfBase.Base.DebugMode != "" {
			ConfBase.DebugMode = ConfBase.Base.DebugMode
		} else {
			ConfBase.DebugMode = "debug"
		}
	}
	if ConfBase.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.TimeLocation = "Asia/Chongqing"
		}
	}
	if ConfBase.Log.Level == "" {
		ConfBase.Log.Level = "trace"
	}

	//配置日志
	logConf := dlog.Config{
		Level: ConfBase.Log.Level,
		FW: dlog.ConfFileWriter{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
		},
		CW: dlog.ConfConsoleWriter{
			On:    ConfBase.Log.CW.On,
			Color: ConfBase.Log.CW.Color,
		},
	}
	if err := dlog.SetupDefaultLogWithConf(logConf); err != nil {
		panic(err)
	}
	dlog.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

//
//func InitLogger(path string) error {
//	if err := dlog.SetupDefaultLogWithFile(path); err != nil {
//		panic(err)
//	}
//	dlog.SetLayout("2006-01-02T15:04:05.000")
//	return nil
//}

func InitRedisConf(path string) error {
	ConfRedis := &RedisMapConf{}
	err := ParseConfig(path, ConfRedis)
	if err != nil {
		return err
	}
	ConfRedisMap = ConfRedis
	return nil
}

// InitViperConf 初始化配置文件
func InitViperConf() error {
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}

// GetStringConf 获取get配置信息
func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
	return confString
}

// GetStringMapConf 获取get配置信息
func GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetConf 获取get配置信息
func GetConf(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.Get(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetBoolConf 获取get配置信息
func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetFloat64Conf 获取get配置信息
func GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetIntConf 获取get配置信息
func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringMapStringConf 获取get配置信息
func GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringSliceConf 获取get配置信息
func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetTimeConf 获取get配置信息
func GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetDurationConf 获取时间阶段长度
func GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// IsSetConf 是否设置了key
func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.IsSet(strings.Join(keys[1:len(keys)], "."))
	return conf
}
