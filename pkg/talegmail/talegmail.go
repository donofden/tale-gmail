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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/donofden/tale-gmail/pkg/gmailparser"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// MailMain exported
func MailMain() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}

	var total int64
	msgs := []message{}
	emailMsgs := []gmailparser.Email{}
	pageToken := ""
	for {
		req := srv.Users.Messages.List("me").Q("larger:5M")
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
