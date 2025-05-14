package cfg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

func NewCfgEnvMap(dir, file string) *CfgEnvMap {
	vals := make(map[string]string, 4)
	vals[T.SL_APP_NAME] = file
	vals[T.SL_LOG_LEVEL] = "INFO" // LOG levels: TRACE, DEBUG, INFO, WARN, ERROR, PANIC, FATAL, NOLOG(default if empty or mess)
	vals[T.SL_HTTP_IP] = "localhost"
	vals[T.SL_HTTP_PORT] = ":8080"
	return &CfgEnvMap{
		vals:  vals,
		fname: filepath.Join(dir, file, ".env"),
	}
}

func (c *CfgEnvMap) Parse() T.ICfg {
	log := L.NewLogFprintf(c)
	c.parseIpFromInterface(log)
	if len(c.fname) != 0 {
		if _, err := os.Stat(c.fname); err == nil {
			c.parseFileDotEnvVars(log)
		}
	}
	c.parseOsEnvVars(log)
	return c
}

func (c *CfgEnvMap) GetVal(key string) string {
	val, _ := c.vals[key]
	return val
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
		log.LogError(fmt.Errorf("%s: %w", "(CfgEnvMap).parseFileDotEnvVars(): error while opening cfg file", err))
		return
	}
	log.LogDebug("load config from file: %s", c.fname)
	defer f.Close()

	pattern := regexp.MustCompile("^[0-9A-Za-z_]+=[0-9A-Za-z_:/.]+")
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

		log.LogError(fmt.Errorf("%s: %w", "(CfgEnvMap).parseFileDotEnvVars(): error while reading cfg file", err))
	}
}

func (c *CfgEnvMap) parseIpFromInterface(log T.ILog) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.LogError(fmt.Errorf("%s: %w", "CfgEnvMap.parseIpFromInterface(): error while getting IP interface", err))
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
