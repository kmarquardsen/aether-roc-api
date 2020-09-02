// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package southbound

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"time"
)

var log = logging.GetLogger("southbound")

const (
	configAddress = "onos-config:5150"
)

// GNMIProvisioner handles provisioning of device configuration via gNMI interface.
type GNMIProvisioner struct {
	gnmi gnmi.GNMIClient
}

// Init initializes the gNMI provisioner
func (p *GNMIProvisioner) Init(opts ...grpc.DialOption) error {
	optsWithRetry := []grpc.DialOption{
		grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
	}
	optsWithRetry = append(opts, optsWithRetry...)
	gnmiConn, err := grpc.Dial(configAddress, optsWithRetry...)
	if err != nil {
		log.Error("Unable to connect to onos-config", err)
		return err
	}
	p.gnmi = gnmi.NewGNMIClient(gnmiConn)
	return nil
}

// Get passes a gNMI SetRequest to the server which synchronously replies with a GetResponse
func (p *GNMIProvisioner) Get(ctx context.Context, request *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	return p.gnmi.Get(ctx, request)
}