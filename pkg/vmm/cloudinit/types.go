package cloudinit

// https://cloudinit.readthedocs.io/en/latest/reference/modules.html
type CLoudConfig struct {
	Hostname         string              `yaml:"hostname"`
	DisableRoot      bool                `yaml:"disable_root"`
	PreserveHostname bool                `yaml:"preserve_hostname"`
	SystemInfo       SystemInfoConfig    `yaml:"system_info"`
	Users            []UserCoinfig       `yaml:"users"`
	GrowPartition    GrowPartitionConfig `yaml:"growpart"`
}

type SystemInfoConfig struct {
	Distro string `yaml:"distro"`
}

type UserCoinfig struct {
	Name              string   `yaml:"name"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

type GrowPartitionConfig struct {
	Mode    GrowPartitionMode `yaml:"mode"`
	Devices []string          `yaml:"devices"`
}

type GrowPartitionMode string

const (
	GrowPartitionMode_Auto GrowPartitionMode = "auto"
)
