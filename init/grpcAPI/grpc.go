package grpcAPI

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sheginabo/go-quick-grpc/internal/pb"
	"github.com/sheginabo/go-quick-grpc/internal/presentation/handlers"
	"github.com/sheginabo/go-quick-grpc/internal/presentation/interceptors"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Module struct {
	GrpcApi    *handlers.GrpcApi
	GrpcServer *grpc.Server
	Listener   net.Listener
}

func NewModule() *Module {

	// TODO register logger Interceptor
	loggerInterceptor := grpc.UnaryInterceptor(interceptors.GrpcLogger)

	gAPIModule := &Module{
		GrpcApi:    handlers.NewGrpcApi(),
		GrpcServer: grpc.NewServer(loggerInterceptor),
	}

	return gAPIModule
}

// Run grpc server
func (module *Module) Run(ctx context.Context, waitGroup *errgroup.Group) {
	var err error
	module.Listener, err = net.Listen("tcp", viper.GetString("GRPC_SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// Register your services here
	pb.RegisterGoQuickGRPCServer(module.GrpcServer, module.GrpcApi)
	if viper.Get("ENV") != "prod" {
		reflection.Register(module.GrpcServer) // 使用反射可以使用 grpcurl debug, prod 不建議用
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", module.Listener.Addr().String())
		err = module.GrpcServer.Serve(module.Listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			module.GrpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Info().Msg("graceful shutdown gRPC server done")
		//case t := <-time.After(10 * time.Second): // 10s timeout 1 缺點只提供時間戳記
		//	log.Warn().Msgf("graceful shutdown timed out at %v, forcing stop", t)
		//	module.GrpcServer.Stop()
		//	return errors.New("gRPC server shutdown timed out")
		case <-shutdownCtx.Done(): // 10s timeout 2 提供更多資訊
			log.Warn().Err(shutdownCtx.Err()).Msg("graceful shutdown timed out, forcing stop")
			module.GrpcServer.Stop()
			return fmt.Errorf("graceful shutdown timed out: %w", shutdownCtx.Err())
		}

		return nil
	})
}
