// Package classification hypercloud Server API.
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

//-go:generate swagger generate spec -o app/docs/swagger.json
//go:generate go-bindata -pkg asset -prefix app/ -o app/asset/asset.go app/schema/ app/docs/
