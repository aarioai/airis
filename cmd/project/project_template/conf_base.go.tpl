package conf

const (
	MongoCfgSection            = "mongodb_{{APP_NAME}}"
	MysqlCfgSection              = "mysql_{{APP_NAME}}"
	RedisCfgSection              = "redis_{{APP_NAME}}"
	RedisNoExpireCfgSection    = "redis_persistent"
)
