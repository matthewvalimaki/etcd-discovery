package etcdisco

import (
	"strings"
	"testing"

	"github.com/couchbaselabs/go.assert"
)

func TestTransformArg(t *testing.T) {

	transformed, err := transformArg(
		"http://{{ .LOCAL_IP }}:2379",
		"LOCAL_IP",
		"10.1.1.1",
	)
	assert.True(t, err == nil)
	assert.Equals(t, transformed, "http://10.1.1.1:2379")

}

func TestTransformArgs(t *testing.T) {

	fakeLocalIp := "10.1.1.1"
	bindings := map[string]string{
		LOCAL_IP: fakeLocalIp,
	}

	args := []string{
		"-listen-client-urls",
		"http://0.0.0.0:2379",
		"-advertise-client-urls",
		"http://{{ .LOCAL_IP }}:2379",
	}

	tranformedArgs, err := tranformArgs(args, bindings)
	assert.True(t, err == nil)
	assert.Equals(t, len(tranformedArgs), len(args))
	lastArgWithLocalIp := tranformedArgs[3]
	assert.True(t, strings.Contains(lastArgWithLocalIp, fakeLocalIp))

}
