package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type result struct {
	header     http.Header
	statusCode int
	callback   func(header http.Header, statusCode int)
}

// Requestor is a module for making asyncronous http requests
type Requestor struct {
	queue chan result
}

func (r *Requestor) Update(dt float64) {
	for {
		select {
		case res := <-r.queue:
			res.callback(res.header, res.statusCode)
		default:
			return
		}
	}
}

func (r *Requestor) Request(req *http.Request, response interface{}, callback func(header http.Header, statusCode int)) {
	go func() {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error making requesting: ", err)
			r.queue <- result{callback: callback, statusCode: 500}
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading data from request: ", err)
			r.queue <- result{callback: callback, statusCode: 500}
			return
		}

		if response != nil {
			err = json.Unmarshal(data, response)
			if err != nil {
				log.Println("Error Unmarshaling response from request: ", err)
				r.queue <- result{callback: callback, statusCode: 500}
				return
			}
		}

		r.queue <- result{callback: callback, header: resp.Header, statusCode: resp.StatusCode}
	}()
}

func (r *Requestor) GetRequest(url string, response interface{}, callback func(header http.Header, statusCode int)) {
	go func() {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating get Request: ", err)
			r.queue <- result{callback: callback, statusCode: 500}
			return
		}

		r.Request(req, response, callback)
	}()
}

func (r *Requestor) PostRequest(url string, request, response interface{}, callback func(header http.Header, statusCode int)) {
	go func() {
		requestData, err := json.Marshal(request)
		if err != nil {
			log.Println("Error Marshaling post Request: ", err)
			r.queue <- result{callback: callback, statusCode: 500}
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
		if err != nil {
			log.Println("Error creating post Request: ", err)
			r.queue <- result{callback: callback, statusCode: 500}
			return
		}

		r.Request(req, response, callback)
	}()
}

func NewRequestor() *Requestor {
	return &Requestor{
		queue: make(chan result, 256),
	}
}
