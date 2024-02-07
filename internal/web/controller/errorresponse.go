package controller

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"time"
)

const errTimeFormat = "02.01.-15:04:05"

func ErrorResponse(ctx context.Context, msg string, retryUrl string) []byte {
	currentTime := time.Now().Format(errTimeFormat)
	retry := ""
	if retryUrl != "" {
		retry = fmt.Sprintf(`<p>You can also <a href="%s">go back to try again</a>.</p>`, retryUrl)
	}
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
  %s
</BODY>
</HTML>`, msg, ctxvalues.RequestId(ctx), currentTime, retry)
	return []byte(response)
}
