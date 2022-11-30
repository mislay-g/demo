package config

type Config struct {
	Db     *DBConfig
	Redis  *RedisConfig
	Log    *LogConfig
	Server *ServerConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	Mode         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DBConfig struct {
	Debug        bool
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

type RedisConfig struct {
	Address  string
	Username string
	Password string
	DB       int
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type CacheConfig struct {
	Expired       int
	CleanInterval int
}
