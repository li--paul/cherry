/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package main

import (
	"flag"
	"fmt"
	"git.sds.co.kr/cherry.git/cherryd/internal/device"
	"git.sds.co.kr/cherry.git/cherryd/internal/log"
	"golang.org/x/net/context"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func initSyslog() (log.Logger, error) {
	log, err := syslog.New(syslog.LOG_ERR|syslog.LOG_DAEMON, "")
	if err != nil {
		return nil, err
	}

	return log, nil
}

func waitSignal(log log.Logger, shutdown context.CancelFunc) {
	c := make(chan os.Signal, 5)
	// All incoming signals will be transferred to the channel
	signal.Notify(c)

	for {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			// Graceful shutdown
			log.Info("Shutting down...")
			shutdown()
			time.Sleep(10 * time.Second) // let cancelation propagate
			log.Info("Halted")
			os.Exit(0)
		} else if s == syscall.SIGHUP {
			// XXX: Do something you need
			log.Debug("SIGHUP")
		}
	}
}

func listen(ctx context.Context, log log.Logger, config *Config) {
	type KeepAlive interface {
		SetKeepAlive(keepalive bool) error
		SetKeepAlivePeriod(d time.Duration) error
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ServerPort))
	if err != nil {
		log.Err(fmt.Sprintf("Failed to listen on %v port: %v", config.ServerPort, err))
		return
	}
	defer listener.Close()

	f := func(c chan<- net.Conn) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Err(fmt.Sprintf("Failed to accept a new connection: %v", err))
				continue
			}
			c <- conn
		}
	}

	// Infinite loop
	for {
		backlog := make(chan net.Conn)
		go f(backlog)

		select {
		case conn := <-backlog:
			if v, ok := conn.(KeepAlive); ok {
				log.Debug("Trying to enable socket keepalive..")
				if err := v.SetKeepAlive(true); err == nil {
					log.Debug("Setting socket keepalive period...")
					v.SetKeepAlivePeriod(time.Duration(30) * time.Second)
				} else {
					log.Err(fmt.Sprintf("Failed to enable socket keepalive: %v", err))
				}
			}

			manager := device.NewManager(log)
			go manager.Run(ctx, conn)

		case <-ctx.Done():
			return
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	conf := NewConfig()
	if err := conf.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configurations: %v\n", err)
		os.Exit(1)
	}

	log, err := initSyslog()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init syslog: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go waitSignal(log, cancel)
	listen(ctx, log, conf)
}
