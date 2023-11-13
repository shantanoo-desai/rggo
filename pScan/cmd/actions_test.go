package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/shantanoo-desai/rggo/pScan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	tf, err := os.CreateTemp("", "pScan")

	if err != nil {
		t.Fatal(err)
	}

	tf.Close()

	if initList {
		hl := &scan.HostsList{}

		for _, h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			var out bytes.Buffer
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("expected no error, go %q\n", err)
			}

			if out.String() != tc.expectedOut {
				t.Errorf("expected output %q\n,got %q\n", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"

	hostsEnd := []string{
		"host1",
		"host3",
	}

	var out bytes.Buffer

	expectedOut := ""

	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}

	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	for _, v := range hostsEnd {
		expectedOut += fmt.Sprintf("%s Host not found\n", v)
		expectedOut += fmt.Sprintln()
	}

	// add hosts
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// list hosts
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// list after delete
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// scan hosts
	if err := scanAction(&out, tf, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	if out.String() != expectedOut {
		t.Errorf("expected output\n %q\n,got %q\n", expectedOut, out.String())
	}

}

func TestScanAction(t *testing.T) {
	hosts := []string{
		"localhost",
		"unknownhostoutthere",
	}

	tf, cleanup := setup(t, hosts, true)
	defer cleanup()

	ports := []int{}

	for i := 0; i < 2; i++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))

		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())

		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if i == 1 {
			ln.Close()
		}

	}
	expectedOut := fmt.Sprintln("localhost")
	expectedOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintln("unknownhostoutthere Host not found")
	expectedOut += fmt.Sprintln()

	var out bytes.Buffer

	if err := scanAction(&out, tf, ports); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	if out.String() != expectedOut {
		t.Errorf("expected output\n %q\n, got %q\n", expectedOut, out.String())
	}
}
