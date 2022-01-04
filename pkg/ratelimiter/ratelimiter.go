/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ratelimiter

import (
	"time"

	"go.uber.org/ratelimit"

	"github.com/k8stopologyawareschedwg/resource-topology-exporter/pkg/notification"
)

type RateLimitedEventSource struct {
	es    notification.EventSource
	inCh  <-chan notification.Event
	outCh chan notification.Event
	rt    ratelimit.Limiter

	doneChan chan struct{}
	stopChan chan struct{}
}

func NewRateLimitedEventSource(es notification.EventSource, maxEventsPerTimeUnit uint64, timeUnit time.Duration) (*RateLimitedEventSource, error) {

	rles := RateLimitedEventSource{
		es:       es,
		inCh:     es.Events(),
		outCh:    make(chan notification.Event),
		doneChan: make(chan struct{}, 1),
		stopChan: make(chan struct{}),
	}

	options := ratelimit.Per(timeUnit)
	rles.rt = ratelimit.New(int(maxEventsPerTimeUnit), options)

	return &rles, nil
}

func (rles *RateLimitedEventSource) Events() <-chan notification.Event {
	return rles.outCh
}

func (rles *RateLimitedEventSource) Run() {
	go rles.run()
	rles.es.Run()
}

func (rles *RateLimitedEventSource) Stop() {
	rles.es.Stop()
	rles.stop()
}

func (rles *RateLimitedEventSource) Wait() {
	rles.wait()
	rles.es.Wait()
}

func (rles *RateLimitedEventSource) Close() {
	//nothing to do here, just call decorated Close
	rles.es.Close()
}

func (rles *RateLimitedEventSource) run() {
	keepgoing := true
	for keepgoing {
		select {
		case event := <-rles.inCh:
			rles.rt.Take()
			rles.outCh <- event
		case <-rles.stopChan:
			keepgoing = false
		}
	}
	rles.doneChan <- struct{}{}
}

// Wait stops the caller until the EventSource is exhausted
func (rles *RateLimitedEventSource) wait() {
	<-rles.doneChan
}

func (rles *RateLimitedEventSource) stop() {
	rles.stopChan <- struct{}{}
}
