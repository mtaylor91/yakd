package cmd

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var DefaultEtcdEndpoints []string = []string{
	"127.0.0.1:2379",
}

type Config struct {
	Etcd            ConfigEtcd    `mapstructure:"etcd"`
	IntervalSeconds time.Duration `mapstructure:"interval_seconds"`
	TimeoutSeconds  time.Duration `mapstructure:"timeout_seconds"`
}

type ConfigEtcd struct {
	Endpoints          []string      `mapstructure:"endpoints"`
	DialTimeoutSeconds time.Duration `mapstructure:"connect"`
}

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
	ctx := context.Background()

	f := cmd.Flags()
	v := viper.New()
	v.BindPFlag("etcd.endpoints", f.Lookup("etcd-endpoints"))
	v.BindPFlag("interval_seconds", f.Lookup("interval"))
	v.BindPFlag("timeout_seconds", f.Lookup("timeout"))

	err := v.Unmarshal(config)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{}).Info("Connecting to etcd cluster")
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Etcd.Endpoints,
		DialTimeout: time.Second * config.Etcd.DialTimeoutSeconds,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer etcd.Close()

	for {
		log.Info("Checking state...")
		timeoutCtx, _ := context.WithTimeout(ctx, time.Second*config.TimeoutSeconds)
		etcd.Sync(timeoutCtx)
		time.Sleep(time.Second * config.IntervalSeconds)
	}
}
