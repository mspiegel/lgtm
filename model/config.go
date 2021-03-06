package model

import (
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/ianschenck/envflag"
)

type Config struct {
	Approvals       int    `json:"approvals"         toml:"approvals"`
	Pattern         string `json:"pattern"           toml:"pattern"`
	Team            string `json:"team"              toml:"team"`
	SelfApprovalOff bool   `json:"self_approval_off" toml:"self_approval_off"`
	DoMerge         bool   `json:"do_merge" toml:"do_merge"`
	DoVersion       bool   `json:"do_version" toml:"do_version"`
	ApprovalAlg     string `json:"approval_algorithm" toml:"approval_algorithm"`
	VersionAlg      string `json:"version_algorithm" toml:"version_algorithm"`
	VersionFormat   string `json:"version_format" toml:"version_format"`
	DoComment       bool   `json:"do_comment" toml:"do_comment"`
	DoDeployment    bool   `json:"do_deploy" toml:"do_deploy"`
	DeploymentMap   DeploymentConfigs
	re              *regexp.Regexp
}

var (
	approvals       = envflag.Int("LGTM_APPROVALS", 2, "")
	pattern         = envflag.String("LGTM_PATTERN", `(?i)^LGTM\s*(\S*)`, "")
	team            = envflag.String("LGTM_TEAM", "MAINTAINERS", "")
	selfApprovalOff = envflag.Bool("LGTM_SELF_APPROVAL_OFF", false, "")
	approvalAlg     = envflag.String("LGTM_APPROVAL_ALGORITHM", "simple", "")
	versionAlg      = envflag.String("LGTM_VERSION_ALGORITHM", "semver", "")
)

// ParseConfig parses a projects .lgtm file
func ParseConfig(configData []byte, deployData []byte) (*Config, error) {
	c, err := ParseConfigStr(string(configData))
	if err != nil {
		return nil, err
	}
	if c.DoDeployment {
		c.DeploymentMap, err = loadDeploymentMap(string(deployData))
	}
	if err != nil {
		return nil, err
	}

	return c, nil
}

// ParseConfigStr parses a projects .lgtm file in string format.
func ParseConfigStr(data string) (*Config, error) {
	c := new(Config)
	_, err := toml.Decode(data, c)
	if err != nil {
		return nil, err
	}
	if c.Approvals == 0 {
		c.Approvals = *approvals
	}
	if len(c.Pattern) == 0 {
		c.Pattern = *pattern
	}
	if len(c.Team) == 0 {
		c.Team = *team
	}
	if len(c.ApprovalAlg) == 0 {
		c.ApprovalAlg = *approvalAlg
	}
	if len(c.VersionAlg) == 0 {
		c.VersionAlg = *versionAlg
	}
	if c.SelfApprovalOff == false {
		c.SelfApprovalOff = *selfApprovalOff
	}
	c.re, err = regexp.Compile(c.Pattern)
	return c, err
}

// IsMatch returns true if the text matches the regular
// epxression pattern.
func (c *Config) IsMatch(text string) bool {
	if c.re == nil {
		// this should never happen
		return false
	}
	return c.re.MatchString(text)
}

func loadDeploymentMap(deployData string) (DeploymentConfigs, error) {
	d := DeploymentConfigs{}
	if len(deployData) == 0 {
		return d, nil
	}
	//todo actually load deployment config info from DEPLOYMENT toml file
	_, err := toml.Decode(deployData, &d)
	return d, err
}
