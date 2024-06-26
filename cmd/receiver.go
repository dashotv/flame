/*
Copyright © 2020 Shawn Catanzarite <me@shawncatz.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/dashotv/mercury"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "run flame receiver",
	Long:  "run flame receiver",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := zap.NewDevelopment()
		if err != nil {
			fmt.Printf("error setting up logger: %s\n", err)
			os.Exit(1)
		}
		defer l.Sync()
		log := l.Sugar().Named("receiver")

		m, err := mercury.New("mercury", nats.DefaultURL)
		if err != nil {
			log.Fatalf("creating mercury: %s", err)
		}

		log.Infof("starting receiver...")

		qbittorrents := make(chan *qbt.Response, 5)
		if err := m.Receiver("flame.qbittorrents", qbittorrents); err != nil {
			log.Fatalf("flame torrents receiver: %s", err)
		}

		nzbs := make(chan *nzbget.GroupResponse, 5)
		if err := m.Receiver("flame.nzbs", nzbs); err != nil {
			log.Fatalf("flame nzbs receiver: %s", err)
		}

		downloads := make(chan string, 5)
		if err := m.Receiver("seer.downloads", downloads); err != nil {
			log.Fatalf("seer downloads receiver: %s", err)
		}

		for {
			select {
			case r := <-qbittorrents:
				log.Named("qbt").Infof("%T %s", r, r.Timestamp)
				for _, t := range r.Torrents {
					log.Named("qbt").Infof("%3d %6.2f%% %10.2fmb %8.8s %s", t.Priority, t.Progress, t.SizeMb(), t.State, t.Name)
					for _, f := range t.Files {
						log.Named("qbt").Infof("%3d %6.2f%% %s", f.Priority, f.Progress, f.Name)
					}
				}
			case r := <-nzbs:
				log.Named("nzb").Infof("%T %s", r, r.Timestamp)
				//for _, g := range r.Result {
				//	log.Infof("%5d %25s %s\n", g.ID, g.Status, g.NZBName)
				//}
			case s := <-downloads:
				log.Named("dls").Infof("%#v\n", s)
			case <-time.After(30 * time.Second):
				log.Warn("timeout")
				os.Exit(0)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
