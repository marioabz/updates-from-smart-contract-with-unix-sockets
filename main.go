package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	LP "./contracts"
)

type Reserves struct {
	Reserve0 *big.Int;
	Reserve1 *big.Int;
}
type Event struct {
	Address string;
	Reserve0 *big.Int;
	Reserve1 *big.Int;
}
var client = &ethclient.Client{}
var reserves map[string]Reserves
var socketAddress = "/tmp/reserves.sock"
var contracts map[string]*LP.LiquidityPool
var contractAddresses = []string{
	"0xCD6bcca48069f8588780dFA274960F15685aEe0e", 
	"0x05f04f112a286c4c551897fb19ed2300272656c8",
	};
var url = "wss://jolly-brattain:class-awning-speech-rule-mangy-gating@ws-nd-780-530-825.p2pify.com";

func main() {
	contracts = make(map[string]*LP.LiquidityPool)
	reserves = make(map[string]Reserves)
	load()
	initContracts()
	
	socket, err := net.Dial("unix", socketAddress)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for _, _address := range contractAddresses {
			updatingSocket := false
			_reserves, err := contracts[_address].GetReserves(nil)
			if err != nil {
				log.Fatal(err)
			}
			reservesChanged := compareReserves(reserves[_address], _reserves.Reserve0, _reserves.Reserve1)
			updatingSocket = updatingSocket || reservesChanged
			if updatingSocket {
				event := Event{
					Address: _address,
					Reserve0: _reserves.Reserve0,
					Reserve1: _reserves.Reserve1,
				}
				marshaledStruct, err := json.Marshal(event)
				fmt.Println(marshaledStruct)
				if err != nil {
					log.Fatal(err)
				}
				reserves[_address] = Reserves{
					Reserve0: _reserves.Reserve0,
					Reserve1: _reserves.Reserve1,
				}
				socket.Write(marshaledStruct)
			}
		}
		time.Sleep(1 * time.Second)
	}	
}

func load() {
	_client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	client = _client
}

func compareReserves(r Reserves, reserve0, reserve1 *big.Int) bool{
	return r.Reserve0.Cmp(reserve0) != 0 || r.Reserve1.Cmp(reserve1) != 0
}

func initContracts() {
	for _, address := range contractAddresses {
		_address := common.HexToAddress(address)
		instance, err := LP.NewLiquidityPool(_address,  client)
		if err != nil {
			log.Fatal(err)
		}
		contracts[address] = instance
		fetchReserves, err := instance.GetReserves(nil)
		if err != nil {
			log.Fatal(err)
		}
		reserves[address] = Reserves{
			Reserve0: fetchReserves.Reserve0,
			Reserve1: fetchReserves.Reserve1,
		}
	}
}

