package server

import (
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/prometheus/common/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/onosproject/aether-roc-api/pkg/southbound"
)

// Implement the Server Interface for access to gNMI
var log = logging.GetLogger("app_gtwy")

// Server -
type Server struct {
	GnmiClient      southbound.GnmiClient
	AnalyticsClient AnalyticsClient
}

// GetDevices /aether/app-gtwy/enterprises/{enterprise-id}/sites/{site-id}/devices
func (i *Server) GetDevices(ctx echo.Context, enterpriseId string, siteId string) error {

	var response interface{}
	//var err error

	si, err := i.AnalyticsClient.Query("subscriber_info")
	if err != nil {
		return err
	}

	imsiMap := buildIMSIMap(si)
	fmt.Println(imsiMap)

	//gnmiCtx, cancel := utils.NewGnmiContext(ctx)
	//defer cancel()
	//
	//// Response GET OK 200
	//response, err = i.gnmiGetDevices(gnmiCtx, "/aether/app-gtwy/enterprises/{enterprise-id}/sites/{site-id}/devices", enterpriseId, siteId)
	//
	//if err != nil {
	//	return utils.ConvertGrpcError(err)
	//}
	//// It's not enough to check if response==nil - see https://medium.com/@glucn/golang-an-interface-holding-a-nil-value-is-not-nil-bb151f472cc7
	//if reflect.ValueOf(response).Kind() == reflect.Ptr && reflect.ValueOf(response).IsNil() {
	//	return echo.NewHTTPError(http.StatusNoContent)
	//}
	//
	//log.Infof("GetDevices")
	return ctx.JSON(http.StatusOK, response)
}

// GetDevice /aether/app-gtwy/enterprises/{enterprise-id}/sites/{site-id}/devices/{device-id}
func (i *Server) GetDevice(ctx echo.Context, enterpriseId string, siteId string, deviceId string) error {

	var response interface{}
	//var err error

	si, err := i.AnalyticsClient.Query("subscriber_info")
	if err != nil {
		return err
	}

	imsiMap := buildIMSIMap(si)
	fmt.Println(imsiMap)

	//gnmiCtx, cancel := utils.NewGnmiContext(ctx)
	//defer cancel()
	//
	//// Response GET OK 200
	//response, err = i.gnmiGetDevice(gnmiCtx, "/aether/app-gtwy/enterprises/{enterprise-id}/sites/{site-id}/devices/{device-id}", enterpriseId, siteId, deviceId)
	//
	//if err != nil {
	//	return utils.ConvertGrpcError(err)
	//}
	//// It's not enough to check if response==nil - see https://medium.com/@glucn/golang-an-interface-holding-a-nil-value-is-not-nil-bb151f472cc7
	//if reflect.ValueOf(response).Kind() == reflect.Ptr && reflect.ValueOf(response).IsNil() {
	//	return echo.NewHTTPError(http.StatusNoContent)
	//}

	log.Infof("GetDevice")
	return ctx.JSON(http.StatusOK, response)
}

func buildIMSIMap(si model.Value) map[string]*model.Sample {
	v := si.(model.Vector)
	imsiMap := make(map[string]*model.Sample)
	for _, val := range v {
		imsiMap[string(val.Metric["imsi"])] = val
	}
	return imsiMap
}

// gnmiGetDevice returns an instance of App_Gtwy_Device.
//func (i *Server) gnmiGetDevice(ctx context.Context,
//	openApiPath string, target externalRef0.Target, args ...string) (*externalRef0.AppGtwyDevice, error) {
//
//	gnmiGet, err := utils.NewGnmiGetRequest(openApiPath, string(target), args...)
//	if err != nil {
//		return nil, err
//	}
//	log.Infof("gnmiGetRequest %s", gnmiGet.String())
//	gnmiVal, err := utils.GetResponseUpdate(i.GnmiClient.Get(ctx, gnmiGet))
//	if err != nil {
//		return nil, err
//	}
//	if gnmiVal == nil {
//		return nil, nil
//	}
//	gnmiJsonVal, ok := gnmiVal.Value.(*gnmi.TypedValue_JsonVal)
//	if !ok {
//		return nil, fmt.Errorf("unexpected type of reply from server %v", gnmiVal.Value)
//	}
//
//	log.Infof("gNMI Json %s", string(gnmiJsonVal.JsonVal))
//	var gnmiResponse externalRef1.Device
//	if err = externalRef1.Unmarshal(gnmiJsonVal.JsonVal, &gnmiResponse); err != nil {
//		return nil, fmt.Errorf("error unmarshalling gnmiResponse %v", err)
//	}
//	mpd := ModelPluginDevice{
//		device: gnmiResponse,
//	}
//
//	return mpd.toAppGtwyDevice(args...)
//}

// gnmiGetAppGtwyDevices returns an instance of App_Gtwy_Devices.
//func (i *Server) gnmiGetDevices(ctx context.Context,
//	openApiPath string, target externalRef0.Target, args ...string) (*externalRef0.AppGtwyDevices, error) {
//
//	gnmiGet, err := utils.NewGnmiGetRequest(openApiPath, string(target), args...)
//	if err != nil {
//		return nil, err
//	}
//	log.Infof("gnmiGetRequest %s", gnmiGet.String())
//	gnmiVal, err := utils.GetResponseUpdate(i.GnmiClient.Get(ctx, gnmiGet))
//	if err != nil {
//		return nil, err
//	}
//	if gnmiVal == nil {
//		return nil, nil
//	}
//	gnmiJsonVal, ok := gnmiVal.Value.(*gnmi.TypedValue_JsonVal)
//	if !ok {
//		return nil, fmt.Errorf("unexpected type of reply from server %v", gnmiVal.Value)
//	}
//
//	log.Infof("gNMI Json %s", string(gnmiJsonVal.JsonVal))
//	var gnmiResponse externalRef1.Device
//	if err = externalRef1.Unmarshal(gnmiJsonVal.JsonVal, &gnmiResponse); err != nil {
//		return nil, fmt.Errorf("error unmarshalling gnmiResponse %v", err)
//	}
//	mpd := ModelPluginDevice{
//		device: gnmiResponse,
//	}
//
//	return mpd.toAppGtwyDevices(args...)
//}

//type Translator interface {
//	toAppGtwyDevice(args ...string) (*externalRef0.AppGtwyDevice, error)
//	toAppGtwyDevices(args ...string) (*externalRef0.AppGtwyDevices, error)
//}
