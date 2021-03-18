// Copyright (c) 2015-2021, NVIDIA CORPORATION.
// SPDX-License-Identifier: Apache-2.0

package imgrpkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/NVIDIA/proxyfs/bucketstats"
)

func startHTTPServer() (err error) {
	var (
		ipAddrTCPPort string
	)

	ipAddrTCPPort = net.JoinHostPort(globals.config.PrivateIPAddr, strconv.Itoa(int(globals.config.HTTPServerPort)))

	globals.httpServer = &http.Server{
		Addr:    ipAddrTCPPort,
		Handler: &globals,
	}

	globals.httpServerWG.Add(1)

	go func() {
		var (
			err error
		)

		err = globals.httpServer.ListenAndServe()
		if http.ErrServerClosed != err {
			log.Fatalf("httpServer.ListenAndServe() exited unexpectedly: %v", err)
		}

		globals.httpServerWG.Done()
	}()

	err = nil
	return
}

func stopHTTPServer() (err error) {
	err = globals.httpServer.Shutdown(context.TODO())
	if nil == err {
		globals.httpServerWG.Wait()
	}

	return
}

func (dummy *globalsStruct) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		err         error
		requestBody []byte
		requestPath string
	)

	requestPath = strings.TrimRight(request.URL.Path, "/")

	requestBody, err = ioutil.ReadAll(request.Body)
	if nil == err {
		err = request.Body.Close()
		if nil != err {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		_ = request.Body.Close()
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodDelete:
		serveHTTPDelete(responseWriter, request, requestPath)
	case http.MethodGet:
		serveHTTPGet(responseWriter, request, requestPath)
	case http.MethodPut:
		serveHTTPPut(responseWriter, request, requestPath, requestBody)
	default:
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func serveHTTPDelete(responseWriter http.ResponseWriter, request *http.Request, requestPath string) {
	switch {
	case strings.HasPrefix(requestPath, "/volume"):
		serveHTTPDeleteOfVolume(responseWriter, request, requestPath)
	default:
		responseWriter.WriteHeader(http.StatusNotFound)
	}
}

func serveHTTPDeleteOfVolume(responseWriter http.ResponseWriter, request *http.Request, requestPath string) {
	var (
		err       error
		pathSplit []string
	)

	pathSplit = strings.Split(requestPath, "/")

	switch len(pathSplit) {
	case 3:
		err = deleteVolume(pathSplit[2])
		if nil == err {
			responseWriter.WriteHeader(http.StatusNoContent)
		} else {
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	default:
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}

func serveHTTPGet(responseWriter http.ResponseWriter, request *http.Request, requestPath string) {
	switch {
	case "/config" == requestPath:
		serveHTTPGetOfConfig(responseWriter, request)
	case "/stats" == requestPath:
		serveHTTPGetOfStats(responseWriter, request)
	case strings.HasPrefix(requestPath, "/volume"):
		serveHTTPGetOfVolume(responseWriter, request, requestPath)
	default:
		responseWriter.WriteHeader(http.StatusNotFound)
	}
}

func serveHTTPGetOfConfig(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		confMapJSON []byte
		err         error
	)

	confMapJSON, err = json.Marshal(globals.config)
	if nil != err {
		logFatalf("json.Marshal(globals.config) failed: %v", err)
	}

	responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(confMapJSON)))
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	_, err = responseWriter.Write(confMapJSON)
	if nil != err {
		logWarnf("responseWriter.Write(confMapJSON) failed: %v", err)
	}
}

func serveHTTPGetOfStats(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		err           error
		statsAsString string
	)

	statsAsString = bucketstats.SprintStats(bucketstats.StatFormatParsable1, "*", "*")

	responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(statsAsString)))
	responseWriter.Header().Set("Content-Type", "text/plain")
	responseWriter.WriteHeader(http.StatusOK)

	_, err = responseWriter.Write([]byte(statsAsString))
	if nil != err {
		logWarnf("responseWriter.Write([]byte(statsAsString)) failed: %v", err)
	}
}

func serveHTTPGetOfVolume(responseWriter http.ResponseWriter, request *http.Request, requestPath string) {
	var (
		err          error
		jsonToReturn []byte
		pathSplit    []string
	)

	pathSplit = strings.Split(requestPath, "/")

	switch len(pathSplit) {
	case 2:
		jsonToReturn = getVolumeListAsJSON()

		responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonToReturn)))
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)

		_, err = responseWriter.Write(jsonToReturn)
		if nil != err {
			logWarnf("responseWriter.Write(jsonToReturn) failed: %v", err)
		}
	case 3:
		jsonToReturn, err = getVolumeAsJSON(pathSplit[2])
		if nil == err {
			responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonToReturn)))
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusOK)

			_, err = responseWriter.Write(jsonToReturn)
			if nil != err {
				logWarnf("responseWriter.Write(jsonToReturn) failed: %v", err)
			}
		} else {
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	default:
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}

func serveHTTPPut(responseWriter http.ResponseWriter, request *http.Request, requestPath string, requestBody []byte) {
	switch {
	case strings.HasPrefix(requestPath, "/volume"):
		serveHTTPPutOfVolume(responseWriter, request, requestPath, requestBody)
	default:
		responseWriter.WriteHeader(http.StatusNotFound)
	}
}

func serveHTTPPutOfVolume(responseWriter http.ResponseWriter, request *http.Request, requestPath string, requestBody []byte) {
	var (
		err        error
		pathSplit  []string
		storageURL string
	)

	pathSplit = strings.Split(requestPath, "/")

	switch len(pathSplit) {
	case 3:
		storageURL = string(requestBody[:])

		err = putVolume(pathSplit[2], storageURL)
		if nil == err {
			responseWriter.WriteHeader(http.StatusCreated)
		} else {
			responseWriter.WriteHeader(http.StatusConflict)
		}
	default:
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
}