package configs

type Middleware struct {
	Recovery struct {
		Skip int `default:"3"` // skip the first n stack frames
	}
	CORS struct {
		Enable                 bool
		AllowAllOrigins        bool
		AllowOrigins           []string
		AllowMethods           []string
		AllowHeaders           []string
		AllowCredentials       bool
		ExposeHeaders          []string
		MaxAge                 int
		AllowWildcard          bool
		AllowBrowserExtensions bool
		AllowWebSockets        bool
		AllowFiles             bool
	}
	Trace struct {
		RequestHeaderKey string `default:"X-Request-Id"`
		ResponseTraceKey string `default:"X-Trace-Id"`
	}
	Logger struct {
		MaxOutputRequestBodyLen  int `default:"4096"`
		MaxOutputResponseBodyLen int `default:"1024"`
	}
	CopyBody struct {
		MaxContentLen int64 `default:"33554432"` // max content length (default 32MB)
	}
	Auth struct {
		Disable       bool
		SigningMethod string `default:"HS512"`    // HS256/HS384/HS512
		SigningKey    string `default:"XnEsT0S@"` // secret key
		RefreshKey    string `default:"TfYOj6Tf"` // secret key for refresh token
		OldSigningKey string // old secret key (for migration)
		Expired       int    `default:"86400"` // seconds
		Store         struct {
			Type      string `default:"memory"` // memory/badger/redis
			Delimiter string `default:":"`      // delimiter for key
			Memory    struct {
				CleanupInterval int `default:"60"` // seconds
			}
			Badger struct {
				Path string `default:"data/auth"`
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	RateLimiter struct {
		Enable             bool
		Period             int // seconds
		MaxRequestsPerIP   int
		MaxRequestsPerUser int
		Store              struct {
			Type   string // memory/redis
			Memory struct {
				Expiration      int `default:"3600"` // seconds
				CleanupInterval int `default:"60"`   // seconds
			}
			Redis struct {
				Addr     string
				Username string
				Password string
				DB       int
			}
		}
	}
	Casbin struct {
		Disable          bool
		LoadThread       int    `default:"2"`
		AutoLoadInterval int    `default:"3"` // seconds
		ModelFile        string `default:"rbac_model.conf"`
		GenPolicyFile    string `default:"data/gen_rbac_policy.csv"`
	}
	Static struct {
		ExcludedPathPrefixes []string // excluded path prefixes
		Root                 string   // static file root
	}
}
