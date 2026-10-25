package acontext

import "sync/atomic"

type ServHealthEnum string

const (
	ServHealthStarting ServHealthEnum = "starting"  // not ready to serve, checking dependencies in bootstrap
	ServHeathReady     ServHealthEnum = "ready"     // ready to serve, all dependencies are ready
	ServHeathStopping  ServHealthEnum = "stopping"  // service is stopping before restarting
	ServHeathUnhealthy ServHealthEnum = "unhealthy" // service is unhealthy
)

var (
	servHealth atomic.Value
	_          = func() struct{} {
		servHealth.Store(ServHealthStarting)
		return struct{}{}
	}()
)

// No need to overthink the issue of concurrent modifications more

func SetServStarting()  { servHealth.Store(ServHealthStarting) }
func SetServReady()     { servHealth.Store(ServHeathReady) }
func SetServStopping()  { servHealth.Store(ServHeathStopping) }
func SetServUnhealthy() { servHealth.Store(ServHeathUnhealthy) }

// ServFallbackReady only if the status is still starting, no any error issues, fallback to set to ready
func ServFallbackReady() {
	if ServHealth().IsStarting() {
		servHealth.Store(ServHeathReady)
	}
}

func ServHealth() ServHealthEnum {
	if val := servHealth.Load(); val != nil {
		return val.(ServHealthEnum)
	}
	return ServHealthStarting
}

func (h ServHealthEnum) String() string    { return string(h) }
func (h ServHealthEnum) IsStarting() bool  { return h == ServHealthStarting }
func (h ServHealthEnum) IsReady() bool     { return h == ServHeathReady }
func (h ServHealthEnum) IsStopping() bool  { return h == ServHeathStopping }
func (h ServHealthEnum) IsUnhealthy() bool { return h == ServHeathUnhealthy }
