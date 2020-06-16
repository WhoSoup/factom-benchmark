package main

import (
	"fmt"

	"github.com/FactomProject/factomd/p2p"
)

type FakeNet struct {
	Load chan *p2p.Parcel
}

func NewFakeNet() *FakeNet {
	f := new(FakeNet)
	f.Load = make(chan *p2p.Parcel, 5000)
	return f
}

func (f *FakeNet) Send(parcel *p2p.Parcel) {
	externalReceived++
}

func (f *FakeNet) Reader() <-chan *p2p.Parcel {
	return f.Load
}

func (f *FakeNet) ApplyNetworkIncomingLoad(size int) {
	for i := 0; i < size; i++ {
		f.Load <- p2p.NewParcel(p2p.Broadcast, []byte(fmt.Sprintf("network load-%d/%d", i, size)))
	}
}
