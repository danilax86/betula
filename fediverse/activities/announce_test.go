package activities

import (
	"errors"
	"reflect"
	"testing"
)

const json1 = `
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "type": "Announce",
    "actor": {
        "type": "Person",
        "id": "https://links.alice",
        "inbox": "https://links.alice/inbox",
        "name": "Alice",
        "preferredUsername": "alice"
    },
    "id": "https://links.alice/84",
    "object": "https://links.bob/42"
}`

const jsonNoID = `
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "type": "Announce",
    "actor": {
        "type": "Person",
        "id": "https://links.alice",
        "inbox": "https://links.alice/inbox",
        "name": "Alice",
        "preferredUsername": "alice"
    },
    "object": "https://links.bob/42"
}`

const jsonBadID = `
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "type": "Announce",
    "actor": {
        "type": "Person",
        "id": "https://links.alice",
        "inbox": "https://links.alice/inbox",
        "name": "Alice",
        "preferredUsername": "alice"
    },
	"id": {"c'était trop beau": 4},
    "object": "https://links.bob/42"
}`

const jsonNoUsername = `
{
    "@context": "https://www.w3.org/ns/activitystreams",
    "type": "Announce",
    "actor": {
        "type": "Person",
        "id": "https://links.alice",
        "inbox": "https://links.alice/inbox",
        "name": "Alice"
    },
    "id": "https://links.alice/84",
    "object": "https://links.bob/42"
}`

const jsonSallyOffered = `
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "summary": "Sally offered the Foo object",
  "type": "Offer",
  "actor": {
    "type": "Person",
    "id": "http://sally.example.org",
    "summary": "Sally"
  },
  "object": "http://example.org/foo"
}`

const badJSON = `Laika`

var table = []struct {
	json   string
	err    error
	report any
}{
	{json1, nil, AnnounceReport{
		ReposterUsername: "alice",
		RepostPage:       "https://links.alice/84",
		OriginalPage:     "https://links.bob/42",
	}},
	{jsonNoID, ErrNoID, nil},
	{jsonBadID, ErrNoID, nil},
	{jsonNoUsername, ErrNoActorUsername, nil},
	{badJSON, errors.New("invalid character 'L' looking for beginning of value"), nil},
	{jsonSallyOffered, ErrUnknownType, nil},
	// one might want to write many more tests
}

func TestGuess(t *testing.T) {
	for i, test := range table {
		report, err := Guess([]byte(test.json))
		if test.err != nil && err.Error() != test.err.Error() {
			t.Errorf("Error failed. Test %d: %q ≠ %q", i+1, err, test.err)
		}
		if report == nil && test.report == nil {
			continue
		}
		if reflect.TypeOf(report) != reflect.TypeOf(test.report) {
			t.Errorf("Report types mismatch. Test %d: %v ≠ %v", i+1, report, test.report)
		}
		switch r := test.report.(type) {
		case AnnounceReport:
			R := report.(AnnounceReport)
			if !reflect.DeepEqual(r, R) {
				t.Errorf("Report failed. Test %d: %v ≠ %v", i+1, report, test.report)
			}
		default:
			panic("how did this happen")
		}
	}
}
