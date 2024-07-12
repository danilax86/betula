// Package signing manages HTTP signatures and managing a pair of private and public keys. This package is a wrapper around Ted of the Honk's httpsig package.
package signing

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"log"
	"net/http"

	"git.sr.ht/~bouncepaw/betula/db"
	"git.sr.ht/~bouncepaw/betula/settings"
	"humungus.tedunangst.com/r/webs/httpsig"
)

// SignRequest signs the request.
func SignRequest(rq *http.Request, content []byte) {
	keyID := settings.SiteURL() + "/@" + settings.AdminUsername() + "#main-key"
	httpsig.SignRequest(keyID, privateKey, rq, content)
}

var (
	privateKey   httpsig.PrivateKey
	publicKey    httpsig.PublicKey
	publicKeyPEM string
)

func PublicKey() string {
	return publicKeyPEM
}

func setKeys(privateKeyPEM string) {
	var err error
	privateKey, publicKey, err = httpsig.DecodeKey(privateKeyPEM)
	if err != nil {
		log.Fatalf("When decoding private key PEM: %s\n", err)
	}

	publicKeyPEM, err = httpsig.EncodeKey(publicKey.Key)
	if err != nil {
		log.Fatalf("When encoding public key PEM: %s\n", err)
	}
}

// EnsureKeysFromDatabase reads the keys from the database and remembers them. If they are not found, it comes up with new ones and saves them. This function might crash the application.
func EnsureKeysFromDatabase() {
	var pem string
	privKeyPEMMaybe := db.MetaEntry[sql.NullString](db.BetulaMetaPrivateKey)
	if !privKeyPEMMaybe.Valid || privKeyPEMMaybe.String == "" {
		log.Println("Generating a new pair of RSA keys")
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("When generating new keys: %s\n", err)
		}

		pem, err = httpsig.EncodeKey(priv)
		if err != nil {
			log.Fatalf("When generating private key PEM: %s\n", err)
		}

		db.SetMetaEntry(db.BetulaMetaPrivateKey, pem)
		setKeys(pem)
	} else {
		setKeys(privKeyPEMMaybe.String)
	}
}
