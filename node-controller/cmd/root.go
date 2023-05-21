package cmd

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/mtaylor91/yakd/node-controller/pkg"
)

var DefaultEtcdEndpoints []string = []string{}

type Config struct {
	Etcd            ConfigEtcd    `mapstructure:"etcd"`
	IntervalSeconds time.Duration `mapstructure:"interval"`
	TimeoutSeconds  time.Duration `mapstructure:"timeout"`
}

type ConfigEtcd struct {
	Endpoints          []string      `mapstructure:"endpoints"`
	DialTimeoutSeconds time.Duration `mapstructure:"connect"`
}

type State struct {
	epoch     int
	tasks     []Task
	scheduler pkg.Scheduler
}

type Task func(context.Context, *State)

var RootCmd = &cobra.Command{
	Use: "node-controller",
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Node Controller",
	Run:   Run,
}

func init() {
	RootCmd.AddCommand(RunCmd)
	runFlags := RunCmd.Flags()
	runFlags.StringSlice("etcd-endpoints", DefaultEtcdEndpoints, "Etcd Endpoints")
	runFlags.Duration("interval", 3, "Controller interval seconds")
	runFlags.Duration("timeout", 3, "Controller timeout seconds")
}

func Main() {
	RootCmd.Execute()
}

func Run(cmd *cobra.Command, args []string) {
	log.Info("Running YAKD node controller")

	config := &Config{}

	f := cmd.Flags()
	v := viper.New()
	v.BindPFlag("etcd.endpoints", f.Lookup("etcd-endpoints"))
	v.BindPFlag("interval", f.Lookup("interval"))
	v.BindPFlag("timeout", f.Lookup("timeout"))

	err := v.Unmarshal(config)
	if err != nil {
		log.Fatal(err)
	}

	state := &State{
		epoch:     0,
		tasks:     []Task{},
		scheduler: pkg.Scheduler{},
	}

	if len(config.Etcd.Endpoints) > 0 {
		log.WithFields(log.Fields{
			"endpoints": strings.Join(config.Etcd.Endpoints, ","),
		}).Info("Connecting to etcd cluster")
		if etcd, err := clientv3.New(clientv3.Config{
			Endpoints:   config.Etcd.Endpoints,
			DialTimeout: time.Second * config.Etcd.DialTimeoutSeconds,
		}); err != nil {
			log.Fatalf("Etcd connection failed: %v", err)
		} else {
			state.tasks = append(state.tasks, func(ctx context.Context, s *State) {
				etcd.Sync(ctx)
				log.WithFields(log.Fields{
					"endpoints": etcd.Endpoints(),
				}).Info("Synced etcd client endpoints")
			})

			defer etcd.Close()
		}
	}

	for {
		state.epoch += 1
		log.WithFields(log.Fields{
			"epoch": state.epoch,
		}).Info("Checking state...")
		timeout := time.Second * config.TimeoutSeconds
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		for _, task := range state.tasks {
			task(ctx, state)
		}

		time.Sleep(time.Second * config.IntervalSeconds)
	}
}
