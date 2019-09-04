package http

import (
	"net/http"
	"strings"
        "log"
	"mail-provider/config"
	"mail-provider/db"
	"mail-provider/smtp"
	"github.com/toolkits/web/param"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		log.Println(content)
                if strings.Contains(strings.Split(content, "\r\n")[1],"P6"){
                    db.AlarmVerificationSuccess(content)
                    return
                }
		s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
		err := s.SendMail(cfg.Smtp.From, tos, subject, content)
                db.AlarmVerificationSuccess(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
