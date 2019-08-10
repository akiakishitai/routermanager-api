package openapi

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

// SysTimedateGet - The time date of this machine
func SysTimedateGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	out, err := exec.Command("date", "+\"%Y/%m/%d %H:%M:%S\"").Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(unexpectedError)
		log.Print("Failed `date` command.")
	} else {
		removeQuote := out[1 : len(out) - 2]
		time := Date{
			Date: string(removeQuote),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time)
	}
}

// SysTimedateSync - Synchronize clock to NTP server on this machine
func SysTimedateSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := exec.Command("chronyc", "makestep").Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(unexpectedError)
		log.Print("Failed `chronyc` command.")
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
