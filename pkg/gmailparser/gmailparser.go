// Package gmailparser implements a way to parse emails.
package gmailparser

import (
	"google.golang.org/api/gmail/v1"
)

// GetMessageBody finds the HTML body of an email.
func GetMessageBody(parts []*gmail.MessagePart) string {
	for _, part := range parts {
		if len(part.Parts) > 0 {
			return GetMessageBody(part.Parts)
		} else {
			if part.MimeType == "text/html" {
				return part.Body.Data
			}
		}
	}

	return ""
}

// GetMessageSender goes through the headers to find the From header.
func GetMessageSender(headers []*gmail.MessagePartHeader) string {
	return GetMessageHeader(headers, "From")
}

// GetMessageSubject goes through the headers to find the Subject header.
func GetMessageSubject(headers []*gmail.MessagePartHeader) string {
	return GetMessageHeader(headers, "Subject")
}

// GetMessageHeader goes through a list of headers and returns the header where
// the name matches the one we want.
func GetMessageHeader(headers []*gmail.MessagePartHeader, wanted string) string {
	for _, header := range headers {
		if header.Name == wanted {
			return header.Value
		}

	}

	return ""
}

// Email represents an email fetched from your gmail account.
type Email struct {
	Size    int64  `json:"size"`
	Subject string `json:"subject"`
	Body    string `json:"body"` // Base64.URLEncoding
	ID      string `json:"id"`
	Sender  string `json:"sender"`
	Date    string `json:"date"`
}
