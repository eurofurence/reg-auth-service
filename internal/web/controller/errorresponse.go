package controller

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"time"
)

const errTimeFormat = "02.01.-15:04:05"

func ErrorResponse(ctx context.Context, msg string) []byte {
	currentTime := time.Now().Format(errTimeFormat)
	response := fmt.Sprintf(`<HTML>
<HEAD>
  <TITLE>Reg Auth Service Error</TITLE>
  <meta name="robots" content="noindex"/>
  <meta http-equiv="expires" content="0"/>
  <style>body { font-family: Arial, sans-serif; }</style>
</HEAD>
<BODY bgcolor="white">
  <p><b>error:</b> %s</p>
  <p><font color="red"><b>code:</b> %s-%s</font></p>
  <p>If you wish us to investigate an error, please provide us with the code.</p>
</BODY>
</HTML>`, msg, ctxvalues.RequestId(ctx), currentTime)
	return []byte(response)
}
