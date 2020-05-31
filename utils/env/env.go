package env

type Config struct {
	LogLevel   string `envconfig:"LOG_LEVEL" default:"info"`
	ServerAddr string `envconfig:"ISSUINGSERVICE_ADDR" default:":8080"`
	Service    Service
	DB         DB
}

type Service struct {
	ProcessEnv   string `envconfig:"PROCESS_ENV" default:"dev"`
	CloudService string `envconfig:"CLOUD_SERVICE" default:"aws"`
}

type DB struct {
	Host     string `envconfig:"DB_HOST" default:"db"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"lastrust"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Database string `envconfig:"DB_NAME" default:"postgres"`

	MaxOpenConns int `envconfig:"DB_MAXOPENCONNS" default:"10"`
	MaxIdleConns int `envconfig:"DB_MAXIDLECONNS" default:"10"`
}
