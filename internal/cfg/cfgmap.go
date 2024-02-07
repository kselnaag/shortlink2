package cfg

import (
	"bufio"
	"net"
	"os"
	"regexp"
	L "shortlink2/internal/log"
	T "shortlink2/internal/types"
	"strings"
)

var _ T.ICfg = (*CfgEnvMap)(nil)

type CfgEnvMap struct {
	vals  map[string]string
	fname string
}

func NewCfgEnvMap(cfgfilename string) *CfgEnvMap {
	vals := make(map[string]string, 6) // default vals
	vals[T.SL_APP_NAME] = "shortlink2"
	vals[T.SL_APP_PROTOCS] = ":http:grpc"
	vals[T.SL_LOG_LEVEL] = "INFO"
	vals[T.SL_HTTP_IP] = "localhost"
	vals[T.SL_HTTP_PORT] = ":8080"
	vals[T.SL_GRPC_PORT] = ":8082"
	return &CfgEnvMap{
		vals:  vals,
		fname: cfgfilename,
	}
}

func (c *CfgEnvMap) Parse() T.ICfg {
	log := L.NewLogFprintf(c)
	c.parseIpFromInterface(log)
	if len(c.fname) != 0 {
		c.parseFileDotEnvVars(log)
	}
	c.parseOsEnvVars(log)
	return c
}

func (c *CfgEnvMap) GetVal(key string) string {
	if val, ok := c.vals[key]; ok {
		return val
	}
	return ""
}

func (c *CfgEnvMap) parseOsEnvVars(log T.ILog) {
	for key := range c.vals {
		if v, ok := os.LookupEnv(key); ok && (len(v) > 0) {
			c.vals[key] = v
			log.LogDebug("OSENV %s=%s\n", key, v)
		}
	}
}

func (c *CfgEnvMap) parseFileDotEnvVars(log T.ILog) {
	f, err := os.Open(c.fname)
	if err != nil {
		log.LogError(err, "(CfgEnvMap).parseFileDotEnvVars(): error while opening cfg file")
		return
	}
	defer f.Close()
	pattern := regexp.MustCompile("^[0-9A-Za-z_:]+=[0-9A-Za-z_:]+")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := pattern.FindString(scanner.Text())
		if len(str) > 0 {
			strarr := strings.Split(str, "=")
			if _, ok := c.vals[strarr[0]]; ok {
				c.vals[strarr[0]] = strarr[1]
				log.LogDebug("CFGFILE %s=%s\n", strarr[0], strarr[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.LogError(err, "(CfgEnvMap).parseFileDotEnvVars(): error while reading cfg file")
	}
}

func (c *CfgEnvMap) parseIpFromInterface(log T.ILog) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.LogError(err, "(CfgEnvMap).parseIpFromInterface(): error while getting IP interface")
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
		c.vals[T.SL_HTTP_IP] = ip
		log.LogDebug("(CfgEnvMap).parseIpFromInterface(): %s", ip)
	}
}
