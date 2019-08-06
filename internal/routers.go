/*
 * Router Manager
 *
 * This is a managing network service.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

	{
		"SysPoweroff",
		strings.ToUpper("Post"),
		"/v1/sys/poweroff",
		SysPoweroff,
	},

	{
		"SysTimedateGet",
		strings.ToUpper("Get"),
		"/v1/sys/timedate",
		SysTimedateGet,
	},

	{
		"SysTimedateSync",
		strings.ToUpper("Post"),
		"/v1/sys/timedate",
		SysTimedateSync,
	},

	{
		"WifiReconnect",
		strings.ToUpper("Post"),
		"/v1/wifi",
		WifiReconnect,
	},

	{
		"WifiStatus",
		strings.ToUpper("Get"),
		"/v1/wifi",
		WifiStatus,
	},
}