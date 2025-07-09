package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	MongoDB     MongoDBConfig     `mapstructure:"mongodb"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Security    SecurityConfig    `mapstructure:"security"`
	Performance PerformanceConfig `mapstructure:"performance"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI         string        `mapstructure:"uri"`
	Database    string        `mapstructure:"database"`
	MaxPoolSize int           `mapstructure:"max_pool_size"`
	MinPoolSize int           `mapstructure:"min_pool_size"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
	Issuer             string        `mapstructure:"issuer"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	MaxLoginAttempts       int           `mapstructure:"max_login_attempts"`
	LockoutDuration        time.Duration `mapstructure:"lockout_duration"`
	PasswordMinLength      int           `mapstructure:"password_min_length"`
	SessionCleanupInterval time.Duration `mapstructure:"session_cleanup_interval"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	EnableTextSearch     bool          `mapstructure:"enable_text_search"`
	CacheUserPermissions bool          `mapstructure:"cache_user_permissions"`
	MaxQueryTime         time.Duration `mapstructure:"max_query_time"`
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 设置环境变量前缀和替换
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 先读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，继续使用默认值和环境变量
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "auth_center")
	viper.SetDefault("mongodb.max_pool_size", 100)
	viper.SetDefault("mongodb.min_pool_size", 10)
	viper.SetDefault("mongodb.timeout", "30s")

	viper.SetDefault("jwt.access_token_expire", "15m")
	viper.SetDefault("jwt.refresh_token_expire", "168h")
	viper.SetDefault("jwt.issuer", "AuthCenter")
	viper.SetDefault("jwt.secret", "change-this-secret-in-production")

	viper.SetDefault("security.max_login_attempts", 5)
	viper.SetDefault("security.lockout_duration", "30m")
	viper.SetDefault("security.password_min_length", 8)
	viper.SetDefault("security.session_cleanup_interval", "1h")
	viper.SetDefault("security.bcrypt_cost", 12)

	viper.SetDefault("performance.enable_text_search", true)
	viper.SetDefault("performance.cache_user_permissions", true)
	viper.SetDefault("performance.max_query_time", "30s")
}
