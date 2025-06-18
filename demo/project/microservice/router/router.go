package router

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/helpers/debug"
	infragrpc "project/microservice/app/infra/grpc"
	"project/microservice/sdk/infra"
)

func Serve(app *aa.App, prof *debug.Profile) {
	// Start service infra gRPC server
	_, _, err1 := infragrpc.New(app).Serve(prof)
	ae.PanicOnErrs(err1)

	// Init service infra gRPC client
	infra.New(app).Init(prof)
}
