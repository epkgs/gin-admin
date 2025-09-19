package config

type Captcha struct {
	Length int `default:"4"`
	Width  int `default:"400"`
	Height int `default:"160"`
}

type Prometheus struct {
	Enable         bool
	Port           int    `default:"9100"`
	BasicUsername  string `default:"admin"`
	BasicPassword  string `default:"admin"`
	LogApis        []string
	LogMethods     []string
	DefaultCollect bool
}

type Swagger struct {
	Enable bool `default:"false"` // disable swagger
}

type Pprof struct {
	Addr string `default:""` // pprof monitor address, e.g., "localhost:6060"
}

type Menu struct {
	File        string // Data to restore model.Menus (JSON/YAML)
	DenyOperate bool   // Deny operate menu
}

type RateLimiter struct {
	Period             int `default:"10"` // seconds
	MaxRequestsPerIP   int `default:"1000"`
	MaxRequestsPerUser int `default:"500"`
}

type Jwt struct {
	SigningMethod string `default:"HS512"`    // HS256/HS384/HS512
	SigningKey    string `default:"XnEsT0S@"` // secret key
	Expired       int    `default:"86400"`    // seconds
}

type Casbin struct {
	Disable          bool
	LoadThread       int    `default:"2"`
	AutoLoadInterval int    `default:"3"` // seconds
	ModelFile        string `default:"rbac_model.conf"`
	GenPolicyFile    string `default:"gen_rbac_policy.csv"`
}

type Logger struct {
	Level      string `default:"info"` // debug/info/warn/error
	CallerSkip int    `default:"1"`
	Console    struct {
		Enable bool `default:"true"`
	}
	File struct {
		Enable     bool   `default:"false"`
		Path       string `default:"log/ginadmin.log"`
		MaxSize    int    `default:"20"` // Maximum number of backup log files
		MaxBackups int    `default:"64"` // Maximum size of each log file in MB
	}
	Database struct {
		Enable bool `default:"false"`
	}
	MaxRequestLen  int `default:"4096"`
	MaxResponseLen int `default:"1024"`
}
