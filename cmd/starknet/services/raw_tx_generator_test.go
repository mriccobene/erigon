package services

import (
	"bytes"
	"encoding/hex"
	"github.com/ledgerwatch/erigon/crypto"
	"testing"
	"testing/fstest"
)

func TestCreate(t *testing.T) {
	privateKey := "26e86e45f6fc45ec6e2ecd128cec80fa1d1505e5507dcd2ae58c3130a7a97b48"

	var cases = []struct {
		name       string
		privateKey string
		fileName   string
		want       string
	}{
		{name: "success", privateKey: privateKey, fileName: "contract_test.json", want: "03f86583127ed80180800180019637623232363136323639323233613230356235643764c080a0ceb955e6039bf37dbf77e4452a10b4a47906bbbd2f6dcf0c15bccb052d3bbb60a03de24d584a0a20523f55a137ebc651e2b092fbc3728d67c9fda09da9f0edd154"},
	}

	fs := fstest.MapFS{
		"contract_test.json": {Data: []byte("{\"abi\": []}")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rawTxGenerator := RawTxGenerator{
				privateKey: tt.privateKey,
			}

			buf := bytes.NewBuffer(nil)
			err := rawTxGenerator.CreateFromFS(fs, tt.fileName, buf)

			assertNoError(t, err)

			if hex.EncodeToString(buf.Bytes()) != tt.want {
				t.Error("got not equals want")
			}
		})
	}
}

func TestErrorCreate(t *testing.T) {
	var cases = []struct {
		name       string
		privateKey string
		fileName   string
		error      error
	}{
		{name: "invalid private key", privateKey: "abc", fileName: "not_exist.json", error: ErrInvalidPrivateKey},
		{name: "contract file not found", privateKey: generatePrivateKey(t), fileName: "not_exist.json", error: ErrReadContract},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fs := fstest.MapFS{}

			rawTxGenerator := RawTxGenerator{
				privateKey: tt.privateKey,
			}

			buf := bytes.NewBuffer(nil)
			err := rawTxGenerator.CreateFromFS(fs, tt.fileName, buf)

			if tt.error != nil {
				assertError(t, err, tt.error)
			}
		})
	}
}

func generatePrivateKey(t testing.TB) string {
	t.Helper()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}

	return hex.EncodeToString(crypto.FromECDSA(privateKey))
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
