package pgmetrics

const ModelSchemaVersion = "1.13"

type Model struct {
	Metadata Metadata `json:"meta"` // metadata about this object

	StartTime        int64  `json:"start_time"`        // of postmaster
	SystemIdentifier string `json:"system_identifier"` // from pg_control

	// Checkpoint information
	CheckpointLSN   string `json:"checkpoint_lsn"`
	PriorLSN        string `json:"prior_lsn"`
	RedoLSN         string `json:"redo_lsn"`
	TimelineID      int    `json:"timeline_id"`
	NextXid         int    `json:"next_xid"`
	OldestXid       int    `json:"oldest_xid"`
	OldestActiveXid int    `json:"oldest_active_xid"`
	CheckpointTime  int64  `json:"checkpoint_time"`

	// wal
	WALFlushLSN  string `json:"wal_flush_lsn"`
	WALInsertLSN string `json:"wal_insert_lsn"`
	WALLSN       string `json:"wal_lsn"`

	// Recovery
	IsInRecovery            bool   `json:"is_in_recovery"`
	IsWalReplayPaused       bool   `json:"is_wal_replay_paused"`
	LastWALReceiveLSN       string `json:"last_wal_receive_lsn"`
	LastWALReplayLSN        string `json:"last_wal_replay_lsn"`
	LastXActReplayTimestamp int64  `json:"last_xact_replay_timestamp"`

	// last committed transaction (needs track_commit_timestamp = on)
	LastXactXid       int   `json:"last_xact_xid"`
	LastXactTimestamp int64 `json:"last_xact_time"`

	// wal - settings, archival stats
	WALArchiving  WALArchiving `json:"wal_archiving"`
	WALCount      int          `json:"wal_count"`
	WALReadyCount int          `json:"wal_ready_count"`

	// NotificationQueueUsage is the fraction of the asynchronous notification
	// queue currently occupied. Postgres v9.6 and above only. Added in
	// schema version 1.1.
	NotificationQueueUsage float64 `json:"notification_queue_usage"`

	// replication
	ReplicationOutgoing []ReplicationOut  `json:"replication_outgoing,omitempty"`
	ReplicationIncoming *ReplicationIn    `json:"replication_incoming,omitempty"`
	ReplicationSlots    []ReplicationSlot `json:"replication_slots,omitempty"`

	// other cluster-level stats
	BGWriter       BGWriter                `json:"bg_writer"`
	Backends       []Backend               `json:"backends"`
	VacuumProgress []VacuumProgressBackend `json:"vacuum_progress,omitempty"`

	// structural cluster-level information
	Roles       []Role       `json:"roles"`
	Databases   []Database   `json:"databases,omitempty"`
	Tablespaces []Tablespace `json:"tablespaces,omitempty"`

	// Database-specific
	Tables           []Table        `json:"tables,omitempty"`
	Indexes          []Index        `json:"indexes,omitempty"`
	Sequences        []Sequence     `json:"sequences,omitempty"`
	UserFunctions    []UserFunction `json:"user_functions,omitempty"`
	Extensions       []Extension    `json:"extensions,omitempty"`
	DisabledTriggers []Trigger      `json:"disabled_triggers,omitempty"`
	Statements       []Statement    `json:"statements,omitempty"`

	// System-level
	System *SystemMetrics `json:"system,omitempty"`

	// settings
	Settings map[string]Setting `json:"settings"` // all settings and their values

	// following fields present only in schema 1.2 and later

	// Logical replication (database-specific)
	Publications  []Publication  `json:"publications,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`

	// following fields present only in schema 1.3 and later

	// Lock information
	Locks        []Lock        `json:"locks,omitempty"`
	BlockingPIDs map[int][]int `json:"blocking_pids,omitempty"`

	// following fields present only in schema 1.4 and later

	PgBouncer *PgBouncer `json:"pgbouncer,omitempty"`

	// following fields are present only in schema 1.6 and later

	// the numerically highest wal segment number
	HighestWALSegment uint64 `json:"highwal,omitempty"`

	// following fields are present only in schema 1.7 and later

	// query execution plans
	Plans []Plan `json:"plans,omitempty"`

	// autovacuum information
	AutoVacuums []AutoVacuum `json:"autovacuums,omitempty"`

	// deadlock information
	Deadlocks []Deadlock `json:"deadlocks,omitempty"`

	// following fields are present only in schema 1.8 and later

	// metrics from AWS RDS
	RDS *RDS `json:"rds,omitempty"`

	// the types of running backends and their counts
	BackendTypeCounts map[string]int `json:"betypecounts,omitempty"`

	// following fields are present only in schema 1.9 and later

	// citus-related information, per db
	Citus map[string]*Citus `json:"citus,omitempty"`

	// following fields are present only in schema 1.11 and later

	// WAL activity info, from pg_stat_wal, pg >= v14
	WAL *WAL `json:"wal,omitempty"`

	// following fields are present only in schema 1.12 and later

	// metrics from Azure PostgreSQL, via Azure Monitor APIs
	Azure *Azure `json:"azure,omitempty"`

	// progress information from pg_stat_progress_* (see above for vacuum)
	AnalyzeProgress     []AnalyzeProgressBackend     `json:"analyze_progress,omitempty"`
	BasebackupProgress  []BasebackupProgressBackend  `json:"basebackup_progress,omitempty"`
	ClusterProgress     []ClusterProgressBackend     `json:"cluster_progress,omitempty"`
	CopyProgress        []CopyProgressBackend        `json:"copy_progress,omitempty"`
	CreateIndexProgress []CreateIndexProgressBackend `json:"create_index_progress,omitempty"`
}

// DatabaseByOID iterates over the databases in the model and returns the reference
// to a Database that has the given oid. If there is no such database, it returns nil.
func (m *Model) DatabaseByOID(oid int) *Database {
	for i, d := range m.Databases {
		if d.OID == oid {
			return &m.Databases[i]
		}
	}
	return nil
}

// RoleByOID iterates over the roles in the model and returns the reference
// to a Role that has the given oid. If there is no such role, it returns nil.
func (m *Model) RoleByOID(oid int) *Role {
	for i, r := range m.Roles {
		if r.OID == oid {
			return &m.Roles[i]
		}
	}
	return nil
}

// TableByName iterates over the tables in the model and returns the reference
// to a Table that has the given database, schema and table names. If there is
// no such table, it returns nil.
func (m *Model) TableByName(db, schema, table string) *Table {
	for i, t := range m.Tables {
		if t.DBName == db && t.SchemaName == schema && t.Name == table {
			return &m.Tables[i]
		}
	}
	return nil
}

// TableByOID iterates over the tables in the model and returns the reference
// to a Table that has the given oid. If there is no such table, it returns nil.
func (m *Model) TableByOID(oid int) *Table {
	for i, t := range m.Tables {
		if t.OID == oid {
			return &m.Tables[i]
		}
	}
	return nil
}

// IndexByName iterates over the indexes in the model and returns the reference
// to an Index that has the given database, schema and index names. If there is
// no such index, it returns nil.
func (m *Model) IndexByName(db, schema, index string) *Index {
	for i, idx := range m.Indexes {
		if idx.DBName == db && idx.SchemaName == schema && idx.Name == index {
			return &m.Indexes[i]
		}
	}
	return nil
}

// IndexByOID iterates over the indexes in the model and returns the reference
// to an Index that has the given oid. If there is no such index, it returns nil.
func (m *Model) IndexByOID(oid int) *Index {
	for i, idx := range m.Indexes {
		if idx.OID == oid {
			return &m.Indexes[i]
		}
	}
	return nil
}

// Metadata contains information about how to interpret the other fields in
// "Model" data structure.
type Metadata struct {
	Version      string   `json:"version"`       // schema version, "semver" format
	At           int64    `json:"at"`            // time when this report was started
	CollectedDBs []string `json:"collected_dbs"` // names of dbs we collected db-level stats from
	Local        bool     `json:"local"`         // was connected to a local postgres server?
	UserAgent    string   `json:"user_agent"`    // "pgmetrics/1.8.1"
	Username     string   `json:"user"`          // user that pgmetrics connected as
}

type SystemMetrics struct {
	CPUModel   string  `json:"cpu_model,omitempty"` // model of the CPU
	NumCores   int     `json:"num_cores"`           // number of cores
	LoadAvg    float64 `json:"loadavg"`             // 1-minute load average from the OS
	MemUsed    int64   `json:"memused"`             // used RAM, in bytes
	MemFree    int64   `json:"memfree"`             // free RAM, in bytes
	MemBuffers int64   `json:"membuffers"`          // RAM used for buffers, in bytes
	MemCached  int64   `json:"memcached"`           // RAM used for cache, in bytes
	SwapUsed   int64   `json:"swapused"`            // used swap memory in bytes, 0 if no swap
	SwapFree   int64   `json:"swapfree"`            // free swap memory in bytes, 0 if no swap
	Hostname   string  `json:"hostname"`            // hostname from the OS
	// following fields present only in schema 1.8 and later
	MemSlab int64 `json:"memslab"` // RAM used for slab in bytes
}

type Backend struct {
	DBName          string `json:"db_name"`
	RoleName        string `json:"role_name"`
	ApplicationName string `json:"application_name"`
	PID             int    `json:"pid"`
	ClientAddr      string `json:"client_addr"`
	BackendStart    int64  `json:"backend_start"`
	XactStart       int64  `json:"xact_start"`
	QueryStart      int64  `json:"query_start"`
	StateChange     int64  `json:"state_change"`
	WaitEventType   string `json:"wait_event_type"`
	WaitEvent       string `json:"wait_event"`
	State           string `json:"state"`
	BackendXid      int    `json:"backend_xid"`
	BackendXmin     int    `json:"backend_xmin"`
	Query           string `json:"query"`
	// following fields present only in schema 1.11 and later
	QueryID int64 `json:"query_id,omitempty"` // >= pg14
}

type ReplicationSlot struct {
	SlotName          string `json:"slot_name"`
	Plugin            string `json:"plugin"`
	SlotType          string `json:"slot_type"`
	DBName            string `json:"db_name"`
	Active            bool   `json:"active"`
	Xmin              int    `json:"xmin"`
	CatalogXmin       int    `json:"catalog_xmin"`
	RestartLSN        string `json:"restart_lsn"`
	ConfirmedFlushLSN string `json:"confirmed_flush_lsn"`
	Temporary         bool   `json:"temporary"`
	// following fields present only in schema 1.11 and later
	WALStatus   string `json:"wal_status,omitempty"`    // >= pg13
	SafeWALSize int64  `json:"safe_wal_size,omitempty"` // >= pg13
	TwoPhase    bool   `json:"two_phase,omitempty"`     // >= pg14
}

type Role struct {
	OID            int      `json:"oid"`
	Name           string   `json:"name"`
	Rolsuper       bool     `json:"rolsuper"`
	Rolinherit     bool     `json:"rolinherit"`
	Rolcreaterole  bool     `json:"rolcreaterole"`
	Rolcreatedb    bool     `json:"rolcreatedb"`
	Rolcanlogin    bool     `json:"rolcanlogin"`
	Rolreplication bool     `json:"rolreplication"`
	Rolbypassrls   bool     `json:"rolbypassrls"`
	Rolconnlimit   int      `json:"rolconnlimit"`
	Rolvaliduntil  int64    `json:"rolvaliduntil"`
	MemberOf       []string `json:"memberof"`
}

type Tablespace struct {
	OID         int    `json:"oid"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Location    string `json:"location"`
	Size        int64  `json:"size"`
	DiskUsed    int64  `json:"disk_used"`
	DiskTotal   int64  `json:"disk_total"`
	InodesUsed  int64  `json:"inodes_used"`
	InodesTotal int64  `json:"inodes_total"`
}

type Database struct {
	OID             int     `json:"oid"`
	Name            string  `json:"name"`
	DatDBA          int     `json:"datdba"`
	DatTablespace   int     `json:"dattablespace"`
	DatConnLimit    int     `json:"datconnlimit"`
	AgeDatFrozenXid int     `json:"age_datfrozenxid"`
	NumBackends     int     `json:"numbackends"`
	XactCommit      int64   `json:"xact_commit"`
	XactRollback    int64   `json:"xact_rollback"`
	BlksRead        int64   `json:"blks_read"`
	BlksHit         int64   `json:"blks_hit"`
	TupReturned     int64   `json:"tup_returned"`
	TupFetched      int64   `json:"tup_fetched"`
	TupInserted     int64   `json:"tup_inserted"`
	TupUpdated      int64   `json:"tup_updated"`
	TupDeleted      int64   `json:"tup_deleted"`
	Conflicts       int64   `json:"conflicts"`
	TempFiles       int64   `json:"temp_files"`
	TempBytes       int64   `json:"temp_bytes"`
	Deadlocks       int64   `json:"deadlocks"`
	BlkReadTime     float64 `json:"blk_read_time"`
	BlkWriteTime    float64 `json:"blk_write_time"`
	StatsReset      int64   `json:"stats_reset"`
	Size            int64   `json:"size"`
	// following fields present only in schema 1.11 and later
	ChecksumFailures    int64   `json:"checksum_failures,omitempty"`        // >= pg12
	ChecksumLastFailure int64   `json:"checksum_last_failure,omitempty"`    // >= pg12
	SessionTime         float64 `json:"session_time,omitempty"`             // in milliseconds, pg >= v14
	ActiveTime          float64 `json:"active_time,omitempty"`              // in milliseconds, pg >= v14
	IdleInTxTime        float64 `json:"idle_in_transaction_time,omitempty"` // in milliseconds, pg >= v14
	Sessions            int64   `json:"sessions,omitempty"`                 // pg >= v14
	SessionsAbandoned   int64   `json:"sessions_abandoned,omitempty"`       // pg >= v14
	SessionsFatal       int64   `json:"sessions_fatal,omitempty"`           // pg >= v14
	SessionsKilled      int64   `json:"sessions_killed,omitempty"`          // pg >= v14
}

type Table struct {
	OID              int    `json:"oid"`
	DBName           string `json:"db_name"`
	SchemaName       string `json:"schema_name"`
	Name             string `json:"name"`
	SeqScan          int64  `json:"seq_scan"`
	SeqTupRead       int64  `json:"seq_tup_read"`
	IdxScan          int64  `json:"idx_scan"`
	IdxTupFetch      int64  `json:"idx_tup_fetch"`
	NTupIns          int64  `json:"n_tup_ins"`
	NTupUpd          int64  `json:"n_tup_upd"`
	NTupDel          int64  `json:"n_tup_del"`
	NTupHotUpd       int64  `json:"n_tup_hot_upd"`
	NLiveTup         int64  `json:"n_live_tup"`
	NDeadTup         int64  `json:"n_dead_tup"`
	NModSinceAnalyze int64  `json:"n_mod_since_analyze"`
	LastVacuum       int64  `json:"last_vacuum"`
	LastAutovacuum   int64  `json:"last_autovacuum"`
	LastAnalyze      int64  `json:"last_analyze"`
	LastAutoanalyze  int64  `json:"last_autoanalyze"`
	VacuumCount      int64  `json:"vacuum_count"`
	AutovacuumCount  int64  `json:"autovacuum_count"`
	AnalyzeCount     int64  `json:"analyze_count"`
	AutoanalyzeCount int64  `json:"autoanalyze_count"`
	HeapBlksRead     int64  `json:"heap_blks_read"`
	HeapBlksHit      int64  `json:"heap_blks_hit"`
	IdxBlksRead      int64  `json:"idx_blks_read"`
	IdxBlksHit       int64  `json:"idx_blks_hit"`
	ToastBlksRead    int64  `json:"toast_blks_read"`
	ToastBlksHit     int64  `json:"toast_blks_hit"`
	TidxBlksRead     int64  `json:"tidx_blks_read"`
	TidxBlksHit      int64  `json:"tidx_blks_hit"`
	Size             int64  `json:"size"`
	Bloat            int64  `json:"bloat"`
	// following fields present only in schema 1.2 and later
	RelKind         string `json:"relkind"`
	RelPersistence  string `json:"relpersistence"`
	RelNAtts        int    `json:"relnatts"`
	AgeRelFrozenXid int    `json:"age_relfrozenxid"`
	RelIsPartition  bool   `json:"relispartition"`
	TablespaceName  string `json:"tablespace_name"`
	ParentName      string `json:"parent_name"`
	PartitionCV     string `json:"partition_cv"` // partition constraint value
	// following fields present only in schema 1.7 and later
	ACL string `json:"acl,omitempty"`
	// following fields present only in schema 1.11 and later
	NInsSinceVacuum int64 `json:"n_ins_since_vacuum,omitempty"` // pg >= v13
}

type Index struct {
	OID         int    `json:"oid"`
	DBName      string `json:"db_name"`
	SchemaName  string `json:"schema_name"`
	TableOID    int    `json:"table_oid"`
	TableName   string `json:"table_name"`
	Name        string `json:"name"`
	IdxScan     int64  `json:"idx_scan"`
	IdxTupRead  int64  `json:"idx_tup_read"`
	IdxTupFetch int64  `json:"idx_tup_fetch"`
	IdxBlksRead int64  `json:"idx_blks_read"`
	IdxBlksHit  int64  `json:"idx_blks_hit"`
	Size        int64  `json:"size"`
	Bloat       int64  `json:"bloat"`
	// following fields present only in schema 1.2 and later
	RelNAtts       int    `json:"relnatts"`
	AMName         string `json:"amname"`
	TablespaceName string `json:"tablespace_name"`
	// following fields present only in schema 1.8 and later
	Definition string `json:"def"`
}

type Sequence struct {
	OID        int    `json:"oid"`
	DBName     string `json:"db_name"`
	SchemaName string `json:"schema_name"`
	Name       string `json:"name"`
	BlksRead   int64  `json:"blks_read"`
	BlksHit    int64  `json:"blks_hit"`
}

type UserFunction struct {
	OID        int     `json:"oid"`
	SchemaName string  `json:"schema_name"`
	DBName     string  `json:"db_name"`
	Name       string  `json:"name"`
	Calls      int64   `json:"calls"`
	TotalTime  float64 `json:"total_time"`
	SelfTime   float64 `json:"self_time"`
}

type VacuumProgressBackend struct {
	DBName           string `json:"db_name"`
	TableOID         int    `json:"table_oid"`
	TableName        string `json:"table_name"`
	Phase            string `json:"phase"`
	HeapBlksTotal    int64  `json:"heap_blks_total"`
	HeapBlksScanned  int64  `json:"heap_blks_scanned"`
	HeapBlksVacuumed int64  `json:"heap_blks_vacuumed"`
	IndexVacuumCount int64  `json:"index_vacuum_count"`
	MaxDeadTuples    int64  `json:"max_dead_tuples"`
	NumDeadTuples    int64  `json:"num_dead_tuples"`
	// following fields present only in schema 1.12 and later
	PID int `json:"pid,omitempty"`
}

type Extension struct {
	Name             string `json:"name"`
	DBName           string `json:"db_name"`
	DefaultVersion   string `json:"default_version"`
	InstalledVersion string `json:"installed_version"`
	Comment          string `json:"comment"`
}

type Setting struct {
	Setting string `json:"setting"`
	BootVal string `json:"bootval,omitempty"`
	Source  string `json:"source,omitempty"`
	// following fields present only in schema 1.11 and later
	Pending bool `json:"pending,omitempty"`
}

type WALArchiving struct {
	ArchivedCount    int    `json:"archived_count"`
	LastArchivedWAL  string `json:"last_archived_wal"`
	LastArchivedTime int64  `json:"last_archived_time"`
	FailedCount      int    `json:"failed_count"`
	LastFailedWAL    string `json:"last_failed_wal"`
	LastFailedTime   int64  `json:"last_failed_time"`
	StatsReset       int64  `json:"stats_reset"`
}

type BGWriter struct {
	CheckpointsTimed     int64   `json:"checkpoints_timed"`
	CheckpointsRequested int64   `json:"checkpoints_req"`
	CheckpointWriteTime  float64 `json:"checkpoint_write_time"`
	CheckpointSyncTime   float64 `json:"checkpoint_sync_time"`
	BuffersCheckpoint    int64   `json:"buffers_checkpoint"`
	BuffersClean         int64   `json:"buffers_clean"`
	MaxWrittenClean      int64   `json:"maxwritten_clean"`
	BuffersBackend       int64   `json:"buffers_backend"`
	BuffersBackendFsync  int64   `json:"buffers_backend_fsync"`
	BuffersAlloc         int64   `json:"buffers_alloc"`
	StatsReset           int64   `json:"stats_reset"`
}

type ReplicationOut struct {
	RoleName        string `json:"role_name"`
	ApplicationName string `json:"application_name"`
	ClientAddr      string `json:"client_addr"`
	BackendStart    int64  `json:"backend_start"`
	BackendXmin     int    `json:"backend_xmin"`
	State           string `json:"state"`
	SentLSN         string `json:"sent_lsn"`
	WriteLSN        string `json:"write_lsn"`
	FlushLSN        string `json:"flush_lsn"`
	ReplayLSN       string `json:"replay_lsn"`
	WriteLag        int    `json:"write_lag"`  // only in 10.x
	FlushLag        int    `json:"flush_lag"`  // only in 10.x
	ReplayLag       int    `json:"replay_lag"` // only in 10.x
	SyncPriority    int    `json:"sync_priority"`
	SyncState       string `json:"sync_state"`
	// following fields present only in schema 1.5 and later
	PID int `json:"pid,omitempty"`
	// following fields present only in schema 1.11 and later
	ReplyTime int64 `json:"reply_time,omitempty"` // >= pg12
}

type ReplicationIn struct {
	Status             string `json:"status"`
	ReceiveStartLSN    string `json:"receive_start_lsn"`
	ReceiveStartTLI    int    `json:"receive_start_tli"`
	ReceivedLSN        string `json:"received_lsn"` // empty in 13.x+
	ReceivedTLI        int    `json:"received_tli"`
	LastMsgSendTime    int64  `json:"last_msg_send_time"`
	LastMsgReceiptTime int64  `json:"last_msg_receipt_time"`
	Latency            int64  `json:"latency_micros"`
	LatestEndLSN       string `json:"latest_end_lsn"`
	LatestEndTime      int64  `json:"latest_end_time"`
	SlotName           string `json:"slot_name"`
	Conninfo           string `json:"conninfo"`
	// following fields present only in schema 1.9 and later (13.x+)
	WrittenLSN string `json:"written_lsn,omitempty"`
	FlushedLSN string `json:"flushed_lsn,omitempty"`
	// following fields present only in schema 1.11 and later
	SenderHost string `json:"sender_host,omitempty"` // >= pg11
}

type Trigger struct {
	OID        int    `json:"oid"`
	DBName     string `json:"db_name"`
	SchemaName string `json:"schema_name"`
	TableName  string `json:"table_name"`
	Name       string `json:"name"`
	ProcName   string `json:"proc_name"`
}

// Statement represents a row of the pg_stat_statements view. Added in schema
// version 1.1.
type Statement struct {
	UserOID           int     `json:"useroid"`             // OID of user who executed the statement
	UserName          string  `json:"user"`                // Name of the user corresponding to useroid (might be empty)
	DBOID             int     `json:"db_oid"`              // OID of database in which the statement was executed
	DBName            string  `json:"db_name"`             // Name of the database corresponding to db_oid
	QueryID           int64   `json:"queryid"`             // Internal hash code, computed from the statement's parse tree
	Query             string  `json:"query"`               // Text of a representative statement
	Calls             int64   `json:"calls"`               // Number of times executed
	TotalTime         float64 `json:"total_time"`          // Total time spent executing the statement, in milliseconds
	MinTime           float64 `json:"min_time"`            // Minimum time spent executing the statement, in milliseconds
	MaxTime           float64 `json:"max_time"`            // Maximum time spent executing the statement, in milliseconds
	StddevTime        float64 `json:"stddev_time"`         // Population standard deviation of time spent executing the statement, in milliseconds
	Rows              int64   `json:"rows"`                // Total number of rows retrieved or affected by the statement
	SharedBlksHit     int64   `json:"shared_blks_hit"`     // Total number of shared block cache hits by the statement
	SharedBlksRead    int64   `json:"shared_blks_read"`    // Total number of shared blocks read by the statement
	SharedBlksDirtied int64   `json:"shared_blks_dirtied"` // Total number of shared blocks dirtied by the statement
	SharedBlksWritten int64   `json:"shared_blks_written"` // Total number of shared blocks written by the statement
	LocalBlksHit      int64   `json:"local_blks_hit"`      // Total number of local block cache hits by the statement
	LocalBlksRead     int64   `json:"local_blks_read"`     // Total number of local blocks read by the statement
	LocalBlksDirtied  int64   `json:"local_blks_dirtied"`  // Total number of local blocks dirtied by the statement
	LocalBlksWritten  int64   `json:"local_blks_written"`  // Total number of local blocks written by the statement
	TempBlksRead      int64   `json:"temp_blks_read"`      // Total number of temp blocks read by the statement
	TempBlksWritten   int64   `json:"temp_blks_written"`   // Total number of temp blocks written by the statement
	BlkReadTime       float64 `json:"blk_read_time"`       // Total time the statement spent reading blocks, in milliseconds (if track_io_timing is enabled, otherwise zero)
	BlkWriteTime      float64 `json:"blk_write_time"`      // Total time the statement spent writing blocks, in milliseconds (if track_io_timing is enabled, otherwise zero)
	// following fields present only in schema 1.10 and later (for Postgres v13+)
	Plans          int64   `json:"plans"`            // Number of times the statement was planned
	TotalPlanTime  float64 `json:"total_plan_time"`  // Total time spent planning the statement, in milliseconds
	MinPlanTime    float64 `json:"min_plan_time"`    // Minimum time spent planning the statement, in milliseconds
	MaxPlanTime    float64 `json:"max_plan_time"`    // Maximum time spent planning the statement, in milliseconds
	StddevPlanTime float64 `json:"stddev_plan_time"` // Population standard deviation of time spent planning the statement, in milliseconds
	WALRecords     int64   `json:"wal_records"`      // Total number of WAL records generated by the statement
	WALFPI         int64   `json:"wal_fpi"`          // Total number of WAL full page images generated by the statement
	WALBytes       int64   `json:"wal_bytes"`        // Total amount of WAL bytes generated by the statement
}

// Publication represents a single v10+ publication. Added in schema 1.2.
type Publication struct {
	OID        int    `json:"oid"`
	Name       string `json:"name"`
	DBName     string `json:"db_name"`
	AllTables  bool   `json:"alltables"`
	Insert     bool   `json:"insert"`
	Update     bool   `json:"update"`
	Delete     bool   `json:"delete"`
	TableCount int    `json:"table_count"`
}

// Subscription represents a single v10+ subscription. Added in schema 1.2.
type Subscription struct {
	OID                int    `json:"oid"`
	Name               string `json:"name"`
	DBName             string `json:"db_name"`
	Enabled            bool   `json:"enabled"`
	PubCount           int    `json:"pub_count"`
	TableCount         int    `json:"table_count"`
	WorkerCount        int    `json:"worker_count"`
	ReceivedLSN        string `json:"received_lsn"`
	LatestEndLSN       string `json:"latest_end_lsn"`
	LastMsgSendTime    int64  `json:"last_msg_send_time"`
	LastMsgReceiptTime int64  `json:"last_msg_receipt_time"`
	LatestEndTime      int64  `json:"latest_end_time"`
	Latency            int64  `json:"latency_micros"`
	// following fields present only in schema 1.13 and later
	ApplyErrorCount int `json:"apply_error_count,omitempty"` // >= pg15
	SyncErrorCount  int `json:"sync_error_count,omitempty"`  // >= pg15
}

// Lock represents a single row from pg_locks. Added in schema 1.3.
type Lock struct {
	LockType    string `json:"locktype"`
	DBName      string `json:"db_name,omitempty"`
	PID         int    `json:"pid"`
	Mode        string `json:"mode"`
	Granted     bool   `json:"granted"`
	RelationOID int    `json:"relation_oid,omitempty"`
	// following fields present only in schema 1.11 and later
	WaitStart int64 `json:"waitstart,omitempty"` // >= pg14
}

// PgBouncer contains information collected from the virtual "pgbouncer"
// database. Added in schema 1.4.
type PgBouncer struct {
	Pools     []PgBouncerPool     `json:"pools,omitempty"`
	Stats     []PgBouncerStat     `json:"stats,omitempty"`
	Databases []PgBouncerDatabase `json:"dbs,omitempty"`

	SCActive  int     `json:"sc_active"`  // no. of active server conns
	SCIdle    int     `json:"sc_idle"`    // no. of idle server conns
	SCUsed    int     `json:"sc_used"`    // no. of used server conns
	SCMaxWait float64 `json:"sc_maxwait"` // max wait time for server conns

	CCActive  int     `json:"cc_active"`  // no. of active client conns
	CCWaiting int     `json:"cc_waiting"` // no. of waiting client conns
	CCIdle    int     `json:"cc_idle"`    // no. of idle client conns
	CCUsed    int     `json:"cc_used"`    // no. of used client conns
	CCMaxWait float64 `json:"cc_maxwait"` // max wait time for *waiting* client conns
	CCAvgWait float64 `json:"cc_avgwait"` // avg wait time for *waiting* client conns
}

// PgBouncerPool contains information about one pool of PgBouncer (one row
// from SHOW POOLS).
type PgBouncerPool struct {
	Database  string  `json:"db_name"`
	UserName  string  `json:"user"`
	ClActive  int     `json:"cl_active"`
	ClWaiting int     `json:"cl_waiting"`
	SvActive  int     `json:"sv_active"`
	SvIdle    int     `json:"sv_idle"`
	SvUsed    int     `json:"sv_used"`
	SvTested  int     `json:"sv_tested"`
	SvLogin   int     `json:"sv_login"`
	MaxWait   float64 `json:"maxwait"` // seconds
	Mode      string  `json:"pool_mode"`

	// following fields present only in schema 1.11 and later
	ClCancelReq int `json:"cl_cancel_req,omitempty"` // only in pgbouncer >= v1.16.0
}

// PgBouncerDatabase contains information about one database of PgBouncer
// (one row from SHOW DATABASES).
type PgBouncerDatabase struct {
	Database       string `json:"db_name"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	SourceDatabase string `json:"srcdb_name"`
	User           string `json:"force_user"`
	MaxConn        int    `json:"max_connections"`
	CurrConn       int    `json:"current_connections"`
	Paused         bool   `json:"paused"`
	Disabled       bool   `json:"disabled"`
}

// PgBouncerStat contains one row from SHOW STATS. Times are in seconds,
// averages are for the last second (as per PgBouncer docs).
type PgBouncerStat struct {
	Database        string  `json:"db_name"`
	TotalXactCount  int64   `json:"total_xact_count"`
	TotalQueryCount int64   `json:"total_query_count"`
	TotalReceived   int64   `json:"total_received"`   // bytes
	TotalSent       int64   `json:"total_sent"`       // bytes
	TotalXactTime   float64 `json:"total_xact_time"`  // seconds
	TotalQueryTime  float64 `json:"total_query_time"` // seconds
	TotalWaitTime   float64 `json:"total_wait_time"`  // seconds
	AvgXactCount    int64   `json:"avg_xact_count"`
	AvgQueryCount   int64   `json:"avg_query_count"`
	AvgReceived     int64   `json:"avg_received"`   // bytes
	AvgSent         int64   `json:"avg_sent"`       // bytes
	AvgXactTime     float64 `json:"avg_xact_time"`  // seconds
	AvgQueryTime    float64 `json:"avg_query_time"` // seconds
	AvgWaitTime     float64 `json:"avg_wait_time"`  // seconds
}

// Plan represents a query execution plan. Added in schema 1.7.
type Plan struct {
	Database string `json:"db_name"` // might be empty
	UserName string `json:"user"`    // might be empty
	Format   string `json:"format"`  // text, json, yaml or xml
	At       int64  `json:"at"`      // time when plan was logged, as seconds since epoch
	Query    string `json:"query"`   // the sql query
	Plan     string `json:"plan"`    // the plan as a string

	// following fields present only in schema 1.12 and later
	QueryID int64 `json:"queryid,omitempty"` // query id
}

// AutoVacuum contains information about a single autovacuum run.
// Added in schema 1.7.
type AutoVacuum struct {
	At      int64   `json:"at"`         // time when activity was logged, as seconds since epoch
	Table   string  `json:"table_name"` // fully qualified, db.schema.table
	Elapsed float64 `json:"elapsed"`    // in seconds
}

// Deadlock contains information about a single deadlock detection log.
// Added in schema 1.7.
type Deadlock struct {
	At     int64  `json:"at"`     // time when activity was logged, as seconds since epoch
	Detail string `json:"detail"` // information about the deadlocking processes
}

// RDS contains metrics collected from AWS RDS (also includes Aurora).
// Added in schema 1.8.
type RDS struct {
	Basic    map[string]float64     `json:"basic"`              // Basic Monitoring Metrics
	Enhanced map[string]interface{} `json:"enhanced,omitempty"` // Enhanced Monitoring
}

// Citus contains metrics collected from Citus extension.
// Added in schema 1.9.
type Citus struct {
	Version        string           `json:"version"`
	Nodes          []CitusNode      `json:"nodes"`
	Statements     []CitusStatement `json:"statements"`
	Backends       []CitusBackend   `json:"dist_activity,omitempty"`   // citus <=10.x
	WorkerBackends []CitusBackend   `json:"worker_activity,omitempty"` // citus <=10.x
	Locks          []CitusLock      `json:"locks"`
	// following fields present only in schema 1.13 and later
	AllBackends       []CitusBackendV11 `json:"activity,omitempty"`           // citus >=11.x
	Tables            []CitusTable      `json:"tables,omitempty"`             // citus >=11.x
	CoordinatorNodeID int               `json:"coordinator_nodeid,omitempty"` // citus >=11.x
	ConnectedNodeID   int               `json:"connected_nodeid,omitempty"`   // citus >=11.x, the node pgmetrics connected to
}

// CitusNode represents a row from the pg_dist_node table. Added in schema 1.9.
type CitusNode struct {
	ID               int    `json:"nodeid"`
	GroupID          int    `json:"groupid"`
	Name             string `json:"nodename"`
	Port             int    `json:"nodeport"`
	Rack             string `json:"noderack"`
	IsActive         bool   `json:"isactive"`
	Role             string `json:"noderole"`
	Cluster          string `json:"nodecluster"`
	ShouldHaveShards bool   `json:"shouldhaveshards"`
}

// CitusStatement represents a row in citus_stat_statements. Added in schema 1.9.
type CitusStatement struct {
	QueryID      int64  `json:"queryid"`       // same as pg_stat_statements.queryid
	UserOID      int    `json:"useroid"`       // user who ran the query
	DBOID        int    `json:"db_oid"`        // database instance of coordinator
	Query        string `json:"query"`         // anonymized query string
	Executor     string `json:"executor"`      // Citus executor used: adaptive, real-time, task-tracker, router, or insert-select
	PartitionKey string `json:"partition_key"` // value of distribution column in router-executed queries, else NULL
	Calls        int64  `json:"calls"`         // number of times the query was run
}

// CitusBackend represents a row from citus_dist_stat_activity or from
// citus_worker_stat_activity; in Citus <=10.x. Added in schema 1.9.
type CitusBackend struct {
	Backend                    // also include all fields from pg_stat_activity
	QueryHostname       string `json:"query_hostname"`
	QueryPort           int    `json:"query_port"`
	MasterQueryHostname string `json:"master_query_hostname"`
	MasterQueryPort     int    `json:"master_query_port"`
	TxNumber            int64  `json:"transaction_number"`
	TxStamp             int64  `json:"transaction_stamp"`
}

// CitusBackendV11 represents a row from citus_stat_activity in Citus >=11.x.
// Added in schema 1.13.
type CitusBackendV11 struct {
	Backend              // also include all fields from pg_stat_activity
	GlobalPID     int64  `json:"global_pid"`
	NodeID        int    `json:"node_id"`
	IsWorkerQuery bool   `json:"is_worker_query"`
	BackendType   string `json:"backend_type"`
}

// CitusLock represents a single row from citus_lock_waits. Added in schema 1.9.
type CitusLock struct {
	WaitingPID       int    `json:"waiting_pid,omitempty"`  // citus <=10.x, 0 otherwise
	BlockingPID      int    `json:"blocking_pid,omitempty"` // citus <=10.x, 0 otherwise
	BlockedStmt      string `json:"blocked_statement"`
	CurrStmt         string `json:"current_statement_in_blocking_process"`
	WaitingNodeID    int    `json:"waiting_node_id"`
	BlockingNodeID   int    `json:"blocking_node_id"`
	WaitingNodeName  string `json:"waiting_node_name,omitempty"`  // citus <=10.x, '' otherwise
	BlockingNodeName string `json:"blocking_node_name,omitempty"` // citus <=10.x, '' otherwise
	WaitingNodePort  int    `json:"waiting_node_port,omitempty"`  // citus <=10.x, 0 otherwise
	BlockingNodePort int    `json:"blocking_node_port,omitempty"` // citus <=10.x, 0 otherwise
	// following fields present only in schema 1.13 and later
	WaitingGPID  int64 `json:"waiting_gpid,omitempty"`  // citus >=11.x
	BlockingGPID int64 `json:"blocking_gpid,omitempty"` // citus >=11.x
}

// CitusTable represents an equivalent of a single row from citus_tables.
// Added in schema 1.13.
type CitusTable struct {
	OID                int    `json:"oid"`
	TableName          string `json:"table_name"`
	TableType          string `json:"citus_table_type"`
	DistributionColumn string `json:"distribution_column"`
	ColocationID       int    `json:"colocation_id"`
	Size               int64  `json:"table_size"`
	ShardCount         int    `json:"shard_count"`
	TableOwner         string `json:"table_owner"`
	AccessMethod       string `json:"access_method"`
}

// WAL represents a single row from pg_stat_wal. Added in schema 1.11.
// pg_stat_wal is available only in pg >= v14.
type WAL struct {
	Records     int64   `json:"records"`
	FPI         int64   `json:"fpi"`
	Bytes       int64   `json:"bytes"`
	BuffersFull int64   `json:"buffers_full"`
	Write       int64   `json:"write"`
	Sync        int64   `json:"sync"`
	WriteTime   float64 `json:"write_time"` // in milliseconds
	SyncTime    float64 `json:"sync_time"`  // in milliseconds
	StatsReset  int64   `json:"stats_reset"`
}

// Azure represents metrics and information collected from Azure PostgreSQL
// via Azure Monitor APIs. Added in schema 1.12.
type Azure struct {
	ResourceName   string             `json:"resource_name"`
	ResourceType   string             `json:"resource_type"`
	ResourceRegion string             `json:"resource_region"`
	Metrics        map[string]float64 `json:"metrics"`
}

// AnalyzeProgressBackend represents a row (and each row represents one
// backend) from pg_stat_progress_analyze.
//
// pg >= 13, schema >= 1.12, pgmetrics >= 1.13.0
type AnalyzeProgressBackend struct {
	PID                     int    `json:"pid"`
	DBName                  string `json:"db_name"`
	TableOID                int    `json:"table_oid"`
	Phase                   string `json:"phase"`
	SampleBlocksTotal       int64  `json:"sample_blks_total"`
	SampleBlocksScanned     int64  `json:"sample_blks_scanned"`
	ExtStatsTotal           int64  `json:"ext_stats_total"`
	ExtStatsComputed        int64  `json:"ext_stats_computed"`
	ChildTablesTotal        int64  `json:"child_tables_total"`
	ChildTablesDone         int64  `json:"child_tables_done"`
	CurrentChildTableRelOID int    `json:"child_oid"`
}

// BasebackupProgressBackend represents a row (and each row represents one
// backend) from pg_stat_progress_basebackup.
//
// pg >= 13, schema >= 1.12, pgmetrics >= 1.13.0
type BasebackupProgressBackend struct {
	PID                 int    `json:"pid"`
	Phase               string `json:"phase"`
	BackupTotal         int64  `json:"backup_total"`
	BackupStreamed      int64  `json:"backup_streamed"`
	TablespacesTotal    int64  `json:"tablespaces_total"`
	TablespacesStreamed int64  `json:"tablespaces_streamed"`
}

// ClusterProgressBackend represents a row (and each row represents one
// backend) from pg_stat_progress_cluster.
//
// pg >= 12, schema >= 1.12, pgmetrics >= 1.13.0
type ClusterProgressBackend struct {
	PID               int    `json:"pid"`
	DBName            string `json:"db_name"`
	TableOID          int    `json:"table_oid"`
	Command           string `json:"command"`
	Phase             string `json:"phase"`
	ClusterIndexOID   int    `json:"cluser_index_oid"`
	HeapTuplesScanned int64  `json:"heap_tuples_scanned"`
	HeapTuplesWritten int64  `json:"heap_tuples_written"`
	HeapBlksTotal     int64  `json:"heap_blks_total"`
	HeapBlksScanned   int64  `json:"heap_blks_scanned"`
	IndexRebuildCount int    `json:"index_rebuild_count"`
}

// CopyProgressBackend represents a row (and each row represents one
// backend) from pg_stat_progress_copy.
//
// pg >= 14, schema >= 1.12, pgmetrics >= 1.13.0
type CopyProgressBackend struct {
	PID             int    `json:"pid"`
	DBName          string `json:"db_name"`
	TableOID        int    `json:"table_oid"`
	Command         string `json:"command"`
	Type            string `json:"type"`
	BytesProcessed  int64  `json:"bytes_processed"`
	BytesTotal      int64  `json:"bytes_total"`
	TuplesProcessed int64  `json:"tuples_processed"`
	TuplesExcluded  int64  `json:"tuples_excluded"`
}

// CreateIndexProgressBackend represents a row (and each row represents one
// backend) from pg_stat_progress_create_index.
//
// pg >= 12, schema >= 1.12, pgmetrics >= 1.13.0
type CreateIndexProgressBackend struct {
	PID              int    `json:"pid"`
	DBName           string `json:"db_name"`
	TableOID         int    `json:"table_oid"`
	IndexOID         int    `json:"index_oid"`
	Command          string `json:"command"`
	Phase            string `json:"phase"`
	LockersTotal     int64  `json:"lockers_total"`
	LockersDone      int64  `json:"lockers_done"`
	CurrentLockerPID int    `json:"current_locker_pid"`
	BlocksTotal      int64  `json:"blocks_total"`
	BlocksDone       int64  `json:"blocks_done"`
	TuplesTotal      int64  `json:"tuples_total"`
	TuplesDone       int64  `json:"tuples_done"`
	PartitionsTotal  int64  `json:"partitions_total"`
	PartitionsDone   int64  `json:"partitions_done"`
}
