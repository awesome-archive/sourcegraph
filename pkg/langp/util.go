package langp

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"sourcegraph.com/sourcegraph/sourcegraph/pkg/conf/feature"
)

var btrfsPresent bool

func init() {
	if !feature.Features.Universe {
		return
	}
	_, err := exec.LookPath("btrfs")
	if err == nil {
		btrfsPresent = true
	} else {
		log.Println("btrfs command not available, assuming filesystem is not btrfs")
	}
}

func btrfsSubvolumeCreate(path string) error {
	if !btrfsPresent {
		return os.Mkdir(path, 0700)
	}
	return Cmd("btrfs", "subvolume", "create", path).Run()
}

func btrfsSubvolumeSnapshot(subvolumePath, snapshotPath string) error {
	if !btrfsPresent {
		// TODO: This isn't portable outside *nix, but it does spare us a lot
		// of complex logic. Maybe find a good package to copy a directory.
		return Cmd("cp", "-r", subvolumePath, snapshotPath).Run()
	}
	return Cmd("btrfs", "subvolume", "snapshot", subvolumePath, snapshotPath).Run()
}

// dirExists tells if the directory p exists or not.
func dirExists(p string) (bool, error) {
	info, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// ResolveRepoAlias returns import path and clone URI of given repository URI,
// it takes special care to sourcegraph main repository.
func ResolveRepoAlias(repo string) (importPath, cloneURI string) {
	// TODO(slimsag): find a way to pass this information from the app instead
	// of hard-coding it here.
	if repo == "sourcegraph/sourcegraph" {
		return "sourcegraph.com/sourcegraph/sourcegraph", "git@github.com:sourcegraph/sourcegraph"
	}
	return repo, "https://" + repo
}

// UnresolveRepoAlias performs the opposite action of ResolveRepoAlias.
func UnresolveRepoAlias(repo string) string {
	if repo == "sourcegraph.com/sourcegraph/sourcegraph" {
		repo = "sourcegraph/sourcegraph"
	}
	return repo
}

// Cmd is a small helper which logs the command name and parameters and returns
// a command with output going to stdout/stderr.
func Cmd(name string, args ...string) *exec.Cmd {
	s := fmt.Sprintf("exec %s", name)
	for _, arg := range args {
		s = fmt.Sprintf("%s %q", s, arg)
	}
	log.Println(s)
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Clone clones the specified repository at the given commit into the specified
// directory. If update is true, this function assumes the git repository
// already exists and can just be fetched / updated.
func Clone(update bool, cloneURI, repoDir, commit string) error {
	if !update {
		c := Cmd("git", "clone", cloneURI, repoDir)
		if err := c.Run(); err != nil {
			return err
		}
	} else {
		// Update our repo to match the remote.
		c := Cmd("git", "remote", "update", "--prune")
		c.Dir = repoDir
		if err := c.Run(); err != nil {
			return err
		}
	}

	// Reset to the specific revision.
	c := Cmd("git", "reset", "--hard", commit)
	c.Dir = repoDir
	return c.Run()
}

// GetReferenceCommitID tries to fetch and return commit ID given reference
// points to from clone URI.
func GetReferenceCommitID(cloneURI, ref string) (string, error) {
	stdout := bytes.NewBuffer(nil)
	c := Cmd("git", "ls-remote", cloneURI, ref)
	c.Stdout = stdout
	if err := c.Run(); err != nil {
		return "", err
	}
	output := stdout.String()
	if len(output) < 40 {
		return "", fmt.Errorf("invalid length of output: %d", len(output))
	}
	return output[:40], nil
}
