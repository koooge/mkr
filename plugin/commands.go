package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

var CommandPlugin = cli.Command{
	Name:        "plugin",
	Usage:       "Manage mackerel plugin",
	Description: `Manage mackerel plugin`,
	Subcommands: []cli.Command{
		{
			Name:        "install",
			Usage:       "install mackerel plugin",
			Description: `WIP`,
			Action:      doPluginInstall,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "plugin install location",
				},
			},
		},
	},
	Hidden: true,
}

func doPluginInstall(c *cli.Context) error {
	argInstallTarget := c.Args().First()
	if argInstallTarget == "" {
		return fmt.Errorf("Specify install name")
	}
	it, err := parseInstallTarget(argInstallTarget)
	if err != nil {
		return errors.Wrap(err, "failed to install plugin")
	}

	pluginDir, err := setupPluginDir(c.String("prefix"))
	if err != nil {
		return errors.Wrap(err, "failed to install plugin")
	}

	u, err := it.makeDownloadURL()
	if err != nil {
		return errors.Wrap(err, "failed to install plugin while making download url")
	}
	_ = pluginDir
	_ = u

	fmt.Println("do plugin install [wip]")
	return nil
}

func setupPluginDir(prefix string) (string, error) {
	if prefix == "" {
		prefix = "/opt/mackerel-agent/plugins"
	}
	err := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup plugin directory")
	}
	return prefix, nil
}

type installTarget struct {
	owner      string
	repo       string
	pluginName string
	releaseTag string
}

func (it *installTarget) makeDownloadURL() (string, error) {
	if it.owner != "" && it.repo != "" {
		if it.releaseTag == "" {
			return "", fmt.Errorf("not implemented")
		}
		filename := fmt.Sprintf("%s_%s_%s.zip", it.repo, runtime.GOOS, runtime.GOARCH)
		return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
			it.owner, it.repo, it.releaseTag, filename), nil
	}
	return "", fmt.Errorf("not implemented")
}

func parseInstallTarget(target string) (*installTarget, error) {
	it := &installTarget{}

	ownerRepoAndReleaseTag := strings.Split(target, "@")
	var ownerRepo string
	switch len(ownerRepoAndReleaseTag) {
	case 1:
		ownerRepo = ownerRepoAndReleaseTag[0]
	case 2:
		ownerRepo = ownerRepoAndReleaseTag[0]
		it.releaseTag = ownerRepoAndReleaseTag[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	ownerAndRepo := strings.Split(ownerRepo, "/")
	switch len(ownerAndRepo) {
	case 1:
		it.pluginName = ownerAndRepo[0]
	case 2:
		it.owner = ownerAndRepo[0]
		it.repo = ownerAndRepo[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	return it, nil
}
