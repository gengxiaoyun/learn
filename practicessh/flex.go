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

type MysqlDMulti struct {
	MysqlD string `toml:"mysqld"`
	Mysqladmin string `toml:"mysqladmin"`
	Log string `toml:"log"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type Mysql struct {
	DefaultCharacterSet string `toml:"default-character-set"`
}

type MysqlD struct {
	Port string `toml:"port"`
	LogTimestamps string `toml:"log_timestamps"`
	BaseDir string `toml:"basedir"`
	DataDir string `toml:"datadir"`
	Socket string `toml:"socket"`
	LogError string `toml:"log_error"`
	PidFile string `toml:"pid-file"`
	SecureFilePriV string `toml:"secure_file_priv"`

	ServerId string `toml:"server-id"`
	TransactionIsolation string `toml:"transaction-isolation"`
	CharacterSetServer string `toml:"character_set_server"`
	OpenFilesLimit int `toml:"open_files_limit"`
	LowerCaseTableNames int `toml:"lower_case_table_names"`
	MaxConnections int `toml:"max_connections"`
	MaxConnectErrors string `toml:"max_connect_errors"`
	ConnectTimeout int `toml:"connect_timeout"`
	LockWaitTimeout int `toml:"lock_wait_timeout"`
	ThreadCacheSize int `toml:"thread_cache_size"`
	SqlMode string `toml:"sql_mode"`

	PluginLoad string `toml:"plugin_load"`
	RplSemiSyncSlaveEnabled int `toml:"rpl_semi_sync_slave_enabled"`
	RplSemiSyncMasterWaitForSlaveCount int `toml:"rpl_semi_sync_master_wait_for_slave_count"`
	RplSemiSyncMasterWaitNoSlave int `toml:"rpl_semi_sync_master_wait_no_slave"`
	RplSemiSyncMasterTimeout int64 `toml:"rpl_semi_sync_master_timeout"`

	BinlogFormat string `toml:"binlog_format"`
	LogBin string `toml:"log_bin"`
	MaxBinlogSize string `toml:"max_binlog_size"`
	ExpireLogsDays int `toml:"expire_logs_days"`
	BinlogErrorAction string `toml:"binlog_error_action"`

	LogSlaveUpdates int `toml:"log_slave_updates"`
	RelayLog string `toml:"relay_log"`
	MaxRelayLogSize string `toml:"max_relay_log_size"`
	RelayLogPurge int `toml:"relay_log_purge"`

	MasterInfoRepository string `toml:"master_info_repository"`
	RelayLogInfoRepository string `toml:"relay_log_info_repository"`
	RelayLogRecovery int `toml:"relay_log_recovery"`
	ReportHost string `toml:"report_host"`
	ReportPort string `toml:"report_port"`

	SyncBinlog int `toml:"sync_binlog"`
	InnodbFlushLogAtTrxCommit int `toml:"innodb_flush_log_at_trx_commit"`

	InnodbBufferPoolSize string `toml:"innodb_buffer_pool_size"`
	InnodbSortBufferSize string `toml:"innodb_sort_buffer_size"`
	InnodbLogBufferSize string `toml:"innodb_log_buffer_size"`
	InnodbLockWaitTimeout int `toml:"innodb_lock_wait_timeout"`
	InnodbLogFileSize string `toml:"innodb_log_file_size"`
	InnodbLogFilesInGroup int `toml:"innodb_log_files_in_group"`
	InnodbIoCapacity int `toml:"innodb_io_capacity"`
	InnodbIoCapacityMax int `toml:"innodb_io_capacity_max"`
	InnodbFilePerTable int `toml:"innodb_file_per_table"`
	InnodbStatsPersistentSamplePages int `toml:"innodb_stats_persistent_sample_pages"`
	InnodbOnlineAlterLogMaxSize string `toml:"innodb_online_alter_log_max_size"`
	InnodbThreadConcurrency int `toml:"innodb_thread_concurrency"`
	InnodbWriteIoThreads int `toml:"innodb_write_io_threads"`
	InnodbReadIoThreads int `toml:"innodb_read_io_threads"`
	InnodbPageCleaners int `toml:"innodb_page_cleaners"`
	InnodbFlushMethod string `toml:"innodb_flush_method"`
	InnodbMonitorEnable string `toml:"innodb_monitor_enable"`
	InnodbPrintAllDeadlocks int `toml:"innodb_print_all_deadlocks"`

	GtiDMode string `toml:"gtid_mode"`
	EnforceGtidConsistency int `toml:"enforce_gtid_consistency"`
	SlaveParallelType string `toml:"slave-parallel-type"`
	SlaveParallelWorkers int `toml:"slave-parallel-workers"`
	SlavePreserveCommitOrder int `toml:"slave_preserve_commit_order"`
	SlaveTransactionRetries int `toml:"slave_transaction_retries"`
	BinlogGtiDSimpleRecovery int `toml:"binlog_gtid_simple_recovery"`

	LooseInnodbNumaInterleave int `toml:"loose_innodb_numa_interleave"`
	InnodbBufferPoolDumpPct int `toml:"innodb_buffer_pool_dump_pct"`
	InnodbUndoLogs int `toml:"innodb_undo_logs"`
	InnodbUndoLogTruncate int `toml:"innodb_undo_log_truncate"`
	InnodbMaxUndoLogSize string `toml:"innodb_max_undo_log_size"`
	InnodbPurgeRseGTruncateFrequency int `toml:"innodb_purge_rseg_truncate_frequency"`

	MaxAllowedPacket string `toml:"max_allowed_packet"`
	TableOpenCache int `toml:"table_open_cache"`
	TmpTableSize string `toml:"tmp_table_size"`
	MaxHeapTableSize string `toml:"max_heap_table_size"`
	SortBufferSize string `toml:"sort_buffer_size"`
	JoinBufferSize string `toml:"join_buffer_size"`
	ReadBufferSize string `toml:"read_buffer_size"`
	ReadRndBufferSize string `toml:"read_rnd_buffer_size"`
	KeyBufferSize string `toml:"key_buffer_size"`
	BulkInsertBufferSize string `toml:"bulk_insert_buffer_size"`
	BinlogCacheSize string `toml:"binlog_cache_size"`

	SlowQueryLogFile string `toml:"slow_query_log_file"`
	SlowQueryLog string `toml:"slow_query_log"`
	LongQueryTime float64 `toml:"long_query_time"`
	LogOutput string `toml:"log_output"`

	PerformanceSchema string `toml:"performance_schema"`
	PerformanceSchemaInstrument string `toml:"performance-schema-instrument"`

	SymbolicLinks int `toml:"symbolic-links"`
}

type Config struct {
	Client *Client `toml:"client"`
	MysqlDMulti *MysqlDMulti `toml:"mysqld_multi"`
	Mysql *Mysql `toml:"mysql"`
	MysqlD *MysqlD `toml:"mysqld"`
}

var (
	err error
	confPath string
	Conf = &Config{}

	m = "skip-host-cache"
	n = "skip-name-resolve"

	sText = "plugin_load=rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so"
	news= "plugin_load=\"rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so\""
	pText = "[mysqld]"

	filename = "./my.cnf"
)

func init() {
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
	logError := "/mysqldata/mysql" + port + "/log/mysqld.log"
	pidFile := "/mysqldata/mysql" + port + "/mysql.pid"
	dataDir := "/mysqldata/mysql" + port + "/data"
	serverId := port + strings.Split(report_host, ".")[2]+strings.Split(report_host, ".")[3]
	logBin := "/mysqllog/mysql" + port + "/binlog/mysql-bin"
	relayLog := "/mysqllog/mysql" + port + "/relaylog/mysql-relay"
	p,_ := strconv.Atoi(port)
	reportPort := strconv.Itoa(p + 1)
	threadCacheSize,_ := strconv.Atoi(b)
	innodbBufferPoolSize := c + "M"
	keyBufferSize := d + "M"
	slowQueryLogFile := "/mysqldata/mysql" + port + "/log/mysql-slow.log"

	v := reflect.ValueOf(Conf.Client).Elem()
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))

	v = reflect.ValueOf(Conf.MysqlD).Elem()
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))
	v.FieldByName("LogError").Set(reflect.ValueOf(logError))
	v.FieldByName("PidFile").Set(reflect.ValueOf(pidFile))
	v.FieldByName("DataDir").Set(reflect.ValueOf(dataDir))
	v.FieldByName("Socket").Set(reflect.ValueOf(socket))
	v.FieldByName("Port").Set(reflect.ValueOf(port))
	v.FieldByName("ReportHost").Set(reflect.ValueOf(report_host))
	v.FieldByName("ServerId").Set(reflect.ValueOf(serverId))
	v.FieldByName("LogBin").Set(reflect.ValueOf(logBin))
	v.FieldByName("RelayLog").Set(reflect.ValueOf(relayLog))
	v.FieldByName("ReportPort").Set(reflect.ValueOf(reportPort))
	v.FieldByName("ThreadCacheSize").Set(reflect.ValueOf(threadCacheSize))
	v.FieldByName("InnodbBufferPoolSize").Set(reflect.ValueOf(innodbBufferPoolSize))
	v.FieldByName("KeyBufferSize").Set(reflect.ValueOf(keyBufferSize))
	v.FieldByName("SlowQueryLogFile").Set(reflect.ValueOf(slowQueryLogFile))

	return Conf
}

func CheckFileIsExist(file string) bool {
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

// get my.cnf
func Flex(address string) ([][]string,error){
	arr := FlagCommand(address)
	var f *os.File
	if CheckFileIsExist(filename) {
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
			encoder.Encode(p.MysqlD)

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
