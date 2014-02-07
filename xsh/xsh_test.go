package xsh

import (
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestAlias(t *testing.T) {
	session := NewSession()
	session.Alias("gr", "echo", "hi")
	ret, err := session.Capture("gr", []string{"sky"})
	if err != nil {
		t.Error(err)
	}
	if ret.Trim() != "hi sky" {
		t.Errorf("expect 'hi sky' but got:%s", ret)
	}
}

func TestCapture(t *testing.T) {
	r, err := Capture("echo", []string{"hello"})
	if err != nil {
		t.Error(err)
	}
	_ = r
	if r.Trim() != "hello" {
		t.Errorf("expect hello, but got %s", r.Trim())
	}
}

func TestSession(t *testing.T) {
	session := NewSession("pwd")
	session.Set(Dir("/"))
	err := session.Call()
	if err != nil {
		t.Error(err)
	}
	ret, err := session.Capture()
	if err != nil {
		t.Error(err)
	}
	if ret.Trim() != "/" {
		t.Errorf("expect /, but got %s", ret.Trim())
	}
}

func TestPipe(t *testing.T) {
	s := NewSession()
	s.Call("echo", []string{"hello"})
	err := s.Command("echo", []string{"hi"}).Command("cat", []string{"-n"}).Start()
	if err != nil {
		t.Error(err)
	}
	err = s.Wait()
	if err != nil {
		t.Error(err)
	}
	out, err := s.Command("echo", []string{"-n", "hello"}).Output()
	if err != nil {
		t.Error(err)
	}
	if out != "hello" {
		t.Error("capture wrong output:", out)
	}
	s.Command("echo", []string{"hello\tworld"}).Command("cut", []string{"-f2"}).Run()
}

func TestPipeCommand(t *testing.T) {
	c1 := exec.Command("echo", "good")
	rd, wr := io.Pipe()
	c1.Stdout = wr
	c2 := exec.Command("cat", "-n")
	c2.Stdout = os.Stdout
	c2.Stdin = rd
	c1.Start()
	c2.Start()

	c1.Wait()
	wc, ok := c1.Stdout.(io.WriteCloser)
	if ok {
		wc.Close()
	}
	c2.Wait()
}

func TestExample(t *testing.T) {
	s := NewSession()
	s.Env["PATH"] = "/usr/bin:/bin"
	s.Alias("ll", "ls", "-l")
	s.Command("ll").Command("awk", []string{"{print $1}"}).Command("grep", []string{"^-rw"}).Run()
}
