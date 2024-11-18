package release

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseCmdRunnerForKappCmd(t *testing.T) {
	ReleaseCmdRunner := NewReleaseCmdRunner(os.Stdout, false, "", false, nil)
	cmd := exec.Command("kapp", "deploy", "-f", "-", "-a", "pkg-test", "-y")
	err := ReleaseCmdRunner.RunWithCancel(cmd, nil)
	require.NoError(t, err)
}
