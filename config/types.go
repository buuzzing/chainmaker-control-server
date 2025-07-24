/*
 * @Author: Justin
 * @Date: 2025-07-21 06:10:55
 * @filename:
 * @version:
 * @Description:
 * @LastEditTime: 2025-07-21 06:13:04
 */
package config

// 这是 server的config
type ServerConfig struct {
	ServerPort string `toml:"server_port"`
}

// 这是接收的 config 文件
type SingleConfig struct {
	ChainConfig          ChainConfig
	ContractConfig       ContractConfig
	RelayerConfig        RelayerConfig
	VerificationServices []VerificationServiceConfig
}

type ChainConfig struct {
	Type    string   `toml:"type"`
	ChainID uint64   `toml:"chain_id"`
	NodeIP  []string `toml:"node_ip"`
	RPCPort uint16   `toml:"rpc_port"`
}

type Protocol struct {
	IsEmpty    bool   `toml:"is_empty"`
	ProtocolID uint64 `toml:"protocol_id"`
}

type ContractConfig struct {
	Application  Protocol `toml:"application"`
	Transaction  Protocol `toml:"transaction"`
	Verification Protocol `toml:"verification"`
	Transport    Protocol `toml:"transport"`
}

type RelayerMetaConfig struct {
	ChainID          uint64 `toml:"chain_id"`
	IP               string `toml:"ip"`
	VerificationPort uint16 `toml:"verification_port"`
	TransportPort    uint16 `toml:"transport_port"`
}

type LocalChain struct {
	IP      string `toml:"ip"`
	RPCPort uint16 `toml:"rpc_port"`
}

type PeerRelayer struct {
	ChainID uint64 `toml:"chain_id"`
	IP      string `toml:"ip"`
	Port    uint16 `toml:"port"`
}

type VerificationPlugin struct {
	VerificationID uint64 `toml:"verification_id"`
	IP             string `toml:"ip"`
	Port           uint16 `toml:"port"`
}

type RelayerConfig struct {
	Self                RelayerMetaConfig    `toml:"self"`
	LocalChain          LocalChain           `toml:"local_chain"`
	Peers               []PeerRelayer        `toml:"peers"`
	VerificationPlugins []VerificationPlugin `toml:"verification_plugins"`
}

type VerificationServiceConfig struct {
	VID  uint64 `toml:"v_id"`
	IP   string `toml:"ip"`
	Port uint16 `toml:"port"`
}
