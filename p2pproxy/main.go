package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FactomProject/factomd/modules/registry"
	"github.com/FactomProject/factomd/modules/worker"
	"github.com/FactomProject/factomd/p2p"
)

var p2pProxy *P2PProxy

var internalReceived int
var externalReceived int

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

	const load = 5000000
	go fakeNet.ApplyNetworkIncomingLoad(load)

	go func() {
		start := time.Now()
		for internalReceived < load {
		}
		fmt.Println("down:", time.Since(start), load/time.Since(start).Seconds(), "/sec")
	}()

	go func() {
		for i := 0; i < load; i++ {
			fm := FactomMessage{PeerHash: p2p.Broadcast, Message: []byte(fmt.Sprintf("factomd load-%d/%d", i, load))}
			p2pProxy.BroadcastOut <- fm
		}
	}()

	go func() {
		start := time.Now()
		for externalReceived < load {
		}
		fmt.Println("up:", time.Since(start), load/time.Since(start).Seconds(), "/sec")

	}()

	select {}
}
