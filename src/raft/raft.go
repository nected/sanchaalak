package raft

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RaftServer struct {
	raft             *raft.Raft
	transportManager *transport.Manager
}

func NewRaftServer(ctx context.Context, nodeID, nodeAddress string, fsm raft.FSM, raftBootstrap bool) (server *RaftServer, err error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(nodeID)
	server = &RaftServer{}

	baseDir := filepath.Join("/tmp/", nodeID)

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	tm := transport.New(raft.ServerAddress(nodeAddress), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if raftBootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(nodeID),
					Address:  raft.ServerAddress(nodeAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}
	server.raft = r
	server.transportManager = tm
	return server, nil
}

func (r *RaftServer) Raft() *raft.Raft {
	return r.raft
}

func (r *RaftServer) TransportManager() *transport.Manager {
	return r.transportManager
}
