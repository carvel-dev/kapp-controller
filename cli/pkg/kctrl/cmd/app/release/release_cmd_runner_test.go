package release

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseCmdRunnerForKappCmd(t *testing.T) {
	var buf bytes.Buffer
	ReleaseCmdRunner := NewReleaseCmdRunner(&buf, false, "", false, nil)
	cmd := exec.Command("kapp", "deploy", "-f", "https://github.com/carvel-dev/kapp-controller/releases/latest/download/release.yml", "-a", "pkg-test", "-y")
	err := ReleaseCmdRunner.RunWithCancel(cmd, nil)
	require.NoError(t, err)
	require.Emptyf(t, buf.Len(), "Expected buf length 0 but got %d", buf.Len())
}
