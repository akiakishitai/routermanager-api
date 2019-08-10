package openapi

import (
	"encoding/json"
	"strings"
	"io/ioutil"
	"log"
	"os/exec"
	"net/http"
)

// SysPoweroff - System shutdown or reboot
func SysPoweroff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cmds := map[Power]*(exec.Cmd){
		SHUTDOWN: exec.Command("sleep", "3s", "&&", "poweroff"),
		REBOOT: exec.Command("sleep", "3s", "&&", "reboot"),
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(unexpectedError)
	}
	
	v, ok := cmds[Power(strings.ToLower(string(bs)))]
	if !ok {
		msg := "Not allowed this request. Please 'shutdown' or 'reboot'"
		res := Error{
			Code: 400,
			Message: msg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		log.Print(msg)
		return
	}

	if err = v.Start(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(unexpectedError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
