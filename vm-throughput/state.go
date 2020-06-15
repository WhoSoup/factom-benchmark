package main

import "github.com/FactomProject/factomd/common/interfaces"

type MockState struct {
}

var _ interfaces.IState = (*MockState)(nil)
