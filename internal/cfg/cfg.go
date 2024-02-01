package cfg

import (
	"bufio"
	"net"
	"os"
	"reflect"
	"regexp"
	L "shortlink2/internal/log"
	T "shortlink2/internal/types"
	"strings"
)

func NewCfgEnv(cfgfilename string) *T.CfgEnv {
	cfg := &T.CfgEnv{ // default vals
		SL_APP_NAME:    "shortlink2",
		SL_APP_PROTOCS: ":http:grpc",
		SL_LOG_LEVEL:   "INFO",
		SL_HTTP_IP:     "localhost",
		SL_HTTP_PORT:   ":8080",
		SL_GRPC_PORT:   ":8181",
	}
	log := L.NewLogFprintf(cfg)
	parseIpFromInterfaces(cfg, log)
	if _, err := os.Stat(cfgfilename); err == nil {
		parseFileDotEnv(cfgfilename, cfg, log)
	}
	parseOsEnvVars(cfg, log)
	return cfg
}

func parseOsEnvVars(cfg *T.CfgEnv, log T.ILog) {
	t := reflect.TypeOf(*cfg) // TODO: check reflect pointers
	v := reflect.ValueOf(cfg).Elem()
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		val := v.Field(i)
		if el, ok := os.LookupEnv(key); ok && (len(el) > 0) {
			val.SetString(el)
			log.LogDebug("OSENV %s=%s\n", key, val)
		}
	}
}

func parseFileDotEnv(filename string, cfg *T.CfgEnv, log T.ILog) {
	f, err := os.Open(filename)
	if err != nil {
		log.LogError(err, "(CfgEnv).parseFileDotEnv(): error with opening cfg file")
		return
	}
	cfgmap := make(map[string]string, 8)
	pattern := regexp.MustCompile("^[0-9A-Za-z_:]+=[0-9A-Za-z_:]+")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := pattern.FindString(scanner.Text())
		if len(str) > 0 {
			strarr := strings.Split(str, "=")
			cfgmap[strarr[0]] = strarr[1]
			log.LogDebug("CFGFILE %s=%s\n", strarr[0], strarr[1])
		}
	}
	if err := scanner.Err(); err != nil {
		log.LogError(err, "(CfgEnv): error while reading cfg file")
	}
	f.Close()

	t := reflect.TypeOf(*cfg) // TODO: check reflect pointers
	v := reflect.ValueOf(cfg).Elem()
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		val := v.Field(i)
		if el, ok := cfgmap[key]; ok && (len(el) > 0) {
			val.SetString(el)
		}
		log.LogDebug("CFGSTRUCT %s=%s\n", key, val)
	}
}

func parseIpFromInterfaces(cfg *T.CfgEnv, log T.ILog) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.LogError(err, "(CfgEnv).parseIpFromInterfaces(): error, can not get IP interface")
		return
	}
	strarr := []string{}
	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				strarr = append(strarr, ipnet.IP.String())
			}
		}
	}
	ip := strings.Join(strarr, ";")
	if len(ip) > 0 {
		cfg.SL_HTTP_IP = ip
		log.LogDebug("(CfgEnv).parseIpFromInterfaces(): %s", ip)
	}
}
