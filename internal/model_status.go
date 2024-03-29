/*
 * Router Manager
 *
 * This is a managing network service.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Status struct {

	Device string `json:"device"`

	Connected bool `json:"connected"`

	Signal int32 `json:"signal"`
}
