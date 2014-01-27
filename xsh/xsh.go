package xsh

import (
	"io"
	"os"
	"os/exec"

	"github.com/codegangsta/inject"
)

func Call(a ...interface{}) error {
	return NewSession().Call(a...)
}

type Dir string

type Session struct {
	inj    inject.Injector
	Env    map[string]string
	Output io.Writer
}

func NewSession(a ...interface{}) *Session {
	env := map[string]string{
		"PATH": "/bin:/usr/bin:/usr/local/bin",
	}
	s := &Session{
		inj:    inject.New(),
		Output: os.Stdout,
		Env:    env,
	}
	dir := Dir("")
	args := []string{}
	s.inj.Map(env).Map(dir).Map(args)
	for _, v := range a {
		s.inj.Map(v)
	}
	return s
}

func (s *Session) Call(a ...interface{}) error {
	for _, v := range a {
		s.inj.Map(v)
	}
	values, err := s.inj.Invoke(s.invokeExec)
	if err != nil {
		return err
	}
	r := values[0]
	if r.IsNil() {
		return nil
	}
	return r.Interface().(error)
}

func (s *Session) invokeExec(cmd string, args []string, cwd Dir) error {
	envs := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		envs = append(envs, k+"="+v)
	}
	//fmt.Println(cmd, args)
	c := exec.Command(cmd, args...)
	c.Env = envs
	c.Dir = string(cwd)
	c.Stdout = s.Output
	c.Stderr = s.Output
	return c.Run()
}
