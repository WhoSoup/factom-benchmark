package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FactomProject/factomd/modules/registry"
	"github.com/FactomProject/factomd/modules/worker"
)

var p2pProxy *P2PProxy

var internalReceived int

func Sender() {
	for range p2pProxy.BroadcastIn {
		//fmt.Println("<-", n)
		internalReceived++
	}
}

func Receiver() {
	for n := range p2pProxy.BroadcastOut {
		fmt.Println(n)
	}
}

func main() {
	worker.AddInterruptHandler(func() { os.Exit(0) })
	fakeNet := NewFakeNet()

	p := registry.New()
	p.Register(func(w *worker.Thread) {
		nodeName := "FNode0"
		p2pProxy = new(P2PProxy).Initialize(nodeName, "P2P Network").(*P2PProxy)
		p2pProxy.StartProxy(w)
		p2pProxy.Network = fakeNet
	})
	go p.Run()

	time.Sleep(time.Second)

	go Sender()

	load := 5000000
	go fakeNet.ApplyNetworkIncomingLoad(load)
	start := time.Now()
	for internalReceived < load {
	}
	fmt.Println(time.Since(start))

	select {}
}
