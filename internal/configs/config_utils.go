package configs

type Captcha struct {
	Length    int    `default:"4"`
	Width     int    `default:"400"`
	Height    int    `default:"160"`
	CacheType string `default:"memory"` // memory/redis
	Redis     struct {
		Addr      string
		Username  string
		Password  string
		DB        int
		KeyPrefix string `default:"captcha:"`
	}
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
	Disable    bool   `default:"false"`                // disable swagger
	StaticFile string `default:"configs/openapi.json"` // static file for openapi.json
}

type Pprof struct {
	Addr string `default:""` // pprof monitor address, e.g., "localhost:6060"
}

type Menu struct {
	File        string // Data to restore model.Menus (JSON/YAML)
	DenyOperate bool   // Deny operate menu
}
