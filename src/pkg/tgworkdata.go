package pkg

import (
	"encoding/base64"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type TelegramStart struct {
	Method string
	Data   string
}

func ConstructReferralLink(userId uuid.UUID) string {
	linkConstruct := "ref|" + userId.String()
	linkEncoded := base64.StdEncoding.EncodeToString([]byte(linkConstruct))
	linkBuild := os.Getenv("TELEGRAM_BOT_LINK") + "?startapp=" + linkEncoded
	return linkBuild
}

func DecodeStartParam(param string) TelegramStart {
	paramDecoded, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return TelegramStart{"", ""}
	}
	paramDecodedString := string(paramDecoded)
	paramCode := strings.SplitN(paramDecodedString, "|", 2)
	log.Infoln(paramCode)
	if len(paramCode) != 2 {
		return TelegramStart{"", ""}
	}
	switch paramCode[0] {
	case "ref":
		return TelegramStart{"ref", paramCode[1]}
	}
	return TelegramStart{"", ""}
}
