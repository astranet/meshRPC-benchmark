package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/astranet/meshRPC/cluster"
	"github.com/gin-gonic/gin"
	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"

	peer "github.com/astranet/meshRPC-benchmark/peer/service"
)

var (
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
		Value:  "0.0.0.0:11999",
	})
	httpListenHost = app.String(cli.StringOpt{
		Name:   "H http-host",
		Desc:   "Specify listen HTTP host.",
		EnvVar: "APP_HTTP_HOST",
		Value:  "0.0.0.0",
	})
	httpListenPort = app.String(cli.StringOpt{
		Name:   "P http-port",
		Desc:   "Specify listen HTTP port.",
		EnvVar: "APP_HTTP_PORT",
		Value:  "8282",
	})
)

var app = cli.App("mesh_api", "An example API Gateway for meshRPC cluster.")

func main() {
	app.Action = func() {
		// Init a cluster client
		c := cluster.NewAstraCluster("mesh_api", &cluster.AstraOptions{
			Tags: []string{
				*clusterName,
			},
			Nodes: *clusterNodes,
			// Debug: true,
		})
		// Listen on a TCP address, this address can be used
		// by other peers to discover each other in this cluster.
		if err := c.ListenAndServe(*netAddr); err != nil {
			closer.Fatalln(err)
		}

		// Start a normal Gin HTTP server that will use cluster endpoints.
		httpListenAndServe(c)
	}
	app.Run(os.Args)
}

func httpListenAndServe(c cluster.Cluster) {
	router := gin.Default()

	serviceSpecs := make(map[string]cluster.HandlerSpec, len(*allPeers))
	for _, serviceName := range *allPeers {
		serviceSpecs[serviceName] = peer.RPCHandlerSpec
	}
	wait(c, serviceSpecs)

	log.Println("wait for peers is over", *allPeers, "initializing client to", (*allPeers)[0])
	peerClient := c.NewClient((*allPeers)[0], peer.RPCHandlerSpec)
	var peerSvc peer.Service = peer.NewServiceClient(peerClient, nil)

	initState := func(hop *peer.HopState) {
		hop.Start = time.Now().UnixNano()
		hop.Current = 0
		hop.Last = ""
		if len(hop.Queue) == 0 {
			hop.Queue = *allPeers
		}
		hop.QueueIdx = 0
		if hop.Limit == 0 {
			hop.Limit = 10
		}
		hop.HopCount = 0
		hop.HopLatency = 0
	}

	allStaticPeers := staticPeers(*allPeers)
	var staticPeerSvc peer.Service = allStaticPeers[(*allPeers)[0]]

	router.POST("/peer/hop", func(c *gin.Context) {
		var hop *peer.HopState
		if err := c.BindJSON(&hop); err != nil {
			c.JSON(500, err.Error())
			return
		}
		initState(hop)
		lastState, err := peerSvc.Hop(hop)
		if err != nil {
			c.JSON(500, &HopResponse{
				State: lastState,
				Error: err,
			})
			return
		}
		c.JSON(200, &HopResponse{
			State: lastState,
		})
		return
	})

	router.POST("/staticPeer/hop", func(c *gin.Context) {
		var hop *peer.HopState
		if err := c.BindJSON(&hop); err != nil {
			c.JSON(500, err.Error())
			return
		}
		initState(hop)
		lastState, err := staticPeerSvc.Hop(hop)
		if err != nil {
			c.JSON(500, &HopResponse{
				State: lastState,
				Error: err,
			})
			return
		}
		c.JSON(200, &HopResponse{
			State: lastState,
		})
		return
	})

	listenAddr := *httpListenHost + ":" + *httpListenPort
	if err := router.Run(listenAddr); err != nil {
		log.Fatalln(err)
	}
}

func staticPeers(names []string) map[string]peer.Service {
	allPeers := make(map[string]peer.Service, len(names))
	for _, name := range names {
		allPeers[name] = peer.NewService(name, allPeers)
	}
	return allPeers
}

type HopResponse struct {
	State *peer.HopState
	Error error
}

func wait(c cluster.Cluster, specs map[string]cluster.HandlerSpec) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	if err := c.Wait(ctx, specs); err != nil {
		log.Println("bench_api: service await failure:", err)
	}
}
