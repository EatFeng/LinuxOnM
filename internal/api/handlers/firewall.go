package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// LoadFirewallBaseInfo
// @Tags Firewall
// @Summary Load firewall base info
// @Description 获取防火墙基础信息
// @Success 200 {object} dto.FirewallBaseInfo
// @Security ApiKeyAuth
// @Router /host/firewall/base [get]
func (b *BaseApi) LoadFirewallBaseInfo(c *gin.Context) {
	data, err := firewallService.LoadBaseInfo()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, data)
}

// OperateFirewall
// @Tags Firewall
// @Summary Page firewall status
// @Description 修改防火墙状态
// @Accept json
// @Param request body dto.FirewallOperation true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /host/firewall/operate [post]
// @x-panel-log {"bodyKeys":["operation"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"[operation] 防火墙","formatEN":"[operation] firewall"}
func (b *BaseApi) OperateFirewall(c *gin.Context) {
	var req dto.FirewallOperation
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.OperateFirewall(req.Operation); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

// SearchFirewallRule
// @Tags Firewall
// @Summary Page firewall rules
// @Description 获取防火墙规则列表分页
// @Accept json
// @Param request body dto.RuleSearch true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /host/firewall/search [post]
func (b *BaseApi) SearchFirewallRule(c *gin.Context) {
	var req dto.RuleSearch
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := firewallService.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// OperatePortRule
// @Tags Firewall
// @Summary Create group
// @Description 创建防火墙端口规则
// @Accept json
// @Param request body dto.PortRuleOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/port [post]
// @x-panel-log {"bodyKeys":["port","strategy"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"添加端口规则 [strategy] [port]","formatEN":"create port rules [strategy][port]"}
func (b *BaseApi) OperatePortRule(c *gin.Context) {
	var req dto.PortRuleOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.OperatePortRule(req, true); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// OperateForwardRule
// @Tags Firewall
// @Summary Create group
// @Description 更新防火墙端口转发规则
// @Accept json
// @Param request body dto.ForwardRuleOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/forward [post]
// @x-panel-log {"bodyKeys":["source_port"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新端口转发规则 [source_port]","formatEN":"update port forward rules [source_port]"}
func (b *BaseApi) OperateForwardRule(c *gin.Context) {
	var req dto.ForwardRuleOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.OperateForwardRule(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// OperateIPRule
// @Tags Firewall
// @Summary Create group
// @Description 创建防火墙 IP 规则
// @Accept json
// @Param request body dto.AddrRuleOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/ip [post]
// @x-panel-log {"bodyKeys":["strategy","address"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"添加 ip 规则 [strategy] [address]","formatEN":"create address rules [strategy][address]"}
func (b *BaseApi) OperateIPRule(c *gin.Context) {
	var req dto.AddrRuleOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.OperateAddressRule(req, true); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// UpdateFirewallDescription
// @Tags Firewall
// @Summary Update rule description
// @Description 更新防火墙描述
// @Accept json
// @Param request body dto.UpdateFirewallDescription true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/update/description [post]
func (b *BaseApi) UpdateFirewallDescription(c *gin.Context) {
	var req dto.UpdateFirewallDescription
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.UpdateDescription(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// UpdatePortRule
// @Tags Firewall
// @Summary Create group
// @Description 更新端口防火墙规则
// @Accept json
// @Param request body dto.PortRuleUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/update/port [post]
func (b *BaseApi) UpdatePortRule(c *gin.Context) {
	var req dto.PortRuleUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.UpdatePortRule(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// UpdateAddrRule
// @Tags Firewall
// @Summary Create group
// @Description 更新 ip 防火墙规则
// @Accept json
// @Param request body dto.AddrRuleUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/update/addr [post]
func (b *BaseApi) UpdateAddrRule(c *gin.Context) {
	var req dto.AddrRuleUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.UpdateAddrRule(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// BatchOperateRule
// @Tags Firewall
// @Summary Create group
// @Description 批量删除防火墙规则
// @Accept json
// @Param request body dto.BatchRuleOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/firewall/batch [post]
func (b *BaseApi) BatchOperateRule(c *gin.Context) {
	var req dto.BatchRuleOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := firewallService.BatchOperateRule(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
