package openapi

import (
	"strings"
	"io/ioutil"
	"log"
	"os/exec"
	"net/http"
)

// SysPoweroff - System shutdown or reboot
func SysPoweroff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	cmds := map[Power]*(exec.Cmd){
		SHUTDOWN: exec.Command("sleep", "3s", "&&", "poweroff"),
		REBOOT: exec.Command("sleep", "3s", "&&", "reboot"),
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	v, ok := cmds[Power(strings.ToLower(string(bs)))]
	if ok {
		v.Start()
	} else {
		log.Print("Not allowed this request. Please 'shutdown' or 'reboot'")
	}
}
