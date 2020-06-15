package main

import (
	"github.com/FactomProject/factomd/common/adminBlock"
	"github.com/FactomProject/factomd/common/directoryBlock"
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	. "github.com/FactomProject/factomd/state"
)

func NewProcessList(s *State, dbheight uint32) *ProcessList {
	// We default to the number of Servers previous.   That's because we always
	// allocate the FUTURE directoryblock, not the current or previous...

	pl := new(state.ProcessList)
	pl.State = s

	// Make a copy of the previous FedServers
	pl.FedServers = make([]interfaces.IServer, 0)
	pl.AuditServers = make([]interfaces.IServer, 0)
	//pl.Requests = make(map[[20]byte]*Request)

	pl.FactoidBalancesTMutex.Lock()
	pl.FactoidBalancesT = map[[32]byte]int64{}
	pl.FactoidBalancesTMutex.Unlock()

	pl.ECBalancesTMutex.Lock()
	pl.ECBalancesT = map[[32]byte]int64{}
	pl.ECBalancesTMutex.Unlock()

	//		pl.FedServers = append(pl.FedServers, previous.FedServers...)
	//		pl.AuditServers = append(pl.AuditServers, previous.AuditServers...)
	for _, auditServer := range pl.AuditServers {
		auditServer.SetOnline(false)
		/*			if state.GetIdentityChainID().IsSameAs(auditServer.GetChainID()) {
					// Always consider yourself "online"
					auditServer.SetOnline(true)
				}*/
	}
	for _, fedServer := range pl.FedServers {
		fedServer.SetOnline(true)
	}
	pl.SortFedServers()

	now := s.GetTimestamp()
	// We just make lots of VMs as they have nearly no impact if not used.
	pl.VMs = make([]*state.VM, 65)
	for i := 0; i < 65; i++ {
		pl.VMs[i] = new(state.VM)
		pl.VMs[i].List = make([]interfaces.IMsg, 0)
		pl.VMs[i].Synced = true
		pl.VMs[i].WhenFaulted = 0
		pl.VMs[i].ProcessTime = now
		pl.VMs[i].VmIndex = i
		//pl.VMs[i].p = pl
	}

	pl.DBHeight = dbheight

	pl.MakeMap()

	pl.PendingChainHeads = NewSafeMsgMap("PendingChainHeads", pl.State)
	pl.OldMsgs = make(map[[32]byte]interfaces.IMsg)
	pl.OldAcks = make(map[[32]byte]interfaces.IMsg)

	pl.NewEBlocks = make(map[[32]byte]interfaces.IEntryBlock)
	pl.NewEntries = make(map[[32]byte]interfaces.IEntry)

	pl.DBSignatures = make([]DBSig, 0)

	// If a federated server, this is the server index, which is our index in the FedServers list

	var err error

	if previous != nil {
		pl.DirectoryBlock = directoryBlock.NewDirectoryBlock(previous.DirectoryBlock)
		pl.AdminBlock = adminBlock.NewAdminBlock(previous.AdminBlock)
		pl.EntryCreditBlock, err = entryCreditBlock.NextECBlock(previous.EntryCreditBlock)
	} else {
		if pl.DBHeight > 0 {
			pl.DirectoryBlock, _ = state.GetDB().FetchDBlockByHeight(pl.DBHeight)
			pl.AdminBlock, _ = state.GetDB().FetchABlockByHeight(pl.DBHeight)
			pl.EntryCreditBlock, _ = state.GetDB().FetchECBlockByHeight(pl.DBHeight)
		} else {
			pl.DirectoryBlock = directoryBlock.NewDirectoryBlock(nil)
			pl.AdminBlock = adminBlock.NewAdminBlock(nil)
			pl.EntryCreditBlock, err = entryCreditBlock.NextECBlock(nil)
		}
	}

	pl.ResetDiffSigTally()

	if pl.DirectoryBlock != nil {
		pl.DirectoryBlock.GetHeader().SetTimestamp(now) // Well this is awkwardly after it's created but ....
	}
	if err != nil {
		panic(err.Error())
	}

	return pl
}
