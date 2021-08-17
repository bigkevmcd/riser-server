package resources

import (
	"crypto/rsa"
	"fmt"
	"io"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/util/cert"

	sealedCrypto "github.com/bitnami-labs/sealed-secrets/pkg/crypto"
)

// SealedSecret represents an encrypted Secret
//
// This is a copy of the Bbitnami SealedSecret object.
type SealedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SealedSecretSpec `json:"spec"`
}

// SealedSecretSpec is the spec for the SealedSecret.
type SealedSecretSpec struct {
	EncryptedData map[string][]byte `json:"encryptedData"`
}

// CreateSealedSecret creates and returns a new sealed secret.
// TODO: Consider using something like https://github.com/awnumar/memguard instead of passing the secret as a string
func CreateSealedSecret(plaintextSecret string, secretMeta *core.SecretMeta, certBytes []byte, rand io.Reader) (*SealedSecret, error) {
	publicKey, err := parsePublicKey(certBytes)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing public key")
	}
	objectMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s-%d", secretMeta.App.Name, secretMeta.Name, secretMeta.Revision),
		Namespace: secretMeta.App.Namespace,
		Annotations: map[string]string{
			riserLabel("revision"):       fmt.Sprintf("%d", secretMeta.Revision),
			riserLabel("server-version"): util.VersionString,
		},
		Labels: map[string]string{
			riserLabel("app"): secretMeta.App.Name,
		},
	}
	ciphertext, err := sealSecret(objectMeta, publicKey, []byte(plaintextSecret), rand)
	if err != nil {
		return nil, errors.Wrap(err, "error sealing secret")
	}
	return &SealedSecret{
		ObjectMeta: objectMeta,
		TypeMeta: metav1.TypeMeta{
			Kind:       "SealedSecret",
			APIVersion: "bitnami.com/v1alpha1",
		},
		Spec: SealedSecretSpec{
			EncryptedData: map[string][]byte{
				"data": ciphertext,
			},
		},
	}, nil
}

// Derived from https://github.com/bitnami-labs/sealed-secrets/blob/d875137740275f7dea36c54f981a90c795e7e681/cmd/kubeseal/main.go#L75
func parsePublicKey(certBytes []byte) (*rsa.PublicKey, error) {
	certs, err := cert.ParseCertsPEM(certBytes)
	if err != nil {
		return nil, err
	}

	// ParseCertsPem returns error if len(certs) == 0, but best to be sure...
	if len(certs) == 0 {
		return nil, errors.New("failed to read any certificates")
	}

	cert, ok := certs[0].PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected RSA public key but found %T", certs[0].PublicKey)
	}

	return cert, nil
}

func sealSecret(secretMeta metav1.ObjectMeta, publicKey *rsa.PublicKey, plaintext []byte, rand io.Reader) ([]byte, error) {
	// Simplified version of labelFor (https://github.com/bitnami-labs/sealed-secrets/blob/d875137740275f7dea36c54f981a90c795e7e681/pkg/apis/sealed-secrets/v1alpha1/sealedsecret_expansion.go#L22)
	// We don't support namespace or cluster wide annotations
	label := []byte(fmt.Sprintf("%s/%s", secretMeta.GetNamespace(), secretMeta.GetName()))
	return sealedCrypto.HybridEncrypt(rand, publicKey, plaintext, label)
}
