package main

import (
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"

	peer "github.com/astranet/meshRPC-benchmark/peer/service"
	"github.com/astranet/meshRPC/cluster"
)

var (
	peerName = app.String(cli.StringOpt{
		Name:   "P peer-name",
		Desc:   "Specify peer name.",
		EnvVar: "APP_PEER_NAME",
	})
	allPeers = app.Strings(cli.StringsOpt{
		Name:   "A all-peers",
		Desc:   "All peer names list.",
		EnvVar: "APP_ALL_PEERS",
		Value:  []string{},
	})
	clusterNodes = app.Strings(cli.StringsOpt{
		Name:   "N nodes",
		Desc:   "A list of cluster nodes to join for service discovery.",
		EnvVar: "MESHRPC_CLUSTER_NODES",
		Value:  []string{},
	})
	clusterName = app.String(cli.StringOpt{
		Name:   "T tag",
		Desc:   "Cluster tag name.",
		EnvVar: "MESHRPC_CLUSTER_TAGNAME",
		Value:  "benchmark",
	})
	netAddr = app.String(cli.StringOpt{
		Name:   "listen-addr",
		Desc:   "Listen address for cluster discovery and private networking.",
		EnvVar: "MESHRPC_LISTEN_ADDR",
		Value:  "0.0.0.0:0",
	})
)

var app = cli.App("peer", "A Peer service server for meshRPC cluster, for benchmark purposes")

func main() {
	app.Action = func() {
		defer closer.Close()

		c := cluster.NewAstraCluster(*peerName, &cluster.AstraOptions{
			Tags: []string{
				*clusterName,
			},
			Nodes: *clusterNodes,
			Debug: true,
		})

		serviceMap := make(map[string]peer.Service, len(*allPeers))
		for _, serviceName := range *allPeers {
			peerClient := c.NewClient(serviceName, peer.RPCHandlerSpec)
			serviceMap[serviceName] = peer.NewServiceClient(peerClient, nil)
			log.Println("initialized RPC client for", serviceName)
		}
		service := peer.NewService(*peerName, serviceMap)
		meshRPC := peer.NewRPCHandler(service, nil)
		log.Println("service published as", *peerName)
		c.Publish(meshRPC)

		if err := c.ListenAndServe(*netAddr); err != nil {
			log.Fatalln(err)
		}
		closer.Hold()
	}
	app.Run(os.Args)
}
