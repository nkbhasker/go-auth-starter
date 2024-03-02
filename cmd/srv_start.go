package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkbhasker/go-auth-starter/config"
	"github.com/nkbhasker/go-auth-starter/internal/api"
	"github.com/nkbhasker/go-auth-starter/internal/comm"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/misc"
	"github.com/nkbhasker/go-auth-starter/internal/repo"
	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
	"github.com/spf13/cobra"

	_ "github.com/nkbhasker/go-auth-starter/internal/errors"
)

var srvStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	RunE: func(cmd *cobra.Command, _args []string) error {
		return SrvStart(cmd.Version)
	},
}

func SrvStart(version string) error {
	cfg, err := config.InitSrvConfig()
	if err != nil {
		return err
	}
	idGenerator := uid.NewIdGenerator()
	dbStore, err := storage.InitDBStore(cfg.PostgresUrl)
	if err != nil {
		return err
	}
	defer dbStore.CloseDB()
	cacheStore, err := storage.InitCacheStore(cfg.RedisUrl)
	if err != nil {
		return err
	}
	defer cacheStore.CloseDB()
	jwtHelper, err := misc.NewJwtHelper(cfg.Host, cfg.JwtPrivateKey)
	if err != nil {
		return err
	}
	repos := repo.NewRepo(repo.RepoOptions{
		DBStore:                    dbStore,
		CacheStore:                 cacheStore,
		IdGenerator:                idGenerator,
		JwtHelper:                  jwtHelper,
		AccessTokenExpiryInMinutes: cfg.AccessTokenExpiryInMinutes,
		OtpExpiryInMinutes:         cfg.OtpExpiryInMinutes,
	})
	awsSession, err := core.NewAwsSession(core.AwsSessionOptions{
		Region:          cfg.AwsRegion,
		AccessKeyId:     cfg.AwsAccessKeyId,
		SecretAccessKey: cfg.AwsSecretAccessKey,
	})
	if err != nil {
		return err
	}
	emailer, err := comm.NewEmailer(comm.NewSES(awsSession.Session, cfg.AwsSesSender))
	if err != nil {
		return err
	}
	app := core.NewApp(core.AppOption{
		Version:     version,
		DBStore:     dbStore,
		CacheStore:  cacheStore,
		Repos:       repos,
		IdGenerator: idGenerator,
		Emailer:     emailer,
		Validate:    core.NewValidate(),
	})
	otpGenerateRateLimiter := core.NewRateLimiter(
		cacheStore,
		core.RateLimiterKindOtpGenerate,
		cfg.OtpGenerateRateLimit,
		cfg.OtpGenerateRateLimitWindow,
	)
	otpVerifyRateLimiter := core.NewRateLimiter(
		cacheStore,
		core.RateLimiterKindOtpVerify,
		cfg.OtpVerifyRateLimit,
		cfg.OtpVerifyRateLimitWindow,
	)
	handler := api.SetupRouter(api.RouterOptions{
		App:                    app,
		JwtHelper:              jwtHelper,
		OtpGenerateRateLimiter: otpGenerateRateLimiter,
		OtpVerifyRateLimiter:   otpVerifyRateLimiter,
	})
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Server run context
	srvCtx, serverStopCtx := context.WithCancel(context.Background())

	errch := make(chan error, 1)
	sigch := make(chan os.Signal, 1)
	// Listen for syscall signals for process to interrupt/quit
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigch
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelShutdownCtx := context.WithTimeout(srvCtx, 30*time.Second)
		defer cancelShutdownCtx()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				errch <- errors.New("graceful shutdown timed out.. forcing exit")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			errch <- err
		}
		serverStopCtx()
	}()

	log.Printf("Server running on port %s\n", cfg.Port)
	log.Println("Press ctrl+c to stop")
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errch <- err
	} else {
		log.Println("Server stopped")
	}
	// Wait to be marked done
	<-srvCtx.Done()
	// Either return received error or nil
	select {
	case err := <-errch:
		return err
	default:
		return nil
	}
}
