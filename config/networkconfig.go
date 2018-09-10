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


package config

import (
	"fmt"
	"time"
	"github.com/hpb-project/go-hpb/network/p2p/nat"
	"github.com/hpb-project/go-hpb/network/p2p/netutil"
	"github.com/hpb-project/go-hpb/network/p2p/discover"

	"strings"
	"path/filepath"
	"os"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8545        // Default TCP port for the HTTP RPC server
	DefaultWSHost   = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort   = 8546        // Default TCP port for the websocket RPC server
)

const (
	clientIdentifier = "ghpb" // Client identifier to advertise over the network
)
const (
	// BloomBitsBlocks is the number of blocks a single bloom bit section vector
	// contains.
	BloomBitsBlocks uint64 = 4096
)

const (
	// This is the amount of time spent waiting in between
	// redialing a certain node.
	dialHistoryExpiration = 30 * time.Second

	// Discovery lookups are throttled and can only run
	// once every few seconds.
	lookupInterval = 4 * time.Second

	// If no peers are found for this amount of time, the initial bootnodes are
	// attempted to be connected.
	fallbackInterval = 20 * time.Second

	// Endpoint resolution is throttled with bounded backoff.
	initialResolveDelay = 60 * time.Second
	maxResolveDelay     = time.Hour
)

const (
	defaultDialTimeout      = 15 * time.Second
	refreshPeersInterval    = 30 * time.Second
	staticPeerCheckInterval = 15 * time.Second


	// Maximum number of concurrently dialing outbound connections.
	maxActiveDialTasks = 16

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second
)
const (
	maxUint24 = ^uint32(0) >> 8

	sskLen = 16 // ecies.MaxSharedKeyLength(pubKey) / 2
	sigLen = 65 // elliptic S256
	pubLen = 64 // 512 bit pubkey in uncompressed representation without format byte
	shaLen = 32 // hash length (for nonce etc)

	authMsgLen  = sigLen + shaLen + pubLen + shaLen + 1
	authRespLen = pubLen + shaLen + 1

	eciesOverhead = 65 /* pubkey */ + 16 /* IV */ + 32 /* MAC */

	encAuthMsgLen  = authMsgLen + eciesOverhead  // size of encrypted pre-EIP-8 initiator handshake
	encAuthRespLen = authRespLen + eciesOverhead // size of encrypted pre-EIP-8 handshake reply

	// total timeout for encryption handshake and protocol
	// handshake in both directions.
	handshakeTimeout = 5 * time.Second

	// This is the timeout for sending the disconnect reason.
	// This is shorter than the usual timeout because we don't want
	// to wait if the connection is known to be bad anyway.
	discWriteTimeout = 1 * time.Second
)

//Protocal
const (
	baseProtocolVersion    = 1
	baseProtocolLength     = uint64(16)
	baseProtocolMaxMsgSize = 2 * 1024

	pingInterval = 15 * time.Second
)

const (
	// devp2p message codes
	handshakeMsg = 0x00
	discMsg      = 0x01
	pingMsg      = 0x02
	pongMsg      = 0x03
	getPeersMsg  = 0x04
	peersMsg     = 0x05
)


// UDP  超时相关
const (
	respTimeout = 500 * time.Millisecond
	sendTimeout = 500 * time.Millisecond
	expiration  = 20 * time.Second

	ntpFailureThreshold = 32               // Continuous timeouts after which to check NTP
	ntpWarningCooldown  = 10 * time.Minute // Minimum amount of time to pass before repeating NTP warning
	driftThreshold      = 10 * time.Second // Allowed clock drift before warning user
)



const NodeIDBits   = 512
const RandNoceSize = 32

type NodeType  uint8
// 节点类型
const(
	LightNode  NodeType = 0x10  //默认节点类型，没有通过硬件认证的节点类型都是默认类型 UnknownNode

	AuthNode   NodeType = 0x30  //经过认证的节点
	PreNode    NodeType = 0x31  //候选节点
	HpNode     NodeType = 0x60  //高性能节点

	BootNode   NodeType = 0xA0  //启动节点
)

//节点数据库相关

var (
	nodeDBNilNodeID      = discover.NodeID{}       // Special node ID to use as a nil element.
	nodeDBNodeExpiration = 24 * time.Hour // Time after which an unseen node should be dropped.
	nodeDBCleanupCycle   = time.Hour      // Time period for running the expiration task.

	nodeDBNodeExpirationOneHour = time.Hour // Time after which an unseen node should be dropped.
)

var defaultNetworkConfig = NetworkConfig{

	HTTPPort:     DefaultHTTPPort,
	HTTPModules:  []string{"net", "web3", "prometheus"},
	WSPort:       DefaultWSPort,
	WSModules:    []string{"net", "web3", "prometheus"},
	ListenAddr:   ":30303",
	MaxPeers:     50,
	NAT:          nat.Any(),
	IpcEndpoint:  DefaultIPCEndpoint(clientIdentifier),
	HttpEndpoint: DefaultHTTPEndpoint(),
	WsEndpoint:   DefaultWSEndpoint(),
}

var MainnetBootnodes = []string{
	// Hpb Foundation Go Bootnodes
	"hnode://34b0c7792d838a7d0859c0e3701f41678731605625045f65b1cd9e764686be0836cb0b0701053a2517c8804f38c25068deaf308686cba7248c46429053b62752@127.0.0.1:3001",
	//"hnode://7d5fdaee2e78dd5085ffbf7c6d96aff10bfbf40eb464f10a10363e9059b15a90b01d99cb43ba16642e23b7aa77739443f1b573d9e9d24e2a40bfa42bfc19e9f3@127.0.0.1:30011",
	//"hnode://43b1ad158f0333fb19e8ccdfa99aa5fe90ab7f77c667a19f561f99024f52f3254452d563b9cc87399076fc9194bb0cf86de10d145897e51c62c396ec8579af2f@47.105.118.89:30001",
	//"hnode://7d5fdaee2e78dd5085ffbf7c6d96aff10bfbf40eb464f10a10363e9059b15a90b01d99cb43ba16642e23b7aa77739443f1b573d9e9d24e2a40bfa42bfc19e9f3@101.132.180.98:30011",
	//"hnode://7d5fdaee2e78dd5085ffbf7c6d96aff10bfbf40eb464f10a103639059b15a90b01d99cb43ba16642e23b7aa77739443f1b573d9e9d24e2a40bfa42bfc19e9f3@47.105.118.89:30101",
	//"hnode://7d5fdaee2e78dd5085ffbf7c6d96aff10bfbf40eb464f10a10363e9059b15a90b01d99cb43ba16642e23b7aa77739443f1b573d9e9d24e2a40bfa42bfc19e9f3@127.0.0.1:30101",
	//"hnode://ce0220052f26d76a1083f56522d8bfcf12b00c7e92b0cab5febe4fb467dbfed5cc616998878dcf23453f26a2b245d1323859d56f80545580425ae08b4365e237@47.100.227.96:20001",
	//"hnode://fb9ba578cb6a47d0794e983be7438f075db136c3a9e170681bf6d99821493c3f3e48d4ecd925f182c71cbf42b2a17517c883db325d469566978b03cea7a2a601@47.99.36.31:20001",
	//"hnode://b0506da3ee83f33e2c076e9623fd22d371aea68ea1dec158981534bc51c6e38e837c779d4f1cbe002ec1751c1a31a2ee260b1336b0141f06369bcbfce4e6aa0c@47.99.46.246:20001",
	//"hnode://af6568c2913a99401fa567182a39f89bad7a0a273d2d7ba5a4ec1d02ad9c790c3be3f17ac92da84c5a9ed604cb7d44482783c85792d587f2bfc42b1dccd3d7e5&1@47.92.26.84:30301",
}

// TestnetBootnodes are the hnode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
}
type NetworkConfig struct {
	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string `toml:",omitempty"`

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort int `toml:",omitempty"`

	// HTTPCors is the Cross-Origin Resource Sharing header to send to requesting
	// clients. Please be aware that CORS is a browser enforced security, it's fully
	// useless for custom HTTP clients.
	HTTPCors []string `toml:",omitempty"`

	// HTTPModules is a list of API modules to expose via the HTTP RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	HTTPModules []string `toml:",omitempty"`

	// WSHost is the host interface on which to start the websocket RPC server. If
	// this field is empty, no websocket API endpoint will be started.
	WSHost string `toml:",omitempty"`

	// WSPort is the TCP port number on which to start the websocket RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	WSPort int `toml:",omitempty"`

	// WSOrigins is the list of domain to accept websocket requests from. Please be
	// aware that the server can only act upon the HTTP request the client sends and
	// cannot verify the validity of the request header.
	WSOrigins []string `toml:",omitempty"`

	// WSModules is a list of API modules to expose via the websocket RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	WSModules []string `toml:",omitempty"`

	// WSExposeAll exposes all API modules via the WebSocket RPC interface rather
	// than just the public ones.
	//
	// *WARNING* Only set this if the node is running in a trusted network, exposing
	// private APIs to untrusted users is a major security risk.
	WSExposeAll bool `toml:",omitempty"`


	// MaxPeers is the maximum number of peers that can be
	// connected. It must be greater than zero.
	MaxPeers int

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DiscoveryV5 specifies whether the the new topic-discovery based V5 discovery
	// protocol should be started or not.
	//DiscoveryV5 bool `toml:",omitempty"`
    NoDiscovery bool
	// Listener address for the V5 discovery protocol UDP traffic.
	//DiscoveryV5Addr string `toml:",omitempty"`

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// RoleType sets the node type of this server.
	// One of hpnode,prenode,access,light.
	RoleType string

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	//Protocols []p2p.Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// If set to a non-nil value, the given NAT port mapper
	// is used to make the listening port available to the
	// Internet.
	NAT nat.Interface `toml:",omitempty"`

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	//Dialer p2p.NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	IpcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	HttpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	WsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)

	BootstrapNodes []*discover.Node

}



func DefaultNetworkConfig() NetworkConfig{
	cfg:= defaultNetworkConfig

	cfg.HTTPModules = append(cfg.HTTPModules, "hpb")
	cfg.WSModules = append(cfg.WSModules, "hpb")


	return cfg
}


// HTTPEndpoint resolves an HTTP endpoint based on the configured host interface
// and port parameters.
func (c *NetworkConfig) HTTPEndpoint() string {
	if c.HTTPHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.HTTPHost, c.HTTPPort)
}

// DefaultHTTPEndpoint returns the HTTP endpoint used by default.
func DefaultHTTPEndpoint() string {
	config := &NetworkConfig{HTTPHost: DefaultHTTPHost, HTTPPort: DefaultHTTPPort}
	return config.HTTPEndpoint()
}

// DefaultIPCEndpoint returns the IPC path used by default.
func DefaultIPCEndpoint(clientIdentifier string) string {
	if clientIdentifier == "" {
		clientIdentifier = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if clientIdentifier == "" {
			panic("empty executable name")
		}
	}
	config := &Nodeconfig{DataDir: DefaultDataDir(), IPCPath: clientIdentifier + ".ipc"}
	return config.IPCEndpoint()
}
// WSEndpoint resolves an websocket endpoint based on the configured host interface
// and port parameters.
func (c *NetworkConfig) WSEndpoint() string {
	if c.WSHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.WSHost, c.WSPort)
}

// DefaultWSEndpoint returns the websocket endpoint used by default.
func DefaultWSEndpoint() string {
	config := &NetworkConfig{WSHost: DefaultWSHost, WSPort: DefaultWSPort}
	return config.WSEndpoint()
}









