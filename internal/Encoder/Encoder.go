package Encoder

import (
	"crypto/sha1"
	"encoding/base64"
)

func EncodeFromKey(String string, Password string) string {
	Salt := "BGuxLWQtKweKEMV4"
	StrLen := len(String)
	Seq := Password
	Gamma := ""
	for len(Gamma) < StrLen {
		SeqBytes := sha1.Sum([]byte(Gamma + Seq + Salt))
		Seq = string(SeqBytes[:])
		Gamma += Seq[:8]
	}

	StringBytes := []byte(String)
	GammaBytes := []byte(Gamma)
	encodedString := make([]byte, StrLen)
	for i := 0; i < StrLen; i++ {
		encodedString[i] = StringBytes[i] ^ GammaBytes[i]
	}

	return base64.StdEncoding.EncodeToString(encodedString)
}
func DecodeFromKey(String string, Password string) string {
	StringDecoded, _ := base64.StdEncoding.DecodeString(String)
	Salt := "BGuxLWQtKweKEMV4"
	StrLen := len(StringDecoded)
	Seq := Password
	Gamma := ""
	for len(Gamma) < StrLen {
		SeqBytes := sha1.Sum([]byte(Gamma + Seq + Salt))
		Seq = string(SeqBytes[:])
		Gamma += Seq[:8]
	}

	StringXOR := make([]byte, StrLen)
	for i := 0; i < StrLen; i++ {
		StringXOR[i] = StringDecoded[i] ^ Gamma[i]
	}

	return string(StringXOR)
}
