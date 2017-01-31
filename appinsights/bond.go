package appinsights

type Data struct {
	BaseType string      `json:"baseType"`
	BaseData interface{} `json:"baseData"`
}

type Envelope struct {
	Name string            `json:"name"`
	Time string            `json:"time"`
	IKey string            `json:"iKey"`
	Tags map[string]string `json:"tags"`
	Data Data              `json:"data"`
}

type DataPointType int

const (
	Measurement DataPointType = iota
	Aggregation
)

type DataPoint struct {
	Name   string        `json:"name"`
	Kind   DataPointType `json:"kind"`
	Value  float32       `json:"value"`
	Count  int           `json:"count"`
	Min    float32       `json:"min"`
	Max    float32       `json:"max"`
	StdDev float32       `json:"stdDev"`
}

type MetricData struct {
	Ver        int               `json:"ver"`
	Properties map[string]string `json:"properties"`
	Metrics    []*DataPoint      `json:"metrics"`
}

type EventData struct {
	Ver          int                `json:"ver"`
	Properties   map[string]string  `json:"properties"`
	Name         string             `json:"name"`
	Measurements map[string]float32 `json:"measurements"`
}

type SeverityLevel int

const (
	Verbose SeverityLevel = iota
	Information
	Warning
	Error
	Critical
)

type MessageData struct {
	Ver           int               `json:"ver"`
	Properties    map[string]string `json:"properties"`
	Message       string            `json:"message"`
	SeverityLevel SeverityLevel     `json:"severityLevel"`
}

type RequestData struct {
	Ver          int                `json:"ver"`
	Properties   map[string]string  `json:"properties"`
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	StartTime    string             `json:"startTime"` // yyyy-mm-ddThh:mm:ss.fffffff-hh:mm
	Duration     string             `json:"duration"`  // d:hh:mm:ss.fffffff
	ResponseCode string             `json:"responseCode"`
	Success      bool               `json:"success"`
	HttpMethod   string             `json:"httpMethod"`
	Url          string             `json:"url"`
	Measurements map[string]float32 `json:"measurements"`
}

type DependencyKind int

const (
	SQL DependencyKind = iota
	Http
	Other
)

type DependencySourceType int

const (
	Undefined DependencySourceType = iota
	Aic
	Apmc
)

type RemoteDependencyData struct {
	Ver              int                  `json:"ver"`
	Id               string               `json:"id"`
	Name             string               `json:"name"`
	ResultCode       int                  `json:"resultCode"`
	CommandName      string               `json:"commandName"`
	Kind             DataPointType        `json:"kind"`
	Duration         string               `json:"duration"`
	Count            int                  `json:"count"`
	Min              float32              `json:"min"`
	Max              float32              `json:"max"`
	StdDev           float32              `json:"stdDev"`
	Type             string               `json:"type"`
	DependencyKind   DependencyKind       `json:"-"` // omitted according to Java reference
	Success          bool                 `json:"success"`
	Async            bool                 `json:"async"`
	DependencySource DependencySourceType `json:"dependencySource"`
	Properties       map[string]string    `json:"properties"`
}

type ContextTagKeys string

const (
	ApplicationVersion         ContextTagKeys = "ai.application.ver"
	ApplicationBuild                          = "ai.application.build"
	DeviceId                                  = "ai.device.id"
	DeviceIp                                  = "ai.device.ip"
	DeviceLanguage                            = "ai.device.language"
	DeviceLocale                              = "ai.device.locale"
	DeviceModel                               = "ai.device.model"
	DeviceNetwork                             = "ai.device.network"
	DeviceOEMName                             = "ai.device.oemName"
	DeviceOS                                  = "ai.device.os"
	DeviceOSVersion                           = "ai.device.osVersion"
	DeviceRoleInstance                        = "ai.device.roleInstance"
	DeviceRoleName                            = "ai.device.roleName"
	DeviceScreenResolution                    = "ai.device.screenResolution"
	DeviceType                                = "ai.device.type"
	DeviceMachineName                         = "ai.device.machineName"
	LocationIp                                = "ai.location.ip"
	OperationId                               = "ai.operation.id"
	OperationName                             = "ai.operation.name"
	OperationParentId                         = "ai.operation.parentId"
	OperationRootId                           = "ai.operation.rootId"
	OperationSyntheticSource                  = "ai.operation.syntheticSource"
	OperationIsSynthetic                      = "ai.operation.isSynthetic"
	SessionId                                 = "ai.session.id"
	SessionIsFirst                            = "ai.session.isFirst"
	SessionIsNew                              = "ai.session.isNew"
	UserAccountAcquisitionDate                = "ai.user.accountAcquisitionDate"
	UserAccountId                             = "ai.user.accountId"
	UserAgent                                 = "ai.user.userAgent"
	UserId                                    = "ai.user.id"
	UserStoreRegion                           = "ai.user.storeRegion"
	SampleRate                                = "ai.sample.sampleRate"
	InternalSdkVersion                        = "ai.internal.sdkVersion"
	InternalAgentVersion                      = "ai.internal.agentVersion"
)
