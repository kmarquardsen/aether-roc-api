package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	externalRef0 "github.com/onosproject/aether-roc-api/pkg/aether_2_0_0/types"
	"github.com/onosproject/aether-roc-api/pkg/app_gtwy/types"
	"github.com/onosproject/aether-roc-api/pkg/southbound"
	"github.com/onosproject/aether-roc-api/pkg/utils"
	externalRef1 "github.com/onosproject/config-models/modelplugin/aether-2.0.0/aether_2_0_0"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/prometheus/common/model"
	"net/http"
)

// Implement the Server Interface for access to gNMI
var log = logging.GetLogger("app_gtwy")

// Server -
type Server struct {
	GnmiClient      southbound.GnmiClient
	AnalyticsClient AnalyticsClient
}

// GetDevices /appgtwy/v1/{target}/enterprises/{enterprise-id}/sites/{site-id}/devices
func (i *Server) GetDevices(ctx echo.Context, target externalRef0.Target, enterpriseId string, siteId string) error {

	var response interface{}
	//var err error

	gnmiCtx, cancel := utils.NewGnmiContext(ctx)
	defer cancel()

	site, err := i.gnmiGetEnterprisesEnterpriseSite(gnmiCtx, "/aether/v2.0.0/{target}/enterprises/enterprise/{enterprise-id}/site/{site-id}", "connectivity-service-v2", enterpriseId, siteId)
	//response, err = i.gnmiGetEnterprisesEnterpriseSite(gnmiCtx, "/aether/v2.0.0/{target}/enterprises/enterprise/{enterprise-id}/site/{site-id}", "connectivity-service-v2", enterpriseId, siteId)
	if err != nil {
		return utils.ConvertGrpcError(err)
	}
	// It's not enough to check if response==nil - see https://medium.com/@glucn/golang-an-interface-holding-a-nil-value-is-not-nil-bb151f472cc7
	//if reflect.ValueOf(response).Kind() == reflect.Ptr && reflect.ValueOf(response).IsNil() {
	//	return ctx.NoContent(http.StatusNotFound)
	//}

	//site := response.(externalRef0.EnterprisesEnterpriseSite)
	fmt.Println(site.Device)
	fmt.Println(site.DeviceGroup)
	fmt.Println(site.SimCard)

	si, err := i.AnalyticsClient.Query("subscriber_info")
	if err != nil {
		return err
	}

	//simMap := mapSim(*site)
	//deviceGroupMap := mapDeviceToDeviceGroups(*site)
	imsiMap := mapIMSI(si)
	//fmt.Println(simMap)
	//fmt.Println(deviceGroupMap)
	fmt.Println(imsiMap)

	log.Infof("GetEnterprisesEnterpriseSite")
	return ctx.JSON(http.StatusOK, response)
}

// GetDevice /appgtwy/v1/{target}/enterprises/{enterprise-id}/sites/{site-id}/devices/{device-id}
func (i *Server) GetDevice(ctx echo.Context, target externalRef0.Target, enterpriseId string, siteId string, deviceId string) error {

	var response interface{}
	//var err error

	gnmiCtx, cancel := utils.NewGnmiContext(ctx)
	defer cancel()

	site, err := i.gnmiGetEnterprisesEnterpriseSite(gnmiCtx, "/aether/v2.0.0/{target}/enterprises/enterprise/{enterprise-id}/site/{site-id}", "connectivity-service-v2", enterpriseId, siteId)
	//response, err = i.gnmiGetEnterprisesEnterpriseSite(gnmiCtx, "/aether/v2.0.0/{target}/enterprises/enterprise/{enterprise-id}/site/{site-id}", "connectivity-service-v2", enterpriseId, siteId)
	if err != nil {
		return utils.ConvertGrpcError(err)
	}
	// It's not enough to check if response==nil - see https://medium.com/@glucn/golang-an-interface-holding-a-nil-value-is-not-nil-bb151f472cc7
	//if reflect.ValueOf(response).Kind() == reflect.Ptr && reflect.ValueOf(response).IsNil() {
	//	return ctx.NoContent(http.StatusNotFound)
	//}

	//site := response.(externalRef0.EnterprisesEnterpriseSite)
	fmt.Println(site.Device)
	fmt.Println(site.DeviceGroup)
	fmt.Println(site.SimCard)

	si, err := i.AnalyticsClient.Query("subscriber_info")
	if err != nil {
		return err
	}

	simMap := mapSim(site)
	deviceGroupMap := mapDeviceToDeviceGroups(site)
	imsiMap := mapIMSI(si)
	fmt.Println(simMap)
	fmt.Println(deviceGroupMap)
	fmt.Println(imsiMap)

	//var device types.AppGtwyDevice
	for _, d := range *site.Device {
		if d.DeviceId != deviceId {
			continue
		}
		s := simMap[*d.SimCard]

		response = types.AppGtwyDevice{
			DeviceId:    &d.DeviceId,
			DisplayName: d.DisplayName,
			Description: d.Description,
			Imei:        d.Imei,
			SimIccid:    s.Iccid,
			//DeviceGroups: deviceGroupMap[deviceId],
		}

		//imsi := strconv.Itoa(s.Imsi)
		////spoof := "00"
		//if _, ok := imsiMap[imsi]; ok {
		//	ip := imsiMap[imsi].Metric["ip"]
		//	device.Ip = string(ip)
		//	attached := imsiMap[imsi].Value.String()
		//	device.Attached, _ = attached
		//}
	}

	log.Infof("GetDevice")
	return ctx.JSON(http.StatusOK, response)
}

// gnmiGetEnterprisesEnterpriseSite returns an instance of Enterprises_Enterprise_Site.
func (i *Server) gnmiGetEnterprisesEnterpriseSite(ctx context.Context,
	openApiPath string, target externalRef0.Target, args ...string) (*externalRef0.EnterprisesEnterpriseSite, error) {

	gnmiGet, err := utils.NewGnmiGetRequest(openApiPath, string(target), args...)
	if err != nil {
		return nil, err
	}
	log.Infof("gnmiGetRequest %s", gnmiGet.String())
	gnmiVal, err := utils.GetResponseUpdate(i.GnmiClient.Get(ctx, gnmiGet))
	if err != nil {
		return nil, err
	}
	if gnmiVal == nil {
		return nil, nil
	}
	gnmiJsonVal, ok := gnmiVal.Value.(*gnmi.TypedValue_JsonVal)
	if !ok {
		return nil, fmt.Errorf("unexpected type of reply from server %v", gnmiVal.Value)
	}

	log.Infof("gNMI Json %s", string(gnmiJsonVal.JsonVal))
	var gnmiResponse externalRef1.Device
	if err = externalRef1.Unmarshal(gnmiJsonVal.JsonVal, &gnmiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling gnmiResponse %v", err)
	}
	mpd := ModelPluginDevice{
		device: gnmiResponse,
	}

	return mpd.toEnterprisesEnterpriseSite(args...)
}

// ModelPluginDevice - a wrapper for the model plugin
type ModelPluginDevice struct {
	device externalRef1.Device
}

// toEnterprisesEnterpriseSite converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSite(params ...string) (*externalRef0.EnterprisesEnterpriseSite, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSite)

	// Property: description string
	//encoding gNMI attribute to OAPI
	reflectDescription, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDescription", params...)
	if err != nil {
		return nil, err
	}
	if reflectDescription != nil {
		attrDescription := reflectDescription.Interface().(string)
		resource.Description = &attrDescription
	}

	// Property: device []EnterprisesEnterpriseSiteDevice
	// Handle []Object
	devices := make([]externalRef0.EnterprisesEnterpriseSiteDevice, 0)
	reflectEnterprisesEnterpriseSiteDevice, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDevice", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnterprisesEnterpriseSiteDevice != nil {
		for _, key := range reflectEnterprisesEnterpriseSiteDevice.MapKeys() {
			v := reflectEnterprisesEnterpriseSiteDevice.MapIndex(key).Interface()
			// Pass down all top level properties as we don't know which one(s) is key
			attribs, err := utils.ExtractGnmiListKeyMap(v)
			if err != nil {
				return nil, err
			}
			childParams := make([]string, len(params))
			copy(childParams, params)
			for _, attribVal := range attribs {
				childParams = append(childParams, fmt.Sprintf("%v", attribVal))
			}
			device, err := d.toEnterprisesEnterpriseSiteDevice(childParams...)
			if err != nil {
				return nil, err
			}
			devices = append(devices, *device)
		}
	}
	resource.Device = &devices

	// Property: device-group []EnterprisesEnterpriseSiteDeviceGroup
	// Handle []Object
	deviceGroups := make([]externalRef0.EnterprisesEnterpriseSiteDeviceGroup, 0)
	reflectEnterprisesEnterpriseSiteDeviceGroup, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroup", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnterprisesEnterpriseSiteDeviceGroup != nil {
		for _, key := range reflectEnterprisesEnterpriseSiteDeviceGroup.MapKeys() {
			v := reflectEnterprisesEnterpriseSiteDeviceGroup.MapIndex(key).Interface()
			// Pass down all top level properties as we don't know which one(s) is key
			attribs, err := utils.ExtractGnmiListKeyMap(v)
			if err != nil {
				return nil, err
			}
			childParams := make([]string, len(params))
			copy(childParams, params)
			for _, attribVal := range attribs {
				childParams = append(childParams, fmt.Sprintf("%v", attribVal))
			}
			deviceGroup, err := d.toEnterprisesEnterpriseSiteDeviceGroup(childParams...)
			if err != nil {
				return nil, err
			}
			deviceGroups = append(deviceGroups, *deviceGroup)
		}
	}
	resource.DeviceGroup = &deviceGroups

	// Property: display-name string
	//encoding gNMI attribute to OAPI
	reflectDisplayName, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDisplayName", params...)
	if err != nil {
		return nil, err
	}
	if reflectDisplayName != nil {
		attrDisplayName := reflectDisplayName.Interface().(string)
		resource.DisplayName = &attrDisplayName
	}

	// Property: imsi-definition EnterprisesEnterpriseSiteImsiDefinition
	//Handle object
	attrImsiDefinition, err := d.toEnterprisesEnterpriseSiteImsiDefinition(params...)
	if err != nil {
		return nil, err
	}
	resource.ImsiDefinition = attrImsiDefinition

	// Property: ip-domain []EnterprisesEnterpriseSiteIpDomain
	// Handle []Object
	//ipDomains := make([]types.EnterprisesEnterpriseSiteIpDomain, 0)
	//reflectEnterprisesEnterpriseSiteIpDomain, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteIpDomain", params...)
	//if err != nil {
	//	return nil, err
	//}
	//if reflectEnterprisesEnterpriseSiteIpDomain != nil {
	//	for _, key := range reflectEnterprisesEnterpriseSiteIpDomain.MapKeys() {
	//		v := reflectEnterprisesEnterpriseSiteIpDomain.MapIndex(key).Interface()
	//		// Pass down all top level properties as we don't know which one(s) is key
	//		attribs, err := utils.ExtractGnmiListKeyMap(v)
	//		if err != nil {
	//			return nil, err
	//		}
	//		childParams := make([]string, len(params))
	//		copy(childParams, params)
	//		for _, attribVal := range attribs {
	//			childParams = append(childParams, fmt.Sprintf("%v", attribVal))
	//		}
	//		ipDomain, err := d.toEnterprisesEnterpriseSiteIpDomain(childParams...)
	//		if err != nil {
	//			return nil, err
	//		}
	//		ipDomains = append(ipDomains, *ipDomain)
	//	}
	//}
	//resource.IpDomain = &ipDomains

	// Property: monitoring EnterprisesEnterpriseSiteMonitoring
	//Handle object
	//attrMonitoring, err := d.toEnterprisesEnterpriseSiteMonitoring(params...)
	//if err != nil {
	//	return nil, err
	//}
	//resource.Monitoring = attrMonitoring

	// Property: sim-card []EnterprisesEnterpriseSiteSimCard
	// Handle []Object
	simCards := make([]externalRef0.EnterprisesEnterpriseSiteSimCard, 0)
	reflectEnterprisesEnterpriseSiteSimCard, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCard", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnterprisesEnterpriseSiteSimCard != nil {
		for _, key := range reflectEnterprisesEnterpriseSiteSimCard.MapKeys() {
			v := reflectEnterprisesEnterpriseSiteSimCard.MapIndex(key).Interface()
			// Pass down all top level properties as we don't know which one(s) is key
			attribs, err := utils.ExtractGnmiListKeyMap(v)
			if err != nil {
				return nil, err
			}
			childParams := make([]string, len(params))
			copy(childParams, params)
			for _, attribVal := range attribs {
				childParams = append(childParams, fmt.Sprintf("%v", attribVal))
			}
			simCard, err := d.toEnterprisesEnterpriseSiteSimCard(childParams...)
			if err != nil {
				return nil, err
			}
			simCards = append(simCards, *simCard)
		}
	}
	resource.SimCard = &simCards

	// Property: site-id string
	//encoding gNMI attribute to OAPI
	reflectSiteId, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSiteId", params...)
	if err != nil {
		return nil, err
	}
	if reflectSiteId != nil {
		attrSiteId := reflectSiteId.Interface().(string)
		resource.SiteId = attrSiteId
	}

	// Property: slice []EnterprisesEnterpriseSiteSlice
	// Handle []Object
	//slices := make([]types.EnterprisesEnterpriseSiteSlice, 0)
	//reflectEnterprisesEnterpriseSiteSlice, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSlice", params...)
	//if err != nil {
	//	return nil, err
	//}
	//if reflectEnterprisesEnterpriseSiteSlice != nil {
	//	for _, key := range reflectEnterprisesEnterpriseSiteSlice.MapKeys() {
	//		v := reflectEnterprisesEnterpriseSiteSlice.MapIndex(key).Interface()
	//		// Pass down all top level properties as we don't know which one(s) is key
	//		attribs, err := utils.ExtractGnmiListKeyMap(v)
	//		if err != nil {
	//			return nil, err
	//		}
	//		childParams := make([]string, len(params))
	//		copy(childParams, params)
	//		for _, attribVal := range attribs {
	//			childParams = append(childParams, fmt.Sprintf("%v", attribVal))
	//		}
	//		slice, err := d.toEnterprisesEnterpriseSiteSlice(childParams...)
	//		if err != nil {
	//			return nil, err
	//		}
	//		slices = append(slices, *slice)
	//	}
	//}
	//resource.Slice = &slices

	// Property: small-cell []EnterprisesEnterpriseSiteSmallCell
	// Handle []Object
	//smallCells := make([]types.EnterprisesEnterpriseSiteSmallCell, 0)
	//reflectEnterprisesEnterpriseSiteSmallCell, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSmallCell", params...)
	//if err != nil {
	//	return nil, err
	//}
	//if reflectEnterprisesEnterpriseSiteSmallCell != nil {
	//	for _, key := range reflectEnterprisesEnterpriseSiteSmallCell.MapKeys() {
	//		v := reflectEnterprisesEnterpriseSiteSmallCell.MapIndex(key).Interface()
	//		// Pass down all top level properties as we don't know which one(s) is key
	//		attribs, err := utils.ExtractGnmiListKeyMap(v)
	//		if err != nil {
	//			return nil, err
	//		}
	//		childParams := make([]string, len(params))
	//		copy(childParams, params)
	//		for _, attribVal := range attribs {
	//			childParams = append(childParams, fmt.Sprintf("%v", attribVal))
	//		}
	//		smallCell, err := d.toEnterprisesEnterpriseSiteSmallCell(childParams...)
	//		if err != nil {
	//			return nil, err
	//		}
	//		smallCells = append(smallCells, *smallCell)
	//	}
	//}
	//resource.SmallCell = &smallCells

	// Property: upf []EnterprisesEnterpriseSiteUpf
	// Handle []Object
	//upfs := make([]types.EnterprisesEnterpriseSiteUpf, 0)
	//reflectEnterprisesEnterpriseSiteUpf, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteUpf", params...)
	//if err != nil {
	//	return nil, err
	//}
	//if reflectEnterprisesEnterpriseSiteUpf != nil {
	//	for _, key := range reflectEnterprisesEnterpriseSiteUpf.MapKeys() {
	//		v := reflectEnterprisesEnterpriseSiteUpf.MapIndex(key).Interface()
	//		// Pass down all top level properties as we don't know which one(s) is key
	//		attribs, err := utils.ExtractGnmiListKeyMap(v)
	//		if err != nil {
	//			return nil, err
	//		}
	//		childParams := make([]string, len(params))
	//		copy(childParams, params)
	//		for _, attribVal := range attribs {
	//			childParams = append(childParams, fmt.Sprintf("%v", attribVal))
	//		}
	//		upf, err := d.toEnterprisesEnterpriseSiteUpf(childParams...)
	//		if err != nil {
	//			return nil, err
	//		}
	//		upfs = append(upfs, *upf)
	//	}
	//}
	//resource.Upf = &upfs

	return resource, nil
}

// toEnterprisesEnterpriseSiteSimCard converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteSimCard(params ...string) (*externalRef0.EnterprisesEnterpriseSiteSimCard, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteSimCard)

	// Property: description string
	//encoding gNMI attribute to OAPI
	reflectDescription, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCardDescription", params...)
	if err != nil {
		return nil, err
	}
	if reflectDescription != nil {
		attrDescription := reflectDescription.Interface().(string)
		resource.Description = &attrDescription
	}

	// Property: display-name string
	//encoding gNMI attribute to OAPI
	reflectDisplayName, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCardDisplayName", params...)
	if err != nil {
		return nil, err
	}
	if reflectDisplayName != nil {
		attrDisplayName := reflectDisplayName.Interface().(string)
		resource.DisplayName = &attrDisplayName
	}

	// Property: iccid string
	//encoding gNMI attribute to OAPI
	reflectIccid, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCardIccid", params...)
	if err != nil {
		return nil, err
	}
	if reflectIccid != nil {
		attrIccid := reflectIccid.Interface().(string)
		resource.Iccid = &attrIccid
	}

	// Property: imsi int64
	//encoding gNMI attribute to OAPI
	reflectImsi, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCardImsi", params...)
	if err != nil {
		return nil, err
	}
	if reflectImsi != nil {
		//OpenAPI does not have unsigned numbers.
		if resource.Imsi, err = utils.ToInt64Ptr(reflectImsi); err != nil {
			return nil, err
		}
	}

	// Property: sim-id string
	//encoding gNMI attribute to OAPI
	reflectSimId, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteSimCardSimId", params...)
	if err != nil {
		return nil, err
	}
	if reflectSimId != nil {
		attrSimId := reflectSimId.Interface().(string)
		resource.SimId = attrSimId
	}

	return resource, nil
}

// toEnterprisesEnterpriseSiteDevice converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteDevice(params ...string) (*externalRef0.EnterprisesEnterpriseSiteDevice, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteDevice)

	// Property: description string
	//encoding gNMI attribute to OAPI
	reflectDescription, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceDescription", params...)
	if err != nil {
		return nil, err
	}
	if reflectDescription != nil {
		attrDescription := reflectDescription.Interface().(string)
		resource.Description = &attrDescription
	}

	// Property: device-id string
	//encoding gNMI attribute to OAPI
	reflectDeviceId, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceDeviceId", params...)
	if err != nil {
		return nil, err
	}
	if reflectDeviceId != nil {
		attrDeviceId := reflectDeviceId.Interface().(string)
		resource.DeviceId = attrDeviceId
	}

	// Property: display-name string
	//encoding gNMI attribute to OAPI
	reflectDisplayName, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceDisplayName", params...)
	if err != nil {
		return nil, err
	}
	if reflectDisplayName != nil {
		attrDisplayName := reflectDisplayName.Interface().(string)
		resource.DisplayName = &attrDisplayName
	}

	// Property: imei string
	//encoding gNMI attribute to OAPI
	reflectImei, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceImei", params...)
	if err != nil {
		return nil, err
	}
	if reflectImei != nil {
		attrImei := reflectImei.Interface().(string)
		resource.Imei = &attrImei
	}

	// Property: sim-card string
	//encoding gNMI attribute to OAPI
	reflectSimCard, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceSimCard", params...)
	if err != nil {
		return nil, err
	}
	if reflectSimCard != nil {
		attrSimCard := reflectSimCard.Interface().(string)
		resource.SimCard = &attrSimCard
	}

	return resource, nil
}

// toEnterprisesEnterpriseSiteDeviceGroup converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteDeviceGroup(params ...string) (*externalRef0.EnterprisesEnterpriseSiteDeviceGroup, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteDeviceGroup)

	// Property: description string
	//encoding gNMI attribute to OAPI
	reflectDescription, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDescription", params...)
	if err != nil {
		return nil, err
	}
	if reflectDescription != nil {
		attrDescription := reflectDescription.Interface().(string)
		resource.Description = &attrDescription
	}

	// Property: device []EnterprisesEnterpriseSiteDeviceGroupDevice
	// Handle []Object
	devices := make([]externalRef0.EnterprisesEnterpriseSiteDeviceGroupDevice, 0)
	reflectEnterprisesEnterpriseSiteDeviceGroupDevice, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDevice", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnterprisesEnterpriseSiteDeviceGroupDevice != nil {
		for _, key := range reflectEnterprisesEnterpriseSiteDeviceGroupDevice.MapKeys() {
			v := reflectEnterprisesEnterpriseSiteDeviceGroupDevice.MapIndex(key).Interface()
			// Pass down all top level properties as we don't know which one(s) is key
			attribs, err := utils.ExtractGnmiListKeyMap(v)
			if err != nil {
				return nil, err
			}
			childParams := make([]string, len(params))
			copy(childParams, params)
			for _, attribVal := range attribs {
				childParams = append(childParams, fmt.Sprintf("%v", attribVal))
			}
			device, err := d.toEnterprisesEnterpriseSiteDeviceGroupDevice(childParams...)
			if err != nil {
				return nil, err
			}
			devices = append(devices, *device)
		}
	}
	resource.Device = &devices

	// Property: device-group-id string
	//encoding gNMI attribute to OAPI
	reflectDeviceGroupId, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDeviceGroupId", params...)
	if err != nil {
		return nil, err
	}
	if reflectDeviceGroupId != nil {
		attrDeviceGroupId := reflectDeviceGroupId.Interface().(string)
		resource.DeviceGroupId = attrDeviceGroupId
	}

	// Property: display-name string
	//encoding gNMI attribute to OAPI
	reflectDisplayName, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDisplayName", params...)
	if err != nil {
		return nil, err
	}
	if reflectDisplayName != nil {
		attrDisplayName := reflectDisplayName.Interface().(string)
		resource.DisplayName = &attrDisplayName
	}

	// Property: ip-domain string
	//encoding gNMI attribute to OAPI
	reflectIpDomain, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupIpDomain", params...)
	if err != nil {
		return nil, err
	}
	if reflectIpDomain != nil {
		attrIpDomain := reflectIpDomain.Interface().(string)
		resource.IpDomain = &attrIpDomain
	}

	// Property: mbr EnterprisesEnterpriseSiteDeviceGroupMbr
	//Handle object
	attrMbr, err := d.toEnterprisesEnterpriseSiteDeviceGroupMbr(params...)
	if err != nil {
		return nil, err
	}
	resource.Mbr = attrMbr

	// Property: traffic-class string
	//encoding gNMI attribute to OAPI
	reflectTrafficClass, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupTrafficClass", params...)
	if err != nil {
		return nil, err
	}
	if reflectTrafficClass != nil {
		attrTrafficClass := reflectTrafficClass.Interface().(string)
		resource.TrafficClass = attrTrafficClass
	}

	return resource, nil
}

// toEnterprisesEnterpriseSiteImsiDefinition converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteImsiDefinition(params ...string) (*externalRef0.EnterprisesEnterpriseSiteImsiDefinition, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteImsiDefinition)

	// Property: enterprise int32
	//encoding gNMI attribute to OAPI
	reflectEnterprise, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteImsiDefinitionEnterprise", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnterprise != nil {
		//OpenAPI does not have unsigned numbers.
		if resource.Enterprise, err = utils.ToInt32(reflectEnterprise); err != nil {
			return nil, err
		}
	}

	// Property: format string
	//encoding gNMI attribute to OAPI
	reflectFormat, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteImsiDefinitionFormat", params...)
	if err != nil {
		return nil, err
	}
	if reflectFormat != nil {
		attrFormat := reflectFormat.Interface().(string)
		resource.Format = attrFormat
	}

	// Property: mcc string
	//encoding gNMI attribute to OAPI
	reflectMcc, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteImsiDefinitionMcc", params...)
	if err != nil {
		return nil, err
	}
	if reflectMcc != nil {
		attrMcc := reflectMcc.Interface().(string)
		resource.Mcc = attrMcc
	}

	// Property: mnc string
	//encoding gNMI attribute to OAPI
	reflectMnc, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteImsiDefinitionMnc", params...)
	if err != nil {
		return nil, err
	}
	if reflectMnc != nil {
		attrMnc := reflectMnc.Interface().(string)
		resource.Mnc = attrMnc
	}

	return resource, nil
}

// toEnterprisesEnterpriseSiteDeviceGroupDevice converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteDeviceGroupDevice(params ...string) (*externalRef0.EnterprisesEnterpriseSiteDeviceGroupDevice, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteDeviceGroupDevice)

	// Property: device-id string
	//encoding gNMI attribute to OAPI
	reflectDeviceId, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDeviceDeviceId", params...)
	if err != nil {
		return nil, err
	}
	if reflectDeviceId != nil {
		attrDeviceId := reflectDeviceId.Interface().(string)
		resource.DeviceId = attrDeviceId
	}

	// Property: enable bool
	//encoding gNMI attribute to OAPI
	reflectEnable, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupDeviceEnable", params...)
	if err != nil {
		return nil, err
	}
	if reflectEnable != nil {
		boolEnable := reflectEnable.Interface().(bool)
		resource.Enable = &boolEnable
	}

	return resource, nil
}

// toEnterprisesEnterpriseSiteDeviceGroupMbr converts gNMI to OAPI.
func (d *ModelPluginDevice) toEnterprisesEnterpriseSiteDeviceGroupMbr(params ...string) (*externalRef0.EnterprisesEnterpriseSiteDeviceGroupMbr, error) {
	resource := new(externalRef0.EnterprisesEnterpriseSiteDeviceGroupMbr)

	// Property: downlink int64
	//encoding gNMI attribute to OAPI
	reflectDownlink, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupMbrDownlink", params...)
	if err != nil {
		return nil, err
	}
	if reflectDownlink != nil {
		//OpenAPI does not have unsigned numbers.
		if resource.Downlink, err = utils.ToInt64(reflectDownlink); err != nil {
			return nil, err
		}
	}

	// Property: uplink int64
	//encoding gNMI attribute to OAPI
	reflectUplink, err := utils.FindModelPluginObject(d.device, "EnterprisesEnterpriseSiteDeviceGroupMbrUplink", params...)
	if err != nil {
		return nil, err
	}
	if reflectUplink != nil {
		//OpenAPI does not have unsigned numbers.
		if resource.Uplink, err = utils.ToInt64(reflectUplink); err != nil {
			return nil, err
		}
	}

	return resource, nil
}

func mapSim(site *externalRef0.EnterprisesEnterpriseSite) map[string]*externalRef0.EnterprisesEnterpriseSiteSimCard {
	simMap := make(map[string]*externalRef0.EnterprisesEnterpriseSiteSimCard)
	for _, s := range *site.SimCard {
		simMap[s.SimId] = &externalRef0.EnterprisesEnterpriseSiteSimCard{
			SimId:       s.SimId,
			Description: s.Description,
			DisplayName: s.DisplayName,
			Iccid:       s.Iccid,
			Imsi:        s.Imsi,
		}
	}
	return simMap
}

func mapDeviceToDeviceGroups(site *externalRef0.EnterprisesEnterpriseSite) map[string]string {
	deviceGroupMap := make(map[string]string)

	for _, dg := range *site.DeviceGroup {
		if len(*dg.Device) > 0 {
			for _, d := range *dg.Device {
				deviceID := d.DeviceId
				//deviceID := fmt.Sprintf("%v", d.(map[string]interface{})["device-id"])
				deviceGroupMap[deviceID] = dg.DeviceGroupId
			}
		}
	}

	return deviceGroupMap
}

func mapIMSI(si model.Value) map[string]*model.Sample {
	v := si.(model.Vector)
	imsiMap := make(map[string]*model.Sample)
	for _, val := range v {
		imsiMap[string(val.Metric["imsi"])] = val
	}
	return imsiMap
}

//type Translator interface {
//	toAppGtwyDevice(args ...string) (*externalRef0.AppGtwyDevice, error)
//	toAppGtwyDevices(args ...string) (*externalRef0.AppGtwyDevices, error)
//}
