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

var (
	unexpectedError = Error {
		Code: 500,
		Message: "Unexpected error",
	}

	notfoundError = Error {
		Code: 404,
		Message: "Not found",
	}
)

// WifiReconnect - Re-connect to wifi
func WifiReconnect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cmd := exec.Command("systemctl", "restart", "network-manager.service")
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(unexpectedError)
		log.Print("Cannot restart network-manager service.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// WifiStatus - Status of wifi
func WifiStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	errFn := func(logMsg string, response Error) {
		log.Print(logMsg)
		w.WriteHeader(response.httpStatus())
		json.NewEncoder(w).Encode(response)
	}

	cmd := exec.Command("nmcli", "-t", "-f", "device,type,state,connection", "dev", "status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errFn(err.Error(), unexpectedError)
		return
	}
	
	if err := cmd.Start(); err != nil {
		errFn(err.Error(), unexpectedError)
		return
	}

	bs, _ := ioutil.ReadAll(stdout)
	cmd.Wait()

	cmd1 := exec.Command("grep", "wifi", "-")
	stdin, err := cmd1.StdinPipe()
	if err != nil {
		errFn(err.Error(), unexpectedError)
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(bs))
	}()
	
	out, _ := cmd1.CombinedOutput()
	if string(out) == "" {
		notfoundError.Message = "Not found a wifi device"
		errFn(notfoundError.Message, notfoundError)
		return
	}

	slice := strings.Split(string(out), ":")
	stat := Status{
		Device: slice[0],
		Connected: (slice[2] == "connected" || slice[2] == "接続済み"),
		Signal: -1,
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
		notfoundError.Message = "Not found an active wifi"
		errFn(notfoundError.Message, notfoundError)
		return
	}

	slice = strings.Split(string(out), ":")
	signal, parseErr := strconv.ParseInt(strings.TrimRight(slice[3], "\n"), 10, 32)
	if parseErr != nil {
		log.Print(parseErr)
		return
	}
	stat.Signal = int32(signal)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stat)
}
