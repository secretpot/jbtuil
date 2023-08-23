package subprocess

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type SubProcess struct {
	name     string
	pid      int
	executor *exec.Cmd
	lock     *sync.Mutex
	exerr    error
}

func Spwan(name string, cmd string, params ...string) *SubProcess {
	proc := &SubProcess{
		name:     name,
		pid:      -1,
		executor: exec.Command(cmd, params...),
		lock:     new(sync.Mutex),
		exerr:    nil,
	}
	return proc
}
func (s *SubProcess) Name() string {
	return s.name
}
func (s *SubProcess) Command() string {
	return s.executor.Path
}
func (s *SubProcess) Args() []string {
	return s.executor.Args
}
func (s *SubProcess) Env() []string {
	return s.executor.Env
}
func (s *SubProcess) PID() int {
	return s.pid
}
func (s *SubProcess) Process() *os.Process {
	return s.executor.Process
}
func (s *SubProcess) IsAlive() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.pid > 0 && s.executor != nil && s.exerr == nil
}
func (s *SubProcess) IsExitByErr() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return (s.pid <= 0 || s.executor == nil) && s.exerr != nil
}

func (s *SubProcess) listenError() {
	// 协程监听导致SubProcess结束的运行错误, Wait表明进程已被杀死, 获取其错误即可
	ec := make(chan error, 1)
	ec <- s.executor.Wait()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.exerr = <-ec
	s.executor = nil
	s.pid = -1
}
func (s *SubProcess) Start() error {
	// 考虑并发场景, 用锁保证一个SubProcess只能被启动一次
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.pid != -1 && s.pid > 0 {
		return fmt.Errorf("subprocess has already started")
	}
	if s.executor == nil {
		s.executor = exec.Command(s.Command(), s.Args()...)
	}
	err := s.executor.Start()
	if err == nil {
		s.pid = s.executor.Process.Pid
	}
	go s.listenError()
	return err
}

func (s *SubProcess) Kill() error {
	// 考虑并发场景, 用锁保证一个SubProcess只能被杀死一次
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.executor == nil {
		return fmt.Errorf("subprocess has been killed or not started")
	}
	err := s.executor.Process.Kill()
	s.executor = nil
	s.pid = -1
	return err
}

func (s *SubProcess) StdoutPipe() (io.ReadCloser, error) {
	return s.executor.StdoutPipe()
}
func (s *SubProcess) StderrPipe() (io.ReadCloser, error) {
	return s.executor.StderrPipe()
}

func (s *SubProcess) CopyOutput(r io.Reader, wf func(*SubProcess, string)) {
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			wf(s, scanner.Text())
		}
	}()
}
func (s *SubProcess) CopyStdout(wf func(*SubProcess, string)) error {
	if stdout, err := s.StdoutPipe(); err != nil {
		return err
	} else {
		s.CopyOutput(stdout, wf)
		return nil
	}
}
func (s *SubProcess) CopyStderr(wf func(*SubProcess, string)) error {
	if stderr, err := s.StderrPipe(); err != nil {
		return err
	} else {
		s.CopyOutput(stderr, wf)
		return nil
	}
}
