/*
Copyright © 2024 Nected<dev@nected.ai>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"net"

	// "github.com/nected/sanchaalak/src/services"

	"github.com/Jille/raft-grpc-example/proto"
	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	"github.com/Jille/raftadmin"
	"github.com/nected/sanchaalak/src/config"
	"github.com/nected/sanchaalak/src/raft"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var bootstrap bool

// electCmd represents the elect command
var electCmd = &cobra.Command{
	Use:   "elect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("elect called")
		runElect()
	},
}

func init() {
	rootCmd.AddCommand(electCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// electCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// electCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	electCmd.Flags().BoolVarP(&bootstrap, "bootstrap", "b", false, "Bootstrap the cluster")
}

func runElect() {
	context := context.Background()
	config := config.GetConfig()
	wt := raft.WordTracker{}

	_, port, err := net.SplitHostPort(config.Raft.NodeInfo.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	sock, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	rSvr, err := raft.NewRaftServer(context, config.Raft.NodeInfo.ID, config.Raft.NodeInfo.Address, &wt, bootstrap)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()

	rpcInterface := raft.NewRpcInterface(&wt, rSvr.Raft())

	proto.RegisterExampleServer(s, rpcInterface)

	rSvr.TransportManager().Register(s)

	leaderhealth.Setup(rSvr.Raft(), s, []string{"Example"})
	raftadmin.Register(s, rSvr.Raft())
	reflection.Register(s)

	if err := s.Serve(sock); err != nil {
		fmt.Println(err)
		return
	}
}
