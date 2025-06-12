env = local
log_level = debug
timezone_id = Asia/Shanghai
time_format = 2006-02-01 15:04:05
mock = 0
text_config_dirs = ./config/rsa
text_config_rsa_checks=app-1024,app-2048,app-4096
pprof = 0

[app]
config_root = ./config
rate_limit = 100,500,1m,5m
view_rate_limit = 100,500,1m,5m
api_rate_limit = 100,500,1m,5m
log_buffer_size = 0
log_dir = ./storage/log
log_symlink = /var/log/app-{{APP_NAME}}.log
views_root = ./frontend/view
emb_root = ./storage/emb

[svc_test]
port = 80
grpc_port = 8000

;[consul]
;scheme = http
;address = aa-consul:8500
;path_prefix =
;data_center =
;username =
;password =
;wait_time =
;token =
;token_file =
;namespace =
;tls_address =
;tls_ca_file =
;tls_ca_path =
;tls_cert_file =
;tls_key_file =
;tls_insecure_skip_verify =

[mongodb]
hosts = aa-mongodb:27017
username =
password =
min_pool_size = 3
max_pool_size = 100

[mongodb_test]
db = test

[mysql]
host = aa-mysql:3306
user =
password =
tls = false
timeout = 5s,5s,5s
pool_max_idle_conns = 0
pool_max_open_conns = 0
pool_conn_max_life_time = 0
pool_conn_max_idle_time = 0

[mysql_test]
schema = test

[redis]
addr = aa-redis:6379
password =
db = 0
dial_timeout = 5s
timeout = 5s,3s,3s
pool_timeout=4s
conn_max_idle_time=30s
max_retries=3
min_retry_backoff=8ms
max_retry_backoff=512ms

[redis_test]
db = 15

[rabbitmq]
host = aa-rabbitmq:5672
user =
password =
vhost = /
channel_max = 64
mqtt_user =
mqtt_password =