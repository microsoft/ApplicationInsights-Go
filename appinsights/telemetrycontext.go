package appinsights

// NOTE: This file was automatically generated.

import (
	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
	"strconv"
)

type TelemetryContext interface {
	InstrumentationKey() string

	Application() ApplicationContext
	Device() DeviceContext
	Location() LocationContext
	Operation() OperationContext
	Session() SessionContext
	User() UserContext
	Cloud() CloudContext
	Internal() InternalContext
}

type ApplicationContext interface {
	GetVer() string
	SetVer(value string)
}

type DeviceContext interface {
	GetId() string
	SetId(value string)
	GetLocale() string
	SetLocale(value string)
	GetModel() string
	SetModel(value string)
	GetOemName() string
	SetOemName(value string)
	GetOsVersion() string
	SetOsVersion(value string)
	GetType() string
	SetType(value string)
}

type LocationContext interface {
	GetIp() string
	SetIp(value string)
}

type OperationContext interface {
	GetId() string
	SetId(value string)
	GetName() string
	SetName(value string)
	GetParentId() string
	SetParentId(value string)
	GetSyntheticSource() string
	SetSyntheticSource(value string)
	GetCorrelationVector() string
	SetCorrelationVector(value string)
}

type SessionContext interface {
	GetId() string
	SetId(value string)
	GetIsFirst() string
	SetIsFirst(value string)
}

type UserContext interface {
	GetAccountId() string
	SetAccountId(value string)
	GetId() string
	SetId(value string)
	GetAuthUserId() string
	SetAuthUserId(value string)
}

type CloudContext interface {
	GetRole() string
	SetRole(value string)
	GetRoleInstance() string
	SetRoleInstance(value string)
}

type InternalContext interface {
	GetSdkVersion() string
	SetSdkVersion(value string)
	GetAgentVersion() string
	SetAgentVersion(value string)
	GetNodeName() string
	SetNodeName(value string)
}

type telemetryContext struct {
	iKey string
	tags map[string]string
}

type applicationContext struct {
	context *telemetryContext
}

type deviceContext struct {
	context *telemetryContext
}

type locationContext struct {
	context *telemetryContext
}

type operationContext struct {
	context *telemetryContext
}

type sessionContext struct {
	context *telemetryContext
}

type userContext struct {
	context *telemetryContext
}

type cloudContext struct {
	context *telemetryContext
}

type internalContext struct {
	context *telemetryContext
}

func (context *telemetryContext) InstrumentationKey() string {
	return context.iKey
}

func (context *telemetryContext) Application() ApplicationContext {
	return &applicationContext{context: context}
}

func (context *telemetryContext) Device() DeviceContext {
	return &deviceContext{context: context}
}

func (context *telemetryContext) Location() LocationContext {
	return &locationContext{context: context}
}

func (context *telemetryContext) Operation() OperationContext {
	return &operationContext{context: context}
}

func (context *telemetryContext) Session() SessionContext {
	return &sessionContext{context: context}
}

func (context *telemetryContext) User() UserContext {
	return &userContext{context: context}
}

func (context *telemetryContext) Cloud() CloudContext {
	return &cloudContext{context: context}
}

func (context *telemetryContext) Internal() InternalContext {
	return &internalContext{context: context}
}

func NewItemTelemetryContext() TelemetryContext {
	return &telemetryContext{
		tags: make(map[string]string),
	}
}

func NewClientTelemetryContext(ikey string) TelemetryContext {
	return &telemetryContext{
		iKey: ikey,
		tags: make(map[string]string),
	}
}

func (context *telemetryContext) getStringTag(key string) string {
	if result, ok := context.tags[key]; ok {
		return result
	}

	return ""
}

func (context *telemetryContext) setStringTag(key, value string) {
	if value != "" {
		context.tags[key] = value
	} else {
		delete(context.tags, key)
	}
}

func (context *telemetryContext) getBoolTag(key string) bool {
	if result, ok := context.tags[key]; ok {
		if value, err := strconv.ParseBool(result); err == nil {
			return value
		}
	}

	return false
}

func (context *telemetryContext) setBoolTag(key string, value bool) {
	if value {
		context.tags[key] = "true"
	} else {
		delete(context.tags, key)
	}
}

func (context *applicationContext) GetVer() string {
	return context.context.getStringTag(contracts.ApplicationVer)
}

func (context *applicationContext) SetVer(value string) {
	context.context.setStringTag(contracts.ApplicationVer, value)
}

func (context *deviceContext) GetId() string {
	return context.context.getStringTag(contracts.DeviceId)
}

func (context *deviceContext) SetId(value string) {
	context.context.setStringTag(contracts.DeviceId, value)
}

func (context *deviceContext) GetLocale() string {
	return context.context.getStringTag(contracts.DeviceLocale)
}

func (context *deviceContext) SetLocale(value string) {
	context.context.setStringTag(contracts.DeviceLocale, value)
}

func (context *deviceContext) GetModel() string {
	return context.context.getStringTag(contracts.DeviceModel)
}

func (context *deviceContext) SetModel(value string) {
	context.context.setStringTag(contracts.DeviceModel, value)
}

func (context *deviceContext) GetOemName() string {
	return context.context.getStringTag(contracts.DeviceOemName)
}

func (context *deviceContext) SetOemName(value string) {
	context.context.setStringTag(contracts.DeviceOemName, value)
}

func (context *deviceContext) GetOsVersion() string {
	return context.context.getStringTag(contracts.DeviceOsVersion)
}

func (context *deviceContext) SetOsVersion(value string) {
	context.context.setStringTag(contracts.DeviceOsVersion, value)
}

func (context *deviceContext) GetType() string {
	return context.context.getStringTag(contracts.DeviceType)
}

func (context *deviceContext) SetType(value string) {
	context.context.setStringTag(contracts.DeviceType, value)
}

func (context *locationContext) GetIp() string {
	return context.context.getStringTag(contracts.LocationIp)
}

func (context *locationContext) SetIp(value string) {
	context.context.setStringTag(contracts.LocationIp, value)
}

func (context *operationContext) GetId() string {
	return context.context.getStringTag(contracts.OperationId)
}

func (context *operationContext) SetId(value string) {
	context.context.setStringTag(contracts.OperationId, value)
}

func (context *operationContext) GetName() string {
	return context.context.getStringTag(contracts.OperationName)
}

func (context *operationContext) SetName(value string) {
	context.context.setStringTag(contracts.OperationName, value)
}

func (context *operationContext) GetParentId() string {
	return context.context.getStringTag(contracts.OperationParentId)
}

func (context *operationContext) SetParentId(value string) {
	context.context.setStringTag(contracts.OperationParentId, value)
}

func (context *operationContext) GetSyntheticSource() string {
	return context.context.getStringTag(contracts.OperationSyntheticSource)
}

func (context *operationContext) SetSyntheticSource(value string) {
	context.context.setStringTag(contracts.OperationSyntheticSource, value)
}

func (context *operationContext) GetCorrelationVector() string {
	return context.context.getStringTag(contracts.OperationCorrelationVector)
}

func (context *operationContext) SetCorrelationVector(value string) {
	context.context.setStringTag(contracts.OperationCorrelationVector, value)
}

func (context *sessionContext) GetId() string {
	return context.context.getStringTag(contracts.SessionId)
}

func (context *sessionContext) SetId(value string) {
	context.context.setStringTag(contracts.SessionId, value)
}

func (context *sessionContext) GetIsFirst() string {
	return context.context.getStringTag(contracts.SessionIsFirst)
}

func (context *sessionContext) SetIsFirst(value string) {
	context.context.setStringTag(contracts.SessionIsFirst, value)
}

func (context *userContext) GetAccountId() string {
	return context.context.getStringTag(contracts.UserAccountId)
}

func (context *userContext) SetAccountId(value string) {
	context.context.setStringTag(contracts.UserAccountId, value)
}

func (context *userContext) GetId() string {
	return context.context.getStringTag(contracts.UserId)
}

func (context *userContext) SetId(value string) {
	context.context.setStringTag(contracts.UserId, value)
}

func (context *userContext) GetAuthUserId() string {
	return context.context.getStringTag(contracts.UserAuthUserId)
}

func (context *userContext) SetAuthUserId(value string) {
	context.context.setStringTag(contracts.UserAuthUserId, value)
}

func (context *cloudContext) GetRole() string {
	return context.context.getStringTag(contracts.CloudRole)
}

func (context *cloudContext) SetRole(value string) {
	context.context.setStringTag(contracts.CloudRole, value)
}

func (context *cloudContext) GetRoleInstance() string {
	return context.context.getStringTag(contracts.CloudRoleInstance)
}

func (context *cloudContext) SetRoleInstance(value string) {
	context.context.setStringTag(contracts.CloudRoleInstance, value)
}

func (context *internalContext) GetSdkVersion() string {
	return context.context.getStringTag(contracts.InternalSdkVersion)
}

func (context *internalContext) SetSdkVersion(value string) {
	context.context.setStringTag(contracts.InternalSdkVersion, value)
}

func (context *internalContext) GetAgentVersion() string {
	return context.context.getStringTag(contracts.InternalAgentVersion)
}

func (context *internalContext) SetAgentVersion(value string) {
	context.context.setStringTag(contracts.InternalAgentVersion, value)
}

func (context *internalContext) GetNodeName() string {
	return context.context.getStringTag(contracts.InternalNodeName)
}

func (context *internalContext) SetNodeName(value string) {
	context.context.setStringTag(contracts.InternalNodeName, value)
}
