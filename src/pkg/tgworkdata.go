package pkg

import (
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

type TelegramStart struct {
	method string
	data   string
}

func ConstructReferralLink(userId uuid.UUID) string {
	linkConstruct := "ref|" + userId.String()
	linkEncoded := base64.StdEncoding.EncodeToString([]byte(linkConstruct))
	return linkEncoded
}

func DecodeStartParam(param string) TelegramStart {
	paramDecoded, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return TelegramStart{"", ""}
	}
	paramDecodedString := string(paramDecoded)
	paramCode := strings.SplitN(paramDecodedString, "|", 1)
	if len(paramCode) != 2 {
		return TelegramStart{"", ""}
	}
	switch paramCode[0] {
	case "ref":
		return TelegramStart{"ref", paramCode[1]}
	}
	return TelegramStart{"", ""}
}
