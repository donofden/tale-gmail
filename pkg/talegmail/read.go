/**
 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// [START tale_gmail]
package talegmail

import (
	"fmt"
	"log"
	"sort"

	"github.com/donofden/tale-gmail/pkg/gmailparser"
)

// ReadMail exported
func ReadMail() {
	srv := gmailService()

	user := "me"

	var total int64
	msgs := []message{}
	emailMsgs := []gmailparser.Email{}
	pageToken := ""
	for {
		req := srv.Users.Messages.List(user).Q("larger:5M")
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}

		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {
			msg, err := srv.Users.Messages.Get("me", m.Id).Do()

			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}
			total += msg.SizeEstimate
			date := ""
			for _, h := range msg.Payload.Headers {
				if h.Name == "Date" {
					date = h.Value
					break
				}
			}
			//body := mail.GetMessageBody(msg.Payload.Parts)
			sender := gmailparser.GetMessageSender(msg.Payload.Headers)
			subject := gmailparser.GetMessageSubject(msg.Payload.Headers)

			msgs = append(msgs, message{
				size:    msg.SizeEstimate,
				gmailID: msg.Id,
				date:    date,
				snippet: msg.Snippet,
			})
			emailMsgs = append(emailMsgs, gmailparser.Email{
				Subject: subject,
				Body:    msg.Id,
				ID:      msg.Id,
				Sender:  sender,
			})
		}
		fmt.Printf("sender <%v> \n", emailMsgs)

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
	log.Printf("total: %v\n", total)

	sortBySize(msgs)
}

type messageSorter struct {
	msg  []message
	less func(i, j message) bool
}

func sortBySize(msg []message) {
	sort.Sort(messageSorter{msg, func(i, j message) bool {
		return i.size > j.size
	}})
}

func (s messageSorter) Len() int {
	return len(s.msg)
}

func (s messageSorter) Swap(i, j int) {
	s.msg[i], s.msg[j] = s.msg[j], s.msg[i]
}

func (s messageSorter) Less(i, j int) bool {
	return s.less(s.msg[i], s.msg[j])
}

type message struct {
	size    int64
	gmailID string
	date    string // retrieved from message header
	snippet string
}

// [END gmail_quickstart]
