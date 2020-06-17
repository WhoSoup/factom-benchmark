package main

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/state"
)

type FNode struct {
	NodeName string
	State    *State
	Peers    []interfaces.IPeer
	P2PIndex int
}

type State struct {
	DBFinished                    bool
	RecentMessage                 *state.RecentMessage
	ChainCommits                  *state.Last100
	Reveals                       *state.Last100
	MissingMessageResponseHandler *state.MissingMessageResponseCache
}

func (s *State) GetTrueLeaderHeight() int              { return 0 }
func (s *State) GetHighestCompletedBlk() int           { return 0 }
func (s *State) InMsgQueue() interfaces.IQueue         { return nil }
func (s *State) InMsgQueue2() interfaces.IQueue        { return nil }
func (s *State) NetworkOutMsgQueue() interfaces.IQueue { return nil }
func (s *State) APIQueue() interfaces.IQueue           { return nil }

func (s *State) LogPrintf(name, format string, more ...interface{})   {}
func (s *State) LogMessage(name, comment string, msg interfaces.IMsg) {}
func (s *State) DataMsgQueue() chan interfaces.IMsg                   { return nil }

func (s *State) NetworkInvalidMsgQueue() chan interfaces.IMsg { return nil }
