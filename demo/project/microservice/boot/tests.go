package boot

import (
	"sync"

	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/alog"
)

func SelfTest(app *aa.App) bool {
	checks := []func(*aa.App, chan<- error){
		//driver.CheckMongodbHealth(conf.MongoCfgSection),
		//driver.CheckMySQLHealth(conf.MysqlCfgSection),
		//driver.CheckRabbitmqHealth(conf.RabbitmqCfgSection),
		//driver.CheckRedisHealth(conf.RedisCfgSection),

		//driver.CheckMongodbHealth(infraconf.MongoCfgSection),
		//driver.CheckMySQLHealth(infraconf.MysqlCfgSection),
		//driver.CheckRabbitmqHealth(infraconf.RabbitmqCfgSection),
		//driver.CheckRedisHealth(infraconf.RedisCfgSection),
	}

	errChan := make(chan error, len(checks))
	var wg sync.WaitGroup
	for _, check := range checks {
		wg.Add(1)
		go func(c func(*aa.App, chan<- error)) {
			defer wg.Done()
			c(app, errChan)
		}(check)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			acontext.SetServUnhealthy()
			alog.Errorf("health check failed: %v", err)
		}
	}

	acontext.ServFallbackReady()
	return acontext.ServHealth().IsReady()
}

// @TODO check posgresql/kafka/rabbitmq... here, and then set to acontext ServiceHealth
