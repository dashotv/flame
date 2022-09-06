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
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/mercury"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "run flame receiver",
	Long:  "run flame receiver",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.Instance()

		m, err := mercury.New("mercury", nats.DefaultURL)
		if err != nil {
			logrus.Fatalf("creating mercury: %w", err)
		}

		app.Log.Infof("starting receiver...")

		qbittorrents := make(chan *qbt.Response, 5)
		if err := m.Receiver("flame.qbittorrents", qbittorrents); err != nil {
			logrus.Fatalf("flame torrents receiver: %w", err)
		}

		nzbs := make(chan *nzbget.GroupResponse, 5)
		if err := m.Receiver("flame.nzbs", nzbs); err != nil {
			logrus.Fatalf("flame nzbs receiver: %w", err)
		}

		downloads := make(chan string, 5)
		if err := m.Receiver("seer.downloads", downloads); err != nil {
			logrus.Fatalf("seer downloads receiver: %w", err)
		}

		for {
			select {
			case r := <-qbittorrents:
				app.Log.WithField("prefix", "qbt").Infof("%T %s", r, r.Timestamp)
				for _, t := range r.Torrents {
					app.Log.WithField("prefix", "qbt").Infof("%3d %6.2f%% %10.2fmb %8.8s %s", t.Priority, t.Progress, t.SizeMb(), t.State, t.Name)
					for _, f := range t.Files {
						app.Log.WithField("prefix", "qbt").Infof("%3d %6.2f%% %s", f.Priority, f.Progress, f.Name)
					}
				}
			case r := <-nzbs:
				app.Log.WithField("prefix", "nzb").Infof("%T %s", r, r.Timestamp)
				//for _, g := range r.Result {
				//	logrus.Infof("%5d %25s %s\n", g.ID, g.Status, g.NZBName)
				//}
			case s := <-downloads:
				app.Log.WithField("prefix", "dls").Infof("%#v\n", s)
			case <-time.After(30 * time.Second):
				app.Log.Warn("timeout")
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
