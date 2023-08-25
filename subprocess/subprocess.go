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
	path     string
	args     []string
	pid      int
	executor *exec.Cmd
	lock     *sync.Mutex
	exerr    error
}

func Spwan(name string, cmd string, params ...string) *SubProcess {
	proc := &SubProcess{
		name:     name,
		path:     cmd,
		args:     params,
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
	return s.path
}
func (s *SubProcess) Args() []string {
	return s.args
}
func (s *SubProcess) PID() int {
	return s.pid
}
func (s *SubProcess) Process() *os.Process {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.executor != nil {
		return s.executor.Process
	}
	return nil
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

func (s *SubProcess) listenExit(listener func(error)) {
	// 协程监听导致SubProcess结束的运行错误, Wait表明进程已被杀死, 获取其错误即可
	ec := s.executor.Wait()
	s.lock.Lock()
	s.exerr = ec
	s.executor = nil
	s.pid = -1
	s.lock.Unlock()
	// 将错误传递给监听器函数,
	// 因此通过监听函数立刻获得该错误, 而不需要再开启协程读取对象的IsExitByErr
	// 此外, 也能更清楚地知悉退出原因
	if listener != nil {
		listener(ec)
	}
}

func (s *SubProcess) ensureProcess() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.executor != nil && s.pid != -1 && s.pid > 0 {
		return fmt.Errorf("subprocess has already started")
	}
	if s.executor == nil {
		s.executor = exec.Command(s.Command(), s.Args()...)
		s.pid = -1
		s.exerr = nil
	}
	return nil
}

func (s *SubProcess) StartWithExitListener(listener func(error)) error {
	// 考虑并发场景, 用锁保证一个SubProcess只能被启动一次
	if err := s.ensureProcess(); err != nil {
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	err := s.executor.Start()
	if err == nil {
		s.pid = s.executor.Process.Pid
	}
	go s.listenExit(listener)
	return err
}
func (s *SubProcess) Start() error {
	return s.StartWithExitListener(nil)
}

func (s *SubProcess) Kill() error {
	// 考虑并发场景, 用锁保证一个SubProcess只能被杀死一次
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pid = -1
	if s.executor == nil {
		return fmt.Errorf("subprocess has been killed or not started")
	}
	err := s.executor.Process.Kill()
	s.executor = nil
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
	if err := s.ensureProcess(); err != nil {
		return err
	}
	if stdout, err := s.StdoutPipe(); err != nil {
		return err
	} else {
		s.CopyOutput(stdout, wf)
		return nil
	}
}
func (s *SubProcess) CopyStderr(wf func(*SubProcess, string)) error {
	if err := s.ensureProcess(); err != nil {
		return err
	}
	if stderr, err := s.StderrPipe(); err != nil {
		return err
	} else {
		s.CopyOutput(stderr, wf)
		return nil
	}
}
