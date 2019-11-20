// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pulsar

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jasonahills/pulsar-client-go/pulsar/internal"
)

func TestSingleMessageIDNoAckTracker(t *testing.T) {
	eventsCh := make(chan interface{}, 1)
	pc := partitionConsumer{
		queueCh:  make(chan []*message, 1),
		eventsCh: eventsCh,
	}

	headersAndPayload := internal.NewBufferWrapper(rawCompatSingleMessage)
	if err := pc.MessageReceived(nil, headersAndPayload); err != nil {
		t.Fatal(err)
	}

	// ensure the tracker was set on the message id
	messages := <-pc.queueCh
	for _, m := range messages {
		assert.Nil(t, m.ID().(*messageID).tracker)
	}

	// ack the message id
	pc.AckID(messages[0].msgID.(*messageID))

	select {
	case <-eventsCh:
	default:
		t.Error("Expected an ack request to be triggered!")
	}
}

func TestBatchMessageIDNoAckTracker(t *testing.T) {
	eventsCh := make(chan interface{}, 1)
	pc := partitionConsumer{
		queueCh:  make(chan []*message, 1),
		eventsCh: eventsCh,
	}

	headersAndPayload := internal.NewBufferWrapper(rawBatchMessage1)
	if err := pc.MessageReceived(nil, headersAndPayload); err != nil {
		t.Fatal(err)
	}

	// ensure the tracker was set on the message id
	messages := <-pc.queueCh
	for _, m := range messages {
		assert.Nil(t, m.ID().(*messageID).tracker)
	}

	// ack the message id
	pc.AckID(messages[0].msgID.(*messageID))

	select {
	case <-eventsCh:
	default:
		t.Error("Expected an ack request to be triggered!")
	}
}

func TestBatchMessageIDWithAckTracker(t *testing.T) {
	eventsCh := make(chan interface{}, 1)
	pc := partitionConsumer{
		queueCh:  make(chan []*message, 1),
		eventsCh: eventsCh,
	}

	headersAndPayload := internal.NewBufferWrapper(rawBatchMessage10)
	if err := pc.MessageReceived(nil, headersAndPayload); err != nil {
		t.Fatal(err)
	}

	// ensure the tracker was set on the message id
	messages := <-pc.queueCh
	for _, m := range messages {
		assert.NotNil(t, m.ID().(*messageID).tracker)
	}

	// ack all message ids except the last one
	for i := 0; i < 9; i++ {
		pc.AckID(messages[i].msgID.(*messageID))
	}

	select {
	case <-eventsCh:
		t.Error("The message id should not be acked!")
	default:
	}

	// ack last message
	pc.AckID(messages[9].msgID.(*messageID))

	select {
	case <-eventsCh:
	default:
		t.Error("Expected an ack request to be triggered!")
	}
}

// Raw single message in old format
// metadata properties:<key:"a" value:"1" > properties:<key:"b" value:"2" >
// payload = "hello"
var rawCompatSingleMessage = []byte{
	0x0e, 0x01, 0x08, 0x36, 0xb4, 0x66, 0x00, 0x00,
	0x00, 0x31, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x6e,
	0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x2d, 0x37,
	0x34, 0x2d, 0x30, 0x10, 0x00, 0x18, 0xac, 0xef,
	0xe8, 0xa0, 0xe2, 0x2d, 0x22, 0x06, 0x0a, 0x01,
	0x61, 0x12, 0x01, 0x31, 0x22, 0x06, 0x0a, 0x01,
	0x62, 0x12, 0x01, 0x32, 0x48, 0x05, 0x60, 0x05,
	0x82, 0x01, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
}

// Message with batch of 1
// singe message metadata properties:<key:"a" value:"1" > properties:<key:"b" value:"2" >
// payload = "hello"
var rawBatchMessage1 = []byte{
	0x0e, 0x01, 0x1f, 0x80, 0x09, 0x68, 0x00, 0x00,
	0x00, 0x1f, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x6e,
	0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x2d, 0x37,
	0x34, 0x2d, 0x31, 0x10, 0x00, 0x18, 0xdb, 0x80,
	0xf4, 0xa0, 0xe2, 0x2d, 0x58, 0x01, 0x82, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x16, 0x0a, 0x06, 0x0a,
	0x01, 0x61, 0x12, 0x01, 0x31, 0x0a, 0x06, 0x0a,
	0x01, 0x62, 0x12, 0x01, 0x32, 0x18, 0x05, 0x28,
	0x05, 0x40, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
}

var rawBatchMessage10 = []byte{
	0x0e, 0x01, 0x7b, 0x28, 0x8c, 0x08,
	0x00, 0x00, 0x00, 0x1f, 0x0a, 0x0f, 0x73, 0x74,
	0x61, 0x6e, 0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65,
	0x2d, 0x37, 0x34, 0x2d, 0x32, 0x10, 0x00, 0x18,
	0xd0, 0xc2, 0xfa, 0xa0, 0xe2, 0x2d, 0x58, 0x0a,
	0x82, 0x01, 0x00, 0x00, 0x00, 0x00, 0x16, 0x0a,
	0x06, 0x0a, 0x01, 0x61, 0x12, 0x01, 0x31, 0x0a,
	0x06, 0x0a, 0x01, 0x62, 0x12, 0x01, 0x32, 0x18,
	0x05, 0x28, 0x05, 0x40, 0x00, 0x68, 0x65, 0x6c,
	0x6c, 0x6f, 0x00, 0x00, 0x00, 0x16, 0x0a, 0x06,
	0x0a, 0x01, 0x61, 0x12, 0x01, 0x31, 0x0a, 0x06,
	0x0a, 0x01, 0x62, 0x12, 0x01, 0x32, 0x18, 0x05,
	0x28, 0x05, 0x40, 0x01, 0x68, 0x65, 0x6c, 0x6c,
	0x6f, 0x00, 0x00, 0x00, 0x16, 0x0a, 0x06, 0x0a,
	0x01, 0x61, 0x12, 0x01, 0x31, 0x0a, 0x06, 0x0a,
	0x01, 0x62, 0x12, 0x01, 0x32, 0x18, 0x05, 0x28,
	0x05, 0x40, 0x02, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	0x00, 0x00, 0x00, 0x16, 0x0a, 0x06, 0x0a, 0x01,
	0x61, 0x12, 0x01, 0x31, 0x0a, 0x06, 0x0a, 0x01,
	0x62, 0x12, 0x01, 0x32, 0x18, 0x05, 0x28, 0x05,
	0x40, 0x03, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00,
	0x00, 0x00, 0x16, 0x0a, 0x06, 0x0a, 0x01, 0x61,
	0x12, 0x01, 0x31, 0x0a, 0x06, 0x0a, 0x01, 0x62,
	0x12, 0x01, 0x32, 0x18, 0x05, 0x28, 0x05, 0x40,
	0x04, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00,
	0x00, 0x16, 0x0a, 0x06, 0x0a, 0x01, 0x61, 0x12,
	0x01, 0x31, 0x0a, 0x06, 0x0a, 0x01, 0x62, 0x12,
	0x01, 0x32, 0x18, 0x05, 0x28, 0x05, 0x40, 0x05,
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x00,
	0x16, 0x0a, 0x06, 0x0a, 0x01, 0x61, 0x12, 0x01,
	0x31, 0x0a, 0x06, 0x0a, 0x01, 0x62, 0x12, 0x01,
	0x32, 0x18, 0x05, 0x28, 0x05, 0x40, 0x06, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x00, 0x16,
	0x0a, 0x06, 0x0a, 0x01, 0x61, 0x12, 0x01, 0x31,
	0x0a, 0x06, 0x0a, 0x01, 0x62, 0x12, 0x01, 0x32,
	0x18, 0x05, 0x28, 0x05, 0x40, 0x07, 0x68, 0x65,
	0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x00, 0x16, 0x0a,
	0x06, 0x0a, 0x01, 0x61, 0x12, 0x01, 0x31, 0x0a,
	0x06, 0x0a, 0x01, 0x62, 0x12, 0x01, 0x32, 0x18,
	0x05, 0x28, 0x05, 0x40, 0x08, 0x68, 0x65, 0x6c,
	0x6c, 0x6f, 0x00, 0x00, 0x00, 0x16, 0x0a, 0x06,
	0x0a, 0x01, 0x61, 0x12, 0x01, 0x31, 0x0a, 0x06,
	0x0a, 0x01, 0x62, 0x12, 0x01, 0x32, 0x18, 0x05,
	0x28, 0x05, 0x40, 0x09, 0x68, 0x65, 0x6c, 0x6c,
	0x6f,
}
