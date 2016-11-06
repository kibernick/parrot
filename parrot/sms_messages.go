package parrot

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// SMSMessage is a container for SMS message data that will be sent to MessageBird.
type SMSMessage struct {
	Recipient  string
	Originator string
	Message    string
	UDHHeader  string
}

// splitMessageIntoParts splits a simple message string into multiple parts, ready to be
// transmitted as concatenated SMS.
func splitMessageIntoParts(message string, partSize int) ([]string, error) {
	if partSize < 1 {
		return nil, fmt.Errorf("Invalid partSize: %v", partSize)
	}
	var msgParts []string

	for len(message) > partSize {
		msgParts = append(msgParts, message[:partSize])
		message = message[partSize:]
	}
	msgParts = append(msgParts, message)
	return msgParts, nil
}

// generateUDHHeader builds a header for sending concatenated SMS messages.
// https://en.wikipedia.org/wiki/Concatenated_SMS#Sending_a_concatenated_SMS_using_a_User_Data_Header
func generateUDHHeader(index, total int, ref byte) (string, error) {
	if index > 255 {
		return "", fmt.Errorf("index (%v) not in [0-255]", index)
	}
	if total < 1 || total > 255 {
		return "", fmt.Errorf("total (%v) not in [1-255]", total)
	}
	if index >= total {
		return "", fmt.Errorf("index (%v) cannot be >= total (%v)", index, total)
	}
	i := byte(index + 1)
	octets := []byte{
		byte(5),     // Field 1 (1 octet): Length of User Data Header, in this case 05.
		byte(0),     // Field 2 (1 octet): Information Element Identifier/reference number, equal to 00 (Concatenated short messages).
		byte(3),     // Field 3 (1 octet): Length of the header, excluding the first two fields; equal to 03
		ref,         // Field 4 (1 octet): 00-FF, CSMS reference number, must be same for all the SMS parts.
		byte(total), // Field 5 (1 octet): 00-FF, total number of parts.
		i,           // Field 6 (1 octet): 00-FF, this part's number in the sequence.
	}
	return hex.EncodeToString(octets), nil
}

// generateUDHRef returns the reference byte to be used in the UDH.
func generateUDHRef() byte {
	b := make([]byte, 1)
	rand.Read(b) // Note that err == nil only if we read len(b) bytes.
	return b[0]
}
