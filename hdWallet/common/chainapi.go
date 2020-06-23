package common

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	rpcclient "github.com/bcbchain/bclib/rpc/lib/client"
	rpctypes "github.com/bcbchain/bclib/rpc/lib/types"
	"github.com/bcbchain/bclib/tendermint/go-amino"
	core_types "github.com/bcbchain/tendermint/rpc/core/types"
)

var ConsumeTimeFlag bool

//网络请求和结果解析-故障队列版
func DoHttpRequestAndParseExBlock(nodeAddrSlice []string, methodName string, params map[string]interface{}, result interface{}) (err error) {

	for {

		RWLock.Lock()
		if len(CorrectUrls) == 0 {
			RWLock.Unlock()
			return errors.New("no available nodes can to connect")
		}

		length := len(CorrectUrls)

		var rnd int
		if length > 1 {
			rnd = rand.Intn(length - 1)
		} else {
			rnd = 0
		}

		url := CorrectUrls[rnd]
		RWLock.Unlock()

		err = CallChainApi(url, methodName, params, result)
		if err == nil {
			break
		} else {
			//RWLock.Lock()
			//FaultCounterMap[url] += 1
			//if FaultCounterMap[url] > 10 {
			//	if rnd == length-1 {
			//		CorrectUrls = append(CorrectUrls[:rnd])
			//	} else {
			//		CorrectUrls = append(CorrectUrls[:rnd], CorrectUrls[rnd+1:]...)
			//	}
			//	length -= 1
			//}
			//RWLock.Unlock()
			//
			//if length <= len(nodeAddrSlice)/3 {
			//	go DealFaultUrls()
			//}
			//
			//if length == 0 {
			splitErr := strings.Split(err.Error(), ":")
			return errors.New(strings.Trim(splitErr[len(splitErr)-1], " "))
			//}
		}
	}

	return
}

func CallChainApi(url string, methodName string, params map[string]interface{}, result interface{}) (err error) {

	if methodName == "broadcast_tx_commit" {
		ConsumeTimeFlag = true
	}
	rpc := NewJSONRPCClientEx(url, "", true)
	_, err = rpc.Call(methodName, params, result)
	return
}

func DealFaultUrls() {
	RWLock.Lock()
	FaultUrls2 := FaultCounterMap
	RWLock.Unlock()
	methodName := "abci_info"
	params := map[string]interface{}{}
	result := new(core_types.ResultABCIInfo)

	for k, _ := range FaultUrls2 {
		err := CallChainApi(k, methodName, params, result)
		if err == nil {
			RWLock.Lock()
			CorrectUrls = append(CorrectUrls, k)
			FaultCounterMap[k] = 0
			RWLock.Unlock()
		}
	}
}

type JSONRPCClient struct {
	address string
	client  *http.Client
	cdc     *amino.Codec
}

func NewJSONRPCClientEx(remote, certFile string, disableKeepAlive bool) *JSONRPCClient {
	var pool *x509.CertPool
	if certFile != "" {
		pool = x509.NewCertPool()
		caCert, err := ioutil.ReadFile(certFile)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		pool.AppendCertsFromPEM(caCert)
	}

	address, client := makeHTTPSClient(remote, pool, disableKeepAlive)

	return &JSONRPCClient{
		address: address,
		client:  client,
		cdc:     rpcclient.CDC,
	}
}

func makeHTTPSClient(remoteAddr string, pool *x509.CertPool, disableKeepAlive bool) (string, *http.Client) {
	//_, dialer := makeHTTPDialer(remoteAddr)

	tr := new(http.Transport)
	tr.DisableKeepAlives = disableKeepAlive
	tr.IdleConnTimeout = time.Second * 120
	if pool != nil {
		tr.TLSClientConfig = &tls.Config{RootCAs: pool}
	} else {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if ConsumeTimeFlag == true {
		return remoteAddr, &http.Client{Transport: tr, Timeout: time.Duration(time.Second * 120)}
	} else {
		return remoteAddr, &http.Client{Transport: tr, Timeout: time.Duration(time.Second * 3)}
	}

}

func (c *JSONRPCClient) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	//request, err := types.MapToRequest("jsonrpc-client", method, params)
	request, err := rpctypes.MapToRequest("jsonrpc-client", method, params)
	if err != nil {
		return nil, err
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("lib client http_client error to json.Marshal(request)")
		return nil, err
	}

	// log.Info(string(requestBytes))
	requestBuf := bytes.NewBuffer(requestBytes)
	// log.Info(Fmt("RPC request to %v (%v): %v", c.remote, method, string(requestBytes)))
	httpResponse, err := c.client.Post(c.address, "text/json", requestBuf)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	// 	log.Info(Fmt("RPC response: %v", string(responseBytes)))
	return unmarshalResponseBytes(c.cdc, responseBytes, result)
}

func unmarshalResponseBytes(cdc *amino.Codec, responseBytes []byte, result interface{}) (interface{}, error) {
	// Read response.  If rpc/core/types is imported, the result will unmarshal
	// into the correct type.
	// log.Notice("response", "response", string(responseBytes))
	var err error
	//response := &types.RPCResponse{}
	response := &rpctypes.RPCResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, errors.New("Error unmarshalling rpc response: " + err.Error())
	}
	if response.Error != nil {
		return nil, errors.New("Response error: " + response.Error.Error())
	}
	// Unmarshal the RawMessage into the result.
	err = cdc.UnmarshalJSON(response.Result, result)
	if err != nil {
		return nil, errors.New("Error unmarshalling rpc response result: " + err.Error())
	}
	return result, nil
}
