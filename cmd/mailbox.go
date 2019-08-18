/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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

	"github.com/donofden/tale-gmail/pkg/talegmail"
	"github.com/spf13/cobra"
)

// ListOptions provides the flags for the `list` command
type ListOptions struct {
	UnreadMails bool
	Drafts      bool
	refs        []string
}

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read All Gmail",
	Long:  `Read All Gmail`,
	Run: func(cmd *cobra.Command, args []string) {
		talegmail.ReadMail()
	},
}

// mailboxCmd represents the mailbox command
var mailboxCmd = &cobra.Command{

	Use:   "mailbox",
	Short: "To Work on the Mailbox",
	Long:  `To Work on the Mailbox`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			talegmail.Mailbox(args[0])
		}
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	var opts ListOptions
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(mailboxCmd)

	flags := mailboxCmd.Flags()
	flags.BoolVar(&opts.UnreadMails, "unread-emails", false, "View all unread emails")
	flags.BoolVar(&opts.Drafts, "drafts", false, "View all drafts")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
