package config

// ServerConfig test server 配置
type ServerConfig struct {
	ServerPort string `toml:"server_port"`
}

// SingleConfig 接收的单个配置
type SingleConfig struct {
	ChainConfig          ChainConfig
	ContractConfig       ContractConfig
	RelayerConfig        RelayerConfig
	VerificationServices []VerificationServiceConfig
}

// ChainConfig 链配置
type ChainConfig struct {
	Type    string   `toml:"type"`
	ChainID uint64   `toml:"chain_id"`
	NodeIP  []string `toml:"node_ip"`
	RPCPort uint16   `toml:"rpc_port"`
}

// Protocol 协议配置
type Protocol struct {
	IsEmpty    bool   `toml:"is_empty"`
	ProtocolID uint64 `toml:"protocol_id"`
}

// ContractConfig 合约配置
type ContractConfig struct {
	Application  Protocol `toml:"application"`
	Transaction  Protocol `toml:"transaction"`
	Verification Protocol `toml:"verification"`
	Transport    Protocol `toml:"transport"`
}

// RelayerMetaConfig relayer 基础配置
type RelayerMetaConfig struct {
	ChainID          uint64 `toml:"chain_id"`
	IP               string `toml:"ip"`
	VerificationPort uint16 `toml:"verification_port"`
	TransportPort    uint16 `toml:"transport_port"`
}

// LocalChain 本地链配置
type LocalChain struct {
	IP      string `toml:"ip"`
	RPCPort uint16 `toml:"rpc_port"`
}

// PeerRelayer 关联的 relayer 配置
type PeerRelayer struct {
	ChainID uint64 `toml:"chain_id"`
	IP      string `toml:"ip"`
	Port    uint16 `toml:"port"`
}

// VerificationPlugin 验证插件配置
type VerificationPlugin struct {
	VerificationID uint64 `toml:"verification_id"`
	IP             string `toml:"ip"`
	Port           uint16 `toml:"port"`
}

// RelayerConfig relayer 完整配置
type RelayerConfig struct {
	Self                RelayerMetaConfig    `toml:"self"`
	LocalChain          LocalChain           `toml:"local_chain"`
	Peers               []PeerRelayer        `toml:"peers"`
	VerificationPlugins []VerificationPlugin `toml:"verification_plugins"`
}

// VerificationServiceConfig 验证服务配置
type VerificationServiceConfig struct {
	VID  uint64 `toml:"v_id"`
	IP   string `toml:"ip"`
	Port uint16 `toml:"port"`
}
