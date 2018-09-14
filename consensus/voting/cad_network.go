// Copyright 2018 The go-hpb Authors
// This file is part of the go-hpb.
//
// The go-hpb is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-hpb is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-hpb. If not, see <http://www.gnu.org/licenses/>.

package voting

import (
	//"math"
	//"strconv"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"

	//"github.com/hpb-project/go-hpb/common"
	//"github.com/hpb-project/go-hpb/consensus"
	// "math/big"
	"github.com/hpb-project/go-hpb/consensus/snapshots"
	//"github.com/hpb-project/go-hpb/blockchain/storage"
	"github.com/hpb-project/go-hpb/blockchain/state"
	"github.com/hpb-project/go-hpb/common"
	"github.com/hpb-project/go-hpb/common/log"
	"github.com/hpb-project/go-hpb/consensus"
	"github.com/hpb-project/go-hpb/network/p2p"
	"github.com/hpb-project/go-hpb/network/p2p/discover"
)

// 从网络中获取最优化的
func GetCadNodeFromNetwork(state *state.StateDB) ([]*snapshots.CadWinner, error) {
	//str := strconv.FormatUint(number, 10)
	// 模拟从外部获取
	//type CadWinners []*snapshots.CadWinner

	bigaddr, _ := new(big.Int).SetString("0000000000000000000000000000000000000000", 16)
	address := common.BigToAddress(bigaddr)

	bestCadWinners := []*snapshots.CadWinner{}
	peers := p2p.PeerMgrInst().PeersAll()
	fmt.Println("######### peers length is:", len(peers))
	if len(peers) == 0 {
		return nil, nil
	}
	for _, peer := range peers {

		//fmt.Println("this is TxsRate:", peer.TxsRate())
		//fmt.Println("this is Bandwidth:", peer.Bandwidth())
		//networkBandwidth := float64(peer.Bandwidth()) * float64(0.3)
		//transactionNum := float64(peer.TxsRate()) * float64(0.7)
		//VoteIndex := networkBandwidth + transactionNum

		if peer.RemoteType() != discover.BootNode && peer.RemoteType() != discover.SynNode {
			//networkBandwidth := float64(rand.Intn(1000)) * float64(0.3)
			//transactionNum := float64(rand.Intn(1000)) * float64(0.7)
			if len(peer.Address()) == 0 || peer.Address() == address {
				continue
			}
			transactionNum := peer.TxsRate() * float64(0.6)
			networkBandwidth := peer.Bandwidth() * float64(0.3)
			log.Error("GetCadNodeFromNetwork print peer addr", "addr", peer.Address().Str())
			bigval := new(big.Float).SetInt(state.GetBalance(peer.Address()))

			onether2weis := big.NewInt(10)
			onether2weis.Exp(onether2weis, big.NewInt(18), nil)
			onether2weisf := new(big.Float).SetInt(onether2weis)
			bigval.Quo(bigval, onether2weisf)

			val64, _ := bigval.Float64()
			balanceIndex := val64 * float64(0.1)

			VoteIndex := networkBandwidth + transactionNum + balanceIndex
			if peer.Address() != address {
				bestCadWinners = append(bestCadWinners, &snapshots.CadWinner{peer.GetID(), peer.Address(), uint64(VoteIndex)})
			}
		}

	}

	if len(bestCadWinners) == 0 {
		return nil, nil
	}

	// 先获取长度，然后进行随机获取
	//lnlen := int(math.Log2(float64(len(bestCadWinners))))

	var lastCadWinners []*snapshots.CadWinner

	//fmt.Println("bestCadWinners - 1:", len(bestCadWinners)-1)
	//fmt.Println("hpbAddresses:", len(hpbAddresses))

	//for i := 0 ; i < lnlen; i++{
	for i := 0; i < len(bestCadWinners); i++ {
		if len(bestCadWinners) > 1 {
			//lastCadWinners = append(lastCadWinners, bestCadWinners[rand.Intn(len(bestCadWinners)-1)])
			lastCadWinners = append(lastCadWinners, bestCadWinners[i])
		} else {
			lastCadWinners = append(lastCadWinners, bestCadWinners[0])
		}
	}

	//开始进行排序获取最大值
	winners := []*snapshots.CadWinner{}
	lastCadWinnerToChain := &snapshots.CadWinner{}
	voteIndexTemp := uint64(0)

	for _, lastCadWinner := range lastCadWinners {
		if lastCadWinner.VoteIndex >= voteIndexTemp {
			voteIndexTemp = lastCadWinner.VoteIndex
			lastCadWinnerToChain = lastCadWinner
		}
	}
	//if selectable candidate nodes not enough, rand select one in all candidate nodes,else rand select one in first consensus.HpbNodenumber numbers.
	if len(lastCadWinners) == 1 {
		lastCadWinnerToChain = lastCadWinners[0]
	} else if len(lastCadWinners) <= consensus.HpbNodenumber {
		lastCadWinnerToChain = lastCadWinners[rand.Intn(len(lastCadWinners))]
	} else {
		//order the candidate nodes
		for i := len(lastCadWinners) - 1; i >= 0; i-- {
			for j := len(lastCadWinners) - 1; j >= len(lastCadWinners)-i; j-- {
				if lastCadWinners[j].VoteIndex > lastCadWinners[j-1].VoteIndex {
					lastCadWinners[j], lastCadWinners[j-1] = lastCadWinners[j-1], lastCadWinners[j]
				}
			}
		}
		lastCadWinnerToChain = lastCadWinners[rand.Intn(31)]
	}

	winners = append(winners, lastCadWinnerToChain) //返回最优的

	if len(bestCadWinners) > 1 {
		winners = append(winners, bestCadWinners[rand.Intn(len(bestCadWinners)-1)]) //返回随机
	} else {
		winners = append(winners, bestCadWinners[0]) // 获取第一个
	}
	return winners, nil
}

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

/*
func GetCadNodeMap(db hpbdb.Database,chain consensus.ChainReader, number uint64, hash common.Hash) (map[string]*snapshots.CadWinner, error) {

	cadWinnerms := make(map[string]*snapshots.CadWinner)

	if cadNodeSnapformap, err  := GetCadNodeSnap(db, chain, number, hash); err == nil{
		for _, cws := range cadNodeSnapformap.CadWinners {
		    cadWinnerms[cws.NetworkId] = &snapshots.CadWinner{cws.NetworkId,cws.Address,cws.VoteIndex}
		}
	}

    return cadWinnerms,nil
}
*/
