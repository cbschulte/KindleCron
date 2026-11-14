package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestShutdownRejectsJobBeforeProcessStart(t *testing.T) {
	dm := newDaemon()

	id, accepted := dm.jobDispatched("shutdown-before-start")
	if !accepted {
		t.Fatal("expected job dispatch to succeed")
	}

	// Simulate shutdown before the dispatched job starts.
	dm.killRunningJobs()

	cmd := exec.Command("sh", "-c", "exit 0")
	callbackCalled := false

	_, err := dm.startJobProcess(id, "shutdown-before-start", cmd, func() {
		callbackCalled = true
	})

	if err == nil {
		t.Fatal("expected process start to be rejected during shutdown")
	}

	if cmd.Process != nil {
		t.Fatal("process was started despite shutdown")
	}

	if callbackCalled {
		t.Fatal("onStarted callback must not be called")
	}

	if dm.hasRunningProcess() {
		t.Fatal("no process should be registered")
	}

	dm.jobFinished(id, "shutdown-before-start")

	if dm.anyJobRunning() {
		t.Fatal("pending count should be zero after jobFinished")
	}
}

func TestShutdownKillsStartedJobProcess(t *testing.T) {
	dm := newDaemon()
	name := "daemon-test-running-process"

	// Clean up a possible leftover PID file.
	defer os.Remove(jobPidPath(name))

	id, accepted := dm.jobDispatched(name)
	if !accepted {
		t.Fatal("expected job dispatch to succeed")
	}

	cmd := exec.Command("sh", "-c", "sleep 30")
	setProcessGroup(cmd)

	callbackCalled := false

	_, err := dm.startJobProcess(id, name, cmd, func() {
		callbackCalled = true
	})
	if err != nil {
		t.Fatalf("startJobProcess failed: %v", err)
	}

	if !callbackCalled {
		t.Fatal("onStarted callback was not called")
	}

	if !dm.hasRunningProcess() {
		t.Fatal("started process was not registered")
	}

	if cmd.Process == nil {
		t.Fatal("expected a started process")
	}

	killed := dm.killRunningJobs()
	if killed != 1 {
		t.Fatalf("expected 1 killed process, got %d", killed)
	}

	// Reap the terminated process.
	if err := cmd.Wait(); err == nil {
		t.Fatal("expected process to have been terminated")
	}

	dm.jobFinished(id, name)

	if dm.hasRunningProcess() {
		t.Fatal("process should no longer be registered")
	}

	if dm.anyJobRunning() {
		t.Fatal("pending count should be zero after jobFinished")
	}
}
