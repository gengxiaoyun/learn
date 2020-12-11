package practicessh

import(
	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/mem"
	"fmt"
	"math"
	"flag"
	"reflect"
	"bytes"
	"strings"
	"strconv"
	"os"
	"bufio"
	"io"
)

type Client struct {
	Socket string `toml:"socket"`
}

type Mysqld_multi struct {
	Mysqld string `toml:"mysqld"`
	Mysqladmin string `toml:"mysqladmin"`
	Log string `toml:"log"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type Mysqld_safe struct {
	Log_error string `toml:"log_error"`
	Pid_file string `toml:"pid-file"`
	Datadir string `toml:"datadir"`
	Socket string `toml:"socket"`
}

type Mysql struct {
	Default_character_set string `toml:"default-character-set"`
}

type Mysqld struct {
	Port string `toml:"port"`
	Log_timestamps string `toml:"log_timestamps"`
	Basedir string `toml:"basedir"`
	Datadir string `toml:"datadir"`
	Socket string `toml:"socket"`
	Log_error string `toml:"log_error"`
	Pid_file string `toml:"pid-file"`
	Secure_file_priv string `toml:"secure_file_priv"`

	Server_id string `toml:"server-id"`
	Transaction_isolation string `toml:"transaction-isolation"`
	Character_set_server string `toml:"character_set_server"`
	Open_files_limit int `toml:"open_files_limit"`
	Lower_case_table_names int `toml:"lower_case_table_names"`
	Max_connections int `toml:"max_connections"`
	Max_connect_errors string `toml:"max_connect_errors"`
	Connect_timeout int `toml:"connect_timeout"`
	Lock_wait_timeout int `toml:"lock_wait_timeout"`
	Thread_cache_size int `toml:"thread_cache_size"`
	Sql_mode string `toml:"sql_mode"`

	Plugin_load string `toml:"plugin_load"`
	Rpl_semi_sync_slave_enabled int `toml:"rpl_semi_sync_slave_enabled"`
	Rpl_semi_sync_master_wait_for_slave_count int `toml:"rpl_semi_sync_master_wait_for_slave_count"`
	Rpl_semi_sync_master_wait_no_slave int `toml:"rpl_semi_sync_master_wait_no_slave"`
	Rpl_semi_sync_master_timeout int64 `toml:"rpl_semi_sync_master_timeout"`

	Binlog_format string `toml:"binlog_format"`
	Log_bin string `toml:"log_bin"`
	Max_binlog_size string `toml:"max_binlog_size"`
	Expire_logs_days int `toml:"expire_logs_days"`
	Binlog_error_action string `toml:"binlog_error_action"`

	Log_slave_updates int `toml:"log_slave_updates"`
	Relay_log string `toml:"relay_log"`
	Max_relay_log_size string `toml:"max_relay_log_size"`
	Relay_log_purge int `toml:"relay_log_purge"`

	Master_info_repository string `toml:"master_info_repository"`
	Relay_log_info_repository string `toml:"relay_log_info_repository"`
	Relay_log_recovery int `toml:"relay_log_recovery"`
	Report_host string `toml:"report_host"`
	Report_port string `toml:"report_port"`

	Sync_binlog int `toml:"sync_binlog"`
	Innodb_flush_log_at_trx_commit int `toml:"innodb_flush_log_at_trx_commit"`

	Innodb_buffer_pool_size string `toml:"innodb_buffer_pool_size"`
	Innodb_sort_buffer_size string `toml:"innodb_sort_buffer_size"`
	Innodb_log_buffer_size string `toml:"innodb_log_buffer_size"`
	Innodb_lock_wait_timeout int `toml:"innodb_lock_wait_timeout"`
	Innodb_log_file_size string `toml:"innodb_log_file_size"`
	Innodb_log_files_in_group int `toml:"innodb_log_files_in_group"`
	Innodb_io_capacity int `toml:"innodb_io_capacity"`
	Innodb_io_capacity_max int `toml:"innodb_io_capacity_max"`
	Innodb_file_per_table int `toml:"innodb_file_per_table"`
	Innodb_stats_persistent_sample_pages int `toml:"innodb_stats_persistent_sample_pages"`
	Innodb_online_alter_log_max_size string `toml:"innodb_online_alter_log_max_size"`
	Innodb_thread_concurrency int `toml:"innodb_thread_concurrency"`
	Innodb_write_io_threads int `toml:"innodb_write_io_threads"`
	Innodb_read_io_threads int `toml:"innodb_read_io_threads"`
	Innodb_page_cleaners int `toml:"innodb_page_cleaners"`
	Innodb_flush_method string `toml:"innodb_flush_method"`
	Innodb_monitor_enable string `toml:"innodb_monitor_enable"`
	Innodb_print_all_deadlocks int `toml:"innodb_print_all_deadlocks"`

	Gtid_mode string `toml:"gtid_mode"`
	Enforce_gtid_consistency int `toml:"enforce_gtid_consistency"`
	Slave_parallel_type string `toml:"slave-parallel-type"`
	Slave_parallel_workers int `toml:"slave-parallel-workers"`
	Slave_preserve_commit_order int `toml:"slave_preserve_commit_order"`
	Slave_transaction_retries int `toml:"slave_transaction_retries"`
	Binlog_gtid_simple_recovery int `toml:"binlog_gtid_simple_recovery"`

	Loose_innodb_numa_interleave int `toml:"loose_innodb_numa_interleave"`
	Innodb_buffer_pool_dump_pct int `toml:"innodb_buffer_pool_dump_pct"`
	Innodb_undo_logs int `toml:"innodb_undo_logs"`
	Innodb_undo_log_truncate int `toml:"innodb_undo_log_truncate"`
	Innodb_max_undo_log_size string `toml:"innodb_max_undo_log_size"`
	Innodb_purge_rseg_truncate_frequency int `toml:"innodb_purge_rseg_truncate_frequency"`

	Max_allowed_packet string `toml:"max_allowed_packet"`
	Table_open_cache int `toml:"table_open_cache"`
	Tmp_table_size string `toml:"tmp_table_size"`
	Max_heap_table_size string `toml:"max_heap_table_size"`
	Sort_buffer_size string `toml:"sort_buffer_size"`
	Join_buffer_size string `toml:"join_buffer_size"`
	Read_buffer_size string `toml:"read_buffer_size"`
	Read_rnd_buffer_size string `toml:"read_rnd_buffer_size"`
	Key_buffer_size string `toml:"key_buffer_size"`
	Bulk_insert_buffer_size string `toml:"bulk_insert_buffer_size"`
	Binlog_cache_size string `toml:"binlog_cache_size"`

	Slow_query_log_file string `toml:"slow_query_log_file"`
	Slow_query_log string `toml:"slow_query_log"`
	Long_query_time float64 `toml:"long_query_time"`
	Log_output string `toml:"log_output"`

	Performance_schema string `toml:"performance_schema"`
	Performance_schema_instrument string `toml:"performance-schema-instrument"`

	Symbolic_links int `toml:"symbolic-links"`
}

type Config struct {
	Client *Client `toml:"client"`
	Mysqld_multi *Mysqld_multi `toml:"mysqld_multi"`
	Mysqld_safe *Mysqld_safe `toml:"mysqld_safe"`
	Mysql *Mysql `toml:"mysql"`
	Mysqld *Mysqld `toml:"mysqld"`
}

var (
	err error
	//address string
	//user string
	//pass string

	confPath string
	Conf = &Config{}

	m = "skip-host-cache"
	n = "skip-name-resolve"

	sText = "plugin_load=rpl_semi_sync_master=semi_sync_master.so;rpl_semi_sync_slave=semi_slave.so"
	news= "plugin_load=\"rpl_semi_sync_master=semi_sync_master.so;rpl_semi_sync_slave=semi_slave.so\""
	pText = "[mysqld]"

	filename = "./my.cnf"
)

func init() {
	//flag.StringVar(&address,"address","192.168.186.132:3306","set ip and port")
	//flag.StringVar(&user,"user","root","set username")
	//flag.StringVar(&pass,"pass","root","set password")

	flag.StringVar(&confPath, "conf", "cnf.toml", "-conf path")
}

// Init init conf
func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func FlagCommand(address string) [][]string{
	str := strings.Split(address,",")
	a := len(str)
	arr := make([][]string,a)
	for i:=0;i<a;i++{
		arr[i] = make([]string,2)
	}
	for i:=0;i<a;i++ {
		fmt.Println(str[i])
		newStr := strings.Split(str[i], ":")
		ip := newStr[0]
		port := newStr[1]
		arr[i][0] = ip
		arr[i][1] = port
	}
	fmt.Println(arr,len(arr))
	return arr
}


func GetMemPercent() (string,string,string,error) {
	memInfo,err := mem.VirtualMemory()
	if err != nil{
		return "","","",err
	}
	var a,b int
	a = int(math.Floor(float64(memInfo.Total)/float64(1024*1024*1024)+0.5))
	switch a {
	case 1:
		b = 8
	case 2:
		b = 16
	case 3:
		b = 32
	default:
		b = 64
	}
	c := strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Total) * 0.5 /float64(1024*1024)))))))
	d := strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Available) * 0.3 / float64(1024*1024)))))))
	return strconv.Itoa(b),c,d,nil
}


func SetValueToStruct(report_host,port,b,c,d string) *Config {
	socket := "/mysqldata/mysql" + port + "/mysql.sock"
	log_error := "/mysqldata/mysql" + port + "/log/mysqld.log"
	pid_file := "/mysqldata/mysql" + port + "/mysql.pid"
	datadir := "/mysqldata/mysql" + port + "/data"
	server_id := port + strings.Split(report_host, ".")[2]+strings.Split(report_host, ".")[3]
	log_bin := "/mysqllog/mysql" + port + "/binlog/mysql-bin"
	relay_log := "/mysqllog/mysql" + port + "/relaylog/mysql-relay"
	p,_ := strconv.Atoi(port)
	report_port := strconv.Itoa(p + 1)
	thread_cache_size,_ := strconv.Atoi(b)
	innodb_buffer_pool_size := c + "M"
	key_buffer_size := d + "M"
	slow_query_log_file := "/mysqldata/mysql" + port + "/log/mysql-slow.log"

	v := reflect.ValueOf(Conf.Client).Elem()
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))

	v = reflect.ValueOf(Conf.Mysqld_safe).Elem()
	v.FieldByName("Log_error").Set(reflect.ValueOf(log_error))
	v.FieldByName("Pid_file").Set(reflect.ValueOf(pid_file))
	v.FieldByName("Datadir").Set(reflect.ValueOf(datadir))
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))

	v = reflect.ValueOf(Conf.Mysqld).Elem()
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))
	v.FieldByName("Log_error").Set(reflect.ValueOf(log_error))
	v.FieldByName("Pid_file").Set(reflect.ValueOf(pid_file))
	v.FieldByName("Datadir").Set(reflect.ValueOf(datadir))
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))
	v.FieldByName("Port").Set(reflect.ValueOf(port))
	v.FieldByName("Report_host").Set(reflect.ValueOf(report_host))
	v.FieldByName("Server_id").Set(reflect.ValueOf(server_id))
	v.FieldByName("Log_bin").Set(reflect.ValueOf(log_bin))
	v.FieldByName("Relay_log").Set(reflect.ValueOf(relay_log))
	v.FieldByName("Report_port").Set(reflect.ValueOf(report_port))
	v.FieldByName("Thread_cache_size").Set(reflect.ValueOf(thread_cache_size))
	v.FieldByName("Innodb_buffer_pool_size").Set(reflect.ValueOf(innodb_buffer_pool_size))
	v.FieldByName("Key_buffer_size").Set(reflect.ValueOf(key_buffer_size))
	v.FieldByName("Slow_query_log_file").Set(reflect.ValueOf(slow_query_log_file))

	return Conf
}

func checkFileIsExist(file string) bool {
	var exist = true
	if _, err := os.Stat(file); os.IsNotExist(err) {
		exist = false
	}
	return exist
}


func ChangeConfFile(file,val string,) error {
	f,err := os.Open(filename)
	if err != nil{
		return err
	}
	defer f.Close()
	out,err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil{
		return err
	}
	defer out.Close()
	buf := bufio.NewReader(f)
	newline := ""
	for {
		line,_,err := buf.ReadLine()
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
		newline = string(line)
		if newline == sText{
			newline = strings.Replace(newline,sText,news,1)
		}
		if newline == pText{
			newline = strings.Replace(newline,pText,"[mysqld"+val+"]",1)
		}

		_,err = out.WriteString(newline+"\n")
		if err != nil{
			return err
		}
	}
	err = os.Remove(filename)
	if err != nil{
		return err
	}
	err = os.Rename(file,filename)
	if err != nil{
		return err
	}

	return nil
}

func Flex(address string) ([][]string,error){
	arr := FlagCommand(address)
	var f *os.File
	if checkFileIsExist(filename) {
		err = os.Remove(filename)
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	} else {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	}
	if err != nil {
		return arr,err
	}
	defer f.Close()

	b,c,d,err := GetMemPercent()
	if err != nil{
		return arr,err
	}

	for i:=0;i<len(arr);i++{
		var str = "[mysqld"+arr[i][1]+"]"
		if i==0 {
			p := SetValueToStruct(arr[i][0],arr[i][1],b,c,d)
			fmt.Println("====================")
			var buffer bytes.Buffer
			encoder := toml.NewEncoder(&buffer)
			encoder.Encode(p)

			_, err = f.WriteString(strings.Replace(strings.Replace(buffer.String()+"\n"+m+"\n"+n+"\n\n"," ","",-1),"\"","",-1)) //写入文件(字符串)
			if err != nil {
				return arr,err
			}

		} else{
			p := SetValueToStruct(arr[i][0],arr[i][1],b,c,d)
			var buffer bytes.Buffer
			encoder := toml.NewEncoder(&buffer)
			encoder.Encode(p.Mysqld)

			_, err = f.WriteString(strings.Replace(strings.Replace(str+"\n"+buffer.String()+"\n"+m+"\n"+n+"\n\n"," ","",-1),"\"","",-1)) //写入文件(字符串)
			if err != nil {
				return arr,err
			}
		}
	}
	val := arr[0][1]
	err = ChangeConfFile("./my001.cnf",val)
	if err != nil {
		return arr,err
	}

	return arr,nil
}



//func main() {
//	if err = Init(); err != nil {
//		//log.Printf("conf.Init() err:%+v", err)
//		fmt.Println("failed")
//	}
//
//	arr := FlagCommand()
//	val := arr[0][1]
//	err = flex(arr)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	err = ChangeConfFile(filename,"./testout0002.txt",val)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
