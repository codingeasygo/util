package xprop

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
)

//Config is parser for properties file
type Config struct {
	xmap.M
	ShowLog bool
	sec     string
	Lines   []string
	Seces   []string
	SecLn   map[string]int
	Base    string
	Masks   map[string]string
}

//NewConfig will return new config
func NewConfig() (config *Config) {
	config = &Config{
		M:       xmap.M{},
		ShowLog: true,
		SecLn:   map[string]int{},
		Masks:   map[string]string{},
	}
	return
}

//LoadConf will load new config by uri
func LoadConf(uri string) (config *Config, err error) {
	config, err = LoadConfWait(uri, true)
	return
}

//LoadConfWait will load new config by uri
func LoadConfWait(uri string, wait bool) (config *Config, err error) {
	config = NewConfig()
	err = config.LoadWait(uri, wait)
	return
}

func (c *Config) slog(fs string, args ...interface{}) {
	if c.ShowLog {
		fmt.Println(fmt.Sprintf(fs, args...))
	}
}

func (c *Config) autoPath(path ...string) (all []string) {
	all = path
	for _, p := range path {
		if strings.Contains(p, "/") {
			continue
		}
		all = append(all, "loc/"+p)
	}
	return
}

func (c *Config) SetValue(path string, val interface{}) (err error) {
	c.M[path] = val
	return
}

func (c *Config) IntDef(def int, path ...string) int {
	return c.M.IntDef(def, c.autoPath(path...)...)
}

func (c *Config) Int64Def(def int64, path ...string) (v int64) {
	return c.M.Int64Def(def, c.autoPath(path...)...)
}

func (c *Config) Uint64Def(def uint64, path ...string) (v uint64) {
	return c.M.Uint64Def(def, c.autoPath(path...)...)
}

func (c *Config) Float64Def(def float64, path ...string) (v float64) {
	return c.M.Float64Def(def, c.autoPath(path...)...)
}

func (c *Config) StrDef(def string, path ...string) (v string) {
	return c.M.StrDef(def, c.autoPath(path...)...)
}

func (c *Config) MapDef(def xmap.M, path ...string) (v xmap.M) {
	v, err := c.MapVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) MapVal(path ...string) (v xmap.M, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		v = xmap.M{}
		err = json.Unmarshal([]byte(data), &v)
	}
	return
}

func (c *Config) ArrayMapDef(def []xmap.M, path ...string) (v []xmap.M) {
	v, err := c.ArrayMapVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayMapVal(path ...string) (v []xmap.M, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		err = json.Unmarshal([]byte(data), &v)
	}
	return
}

func (c *Config) ArrayStrDef(def []string, path ...string) (v []string) {
	v, err := c.ArrayStrVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayStrVal(path ...string) (v []string, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		v = strings.Split(data, ",")
	}
	return
}

func (c *Config) ArrayIntDef(def []int, path ...string) (v []int) {
	v, err := c.ArrayIntVal(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayIntVal(path ...string) (v []int, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		v, err = converter.ArrayIntVal(strings.Split(data, ","))
	}
	return
}

func (c *Config) ArrayInt64Def(def []int64, path ...string) (v []int64) {
	v, err := c.ArrayInt64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayInt64Val(path ...string) (v []int64, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		v, err = converter.ArrayInt64Val(strings.Split(data, ","))
	}
	return
}

func (c *Config) ArrayUint64Def(def []uint64, path ...string) (v []uint64) {
	v, err := c.ArrayUint64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayUint64Val(path ...string) (v []uint64, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		v, err = converter.ArrayUint64Val(strings.Split(data, ","))
	}
	return
}

func (c *Config) ArrayFloat64Def(def []float64, path ...string) (v []float64) {
	v, err := c.ArrayFloat64Val(path...)
	if err != nil {
		v = def
	}
	return
}

func (c *Config) ArrayFloat64Val(path ...string) (v []float64, err error) {
	data, err := c.M.StrVal(c.autoPath(path...)...)
	if err == nil {
		data = strings.TrimSpace(data)
		v, err = converter.ArrayFloat64Val(strings.Split(data, ","))
	}
	return
}

//FileModeDef read file mode value
func (c *Config) FileModeDef(def os.FileMode, path ...string) (mode os.FileMode) {
	mode, err := c.FileModeVal(path...)
	if err != nil {
		mode = def
	}
	return
}

//FileModeVal read file mode value
func (c *Config) FileModeVal(path ...string) (mode os.FileMode, err error) {
	data, err := c.StrVal(path...)
	if err != nil {
		return
	}
	data = strings.TrimSpace(data)
	val, err := strconv.ParseUint(data, 8, 32)
	if err != nil {
		return
	}
	mode = os.FileMode(val)
	return
}

//Print all configure
func (c *Config) Print() {
	fmt.Println(c.String())
}

//PrintSection print session to stdout
func (c *Config) PrintSection(section string) {
	mask := map[*regexp.Regexp]*regexp.Regexp{}
	for k, v := range c.Masks {
		mask[regexp.MustCompile(k)] = regexp.MustCompile(v)
	}
	sdata := ""
	for k, v := range c.M {
		if !strings.HasPrefix(k, section) {
			continue
		}
		val := fmt.Sprintf("%v", v)
		for maskKey, maskVal := range mask {
			if maskKey.MatchString(k) {
				val = maskVal.ReplaceAllString(val, "***")
			}
		}
		val = strings.Replace(val, "\n", "\\n", -1)
		sdata = fmt.Sprintf("%v %v=%v\n", sdata, k, val)
	}
	fmt.Println(sdata)
}

//Range the section key-value by callback
func (c *Config) Range(section string, callback func(key string, val interface{})) {
	for k, v := range c.M {
		if strings.HasPrefix(k, section) {
			callback(strings.TrimPrefix(k, section+"/"), v)
		}
	}
}

func (c *Config) exec(base, line string, wait bool) error {
	ps := strings.Split(line, "#")
	if len(ps) < 1 || len(ps[0]) < 1 {
		return nil
	}
	line = strings.TrimSpace(ps[0])
	if len(line) < 1 {
		return nil
	}
	if regexp.MustCompile("^\\[[^\\]]*\\][\t ]*$").MatchString(line) {
		sec := strings.Trim(line, "\t []")
		c.sec = sec + "/"
		c.Seces = append(c.Seces, sec)
		c.SecLn[sec] = len(c.Lines)
		return nil
	}
	if !strings.HasPrefix(line, "@") {
		ps = strings.SplitN(line, "=", 2)
		if len(ps) < 2 {
			c.slog("not value key found:%v", ps[0])
		} else {
			key := c.sec + c.EnvReplace(strings.Trim(ps[0], " "))
			val := c.EnvReplace(strings.Trim(ps[1], " "))
			c.M[key] = val
		}
		return nil
	}
	line = strings.TrimPrefix(line, "@")
	ps = strings.SplitN(line, ":", 2)
	if len(ps) < 2 || len(ps[1]) < 1 {
		c.slog("%v", c.EnvReplace(line))
		return nil
	}
	ps[0] = strings.ToLower(strings.Trim(ps[0], " \t"))
	ps[0] = c.EnvReplace(ps[0])
	if ps[0] == "l" {
		ps[1] = strings.Trim(ps[1], " \t")
		return c.load(base, ps[1], wait)
	}
	if cs := strings.SplitN(ps[0], "==", 2); len(cs) == 2 {
		if cs[0] == cs[1] {
			return c.exec(base, ps[1], wait)
		}
		return nil
	}
	if cs := strings.SplitN(ps[0], "!=", 2); len(cs) == 2 {
		if cs[0] != cs[1] {
			return c.exec(base, ps[1], wait)
		}
		return nil
	}
	//all other will print line.
	c.slog("%v", c.EnvReplace(line))
	return nil
}

func (c *Config) load(base, line string, wait bool) error {
	line = c.EnvReplaceEmpty(line, true)
	line = strings.Trim(line, "\t ")
	if len(line) < 1 {
		return nil
	}
	if !(strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") || filepath.IsAbs(line)) {
		line = path.Join(base, line)
	}
	config := NewConfig()
	err := config.LoadWait(line, wait)
	if err == nil {
		c.Merge(config)
	}
	return err
}

//Load will load config by uri
func (c *Config) Load(uri string) error {
	return c.LoadWait(uri, false)
}

//LoadWait will load config by uri and wait when uri is not found
func (c *Config) LoadWait(uri string, wait bool) error {
	if strings.HasPrefix(uri, "http://") {
		return c.LoadWebWait(uri, wait)
	} else if strings.HasPrefix(uri, "https://") {
		return c.LoadWebWait(uri, wait)
	} else if strings.HasPrefix(uri, "data:text/conf,") {
		return c.LoadConfString(strings.TrimPrefix(uri, "data:text/conf,"))
	} else if strings.HasPrefix(uri, "data:text/prop,") {
		return c.LoadPropStringWait(strings.TrimPrefix(uri, "data:text/prop,"), wait)
	} else {
		return c.LoadFileWait(uri, wait)
	}
}

//LoadFile will load the configure by .properties file.
func (c *Config) LoadFile(filename string) error {
	return c.LoadFileWait(filename, true)
}

//LoadFileWait will load the configure by .properties file and wait when file not exist
func (c *Config) LoadFileWait(filename string, wait bool) error {
	c.slog("loading local configure->%v", filename)
	var parts = strings.Split(filename, "?")
	filename = parts[0]
	if len(parts) > 1 {
		furl, err := url.Parse("/?" + parts[1])
		if err == nil {
			query := furl.Query()
			for k := range query {
				c.M.SetValue(k, query.Get(k))
			}
		}
	}
	var delay = 50 * time.Millisecond
	for {
		_, xerr := os.Stat(filename)
		if xerr == nil {
			break
		}
		if wait {
			c.slog("file(%v) not found", filename)
			if delay < 2*time.Second {
				delay *= 2
			}
			time.Sleep(delay)
			continue
		} else {
			return fmt.Errorf("file(%v) not found", filename)
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	dir, _ := filepath.Split(filename)
	if len(dir) < 1 {
		dir = "."
	}
	dir, _ = filepath.Abs(dir)
	if strings.HasSuffix(filename, ".conf") {
		return c.LoadConfReader(dir, file)
	}
	return c.LoadPropReaderWait(dir, file, wait)
}

//LoadPropReader will load properties config by reader
func (c *Config) LoadPropReader(base string, reader io.Reader) error {
	return c.LoadPropReaderWait(base, reader, true)
}

//LoadPropReaderWait will load properties config by reader
func (c *Config) LoadPropReaderWait(base string, reader io.Reader, wait bool) error {
	if len(base) > 0 {
		c.Base = base
	}
	buf := bufio.NewReaderSize(reader, 64*1024)
	for {
		//read one line
		line, err := readLine(buf)
		if err != nil {
			break
		}
		c.Lines = append(c.Lines, line)
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}
		err = c.exec(base, line, wait)
		if err != nil {
			return err
		}
	}
	return nil
}

//LoadConfReader will load conf config by reader
func (c *Config) LoadConfReader(base string, reader io.Reader) error {
	var key, val string
	buf := bufio.NewReaderSize(reader, 64*1024)
	for {
		//read one line
		line, err := readLine(buf)
		if err != nil {
			if len(key) > 0 {
				c.M[key] = strings.Trim(val, "\n")
				key, val = "", ""
			}
			break
		}
		if regexp.MustCompile("^\\[[^\\]]*\\][\t ]*$").MatchString(line) {
			sec := strings.Trim(line, "\t []")
			if len(key) > 0 {
				c.M[key] = strings.Trim(val, "\n")
				key, val = "", ""
			}
			key = sec
		} else {
			val += line + "\n"
		}
	}
	return nil
}

func (c *Config) webGet(url string) (data string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code(%v)", res.StatusCode)
		return
	}
	bys, err := ioutil.ReadAll(res.Body)
	if err == nil {
		data = string(bys)
	}
	return
}

//LoadWeb will load the configure by network .properties URL.
func (c *Config) LoadWeb(remote string) error {
	return c.LoadWebWait(remote, true)
}

//LoadWebWait will load the configure by network .properties URL.
func (c *Config) LoadWebWait(remote string, wait bool) (err error) {
	c.slog("loading remote configure->%v", remote)
	var data string
	var delay = 50 * time.Millisecond
	for {
		data, err = c.webGet(remote)
		if err == nil {
			c.slog("loading remote configure(%v) success", remote)
			break
		}
		c.slog("loading remote configure(%v):%v", remote, err.Error())
		if wait {
			if delay < 2*time.Second {
				delay *= 2
			}
			time.Sleep(delay)
			continue
		} else {
			break
		}
	}
	if err != nil {
		return
	}
	var filename string
	rurl, _ := url.Parse(remote)
	rurl.Path, filename = path.Split(rurl.Path)
	if strings.HasSuffix(filename, ".conf") {
		return c.LoadConfReader(rurl.RequestURI(), bytes.NewBufferString(data))
	}
	return c.LoadPropReaderWait(rurl.RequestURI(), bytes.NewBufferString(data), wait)
}

//LoadPropString will load properties config by string config
func (c *Config) LoadPropString(data string) error {
	return c.LoadPropStringWait(data, true)
}

//LoadPropStringWait will load properties config by string config
func (c *Config) LoadPropStringWait(data string, wait bool) error {
	return c.LoadPropReaderWait("", bytes.NewBufferString(data), wait)
}

//LoadConfString will load conf config by string config
func (c *Config) LoadConfString(data string) error {
	return c.LoadConfReader("", bytes.NewBufferString(data))
}

//EnvReplace will replace tartget patter by ${key} with value in configure map or system environment value.
func (c *Config) EnvReplace(val string) string {
	return c.EnvReplaceEmpty(val, false)
}

//EnvReplaceEmpty will replace tartget patter by ${key} with value in configure map or system environment value.
func (c *Config) EnvReplaceEmpty(val string, empty bool) string {
	reg := regexp.MustCompile("\\$\\{[^\\}]*\\}")
	var rval string
	val = reg.ReplaceAllStringFunc(val, func(m string) string {
		keys := strings.Split(strings.Trim(m, "${}\t "), ",")
		for _, key := range keys {
			if c.Exist(key) {
				rval = c.Str(key)
			} else if key == "C_PWD" {
				rval = c.Base
			} else {
				rval = os.Getenv(key)
			}
			if len(rval) > 0 {
				break
			}
		}
		if len(rval) > 0 {
			return rval
		}
		if empty {
			return ""
		}
		return m
	})
	return val
}

//Merge merge another configure.
func (c *Config) Merge(config *Config) {
	if config == nil {
		return
	}
	for k, v := range config.M {
		c.M[k] = v
	}
	for _, s := range config.Seces {
		if _, ok := c.SecLn[s]; ok {
			continue
		}
		c.Seces = append(c.Seces, s)
	}
}

//MergeSection merge section on another configure
func (c *Config) MergeSection(section string, config *Config) {
	for k, v := range config.M {
		if strings.HasPrefix(k, section) {
			continue
		}
		c.M[k] = v
	}
	if _, ok := c.SecLn[section]; !ok {
		c.Seces = append(c.Seces, section)
	}
}

//Clone will clone the configure
func (c *Config) Clone() (conf *Config) {
	conf = NewConfig()
	for k, v := range c.M {
		conf.M[k] = v
	}
	conf.ShowLog = c.ShowLog
	conf.sec = c.sec
	conf.Lines = append(conf.Lines, c.Lines...)
	conf.Seces = append(conf.Seces, c.Seces...)
	for k, v := range c.SecLn {
		conf.SecLn[k] = v
	}
	conf.Base = c.Base
	for k, v := range c.Masks {
		conf.Masks[k] = v
	}
	return
}

// //Strip will strip one section to new Config
// func (c *Config) Strip(section string) (config *Config) {
// 	config = NewConfig()
// 	for k, v := range c.M {
// 		if !strings.HasPrefix(k, section) {
// 			continue
// 		}
// 		config.M["loc"+strings.TrimPrefix(k, section)] = v
// 	}
// 	return
// }

func (c *Config) String() string {
	mask := map[*regexp.Regexp]*regexp.Regexp{}
	for k, v := range c.Masks {
		mask[regexp.MustCompile(k)] = regexp.MustCompile(v)
	}
	buf := bytes.NewBuffer(nil)
	keys, locs := []string{}, []string{}
	for k := range c.M {
		if strings.HasPrefix(k, "loc/") {
			locs = append(locs, k)
		} else {
			keys = append(keys, k)
		}
	}
	sort.Sort(sort.StringSlice(keys))
	for _, k := range keys {
		val := fmt.Sprintf("%v", c.M[k])
		for maskKey, maskVal := range mask {
			if maskKey.MatchString(k) {
				val = maskVal.ReplaceAllString(val, "***")
			}
		}
		val = strings.Replace(val, "\n", "\\n", -1)
		buf.WriteString(fmt.Sprintf("%v=%v\n", k, val))
	}
	for _, k := range locs {
		val := fmt.Sprintf("%v", c.M[k])
		for maskKey, maskVal := range mask {
			if maskKey.MatchString(k) {
				val = maskVal.ReplaceAllString(val, "***")
			}
		}
		val = strings.Replace(val, "\n", "\\n", -1)
		buf.WriteString(fmt.Sprintf("%v=%v\n", k, val))
	}
	return buf.String()
}

// func (c *Config) Store(sec, fp, tsec string) error {
// 	var seci int = -1
// 	for idx, s := range c.Seces {
// 		if s == sec {
// 			seci = idx
// 		}
// 	}
// 	if seci < 0 {
// 		return fmt.Errorf("section not found by %v", sec)
// 	}
// 	var beg, end int = c.SecLn[sec], len(c.Lines)
// 	if seci < len(c.Seces)-1 {
// 		end = c.SecLn[c.Seces[seci+1]] - 1
// 	}
// 	tf, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	defer tf.Close()
// 	buf := bufio.NewWriter(tf)
// 	buc.WriteString("[" + tsec + "]\n")
// 	for i := beg; i < end; i++ {
// 		buc.WriteString(c.Lines[i])
// 		buc.WriteString("\n")
// 	}
// 	return buc.Flush()
// }

func readLine(buf *bufio.Reader) (line string, err error) {
	var bys []byte
	var prefix bool
	for {
		bys, prefix, err = buf.ReadLine()
		if err != nil {
			break
		}
		line += string(bys)
		if !prefix {
			break
		}
	}
	return
}
