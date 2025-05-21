package main

import (
	"bringauto/modules/bringauto_log"
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_process"
	"bringauto/modules/bringauto_error"
	"os"
	"time"
	"syscall"
)

func main() {
	var err error
	var args CmdLineArgs
	logger := bringauto_prerequisites.CreateAndInitialize[bringauto_log.Logger](time.Now(), "./log")

	args.InitFlags()
	err = args.ParseArgs(os.Args)
	if err != nil {
		logger.Error("Can't parse cmd line arguments - %s",)
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
			if err != bringauto_error.GitLfsErr {
				os.Exit(bringauto_error.CREATING_SYSROOT_ERROR)
			}
			os.Exit(bringauto_error.GIT_LFS_ERROR)
		}
		return
	}

	return
}
