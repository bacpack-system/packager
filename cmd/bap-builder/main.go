package main

import (
	"github.com/bacpack-system/packager/internal/bringauto_log"
	"github.com/bacpack-system/packager/internal/bringauto_prerequisites"
	"github.com/bacpack-system/packager/internal/bringauto_process"
	"github.com/bacpack-system/packager/internal/bringauto_error"
	"os"
	"time"
	"syscall"
	"fmt"
)

func main() {
	var args CmdLineArgs
	logger, err := bringauto_prerequisites.CreateAndInitialize[bringauto_log.Logger](time.Now(), "./log")
	if err != nil {
		panic(fmt.Errorf("cannot initialize Logger - %w", err))
	}

	args.InitFlags()
	err = args.ParseArgs(os.Args)
	if err != nil {
		logger.Error("Can't parse cmd line arguments - %s", err)
		os.Exit(bringauto_error.CMD_LINE_ERROR)
	}
	bringauto_process.SignalHandlerRegisterSignal(syscall.SIGINT)

	if args.BuildImage {
		err = BuildDockerImage(&args.BuildImagesArgs, *args.Context)
		if err != nil {
			logger.Error("Failed to build Docker image: %s", err)
			os.Exit(bringauto_error.GetReturnCode(err))
		}
		return
	}

	if args.BuildPackage {
		err = BuildPackage(&args.BuildPackageArgs, *args.Context)
		if err != nil {
			logger.Error("Failed to build package: %s", err)
			os.Exit(bringauto_error.GetReturnCode(err))
		}
		return
	}
	if args.BuildApp {
		err = BuildApp(&args.BuildAppArgs, *args.Context)
		if err != nil {
			logger.Error("Failed to build App: %s", err)
			os.Exit(bringauto_error.GetReturnCode(err))
		}
		return
	}
	if args.CreateSysroot {
		err = CreateSysroot(&args.CreateSysrootArgs, *args.Context)
		if err != nil {
			logger.Error("Failed to create sysroot: %s", err)
			os.Exit(bringauto_error.GetReturnCode(err))
		}
		return
	}

	return
}
