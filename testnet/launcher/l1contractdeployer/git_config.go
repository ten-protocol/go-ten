package l1contractdeployer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/ten-protocol/go-ten/go/node"
	"gopkg.in/yaml.v3"
)

const (
	_githubRepoUrl = "github.com/ten-protocol/ten-apps"
	// path to the file in the repo (env has to be substituted with the testnet env)
	_githubFilePath = "nonprod-argocd-config/apps/envs/%s/valuesFile/l1-values.yaml"
)

// GitL1ConfigWrapper matches a yaml file structure in the ten-apps config repo which looks like this:
// l1Config:
//	networkConfig:  "0x5E13D7dC06C534a3dED5Fd74cd8020067F203c41"
//	messagebus:  "0x7548DB7ceD159D0FA59b9897b4F043E68a0FBf19"
//	bridge:  "0xA87CAB536dB553F7D06adf2886304390108e70bb"
//	crosschain:  "0x1ADC60E42ca33b41eb91a95DA25aD83e59886966"
//	rollup:  "0x3ECA99c12f8252378E1DE0fA717b527FE74DEf12"
//	enclaveRegistry:  "0x9a73334FFb1d3942817DEaC7f96bEfd9230E77aA"
//	starthash:  "0xfeddfd4719e39270913071f12127096397f8cba57eb800881fcfb4c367510551"

type GitL1ConfigWrapper struct {
	L1Config GitL1Config `yaml:"l1Config"`
}

type GitL1Config struct {
	NetworkConfig   string `yaml:"networkConfig"`
	MessageBus      string `yaml:"messagebus"`
	Bridge          string `yaml:"bridge"`
	CrossChain      string `yaml:"crosschain"`
	Rollup          string `yaml:"rollup"`
	EnclaveRegistry string `yaml:"enclaveRegistry"`
	Starthash       string `yaml:"starthash"`
}

func StoreNetworkCfgInGithub(githubPAT string, networkName string, networkConfig *node.NetworkConfig) error {
	repoUrl := fmt.Sprintf("https://%s@%s", githubPAT, _githubRepoUrl)
	// create a temporary directory to clone the repo
	tmpDir, err := os.MkdirTemp("", "ten-apps")
	if err != nil {
		return fmt.Errorf("failed to create temp repo dir: %w", err)
	}
	// checkout the repo
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: repoUrl,
		Auth: &http.BasicAuth{
			Username: "github-user", // can be anything non-empty
			Password: githubPAT,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	relFilePath := fmt.Sprintf(_githubFilePath, networkName)
	absFilePath := filepath.Join(tmpDir, relFilePath)

	f, err := os.ReadFile(absFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", absFilePath, err)
	}
	var data GitL1ConfigWrapper
	err = yaml.Unmarshal(f, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file %s: %w", absFilePath, err)
	}

	data.L1Config.NetworkConfig = networkConfig.NetworkConfigAddress
	data.L1Config.MessageBus = networkConfig.MessageBusAddress
	data.L1Config.Bridge = networkConfig.BridgeAddress
	data.L1Config.CrossChain = networkConfig.CrossChainAddress
	data.L1Config.Rollup = networkConfig.DataAvailabilityRegistryAddress
	data.L1Config.EnclaveRegistry = networkConfig.EnclaveRegistryAddress
	data.L1Config.Starthash = networkConfig.L1StartHash

	// marshal the struct back to yaml
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal file %s: %w", absFilePath, err)
	}

	// write the changes to the file
	err = os.WriteFile(absFilePath, yamlData, 0o644) //nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to write updated config file %s: %w", absFilePath, err)
	}

	// stage and commit the changes
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = wt.Add(relFilePath)
	if err != nil {
		return fmt.Errorf("failed to add file to worktree: %w", err)
	}

	_, err = wt.Commit(fmt.Sprintf("[GH actions] prepare '%s' network - update l1 config", networkName), &git.CommitOptions{
		Author: &object.Signature{
			Name: "github-actions",
			When: time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	// push the changes to the repo
	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "github-actions",
			Password: githubPAT,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	// clean up the temp directory
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to remove temp dir: %w", err)
	}

	return nil
}
