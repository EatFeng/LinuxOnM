package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/common"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
	"sort"
	"time"
)

// LoadMonitor
// @Tags Monitor
// @Summary Load monitor data
// @Description This function is responsible for retrieving monitor data based on the provided request parameters.
//
//	It first validates and binds the incoming JSON request body of type dto.MonitorSearch. This validation ensures that the incoming data adheres to the expected format and constraints defined for the MonitorSearch structure.
//	After successful validation and binding, it proceeds to handle the time zone conversion for the provided start and end time in the request. It uses the time zone obtained from the common.LoadTimeZoneByCmd function and applies it to the req.StartTime and req.EndTime fields, ensuring that the time-based queries later in the function are performed in the correct time zone context.
//	Then, depending on the value of the 'Param' field in the request, it queries different database tables to fetch relevant monitor data.
//	If the 'Param' value is "all", "cpu", "memory", or "load", it queries the model.MonitorBase table. It fetches records within the specified time range (between req.StartTime and req.EndTime) from the database using the global.MonitorDB object. For each retrieved record, it constructs a dto.MonitorData object with the 'Param' set to "base", populates its 'Date' and 'Value' fields with the relevant data from the retrieved record, and appends this object to the backdatas slice.
//	When the 'Param' value is "all" or "io", it queries the model.MonitorIO table following a similar process. It retrieves records within the given time range, constructs a dto.MonitorData object with 'Param' as "io", populates its fields with the retrieved data, and adds it to the backdatas slice. In case of any database query errors during this process, it calls the helper.ErrorWithDetail function to send back an error response with appropriate error code and type, along with the detailed error message.
//	For the case where the 'Param' value is "all" or "network", it queries the model.MonitorNetwork table with an additional condition on the 'name' field (name = req.Info) along with the time range check. It follows the same pattern of constructing a dto.MonitorData object with 'Param' set to "network", populating its fields, and appending it to the backdatas slice. Again, if any errors occur during the database query, an error response is sent.
//	Finally, if all the data retrieval operations are completed without errors, it sends back a success response with the collected monitor data in the backdatas slice using the helper.SuccessWithData function.
//
// @Param request body dto.MonitorSearch true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/monitor/search [post]
func (b *BaseApi) LoadMonitor(c *gin.Context) {
	var req dto.MonitorSearch
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	loc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	req.StartTime = req.StartTime.In(loc)
	req.EndTime = req.EndTime.In(loc)

	var backdatas []dto.MonitorData
	if req.Param == "all" || req.Param == "cpu" || req.Param == "memory" || req.Param == "load" {
		var bases []models.MonitorBase
		if err := global.MonitorDB.
			Where("created_at > ? AND created_at < ?", req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "base"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	if req.Param == "all" || req.Param == "io" {
		var bases []models.MonitorIO
		if err := global.MonitorDB.
			Where("created_at > ? AND created_at < ?", req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "io"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	if req.Param == "all" || req.Param == "network" {
		var bases []models.MonitorNetwork
		if err := global.MonitorDB.
			Where("name = ? AND created_at > ? AND created_at < ?", req.Info, req.StartTime, req.EndTime).
			Find(&bases).Error; err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		var itemData dto.MonitorData
		itemData.Param = "network"
		for _, base := range bases {
			itemData.Date = append(itemData.Date, base.CreatedAt)
			itemData.Value = append(itemData.Value, base)
		}
		backdatas = append(backdatas, itemData)
	}
	helper.SuccessWithData(c, backdatas)
}

// CleanMonitor
// @Tags Monitor
// @Summary Clean monitor datas
// @Description Delete all the existing monitor data.
// @Success 200
// @Security ApiKeyAuth
// @Router /host/monitor/clean [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"清空监控数据","formatEN":"clean monitor datas"}
func (b *BaseApi) CleanMonitor(c *gin.Context) {
	if err := global.MonitorDB.Exec("DELETE FROM monitor_bases").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	if err := global.MonitorDB.Exec("DELETE FROM monitor_ios").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	if err := global.MonitorDB.Exec("DELETE FROM monitor_networks").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

// GetNetworkOptions
// @Tags Monitor
// @Summary Get network options.
// @Description Retrieve available network options which include "all" and the names of network interfaces obtained from network statistics.
// @Success 200
// @Security ApiKeyAuth
// @Router /host/monitor/net_options [get]
func (b *BaseApi) GetNetworkOptions(c *gin.Context) {
	netStat, _ := net.IOCounters(true)
	var options []string
	options = append(options, "all")
	for _, net := range netStat {
		options = append(options, net.Name)
	}
	sort.Strings(options)
	helper.SuccessWithData(c, options)
}

// GetIOOptions
// @Tags Monitor
// @Summary Get I/O options.
// @Description Retrieve available I/O options which include "all" and the names of disk devices obtained from disk I/O counters.
// @Success 200
// @Security ApiKeyAuth
// @Router /host/monitor/io_options [get]
func (b *BaseApi) GetIOOptions(c *gin.Context) {
	diskStat, _ := disk.IOCounters()
	var options []string
	options = append(options, "all")
	for _, net := range diskStat {
		options = append(options, net.Name)
	}
	sort.Strings(options)
	helper.SuccessWithData(c, options)
}
