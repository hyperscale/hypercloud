// Package classification HyperPaaS Server API.
//
//     Schemes: https
//     Host: localhost
//     BasePath: /v1
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Axel Etcheverry <axel@etcheverry.biz>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

//go:generate swagger generate spec -o docs/swagger.json
//go:generate go-bindata -pkg assets -o assets/assets.go schema/ docs/
