package configs

type Cache struct {
	Type      string `default:"memory"` // memory/badger/redis
	Delimiter string `default:":"`      // delimiter for key
	Memory    struct {
		CleanupInterval int `default:"60"` // seconds
	}
	Badger struct {
		Path string `default:"data/cache"`
	}
	Redis struct {
		Addr     string
		Username string
		Password string
		DB       int
	}

	Expiration struct { // Expiration times for various cache entries
		User int `default:"4"` // User cache expiration time in hours
	}
}

type DB struct {
	Debug        bool
	Type         string `default:"sqlite3"`     // sqlite3/mysql/postgres
	DSN          string `default:"data/app.db"` // database source name
	MaxLifetime  int    `default:"86400"`       // seconds
	MaxIdleTime  int    `default:"3600"`        // seconds
	MaxOpenConns int    `default:"100"`         // connections
	MaxIdleConns int    `default:"50"`          // connections
	TablePrefix  string `default:""`
	AutoMigrate  bool
	PrepareStmt  bool
	Resolver     []struct {
		DBType   string   // sqlite3/mysql/postgres
		Sources  []string // DSN
		Replicas []string // DSN
		Tables   []string
	}
}

type Upload struct {
	Path       string `default:"uploads"`
	UseDateDir bool   `default:"true"`
}
