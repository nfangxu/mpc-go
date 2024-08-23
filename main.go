package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nfangxu/mpc-go/internal/conf"
	"github.com/nfangxu/mpc-go/internal/datax"
	"github.com/nfangxu/mpc-go/internal/psi/ecdh"
	"github.com/nfangxu/mpc-go/internal/server"
	"github.com/nfangxu/mpc-go/internal/utils"
	"github.com/samber/lo"
	rpcxclient "github.com/smallnest/rpcx/client"
	rpcxserver "github.com/smallnest/rpcx/server"
	"log"
	"math/big"
)

var configPath string
var selfPartyId string

// rpcx Server
var root = rpcxserver.NewServer()
var psiData = &datax.Data{}

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Config path")
	flag.StringVar(&selfPartyId, "party", "", "self party id")
	flag.Parse()

	_ = root.RegisterName("server", new(server.Server), "")
	_ = root.RegisterName("psiData", psiData, "")
}

func main() {
	config := &conf.Conf{}
	if err := config.Load(configPath); err != nil {
		log.Fatal("Cannot parse config:", configPath)
		return
	}
	if len(selfPartyId) > 0 {
		config.Self = selfPartyId
	}

	fmt.Printf("Nodes: %+v\n", config.Nodes)
	fmt.Printf("Parties: %+v\n", config.Parties)

	selfNode, ok := lo.Find(config.Nodes, func(item conf.Node) bool {
		return item.PartyId == selfPartyId
	})
	if !ok {
		log.Fatal("Cannot find self party in config.Nodes")
		return
	}

	go func() {
		// defer psiData.CleanAll(nil, nil, nil)
		remoteNode, ok := lo.Find(config.Nodes, func(item conf.Node) bool {
			return item.PartyId != selfPartyId
		})
		if !ok {
			fmt.Println("Cannot find remote party in config.Nodes")
			return
		}
		d, err := rpcxclient.NewPeer2PeerDiscovery("tcp@"+remoteNode.Address, "")
		if err != nil {
			fmt.Println("P2P error:", err)
			return
		}

		if err := server.Ping(d); err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		psiClient := rpcxclient.NewXClient("psi", rpcxclient.Failtry, rpcxclient.RandomSelect, d, rpcxclient.DefaultOption)
		defer psiClient.Close()
		psiDataClient := rpcxclient.NewXClient("psiData", rpcxclient.Failtry, rpcxclient.RandomSelect, d, rpcxclient.DefaultOption)
		defer psiDataClient.Close()

		lfilename := fmt.Sprintf("runtime/datasets/%s.csv", selfNode.PartyId)
		rfilename := fmt.Sprintf("runtime/datasets/%s.csv", remoteNode.PartyId)

		origin, err := utils.Readfile(lfilename)
		xs, ys, err := ecdh.GetPoints(origin)
		if err != nil {
			fmt.Println("Ecdh get points failed:", err)
			return
		}
		lxs, lys := ecdh.Exp(xs, ys, ecdh.Key())

		// 第一次数据交换 Start
		lKey1 := fmt.Sprintf("%s:hash(%s)", remoteNode.PartyId, lfilename)
		_ = psiData.Push(lKey1, utils.Json([][]*big.Int{lxs, lys}))
		rKey1 := fmt.Sprintf("%s:hash(%s)", selfPartyId, rfilename)
		fmt.Println("Exchange #1:", rKey1)
		rdata, err := psiData.Pull(psiDataClient, rKey1)
		if err != nil {
			fmt.Println("PsiData.Pull error:", err)
			return
		}
		_rdata := make([][]*big.Int, 0)
		if err := json.Unmarshal(rdata, &_rdata); err != nil {
			fmt.Println("Unmarshal rdata failed:", err)
			return
		}
		rxs, rys := _rdata[0], _rdata[1]
		// 第一次数据交换 End

		_rxs, _rys := ecdh.Exp(rxs, rys, ecdh.Key())

		// 第二次数据交换 Start
		lKey2 := fmt.Sprintf("%s:hash(hash(%s))", remoteNode.PartyId, rfilename)
		_ = psiData.Push(lKey2, utils.Json([][]*big.Int{_rxs, _rys}))
		rKey2 := fmt.Sprintf("%s:hash(hash(%s))", selfPartyId, lfilename)
		fmt.Println("Exchange #2:", rKey2)
		ldata, err := psiData.Pull(psiDataClient, rKey2)
		if err != nil {
			fmt.Println("PsiData.Pull error:", err)
			return
		}
		_ldata := make([][]*big.Int, 0)
		if err := json.Unmarshal(ldata, &_ldata); err != nil {
			fmt.Println("Unmarshal rd1 failed:", err)
			return
		}
		_lxs := _ldata[0]
		// 第二次数据交换 End

		fmt.Println("Intersection")
		idx := ecdh.Intersection(_lxs, _rxs)
		for _, id := range idx {
			fmt.Printf("IDX: %d, Value: %s\n", id, origin[id])
		}
	}()

	fmt.Println("Start Serve")
	if err := root.Serve("tcp", selfNode.ListenAddr); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}
