package test

import "os"

func EnvLoad() {
	// 配置环境变量
	// IKUBEOPS_APP_HTTP_PORT
	os.Setenv("IKUBEOPS_APP_HTTP_ADDR", "127.0.0.1")
	os.Setenv("IKUBEOPS_APP_HTTP_PORT", "1111")
	os.Setenv("IKUBEOPS_APP_GRPC_ADDR", "127.0.0.1")
	os.Setenv("IKUBEOPS_APP_GRPC_PORT", "2222")
	os.Setenv("IKUBEOPS_APP_LANGUAGE", "zh-cn")

	os.Setenv("IKUBEOPS_LOG_OUTPUT", "console")
	os.Setenv("IKUBEOPS_LOG_FORMAT", "text")
	os.Setenv("IKUBEOPS_LOG_LEVEL", "info")
	os.Setenv("IKUBEOPS_LOG_DEV", "true")
	os.Setenv("IKUBEOPS_LOG_FILE_PATH", "logs")
	os.Setenv("IKUBEOPS_LOG_MAX_SIZE", "100")
	os.Setenv("IKUBEOPS_LOG_MAX_AGE", "30")
	os.Setenv("IKUBEOPS_LOG_MAX_BACKUPS", "10")

	os.Setenv("IKUBEOPS_MYSQL_HOST", "127.0.0.1")
	os.Setenv("IKUBEOPS_MYSQL_PORT", "3306")
	os.Setenv("IKUBEOPS_MYSQL_USER", "root")
	os.Setenv("IKUBEOPS_MYSQL_PASSWORD", "123456")
	os.Setenv("IKUBEOPS_MYSQL_DB_NAME", "ikubeops")
	os.Setenv("IKUBEOPS_MYSQL_MAX_OPEN_CONNS", "100")
	os.Setenv("IKUBEOPS_MYSQL_MAX_IDLE_CONNS", "20")
	os.Setenv("IKUBEOPS_MYSQL_CONFIG", "charset=utf8mb4&parseTime=True&loc=Local")
	os.Setenv("IKUBEOPS_MYSQL_LOG_MODE", "false")
	os.Setenv("IKUBEOPS_MYSQL_LOG_ZAP", "false")

	os.Setenv("IKUBEOPS_REDIS_HOST", "127.0.0.1")
	os.Setenv("IKUBEOPS_REDIS_PORT", "6379")
	os.Setenv("IKUBEOPS_REDIS_DB", "0")
	os.Setenv("IKUBEOPS_REDIS_PASSWORD", "")
	os.Setenv("IKUBEOPS_REDIS_POOL_SIZE", "100")
}
