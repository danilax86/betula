package fediverse

import (
	"encoding/json"
	"fmt"
	"git.sr.ht/~bouncepaw/betula/fediverse/signing"
	"git.sr.ht/~bouncepaw/betula/types"
	"io"
	"log"
	"net/http"
	"strings"
)

// RequestActorByNickname returns actor by string like @bouncepaw@links.bouncepaw.com or bouncepaw@links.bouncepaw.com.
func RequestActorByNickname(nickname string) (*types.Actor, error) {
	user, host, ok := strings.Cut(strings.TrimPrefix(nickname, "@"), "@")

	if !ok {
		return nil, fmt.Errorf("bad username: %s", nickname)
	}

	wa, found, err := RequestWebFinger(user, host)
	if !found {
		return nil, fmt.Errorf("user not found 404: %s", nickname)
	}
	if err != nil {
		return nil, err
	}

	actor, err := RequestActor(wa.ActorURL)
	if err != nil {
		return nil, fmt.Errorf("while fetching actor %s: %w", wa.ActorURL, err)
	}

	return actor, nil
}

// RequestActor fetches the actor activity on the specified address.
func RequestActor(actorID string) (*types.Actor, error) {
	if cachedActor, ok := ActorStorage[actorID]; ok {
		log.Printf("Returning cached author %s\n", actorID)
		return cachedActor, nil
	}

	cope := func(err error) error {
		return fmt.Errorf("requesting actor: %w", err)
	}

	req, err := http.NewRequest("GET", actorID, nil)
	if err != nil {
		return nil, cope(err)
	}
	req.Header.Set("Accept", types.ActivityType)
	signing.SignRequest(req, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, cope(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting actor: status not 200, id est %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, cope(err)
	}

	var a types.Actor
	if err = json.Unmarshal(data, &a); err != nil {
		return nil, cope(err)
	}

	ActorStorage[a.ID] = &a
	KeyPEMStorage[a.PublicKey.ID] = a.PublicKey.PublicKeyPEM

	return &a, nil
}

func RequestActorInbox(actorID string) string {
	actor, err := RequestActor(actorID)
	if err != nil {
		log.Printf("When requesting actor %s inbox: %s\n", actorID, err)
		return ""
	}
	return actor.Inbox
}
