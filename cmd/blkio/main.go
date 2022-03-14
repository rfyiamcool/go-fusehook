package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] MOUNTPOINT ORIGINAL...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	mountpoint := flag.Arg(0)
	original := flag.Arg(1)

	// setting log
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.SetLevel(level)

	serve(original, mountpoint)
}

func serve(original string, mountpoint string) {
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

	fs, err := hookfs.NewHookFs(original, mountpoint, &hs)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Serving %s", fs)
	log.Infof("Please run `fusermount -u %s` after using this, manually", mountpoint)
	if err = fs.Serve(); err != nil {
		log.Fatal(err)
	}
}
