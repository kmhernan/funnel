package pbs

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/tests"
	"os"
	"testing"
)

var fun *tests.Funnel

func TestMain(m *testing.M) {
	conf := tests.DefaultConfig()
	if conf.Compute != "pbs" {
		logger.Debug("Skipping PBS/Torque e2e tests...")
		os.Exit(0)
	}

	fun = tests.NewFunnel(conf)
	fun.StartServerInDocker("ohsucompbio/pbs-torque:latest", []string{"--hostname", "docker", "--privileged"})
	defer fun.CleanupTestServerContainer()

	m.Run()
	return
}

func TestHelloWorld(t *testing.T) {
	id := fun.Run(`
    --sh 'echo hello world'
  `)
	task := fun.Wait(id)

	if task.Logs[0].Logs[0].Stdout != "hello world\n" {
		t.Fatal("Missing stdout")
	}
}

func TestResourceRequest(t *testing.T) {
	id := fun.Run(`
    --sh 'echo I need resources!' --cpu 1 --ram 2 --disk 5
  `)
	task := fun.Wait(id)

	if task.Logs[0].Logs[0].Stdout != "I need resources!\n" {
		t.Fatal("Missing stdout")
	}
}
