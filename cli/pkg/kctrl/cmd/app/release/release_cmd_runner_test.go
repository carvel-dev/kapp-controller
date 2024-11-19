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
	cmd := exec.Command("kapp", "deploy", "-f", "-", "-a", "pkg-test", "-y")
	err := ReleaseCmdRunner.RunWithCancel(cmd, nil)
	require.NoError(t, err)
	expectedLength := 0
	if actualLength := buf.Len(); actualLength != expectedLength {
		t.Errorf("Got Buffer length = %d, Expected empty", actualLength)
	}
}
