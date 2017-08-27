package appinsights

// NOTE: This file was automatically generated.

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

type applicationContext struct {
	context *TelemetryContext
}

type deviceContext struct {
	context *TelemetryContext
}

type locationContext struct {
	context *TelemetryContext
}

type operationContext struct {
	context *TelemetryContext
}

type sessionContext struct {
	context *TelemetryContext
}

type userContext struct {
	context *TelemetryContext
}

type cloudContext struct {
	context *TelemetryContext
}

type internalContext struct {
	context *TelemetryContext
}

func (context *TelemetryContext) Application() ApplicationContext {
	return &applicationContext{context: context}
}

func (context *TelemetryContext) Device() DeviceContext {
	return &deviceContext{context: context}
}

func (context *TelemetryContext) Location() LocationContext {
	return &locationContext{context: context}
}

func (context *TelemetryContext) Operation() OperationContext {
	return &operationContext{context: context}
}

func (context *TelemetryContext) Session() SessionContext {
	return &sessionContext{context: context}
}

func (context *TelemetryContext) User() UserContext {
	return &userContext{context: context}
}

func (context *TelemetryContext) Cloud() CloudContext {
	return &cloudContext{context: context}
}

func (context *TelemetryContext) Internal() InternalContext {
	return &internalContext{context: context}
}

func (context *applicationContext) GetVer() string {
	return context.context.getStringTag("ai.application.ver")
}

func (context *applicationContext) SetVer(value string) {
	context.context.setStringTag("ai.application.ver", value)
}

func (context *deviceContext) GetId() string {
	return context.context.getStringTag("ai.device.id")
}

func (context *deviceContext) SetId(value string) {
	context.context.setStringTag("ai.device.id", value)
}

func (context *deviceContext) GetLocale() string {
	return context.context.getStringTag("ai.device.locale")
}

func (context *deviceContext) SetLocale(value string) {
	context.context.setStringTag("ai.device.locale", value)
}

func (context *deviceContext) GetModel() string {
	return context.context.getStringTag("ai.device.model")
}

func (context *deviceContext) SetModel(value string) {
	context.context.setStringTag("ai.device.model", value)
}

func (context *deviceContext) GetOemName() string {
	return context.context.getStringTag("ai.device.oemName")
}

func (context *deviceContext) SetOemName(value string) {
	context.context.setStringTag("ai.device.oemName", value)
}

func (context *deviceContext) GetOsVersion() string {
	return context.context.getStringTag("ai.device.osVersion")
}

func (context *deviceContext) SetOsVersion(value string) {
	context.context.setStringTag("ai.device.osVersion", value)
}

func (context *deviceContext) GetType() string {
	return context.context.getStringTag("ai.device.type")
}

func (context *deviceContext) SetType(value string) {
	context.context.setStringTag("ai.device.type", value)
}

func (context *locationContext) GetIp() string {
	return context.context.getStringTag("ai.location.ip")
}

func (context *locationContext) SetIp(value string) {
	context.context.setStringTag("ai.location.ip", value)
}

func (context *operationContext) GetId() string {
	return context.context.getStringTag("ai.operation.id")
}

func (context *operationContext) SetId(value string) {
	context.context.setStringTag("ai.operation.id", value)
}

func (context *operationContext) GetName() string {
	return context.context.getStringTag("ai.operation.name")
}

func (context *operationContext) SetName(value string) {
	context.context.setStringTag("ai.operation.name", value)
}

func (context *operationContext) GetParentId() string {
	return context.context.getStringTag("ai.operation.parentId")
}

func (context *operationContext) SetParentId(value string) {
	context.context.setStringTag("ai.operation.parentId", value)
}

func (context *operationContext) GetSyntheticSource() string {
	return context.context.getStringTag("ai.operation.syntheticSource")
}

func (context *operationContext) SetSyntheticSource(value string) {
	context.context.setStringTag("ai.operation.syntheticSource", value)
}

func (context *operationContext) GetCorrelationVector() string {
	return context.context.getStringTag("ai.operation.correlationVector")
}

func (context *operationContext) SetCorrelationVector(value string) {
	context.context.setStringTag("ai.operation.correlationVector", value)
}

func (context *sessionContext) GetId() string {
	return context.context.getStringTag("ai.session.id")
}

func (context *sessionContext) SetId(value string) {
	context.context.setStringTag("ai.session.id", value)
}

func (context *sessionContext) GetIsFirst() string {
	return context.context.getStringTag("ai.session.isFirst")
}

func (context *sessionContext) SetIsFirst(value string) {
	context.context.setStringTag("ai.session.isFirst", value)
}

func (context *userContext) GetAccountId() string {
	return context.context.getStringTag("ai.user.accountId")
}

func (context *userContext) SetAccountId(value string) {
	context.context.setStringTag("ai.user.accountId", value)
}

func (context *userContext) GetId() string {
	return context.context.getStringTag("ai.user.id")
}

func (context *userContext) SetId(value string) {
	context.context.setStringTag("ai.user.id", value)
}

func (context *userContext) GetAuthUserId() string {
	return context.context.getStringTag("ai.user.authUserId")
}

func (context *userContext) SetAuthUserId(value string) {
	context.context.setStringTag("ai.user.authUserId", value)
}

func (context *cloudContext) GetRole() string {
	return context.context.getStringTag("ai.cloud.role")
}

func (context *cloudContext) SetRole(value string) {
	context.context.setStringTag("ai.cloud.role", value)
}

func (context *cloudContext) GetRoleInstance() string {
	return context.context.getStringTag("ai.cloud.roleInstance")
}

func (context *cloudContext) SetRoleInstance(value string) {
	context.context.setStringTag("ai.cloud.roleInstance", value)
}

func (context *internalContext) GetSdkVersion() string {
	return context.context.getStringTag("ai.internal.sdkVersion")
}

func (context *internalContext) SetSdkVersion(value string) {
	context.context.setStringTag("ai.internal.sdkVersion", value)
}

func (context *internalContext) GetAgentVersion() string {
	return context.context.getStringTag("ai.internal.agentVersion")
}

func (context *internalContext) SetAgentVersion(value string) {
	context.context.setStringTag("ai.internal.agentVersion", value)
}

func (context *internalContext) GetNodeName() string {
	return context.context.getStringTag("ai.internal.nodeName")
}

func (context *internalContext) SetNodeName(value string) {
	context.context.setStringTag("ai.internal.nodeName", value)
}
