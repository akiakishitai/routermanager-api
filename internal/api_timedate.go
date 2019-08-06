package openapi

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

// SysTimedateGet - The time date of this machine
func SysTimedateGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	out, err := exec.Command("date", "+\"%Y/%m/%d %H:%M:%S\"").Output()
	if err != nil {
		message := Error {
			Code: 400,
			Message: "Failed to exec `date` command.",
		}
		json.NewEncoder(w).Encode(message)
	} else {
		removeQuote := out[1 : len(out) - 2]
		time := Date{
			Date: string(removeQuote),
		}
		json.NewEncoder(w).Encode(time)
	}
}

// SysTimedateSync - Synchronize clock to NTP server on this machine
func SysTimedateSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := exec.Command("chronyc", "makestep").Run()
	if err != nil {
		message := Error {
			Code: 400,
			Message: "Failed to exec `chronyc` command.",
		}
		json.NewEncoder(w).Encode(message)
	}
}
