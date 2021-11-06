package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"gopkg.in/check.v1"

	"github.com/containers/image/v5/manifest"
)

// This image is known to be x86_64 only right now
const knownNotManifestListedImage_x8664 = "docker://quay.io/coreos/11bot"

const expectedProxySemverMajor = "0.2"

// request is copied from proxy.go
// We intentionally copy to ensure that we catch any unexpected "API" changes
// in the JSON.
type request struct {
	// Method is the name of the function
	Method string `json:"method"`
	// Args is the arguments (parsed inside the fuction)
	Args []interface{} `json:"args"`
}

// reply is copied from proxy.go
type reply struct {
	// Success is true if and only if the call succeeded.
	Success bool `json:"success"`
	// Value is an arbitrary value (or values, as array/map) returned from the call.
	Value interface{} `json:"value"`
	// PipeID is an index into open pipes, and should be passed to FinishPipe
	PipeID uint32 `json:"pipeid"`
	// Error should be non-empty if Success == false
	Error string `json:"error"`
}

// maxMsgSize is also copied from proxy.go
const maxMsgSize = 32 * 1024

type proxy struct {
	c *net.UnixConn
}

type pipefd struct {
	// id is the remote identifier "pipeid"
	id uint
	fd *os.File
}

func (self *proxy) call(method string, args []interface{}) (rval interface{}, fd *pipefd, err error) {
	req := request{
		Method: method,
		Args:   args,
	}
	reqbuf, err := json.Marshal(&req)
	if err != nil {
		return
	}
	n, err := self.c.Write(reqbuf)
	if err != nil {
		return
	}
	if n != len(reqbuf) {
		err = fmt.Errorf("short write during call of %d bytes", n)
		return
	}
	oob := make([]byte, syscall.CmsgSpace(1))
	replybuf := make([]byte, maxMsgSize)
	n, oobn, _, _, err := self.c.ReadMsgUnix(replybuf, oob)
	if err != nil {
		err = fmt.Errorf("reading reply: %v", err)
		return
	}
	var reply reply
	err = json.Unmarshal(replybuf[0:n], &reply)
	if err != nil {
		err = fmt.Errorf("Failed to parse reply: %w", err)
		return
	}
	if !reply.Success {
		err = fmt.Errorf("remote error: %s", reply.Error)
		return
	}

	if reply.PipeID > 0 {
		var scms []syscall.SocketControlMessage
		scms, err = syscall.ParseSocketControlMessage(oob[:oobn])
		if err != nil {
			err = fmt.Errorf("failed to parse control message: %v", err)
			return
		}
		if len(scms) != 1 {
			err = fmt.Errorf("Expected 1 received fd, found %d", len(scms))
			return
		}
		var fds []int
		fds, err = syscall.ParseUnixRights(&scms[0])
		if err != nil {
			err = fmt.Errorf("failed to parse unix rights: %v", err)
			return
		}
		fd = &pipefd{
			fd: os.NewFile(uintptr(fds[0]), "replyfd"),
			id: uint(reply.PipeID),
		}
	}

	rval = reply.Value
	return
}

func newProxy() (*proxy, error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_SEQPACKET, 0)
	if err != nil {
		return nil, err
	}
	myfd := os.NewFile(uintptr(fds[0]), "myfd")
	defer myfd.Close()
	theirfd := os.NewFile(uintptr(fds[1]), "theirfd")
	defer theirfd.Close()

	mysock, err := net.FileConn(myfd)
	if err != nil {
		return nil, err
	}

	// Note ExtraFiles starts at 3
	proc := exec.Command("skopeo", "experimental-image-proxy", "--sockfd", "3")
	proc.Stderr = os.Stderr
	proc.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	proc.ExtraFiles = append(proc.ExtraFiles, theirfd)

	if err = proc.Start(); err != nil {
		return nil, err
	}

	return &proxy{
		c: mysock.(*net.UnixConn),
	}, nil
}

func init() {
	check.Suite(&ProxySuite{})
}

type ProxySuite struct {
}

func (s *ProxySuite) SetUpSuite(c *check.C) {
}

func (s *ProxySuite) TearDownSuite(c *check.C) {
}

func initOci(p string) error {
	err := ioutil.WriteFile(filepath.Join(p, "oci-layout"), []byte("{\"imageLayoutVersion\":\"1.0.0\"}"), 0644)
	if err != nil {
		return err
	}

	blobdir := filepath.Join(p, "blobs/sha256")
	err = os.MkdirAll(blobdir, 0755)
	if err != nil {
		return err
	}
	return nil
}

type byteFetch struct {
	content []byte
	err     error
}

func (s *ProxySuite) TestProxy(c *check.C) {
	p, err := newProxy()
	c.Assert(err, check.IsNil)

	v, fd, err := p.call("Initialize", nil)
	c.Assert(err, check.IsNil)
	semver, ok := v.(string)
	if !ok {
		c.Fatalf("Unexpected value %T", v)
	}
	if !strings.HasPrefix(semver, expectedProxySemverMajor) {
		c.Fatalf("Unexpected semver %s", semver)
	}
	c.Assert(fd, check.IsNil)

	v, fd, err = p.call("OpenImage", []interface{}{knownNotManifestListedImage_x8664})
	c.Assert(err, check.IsNil)
	c.Assert(fd, check.IsNil)

	imgidv, ok := v.(float64)
	c.Assert(ok, check.Equals, true)
	imgid := uint32(imgidv)

	v, fd, err = p.call("GetManifest", []interface{}{imgid})
	c.Assert(err, check.IsNil)
	c.Assert(fd, check.NotNil)
	fetchchan := make(chan byteFetch)
	go func() {
		manifestBytes, err := ioutil.ReadAll(fd.fd)
		fetchchan <- byteFetch{
			content: manifestBytes,
			err:     err,
		}
	}()
	_, _, err = p.call("FinishPipe", []interface{}{fd.id})
	c.Assert(err, check.IsNil)
	fetchRes := <-fetchchan
	c.Assert(fetchRes.err, check.IsNil)

	_, err = manifest.OCI1FromManifest(fetchRes.content)
	c.Assert(err, check.IsNil)

	td, err := ioutil.TempDir("", "skopeo-proxy")
	defer os.RemoveAll(td)

	c.Assert(initOci(td), check.IsNil)
}
