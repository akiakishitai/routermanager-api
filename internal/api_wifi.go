package openapi

import (
	"strconv"
	"strings"
	"io"
	"io/ioutil"
	"log"
	"encoding/json"
	"net/http"
	"os/exec"
)

// WifiReconnect - Re-connect to wifi
func WifiReconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// WifiStatus - Status of wifi
func WifiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	cmd := exec.Command("nmcli", "-t", "-f", "device,type,state,connection", "dev", "status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print(err)
		json.NewEncoder(w).Encode(UnexpectedError("Failed to stdout pipe."))
		return
	}
	
	if err := cmd.Start(); err != nil {
		log.Print(err)
		json.NewEncoder(w).Encode(UnexpectedError("Failed to get wifi status."))
		return
	}

	bs, _ := ioutil.ReadAll(stdout)
	cmd.Wait()

	cmd1 := exec.Command("grep", "wifi", "-")
	stdin, err := cmd1.StdinPipe()
	if err != nil {
		log.Print(err)
		json.NewEncoder(w).Encode(UnexpectedError("Failed to stdin pipe."))
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(bs))
	}()
	
	out, _ := cmd1.CombinedOutput()
	if string(out) == "" {
		json.NewEncoder(w).Encode(UnexpectedError("Not found wifi device."))
		return
	}

	slice := strings.Split(string(out), ":")
	stat := Status{
		Device: slice[0],
		Connected: (slice[2] == "connected" || slice[2] == "接続済み"),
		Signal: 23,
	}

	cmd = exec.Command("nmcli", "-t", "-f", "active,device,ssid,signal", "dev", "wifi")
	stdout, _ = cmd.StdoutPipe()
	_ = cmd.Start()
	bs, _ = ioutil.ReadAll(stdout)
	cmd.Wait()

	cmd = exec.Command("grep", "^yes", "-")
	stdin, err = cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(bs))
	}()

	out, _ = cmd.CombinedOutput()
	if string(out) == "" {
		json.NewEncoder(w).Encode(UnexpectedError("Not found active wifi"))
		return
	}

	slice = strings.Split(string(out), ":")
	signal, parseErr := strconv.ParseInt(strings.TrimRight(slice[3], "\n"), 10, 32)
	if parseErr != nil {
		json.NewEncoder(w).Encode(UnexpectedError(parseErr.Error()))
		return
	}
	stat.Signal = int32(signal)
	json.NewEncoder(w).Encode(stat)
}

func UnexpectedError(msg string) Error {
	return Error {
		Code: 400,
		Message: msg,
	}
}
