package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/rfyiamcool/go-fusehook/pkg/blkio"
	"github.com/rfyiamcool/go-fusehook/pkg/hookfs"
)

var (
	readRate  = flag.String("read-rate", "", "disk read bytes ratelimit ( 200kb/s, 1mb/2s, 100mb/m ), default no limit ")
	writeRate = flag.String("write-rate", "", "disk write bytes ratelimit ( 200kb/s, 1mb/s, 100mb/m ) ")
	logLevel  = flag.String("level", "info", "setting log print level")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func checkFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "blkio")
		fmt.Fprintf(os.Stderr, "blkio  [OPTIONS]  MOUNTPOINT  ORIGINAL\n")
		fmt.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
}

type blkioHook struct {
	original   string
	mountpoint string
	fuser      *fuse.Server
}

func newBlkioHook(original, mountpoint string) blkioHook {
	return blkioHook{
		original:   original,
		mountpoint: mountpoint,
	}
}

func (v *blkioHook) start() {
	hs := blkio.Hook{}
	rrater, rdur := blkio.ParseRateParam(*readRate)
	if rrater > 0 && rdur > 0 {
		hs.ReadLimiter = rate.NewLimiter(rate.Every(rdur), int(rrater))
		log.Infof("read ratelimit, every: %s, size: %v", rdur.String(), rrater)
	}

	wrater, wdur := blkio.ParseRateParam(*writeRate)
	if wrater > 0 && wdur > 0 {
		hs.WriteLimiter = rate.NewLimiter(rate.Every(wdur), int(wrater))
		log.Infof("write ratelimit, every: %s, size: %v", wdur.String(), wrater)
	}

	fs, err := hookfs.NewHookFs(v.original, v.mountpoint, &hs)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("start hookfs %s", fs)
	log.Infof("Please run `fusermount -u %s` after using this, manually", v.mountpoint)
	v.fuser, err = fs.ServeAsync()
	if err != nil {
		log.Errorf("failed to hookfs serve, err: %s", err.Error())
	}
}

func (v *blkioHook) stop() {
	// umount mountpoint path
	if v.fuser == nil {
		return
	}

	err := v.fuser.Unmount()
	if err == nil {
		log.Infof("umount path %s successfully", v.mountpoint)
		return
	}

	// retry by fusermount3
	exec.Command("fusermount3", "-u", v.mountpoint).Run()
}

func (v *blkioHook) wait() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP)
	for {
		s := <-sigch
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGHUP:
			log.Infof("receive signal %s", s)
			v.stop()
			return

		default:
			log.Infof("receive unknown signal %s", s)
			continue
		}
	}
}

func main() {
	checkFlags()

	mountpoint := flag.Arg(0)
	original := flag.Arg(1)
	if mountpoint == "" || original == "" {
		log.Fatal("mountpoint or original must not be empty")
	}

	// setting log
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetLevel(level)

	// check fuse3
	cmd := exec.Command("which", "fusermount")
	err = cmd.Run()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		log.Fatal("can't find fusermount command in path, please install fuse3")
	}

	hook := newBlkioHook(original, mountpoint)
	hook.start()
	hook.wait()
}
