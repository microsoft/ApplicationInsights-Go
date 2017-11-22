package appinsights

// NOTE: This file was automatically generated.

// Helper type that provides access to context fields grouped under 'application'.
// This is returned by TelemetryContext.Application()
type ApplicationContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'device'.
// This is returned by TelemetryContext.Device()
type DeviceContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'location'.
// This is returned by TelemetryContext.Location()
type LocationContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'operation'.
// This is returned by TelemetryContext.Operation()
type OperationContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'session'.
// This is returned by TelemetryContext.Session()
type SessionContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'user'.
// This is returned by TelemetryContext.User()
type UserContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'cloud'.
// This is returned by TelemetryContext.Cloud()
type CloudContext struct {
	context *TelemetryContext
}

// Helper type that provides access to context fields grouped under 'internal'.
// This is returned by TelemetryContext.Internal()
type InternalContext struct {
	context *TelemetryContext
}

// Returns a helper to access context fields grouped under 'application'.
func (context *TelemetryContext) Application() *ApplicationContext {
	return &ApplicationContext{context: context}
}

// Returns a helper to access context fields grouped under 'device'.
func (context *TelemetryContext) Device() *DeviceContext {
	return &DeviceContext{context: context}
}

// Returns a helper to access context fields grouped under 'location'.
func (context *TelemetryContext) Location() *LocationContext {
	return &LocationContext{context: context}
}

// Returns a helper to access context fields grouped under 'operation'.
func (context *TelemetryContext) Operation() *OperationContext {
	return &OperationContext{context: context}
}

// Returns a helper to access context fields grouped under 'session'.
func (context *TelemetryContext) Session() *SessionContext {
	return &SessionContext{context: context}
}

// Returns a helper to access context fields grouped under 'user'.
func (context *TelemetryContext) User() *UserContext {
	return &UserContext{context: context}
}

// Returns a helper to access context fields grouped under 'cloud'.
func (context *TelemetryContext) Cloud() *CloudContext {
	return &CloudContext{context: context}
}

// Returns a helper to access context fields grouped under 'internal'.
func (context *TelemetryContext) Internal() *InternalContext {
	return &InternalContext{context: context}
}

// Application version. Information in the application context fields is
// always about the application that is sending the telemetry.
func (context *ApplicationContext) GetVer() string {
	return context.context.getStringTag("ai.application.ver")
}

// Application version. Information in the application context fields is
// always about the application that is sending the telemetry.
func (context *ApplicationContext) SetVer(value string) {
	context.context.setStringTag("ai.application.ver", value)
}

// Unique client device id. Computer name in most cases.
func (context *DeviceContext) GetId() string {
	return context.context.getStringTag("ai.device.id")
}

// Unique client device id. Computer name in most cases.
func (context *DeviceContext) SetId(value string) {
	context.context.setStringTag("ai.device.id", value)
}

// Device locale using <language>-<REGION> pattern, following RFC 5646.
// Example 'en-US'.
func (context *DeviceContext) GetLocale() string {
	return context.context.getStringTag("ai.device.locale")
}

// Device locale using <language>-<REGION> pattern, following RFC 5646.
// Example 'en-US'.
func (context *DeviceContext) SetLocale(value string) {
	context.context.setStringTag("ai.device.locale", value)
}

// Model of the device the end user of the application is using. Used for
// client scenarios. If this field is empty then it is derived from the user
// agent.
func (context *DeviceContext) GetModel() string {
	return context.context.getStringTag("ai.device.model")
}

// Model of the device the end user of the application is using. Used for
// client scenarios. If this field is empty then it is derived from the user
// agent.
func (context *DeviceContext) SetModel(value string) {
	context.context.setStringTag("ai.device.model", value)
}

// Client device OEM name taken from the browser.
func (context *DeviceContext) GetOemName() string {
	return context.context.getStringTag("ai.device.oemName")
}

// Client device OEM name taken from the browser.
func (context *DeviceContext) SetOemName(value string) {
	context.context.setStringTag("ai.device.oemName", value)
}

// Operating system name and version of the device the end user of the
// application is using. If this field is empty then it is derived from the
// user agent. Example 'Windows 10 Pro 10.0.10586.0'
func (context *DeviceContext) GetOsVersion() string {
	return context.context.getStringTag("ai.device.osVersion")
}

// Operating system name and version of the device the end user of the
// application is using. If this field is empty then it is derived from the
// user agent. Example 'Windows 10 Pro 10.0.10586.0'
func (context *DeviceContext) SetOsVersion(value string) {
	context.context.setStringTag("ai.device.osVersion", value)
}

// The type of the device the end user of the application is using. Used
// primarily to distinguish JavaScript telemetry from server side telemetry.
// Examples: 'PC', 'Phone', 'Browser'. 'PC' is the default value.
func (context *DeviceContext) GetType() string {
	return context.context.getStringTag("ai.device.type")
}

// The type of the device the end user of the application is using. Used
// primarily to distinguish JavaScript telemetry from server side telemetry.
// Examples: 'PC', 'Phone', 'Browser'. 'PC' is the default value.
func (context *DeviceContext) SetType(value string) {
	context.context.setStringTag("ai.device.type", value)
}

// The IP address of the client device. IPv4 and IPv6 are supported.
// Information in the location context fields is always about the end user.
// When telemetry is sent from a service, the location context is about the
// user that initiated the operation in the service.
func (context *LocationContext) GetIp() string {
	return context.context.getStringTag("ai.location.ip")
}

// The IP address of the client device. IPv4 and IPv6 are supported.
// Information in the location context fields is always about the end user.
// When telemetry is sent from a service, the location context is about the
// user that initiated the operation in the service.
func (context *LocationContext) SetIp(value string) {
	context.context.setStringTag("ai.location.ip", value)
}

// A unique identifier for the operation instance. The operation.id is created
// by either a request or a page view. All other telemetry sets this to the
// value for the containing request or page view. Operation.id is used for
// finding all the telemetry items for a specific operation instance.
func (context *OperationContext) GetId() string {
	return context.context.getStringTag("ai.operation.id")
}

// A unique identifier for the operation instance. The operation.id is created
// by either a request or a page view. All other telemetry sets this to the
// value for the containing request or page view. Operation.id is used for
// finding all the telemetry items for a specific operation instance.
func (context *OperationContext) SetId(value string) {
	context.context.setStringTag("ai.operation.id", value)
}

// The name (group) of the operation. The operation.name is created by either
// a request or a page view. All other telemetry items set this to the value
// for the containing request or page view. Operation.name is used for finding
// all the telemetry items for a group of operations (i.e. 'GET Home/Index').
func (context *OperationContext) GetName() string {
	return context.context.getStringTag("ai.operation.name")
}

// The name (group) of the operation. The operation.name is created by either
// a request or a page view. All other telemetry items set this to the value
// for the containing request or page view. Operation.name is used for finding
// all the telemetry items for a group of operations (i.e. 'GET Home/Index').
func (context *OperationContext) SetName(value string) {
	context.context.setStringTag("ai.operation.name", value)
}

// The unique identifier of the telemetry item's immediate parent.
func (context *OperationContext) GetParentId() string {
	return context.context.getStringTag("ai.operation.parentId")
}

// The unique identifier of the telemetry item's immediate parent.
func (context *OperationContext) SetParentId(value string) {
	context.context.setStringTag("ai.operation.parentId", value)
}

// Name of synthetic source. Some telemetry from the application may represent
// a synthetic traffic. It may be web crawler indexing the web site, site
// availability tests or traces from diagnostic libraries like Application
// Insights SDK itself.
func (context *OperationContext) GetSyntheticSource() string {
	return context.context.getStringTag("ai.operation.syntheticSource")
}

// Name of synthetic source. Some telemetry from the application may represent
// a synthetic traffic. It may be web crawler indexing the web site, site
// availability tests or traces from diagnostic libraries like Application
// Insights SDK itself.
func (context *OperationContext) SetSyntheticSource(value string) {
	context.context.setStringTag("ai.operation.syntheticSource", value)
}

// The correlation vector is a light weight vector clock which can be used to
// identify and order related events across clients and services.
func (context *OperationContext) GetCorrelationVector() string {
	return context.context.getStringTag("ai.operation.correlationVector")
}

// The correlation vector is a light weight vector clock which can be used to
// identify and order related events across clients and services.
func (context *OperationContext) SetCorrelationVector(value string) {
	context.context.setStringTag("ai.operation.correlationVector", value)
}

// Session ID - the instance of the user's interaction with the app.
// Information in the session context fields is always about the end user.
// When telemetry is sent from a service, the session context is about the
// user that initiated the operation in the service.
func (context *SessionContext) GetId() string {
	return context.context.getStringTag("ai.session.id")
}

// Session ID - the instance of the user's interaction with the app.
// Information in the session context fields is always about the end user.
// When telemetry is sent from a service, the session context is about the
// user that initiated the operation in the service.
func (context *SessionContext) SetId(value string) {
	context.context.setStringTag("ai.session.id", value)
}

// Boolean value indicating whether the session identified by ai.session.id is
// first for the user or not.
func (context *SessionContext) GetIsFirst() string {
	return context.context.getStringTag("ai.session.isFirst")
}

// Boolean value indicating whether the session identified by ai.session.id is
// first for the user or not.
func (context *SessionContext) SetIsFirst(value string) {
	context.context.setStringTag("ai.session.isFirst", value)
}

// In multi-tenant applications this is the account ID or name which the user
// is acting with. Examples may be subscription ID for Azure portal or blog
// name blogging platform.
func (context *UserContext) GetAccountId() string {
	return context.context.getStringTag("ai.user.accountId")
}

// In multi-tenant applications this is the account ID or name which the user
// is acting with. Examples may be subscription ID for Azure portal or blog
// name blogging platform.
func (context *UserContext) SetAccountId(value string) {
	context.context.setStringTag("ai.user.accountId", value)
}

// Anonymous user id. Represents the end user of the application. When
// telemetry is sent from a service, the user context is about the user that
// initiated the operation in the service.
func (context *UserContext) GetId() string {
	return context.context.getStringTag("ai.user.id")
}

// Anonymous user id. Represents the end user of the application. When
// telemetry is sent from a service, the user context is about the user that
// initiated the operation in the service.
func (context *UserContext) SetId(value string) {
	context.context.setStringTag("ai.user.id", value)
}

// Authenticated user id. The opposite of ai.user.id, this represents the user
// with a friendly name. Since it's PII information it is not collected by
// default by most SDKs.
func (context *UserContext) GetAuthUserId() string {
	return context.context.getStringTag("ai.user.authUserId")
}

// Authenticated user id. The opposite of ai.user.id, this represents the user
// with a friendly name. Since it's PII information it is not collected by
// default by most SDKs.
func (context *UserContext) SetAuthUserId(value string) {
	context.context.setStringTag("ai.user.authUserId", value)
}

// Name of the role the application is a part of. Maps directly to the role
// name in azure.
func (context *CloudContext) GetRole() string {
	return context.context.getStringTag("ai.cloud.role")
}

// Name of the role the application is a part of. Maps directly to the role
// name in azure.
func (context *CloudContext) SetRole(value string) {
	context.context.setStringTag("ai.cloud.role", value)
}

// Name of the instance where the application is running. Computer name for
// on-premisis, instance name for Azure.
func (context *CloudContext) GetRoleInstance() string {
	return context.context.getStringTag("ai.cloud.roleInstance")
}

// Name of the instance where the application is running. Computer name for
// on-premisis, instance name for Azure.
func (context *CloudContext) SetRoleInstance(value string) {
	context.context.setStringTag("ai.cloud.roleInstance", value)
}

// SDK version. See
// https://github.com/Microsoft/ApplicationInsights-Home/blob/master/SDK-AUTHORING.md#sdk-version-specification
// for information.
func (context *InternalContext) GetSdkVersion() string {
	return context.context.getStringTag("ai.internal.sdkVersion")
}

// SDK version. See
// https://github.com/Microsoft/ApplicationInsights-Home/blob/master/SDK-AUTHORING.md#sdk-version-specification
// for information.
func (context *InternalContext) SetSdkVersion(value string) {
	context.context.setStringTag("ai.internal.sdkVersion", value)
}

// Agent version. Used to indicate the version of StatusMonitor installed on
// the computer if it is used for data collection.
func (context *InternalContext) GetAgentVersion() string {
	return context.context.getStringTag("ai.internal.agentVersion")
}

// Agent version. Used to indicate the version of StatusMonitor installed on
// the computer if it is used for data collection.
func (context *InternalContext) SetAgentVersion(value string) {
	context.context.setStringTag("ai.internal.agentVersion", value)
}

// This is the node name used for billing purposes. Use it to override the
// standard detection of nodes.
func (context *InternalContext) GetNodeName() string {
	return context.context.getStringTag("ai.internal.nodeName")
}

// This is the node name used for billing purposes. Use it to override the
// standard detection of nodes.
func (context *InternalContext) SetNodeName(value string) {
	context.context.setStringTag("ai.internal.nodeName", value)
}
