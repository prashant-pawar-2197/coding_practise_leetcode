package standByCB

import (
	config "commonconfig"
	"errors"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	bulkratetrie "bulkratetrie"
	commconsts "commconsts"
	mlog "mwp-appcommon/mavcimclient"

	"github.com/couchbase/gocb/v2"
	comPb "protobuf.com/commondef"
)

var (
	Cluster            *gocb.Cluster
	CollectionObj      *gocb.Collection
	BillingBucket      *gocb.Collection
	connflag           bool
	transID            string
	CBDownAlarmRaised  bool
	gCDbConnectionList map[string]*CouchConnectionStruct
	gCommonCbCfgData   CommonCbCfgStruct
)

var (
	PkgInfoSid         string
	SubInfoSid         string
	RatingGroupKey     string
	SpInfoSid          string
	SpListKey          string
	ProductInfoSid     string
	SubAliasSid        string
	SubBalSid          string
	BillAccSid         string
	DunAccSid          string
	BillCycleBucketSid string
	BillPlanSid        string
	UnbillAcctSid      string
	UnbillAcctDataSid  string
	SessionContextSid  string
	RenewalSid         string
	AlertSid           string
	//RenewalBucketKey   string
	BillCycleSid           string
	SpNumSid               string
	ThresholdPolicySid     string
	AlertPolicySid         string
	RateProfileSid         string
	DiscountProfileSid     string
	TaxProfileSid          string
	FUPProfileSid          string
	ConnectionProfileSid   string
	NetworkUsageProfileSid string
	SubLifeCycleSid        string
	CosInfoSid             string
	SubCounterSid          string
	SpendingLimitPolicySid string
	SlcpAliasSid           string
	SlcpSid                string
	ZoneProfileInfoSid     string
	MccMncInfoSid          string
	QuotaProfileSid        string
	RPTablesSid            string
	RPTableSid             string
	CalendarProfileSid     string
	SubInfoHistorySid      string
	SubAliasHistorySid     string
	SubBalHistorySid       string
)

const (
	SEPARATOR                          = "::"
	SEPARATOR_END_BILLING_DOC          = "_"
	RATING_GROUP_KEY_SUFFIX            = "Rating-Groups"
	SpListKey_SUFFIX                   = "MapData"
	ALARM_TOPIC                        = "EVENT"
	SERVICE_TYPE_PATH                  = ".serviceType"
	ACTIVE_RESERV                      = "activeReservations"
	RENEWAL_BUCKET_KEY                 = "MaxBucketCount"
	ALERT_BUCKET_KEY                   = "MaxBucketCount"
	SUB_LIFE_CYCLE_MASTER              = "Master"
	SUB_LIFE_CYCLE_LIST                = "SubLifeCycleList"
	SUB_LIFE_CYCLE_PATH                = "lifeCycles"
	SUB_LIFE_CYCLE_URL_PATH            = "lifecycle"
	COS_INFO_PATH                      = "cosDefinitions"
	COS_INFO_SUFFIX                    = "COSInformation"
	PROFILES_PATH                      = "profiles"
	SP_PATH                            = "serviceProviders"
	ALERT_POLICY_PATH                  = "AlertPolicies"
	THRESHOLD_POLICY_PATH              = "ThresholdPolicies"
	SPENDING_LIMIT_POLICY_LIST         = "SpendLimitPolicies"
	SPENDING_LIMIT_POLICY_PATH         = "policies"
	SLP_URL_PATH                       = "spendingLimitPolicies"
	COLON                              = ":"
	COUNTER_LIST                       = "SubscriberCounters"
	SUB_SERVICE_MAPPING_PATH           = "subServiceMapping"
	SUBSCRIBER_COUNTERS                = "subscriberCounters"
	SUB_SERVICE_COUNTER_PATH           = "serviceLevelCounter"
	SUB_PACKAGE_COUNTER_PATH           = "packageLevelCounter"
	SUB_SUBSCRIBER_COUNTER_PATH        = "subscriberLevelCounter"
	SUB_WALLET_COUNTER_PATH            = "walletLevelCounter"
	PKG_SERVICE_COUNTER_PATH           = "serviceCounters"
	PKG_PACKAGE_COUNTER_PATH           = "packageCounters"
	PKG_WALLET_COUNTER_PATH            = "walletCounters"
	MCCMNC_PATH                        = "mccmnc"
	REACTIVE_TRIGGER_BATCH_RENEWAL_KEY = "ReactiveTriggerBatchRenewal"
)

func initializeKTAB() {
	PkgInfoSid = config.CommonStaticConf.Config.Ktab["package"] + SEPARATOR    //"PKG::"     //"MAV-CCS-PKG-Information::"
	SubInfoSid = config.CommonStaticConf.Config.Ktab["subscriber"] + SEPARATOR // "Line::"    //"MAV-CCS-Subscriber-Information::"
	RatingGroupKey = config.CommonStaticConf.Config.Ktab["ratingGroups"] + SEPARATOR + RATING_GROUP_KEY_SUFFIX
	SpInfoSid = config.CommonStaticConf.Config.Ktab["serviceProvider"] + SEPARATOR
	SpListKey = config.CommonStaticConf.Config.Ktab["serviceProviderMap"] + SEPARATOR + SpListKey_SUFFIX
	ProductInfoSid = config.CommonStaticConf.Config.Ktab["product"] + SEPARATOR
	SubAliasSid = config.CommonStaticConf.Config.Ktab["subscriberAlias"] + SEPARATOR
	SubBalSid = config.CommonStaticConf.Config.Ktab["balance"] + SEPARATOR
	BillAccSid = config.CommonStaticConf.Config.Ktab["billingAccount"] + SEPARATOR
	DunAccSid = config.CommonStaticConf.Config.Ktab["dunningAccount"] + SEPARATOR
	BillCycleBucketSid = config.CommonStaticConf.Config.Ktab["billCycleBucket"] + SEPARATOR
	BillPlanSid = config.CommonStaticConf.Config.Ktab["billPlan"] + SEPARATOR
	UnbillAcctSid = config.CommonStaticConf.Config.Ktab["unbillAccount"] + SEPARATOR
	UnbillAcctDataSid = config.CommonStaticConf.Config.Ktab["unbillAccountData"] + SEPARATOR
	SessionContextSid = config.CommonStaticConf.Config.Ktab["sessionContext"] + SEPARATOR
	RenewalSid = config.CommonStaticConf.Config.Ktab["renewal"] + SEPARATOR
	AlertSid = config.CommonStaticConf.Config.Ktab["alert"] + SEPARATOR
	//RenewalBucketKey = config.CommonStaticConf.Config.Ktab[""] + SEPARATOR
	BillCycleSid = config.CommonStaticConf.Config.Ktab["billCycle"]
	SpNumSid = config.CommonStaticConf.Config.Ktab["specialNumber"] + SEPARATOR
	ThresholdPolicySid = config.CommonStaticConf.Config.Ktab["thresholdPolicy"] + SEPARATOR
	AlertPolicySid = config.CommonStaticConf.Config.Ktab["alertPolicy"] + SEPARATOR
	RateProfileSid = config.CommonStaticConf.Config.Ktab["rateProfile"] + SEPARATOR
	DiscountProfileSid = config.CommonStaticConf.Config.Ktab["discountProfile"] + SEPARATOR
	FUPProfileSid = config.CommonStaticConf.Config.Ktab["fupProfile"] + SEPARATOR
	ConnectionProfileSid = config.CommonStaticConf.Config.Ktab["connectionProfile"] + SEPARATOR
	NetworkUsageProfileSid = config.CommonStaticConf.Config.Ktab["networkUsageProfile"] + SEPARATOR
	TaxProfileSid = config.CommonStaticConf.Config.Ktab["taxProfile"] + SEPARATOR
	SubLifeCycleSid = config.CommonStaticConf.Config.Ktab["subscriberLifeCyle"] + SEPARATOR
	CosInfoSid = config.CommonStaticConf.Config.Ktab["classOfService"] + SEPARATOR
	SubCounterSid = config.CommonStaticConf.Config.Ktab["subCounters"] + SEPARATOR
	SpendingLimitPolicySid = config.CommonStaticConf.Config.Ktab["spendingLimitPolicy"] + SEPARATOR
	SlcpAliasSid = config.CommonStaticConf.Config.Ktab["slcpAlias"] + SEPARATOR
	SlcpSid = config.CommonStaticConf.Config.Ktab["slcp"] + SEPARATOR
	ZoneProfileInfoSid = config.CommonStaticConf.Config.Ktab["zoneProfile"] + SEPARATOR
	MccMncInfoSid = config.CommonStaticConf.Config.Ktab["mccmnc"] + SEPARATOR
	QuotaProfileSid = config.CommonStaticConf.Config.Ktab["quotaProfileSid"] + SEPARATOR
	RPTablesSid = config.CommonStaticConf.Config.Ktab["rpTableMaster"] + SEPARATOR
	RPTableSid = config.CommonStaticConf.Config.Ktab["rpTable"] + SEPARATOR
	CalendarProfileSid = config.CommonStaticConf.Config.Ktab["CalendarProfile"] + SEPARATOR
	SubInfoHistorySid = config.CommonStaticConf.Config.Ktab["subscriberHistory"] + SEPARATOR
	SubAliasHistorySid = config.CommonStaticConf.Config.Ktab["subscriberAliasHistory"] + SEPARATOR
	SubBalHistorySid = config.CommonStaticConf.Config.Ktab["balanceHistory"] + SEPARATOR

}

/*
//constants for couchbase document SID prefixes
const (
	PKG_INFO_SID          = "PKG::"
	SubInfoSid          = "Line::"
	SP_INFO_SID           = "SP::"      //"MAV-CCS-SP-Information::"
	PRODUCT_INFO_SID      = "Product::" //"MAV-CCS-Product-Information::"
	SUB_INFO_SID          = "Line::"    //"MAV-CCS-Subscriber-Information::"
	SUB_ALIAS_SID         = "Line-Alias::"
	SUB_BAL_SID           = "Line-Bal::"             //"MAV-CCS-Subscriber-Balance-Information::"
	BILL_ACC_SID          = "3001_bsBillAcc::"       //"BILLING ACCOUNT KEY"
	DUN_ACC_SID           = "1001_bsCollections::"   //"DUNNING ACCOUNT KEY"
	BILL_CYLCE_BUCKET_SID = "3011_bsBillScheduler::" //"BILL CYCLE ACCOUNT KEY"
	BILL_PLAN_SID         = "3002_bsBillPlan::"      //"BILL PLAN KEY"
	UNBILL_ACC_SID        = "2002_bsUBL::"           //"UNBILL ACC KEY"
	UNBILL_ACCDATA_SID    = "3006_bsUAD::"
	//RATE_PROFILE_SID     = "MAV-CCS-Service-Tariff-Information::RateProfiles"
	//SVC_USAGE_SID        = "MAV-CCS-ServiceUsage-Information::"
	SESSION_CONTEXT_SID = "Session::" //"MAV-CCS-Charging-Session-Information::"
	//DISCOUNT_PROFILE_KEY = "MAV-CCS-Discount-Information::DiscountProfiles"
	RENEWAL_SID          = "Renewal::"
	ALERT_SID            = "Alert::"
	BILL_BILLCYCLE_SID   = "3001_bsBillCycle"
	SP_NUM_SID           = "SpNum-Info::" //Special Number doc key
	THRESHOLD_POLICY_SID = "TP::"         //Threshold Policy doc key
	ALERT_POLICY_SID     = "AP::"         //Alert Policy doc key
	RP_SID               = "RP::"
	DP_SID               = "DP::"
	SubLifeCycleSid   = "SLC::"
	COS_INFO_SID         = "COSInfo::"
)
*/

//constants for supi to normalized subscriber identity
const (
	IMSI_PREFIX = "imsi-"
	NAI_PREFIX  = "nai-"
)

const (
	VOICE          = "VOICE"
	SMS            = "SMS"
	MMS            = "MMS"
	DATA           = "DATA"
	GENERAL        = "GENERAL"
	COMMON_PRODUCT = "ALL"
)

const (
	FUP          = "FUP"
	Connection   = "Connection"
	NetworkUsage = "NetworkUsage"
)

const (
	ADD    = "ADD"
	DELETE = "DELETE"
)

const (
	SUPI_NO_OP = iota
	SUPI_IMSI
	SUPI_NAI
)

const (
	BTP_INSTANCE string = "1"
)

// const - COSID (comes in specialNumber provisioning) constant for default value i.e "ALL"
const COS_ID_ALL string = "ALL"

const (
	LIFETIME_VALIDITY = "LifeTimeValidity"
	LIMITED_VALIDITY  = "LimitedValidity"
)

const (
	ROAMING_ZONE       = "ROAMUSAGEZONE"
	INTERNATIONAL_ZONE = "INTLUSAGEZONE"
)

const (
	PKGLEVELRENCOUNT = "PKGLEVELRENCOUNT"
)

const (
	COUNT_SPEC  = "countSpec"
	EXISTS_SPEC = "existsSpec"
	GET_SPEC    = "getSpec"
)

type ServiceProviderList struct {
	KTab             string            `json:"KTAB"`
	ServiceProviders map[string]string `json:"serviceProviders"`
}

type SubscriberTypeStates struct {
	Pre  []string `json:"pre,omitempty"`
	Post []string `json:"post,omitempty"`
	Hyd  []string `json:"hybrid,omitempty"`
}

type GeneralInfo struct {
	BrandId              string `json:"brandId" validate:"required"`
	NotificationCategory string `json:"notificationCategory,omitempty"`
}
type TableInfo struct {
	TableID      string `json:"tableId,omitempty"`
	Description  string `json:"description,omitempty"`
	Version      int32  `json:"version,omitempty"`
	NumberPrefix string `json:"numberPrefix,omitempty"`
}
type KafkaMessage struct {
	Header      HeaderType      `json:"header"`
	Information InformationType `json:"information"`
}

type HeaderType struct {
	EventType     string `json:"eventType"`
	EventTime     string `json:"eventTime"`
	OutputType    string `json:"outputType"`
	TransactionId string `json:"transactionid"`
}

type InformationType struct {
	CouchBase CouchBaseType `json:"couchbase"`
}

type CouchBaseType struct {
	BucketName   string `json:"bucketName"`
	DocumentName string `json:"documentName"`
}

type KafkaEvent struct {
	ServiceProviderId string           `json:"serviceProviderId"`
	EventType         string           `json:"eventType"`
	Subscribers       []SubscriberInfo `json:"output"`
}

type SubscriberHistory struct {
	KTab        string                  `json:"KTAB"`
	HistoryInfo []SubInfoHistoryDetails `json:"historyInfo,omitempty"`
}

type SubInfoHistoryDetails struct {
	ActivationDate   string      `json:"activationDate,omitempty"`
	DeactivationDate string      `json:"deactivationDate,omitempty"`
	LineDoc          *Subscriber `json:"lineInfo,omitempty"`
}

type SubscriberInfo struct {
	Imsi                 string `json:"imsi"`
	Msisdn               string `json:"msisdn,omitempty"`
	State                string `json:"state,omitempty"`
	LastStatusChangeDate string `json:"lastStatusChangeDate,omitempty"`
	DeactivationDate     string `json:"deactivationDate,omitempty"`
	UsageEndDate         string `json:"usageEndDate,omitempty"`
	SubscriptionId       string `json:"subscriptionId,omitempty"`
	InstanceId           string `json:"instanceID,omitempty"`
}

type Subscriber struct {
	SubscriberId              string                `json:"subscriberId"`
	Imsi                      string                `json:"imsi"`
	NewIMSI                   string                `json:"newIMSI"`
	ImsiChangeDate            string                `json:"imsiChangeDate"`
	ServiceProviderId         string                `json:"serviceProviderId"`
	BillingAccountId          string                `json:"billingAccountId"`
	SubscriberAccountType     string                `json:"subscriberAccountType"`
	SubscriberType            int32                 `json:"subscriberType"`
	Status                    SubscriberState       `json:"status"`
	State                     string                `json:"state"`
	Msisdn                    string                `json:"msisdn"`
	SharedAccountOwnerId      string                `json:"sharedAccountOwnerId,omitempty"`
	Groups                    []string              `json:"groups,omitempty"`
	ServiceClass              string                `json:"serviceClass,omitempty"`
	Class                     string                `json:"class,omitempty"`
	Language                  string                `json:"language"`
	ActivationDate            string                `json:"activationDate"`
	CurrentBillCycleStartDate string                `json:"currentBillCycleStartDate,omitempty"`
	LastStatusChangeDate      string                `json:"lastStatusChangeDate"`
	LastProfileUpdateDate     string                `json:"lastProfileUpdateDate"`
	CreditLimit               *float64              `json:"creditLimit,omitempty"`
	ChargeCycle               int64                 `json:"chargeCycle,omitempty"`
	Subscriptions             []Subscription        `json:"subscriptions,omitempty"`
	SubscriptionHistory       []SubscriptionHistory `json:"subscriptionHistory,omitempty"`
	GeneralInfos              *GeneralInfo          `json:"generalInfo,omitempty"`
	SubTaxExemptionFlag       bool                  `json:"taxExemption,omitempty"`
	PrimaryImsi               string                `json:"primaryImsi,omitempty"`
	MultiCard                 []MultiCard           `json:"multiCard,omitempty"`
	ApplyFUPCharges           bool                  `json:"applyFupCharges,omitempty"`
}

type MultiCard struct {
	Msisdn       string          `json:"msisdn"`
	Imsi         string          `json:"imsi"`
	SubscriberId string          `json:"subscriberId"`
	Status       SubscriberState `json:"status"`
}

type DebitedWalletDetails struct {
	BalanceId string
	Fee       float64
}

type SubscriptionState int32

const (
	SubscriptionState_Active              SubscriptionState = 0
	SubscriptionState_Deactivated         SubscriptionState = 1
	SubscriptionState_Suspended           SubscriptionState = 2
	SubscriptionState_Deleted             SubscriptionState = 3
	SubscriptionState_PendingActivation   SubscriptionState = 4
	SubscriptionState_PendingDeactivation SubscriptionState = 5
)

type WalletBalance struct {
	BalanceType            string     `json:"balanceType,omitempty"`
	BalanceName            string     `json:"balanceName" validate:"required"`
	PreviousBalance        *float64   `json:"previousBalance,omitempty"`
	ChangedBalance         *float64   `json:"changedBalance,omitempty"`
	CurrentBalance         *float64   `json:"currentBalance,omitempty"`
	EffectiveStartDate     string     `json:"effectiveStartDate,omitempty"`
	EffectiveExpiryDate    string     `json:"effectiveExpiryDate,omitempty"`
	PreEffectiveStartDate  string     `json:"preEffectiveStartDate,omitempty"`
	PreEffectiveExpiryDate string     `json:"preEffectiveExpiryDate,omitempty"`
	ThresholdPolicy        *Threshold `json:"thresholdPolicy,omitempty"`
	FeeType                string     `json:"feeType,omitempty"`
	PackageID              string     `json:"subscriptionId,omitempty"`
	InstanceID             string     `json:"instanceID,omitempty"`
}

type Threshold struct {
	PolicyId             string `json:"policyId" validate:"required"`
	ThresholdValue       int64  `json:"thresholdValue" validate:"required"`
	ThresholdUnit        string `json:"thresholdUnit" validate:"required"`
	ThresholdDisplayName string `json:"thresholdDisplayName" validate:"required"`
}

var SubscriptionState_name = map[int32]string{
	0: "Active",
	1: "Deactivated",
	2: "Suspended",
	3: "Deleted",
	4: "PendingActivation",
	5: "PendingDeactivation",
}

const (
	PKG int32 = iota + 1
	PRODUCT
)

func (state SubscriptionState) String() string {
	return SubscriptionState_name[int32(state)]
}

func GetSubscriptionState(state int32) SubscriptionState {
	var subscriptionState SubscriptionState
	switch state {
	case 0:
		subscriptionState = SubscriptionState_Active
	case 1:
		subscriptionState = SubscriptionState_Deactivated
	case 2:
		subscriptionState = SubscriptionState_Suspended
	case 3:
		subscriptionState = SubscriptionState_Deleted
	case 4:
		subscriptionState = SubscriptionState_PendingActivation
	case 5:
		subscriptionState = SubscriptionState_PendingDeactivation
	}
	return subscriptionState
}

const (
	CUMULATIVE int32 = iota + 1
)

type FeeType int32

const (
	FeeType_Other     FeeType = 0 // typo need to be corrected
	FeeType_Renewal   FeeType = 1
	FeeType_Both      FeeType = 2
	FeeType_Recharge  FeeType = 3
	FeeType_Purchanse FeeType = 4
)

var FeeType_name = map[int32]string{
	0: "Other",
	1: "RenewalFee",
	2: "Both",
	3: "Recharge",
	4: "PurchaseFee",
}

func (state FeeType) String() string {
	return FeeType_name[int32(state)]
}

func GetFeeType(state int32) FeeType {
	var feeType FeeType
	switch state {
	case 0:
		feeType = FeeType_Other
	case 1:
		feeType = FeeType_Renewal
	case 2:
		feeType = FeeType_Both
	case 3:
		feeType = FeeType_Recharge
	case 4:
		feeType = FeeType_Purchanse
	}
	return feeType
}

type SubscriptionFeeCharged struct {
	RenewalFee             float64 `json:"renewalFee"`
	TaxAmountOnPurchaseFee float64 `json:"taxOnPurchaseFee,omitempty"`
	TaxAmountOnRenewalFee  float64 `json:"taxOnRenewalFee,omitempty"`
	//Added for EDR purpose
	PurchaseFee                     float64 `json:"purchaseFee"`
	RenewalFeeDeductionInQueuedPlan string  `json:"renewalFeeDeductionInQueuedPlan,omitempty"`
}

type QuotaSlab struct {
	Counter       []*QuotaSlabCounter `json:"counter" validate:"required"`
	SlabInputType string              `json:"slabInputType"`
}

type QuotaSlabCounter struct {
	Low   *int32 `json:"low" validate:"required"`
	High  *int32 `json:"high" validate:"required"`
	Quota *int64 `json:"quota" validate:"required"`
}

type RenewalFeeSlab struct {
	Low  int32   `json:"low" validate:"required"`
	High int32   `json:"high" validate:"required"`
	Fee  float64 `json:"fee" validate:"required"`
}
type DiscountInfo struct {
	Discount []DiscountTierData `json:"discount" validate:"required"`
}

type DiscountTierData struct {
	Low           int64    `json:"low" validate:"required"`
	High          int64    `json:"high" validate:"required"`
	DiscountType  string   `json:"discountType" validate:"required"`
	DiscountVal   *float64 `json:"discount" validate:"required"`
	SlabInputType string   `json:"slabInputType,omitempty"`
}

type SubscriptionType int32

const (
	SubscriptionType_BTP        SubscriptionType = 1
	SubscriptionType_ADDON      SubscriptionType = 2
	SubscriptionType_AUTOMATIC  SubscriptionType = 3
	SubscriptionType_FREECHARGE SubscriptionType = 5
	SubscriptionType_GLOBAL     SubscriptionType = 6
)

var SubscriptionType_name = map[int32]string{
	1: "BTP",
	2: "AddOn",
	3: "AUTOMATIC",
	5: "FREECHARGE",
	6: "GLOBAL",
}

func (sType SubscriptionType) String() string {
	return SubscriptionType_name[int32(sType)]
}

func GetSubscriptionType(subsType int32) SubscriptionType {
	var subscriptionType SubscriptionType
	switch subsType {
	case 1:
		subscriptionType = SubscriptionType_BTP
	case 2:
		subscriptionType = SubscriptionType_ADDON
	case 3:
		subscriptionType = SubscriptionType_AUTOMATIC
	case 5:
		subscriptionType = SubscriptionType_FREECHARGE
	case 6:
		subscriptionType = SubscriptionType_GLOBAL
	}
	return subscriptionType
}

func GetSubscriberType(subType int32) string {
	var subTypeString string
	switch subType {
	case 0:
		subTypeString = PREPAID_STR
	case 1:
		subTypeString = POSTPAID_STR
	case 2:
		subTypeString = HYBRID_STR
	}
	return subTypeString
}

func GetSubscriberTypeOperationForKPI(subType int32) string {
	var subTypeOp string
	switch subType {
	case 0:
		subTypeOp = PREPAID_OPERATION
	case 1:
		subTypeOp = POSTPAID_OPERATION
	case 2:
		subTypeOp = HYBRID_OPERATION
	}
	return subTypeOp
}

//TODO: Groups need to be defined
type Groups struct {
}

type FnInfoReq struct {
	FnNumber      string `json:"fnNumber,omitempty"`
	OperationType string `json:"operationType,omitempty"`
}

type FnFeatureReq struct {
	FnGroupId        string      `json:"fnGroupId"`
	MaxMemberInGroup *int32      `json:"maxMemberInGroup"`
	FnInfo           []FnInfoReq `json:"fnInfo"`
}

type FnInfo struct {
	FnNumber []string `json:"fnNumber,omitempty"`
}
type FnFeature struct {
	FnGroupId        string `json:"fnGroupId"`
	MaxMemberInGroup int32  `json:"maxMemberInGroup"`
	FnInfo           FnInfo `json:"fnInfo"`
}

type Subscription struct {
	SubscriptionId           string                  `json:"subscriptionId"`
	PackageName              string                  `json:"packageName"`
	State                    SubscriptionState       `json:"state"`
	PreviousState            SubscriptionState       `json:"previousState"`
	RenewalPolicy            string                  `json:"renewalPolicy"`
	RenewalCycleDay          *int32                  `json:"renewalCycleDay,omitempty"`
	InstanceId               string                  `json:"instanceID"`
	Validity                 int32                   `json:"validity"`
	ValidityUnit             string                  `json:"validityUnit"`
	Type                     SubscriptionType        `json:"type"`
	ExpiryDate               string                  `json:"expiryDate,omitempty"`
	UsageStartDate           string                  `json:"usageStartDate"`
	UsageEndDate             string                  `json:"usageEndDate"`
	NextChargeDate           string                  `json:"nextChargeDate"`
	LastChargeDate           string                  `json:"lastChargeDate"`
	LastRenewRunDate         string                  `json:"lastRenewRunDate"`
	AlertPolicyID            string                  `json:"alertPolicyID,omitempty"`
	Alerts                   map[string]AlertStatus  `json:"alerts,omitempty"`
	RenewalFee               float64                 `json:"renewalFee,omitempty"`
	RenewalFeeSlab           []*RenewalFeeSlab       `json:"renewalFeeSlab,omitempty"`
	DiscountSlab             *DiscountInfo           `json:"discountSlab,omitempty"`
	PurchaseFee              float64                 `json:"purchaseFee"`
	Priority                 int32                   `json:"priority"`
	RollOverPending          bool                    `json:"rollOverPending"`
	ExternalRenewal          bool                    `json:"externalRenewal"`
	DisableRollOver          bool                    `json:"disableRollOver"`
	OldBTPID                 string                  `json:"oldBTPID,omitempty"`
	PreviousUsageEndDate     string                  `json:"previousUsageEndDate,omitempty"`
	FeeRefund                float64                 `json:"feeRefund,omitempty"`
	PurchaseFeeDebitFlag     bool                    `json:"purchaseFeeDebitFlag"`
	ActivationDate           string                  `json:"activationDate,omitempty"`
	MaxRenewalCycle          *int32                  `json:"maxRenewalCycle,omitempty"`
	RenewalCycleCount        *int32                  `json:"renewalCycleCount,omitempty"`
	TaxOnPurchaseFee         string                  `json:"taxOnPurchaseFee,omitempty"`
	TaxOnRenewalFee          string                  `json:"taxOnRenewalFee,omitempty"`
	DiscountOnRenewalFee     float64                 `json:"discountOnRenewalFee,omitempty"`
	DiscountOnPurchaseFee    float64                 `json:"discountOnPurchaseFee,omitempty"`
	FeeTypeDeducted          FeeType                 `json:"feeType,omitempty"`
	ProrationFactor          float64                 `json:"prorationFactor"`
	WalletForPurchaseFee     []string                `json:"walletForPurchaseFee,omitempty"`
	WalletForRenewalFee      []string                `json:"walletForRenewalFee,omitempty"`
	FeeCharged               *SubscriptionFeeCharged `json:"feeCharged,omitempty"`
	RefundDetails            *RefundDetails          `json:"refundDetails,omitempty"`
	ProrateOnCommencement    string                  `json:"prorateOnCommencement,omitempty"`
	ProrateOnTermination     string                  `json:"prorateOnTermination,omitempty"`
	CUGID                    []string                `json:"cugId,omitempty"`
	FeatureType              string                  `json:"featureType,omitempty"`
	FnFeature                *FnFeature              `json:"fnFeature,omitempty"`
	PurFeeDebitedDetails     []*DebitedWalletDetails `json:"purFeeDebitedDetails,omitempty"`
	RenewalFeeDebitedDetails []*DebitedWalletDetails `json:"renewalFeeDebitedDetails,omitempty"`
	Counters                 *CountersProvision      `json:"counters,omitempty"`
	RoamingCC                []string                `json:"roamingCC,omitempty"`
	DestinationCC            []string                `json:"destinationCC,omitempty"`
	RoamingZoneIds           []string                `json:"roamingZoneIds,omitempty"`
	DestinationZoneIds       []string                `json:"destinationZoneIds,omitempty"`
}

type RefundDetails struct {
	UsageEndDate string `json:"usageEndDate"`
}

type SubscriptionHistory struct {
	SubscriptionId string `json:"subscriptionId"`
	UsageStartDate string `json:"usageStartDate"`
	UsageEndDate   string `json:"usageEndDate"`
	LastChargeDate string `json:"lastChargeDate"`
}

//SubscriberAccountParams - structure to store cash Balance
type SubscriberAccountParams struct {
	BalanceId              string                  `json:"balanceId,omitempty"`
	Units                  string                  `json:"units"`
	Priority               int32                   `json:"priority"`
	EffectiveStartDate     string                  `json:"effectiveStartDate"`
	EffectiveExpiryDate    string                  `json:"effectiveExpiryDate"`
	EffectiveExpiryEpoc    int64                   `json:"effectiveExpiryEpoc,omitempty"`
	BalanceType            string                  `json:"balanceType"`
	State                  WalletState             `json:"state"`
	LastRechargeDate       string                  `json:"lastRechargeDate"`
	LastRechargeType       string                  `json:"lastRechargeType"`
	CreditStatus           CreditStatus            `json:"creditStatus"`
	PreviousCycleStatus    *CreditStatus           `json:"previousCycleStatus,omitempty"`
	TransactionID          string                  `json:"transactionId,omitempty"`
	TransactionDescription string                  `json:"transactionDescription,omitempty"`
	ThresholdPolicy        *CashBalThresholdPolicy `json:"thresholdPolicy,omitempty"`
	FeeTypeDeducted        FeeType                 `json:"feeType,omitempty"`
}

//CreditStatus - structure to store credit in cash Balance
type CreditStatus struct {
	Credit                   float64 `json:"credit"`
	LastUpdateDate           string  `json:"lastUpdateDate"`
	CurrentActiveReservation float64 `json:"currentActiveReservation"`
}

type CashBalThresholdPolicy struct {
	ThresholdPolicyID string       `json:"thresholdPolicyID"`
	NotifStatus       *NotifStatus `json:"notifStatus"`
}

type CounterThresholdPolicy struct {
	ThresholdPolicyID string       `json:"thresholdPolicyID"`
	NotifStatus       *NotifStatus `json:"notifStatus"`
}

type WalletState int32

const (
	WalletState_Pre_Use    WalletState = 0
	WalletState_Active     WalletState = 1
	WalletState_Dormant    WalletState = 2
	WalletState_Frozen     WalletState = 3
	WalletState_Suspended  WalletState = 4
	WalletState_Terminated WalletState = 5
)

var WalletState_name = map[int32]string{
	0: "Pre-use",
	1: "Active",
	2: "Dormant",
	3: "Frozen",
	4: "Suspended",
	5: "Terminated",
}

func (state WalletState) String() string {
	return WalletState_name[int32(state)]
}

//Product - structure to store product details
type Product struct {
	PackageID                  string                      `json:"packageId"`
	PackageName                string                      `json:"packageName,omitempty"`
	Name                       string                      `json:"name,omitempty"`
	SubServiceType             []string                    `json:"subServiceType,omitempty"`
	BalanceType                string                      `json:"balanceType"`
	Priority                   int32                       `json:"priority"`
	RenewalPolicy              string                      `json:"renewalPolicy,omitempty"`
	Validity                   *int32                      `json:"validity,omitempty"`
	ValidityUnit               string                      `json:"validityUnit,omitempty"`
	BalanceRollOver            bool                        `json:"balanceRollOver"`
	BalanceRollOverLimit       int64                       `json:"balanceRollOverLimit,omitempty"`
	EffectiveStartDate         string                      `json:"effectiveStartDate"`
	EffectiveExpiryDate        string                      `json:"effectiveExpiryDate"`
	NextResetDate              string                      `json:"nextResetDate,omitempty"`
	LastResetDate              string                      `json:"lastResetDate,omitempty"`
	LastRenewRunDate           string                      `json:"lastRenewRunDate"`
	ProrationFactor            float64                     `json:"ProrationFactor"`
	ThresholdPolicyID          string                      `json:"thresholdPolicyID,omitempty"`
	AlertPolicyID              string                      `json:"alertPolicyID,omitempty"`
	Alerts                     map[string]AlertStatus      `json:"alerts,omitempty"`
	RenewalCycleCount          *int32                      `json:"renewalCycleCount,omitempty"`
	QuotaDetails               *QuotaDetails               `json:"quotaDetails,omitempty"`
	RateProfiles               map[string]*RateProfile     `json:"rateProfiles,omitempty"`
	DiscountProfiles           map[string]*DiscountProfile `json:"discountProfiles,omitempty"`
	RolloverType               string                      `json:"rolloverType,omitempty"`
	BalRolloverValidity        int32                       `json:"balRolloverValidity,omitempty"`
	BalRolloverValidityUnit    string                      `json:"balRolloverValidityUnit,omitempty"`
	AllowedIntervalForRollover int32                       `json:"allowedIntervalForRollover,omitempty"`
	RollOverInstance           bool                        `json:"rollOverInstance"`
	PolicyCounterDetail        *PolicyCounterDetail        `json:"policyCounterDetail,omitempty"`
	SubscriptionPendingUpdate  map[string]interface{}      `json:"subscriptionPendingUpdate,omitempty"`
	RenewalCycleDay            *int32                      `json:"renewalCycleDay,omitempty"`
	AppliedRenewalCycleDay     *int32                      `json:"appliedRenewalCycleDay,omitempty"`
	QuotaExhausted             bool                        `json:"quotaExhausted,omitempty"`
	ProdIdInstanceIdHash       uint32                      `json:"prodIdInstanceIdHash,omitempty"`
}

type PolicyCounterDetail struct {
	PCID       string      `json:"pcId"`
	UnitType   string      `json:"unitType,omitempty"`
	Overridden bool        `json:"overridden"`
	Tiers      []TierRange `json:"tiers"`
}

type TierRange struct {
	StartRange *int64 `json:"startRange,omitempty"`
	EndRange   int64  `json:"endRange"`
	Status     string `json:"status"`
}

//Policy definitions
type SpendingLimitPolicyDetail struct {
	Unit        string       `json:"unit"`
	Overridable bool         `json:"overridable"`
	Tiers       []TierDetail `json:"tiers"`
}

type TierDetail struct {
	StartRange   int64  `json:"startRange"`
	EndRange     int64  `json:"endRange"`
	Status       string `json:"status"`
	CreationDate string `json:"creationDate"`
	ModifiedDate string `json:"modifiedDate"`
}

//QuotaDetails- strtucture to store quota Balance
type QuotaDetails struct {
	Quota                int64                  `json:"quota"`
	QuotaDisplayUnit     string                 `json:"quotaDisplayUnit,omitempty"`
	BalanceStatus        BalanceStatus          `json:"balanceStatus"`
	NotifStatus          *NotifStatus           `json:"notifStatus,omitempty"`
	QuotaSlabInfo        *QuotaSlab             `json:"quotaSlabInfo,omitempty"`
	PreviousCycleBalance *PreviousBalanceStatus `json:"previousCycleBalance,omitempty"`
}

type NotifStatus struct {
	Thresholds []ThresholdStatus `json:"thresholds,omitempty"`
}
type ThresholdStatus struct {
	ThresholdDisplayName string `json:"thresholdDisplayName"`
	Value                int64  `json:"value"`
	Status               bool   `json:"status"`
}
type PreviousBalanceStatus struct {
	Usage   int64 `json:"usage"`
	Balance int64 `json:"balance"`
}

//BalanceStatus - structure to store Balance in quota Balance
type BalanceStatus struct {
	Usage                    int64  `json:"usage"`
	Balance                  int64  `json:"balance"`
	RollOverTruncatedBalance int64  `json:"rollOverTruncatedBalance"`
	LastUpdateDate           string `json:"lastUpdateDate"`
	CurrentActiveReservation int64  `json:"currentActiveReservation"`
	NonUsagePeriodBalance    int64  `json:"nonUsagePeriodBalance,omitempty"`
	OldBalance               int64  `json:"oldBalance,omitempty"`
}

//SubBalances - structure to fetch balances subdoc from subscriber balance doc
type SubBalances struct {
	ProductBalances map[string](map[string]*Product)    `json:"productBalances"`
	CashBalances    map[string]*SubscriberAccountParams `json:"cashBalances"`
}

type Rates struct {
	Low               int64   `json:"low"`
	High              int64   `json:"high"`
	Rate              float64 `json:"rate"`
	MaxAllowedFUPRate float64 `json:"maxAllowedFUPRate,omitempty"`
	TaxAmount         float64 `json:"taxAmount,omitempty"`
	RatedUnit         int64   `json:"ratedUnit,omitempty"`
	UnitSize          *int64  `json:"unitSize,omitempty"`
	RateProfileId     string  `json:"rateProfileId,omitempty"`
	Units             int64   `json:"units,omitempty"`
}

type RateProfile struct {
	SubServiceType     []string            `json:"subServiceType,omitempty"`
	RatType            []string            `json:"ratType,omitempty"`
	Rates              []Rates             `json:"rates,omitempty"`
	UsageDetails       UsageDetails        `json:"usageDetails,omitempty"`
	PreviousCycleUsage *PreviousCycleUsage `json:"previousCycleUsage,omitempty"`
	TaxId              string              `json:"taxId,omitempty"`
	RatingType         RatingType          `json:"ratingType,omitempty"`
	RateProfileVersion string              `json:"rateProfileVersion,omitempty"`
	PremiumRateType    string              `json:"premiumRateType,omitempty"`
}

type RateProfileWithMutex struct {
	RateProf      *RateProfile
	RateProfMutex sync.RWMutex
}

type PreviousCycleUsage struct {
	CumulativeUsage  int64 `json:"cumulativeUsage"`
	TransactionCount int64 `json:"transactionCount"`
}

type TransactionHistoryCDB struct {
}

type SubBalInfoHistory struct {
	KTab        string                     `json:"KTAB"`
	HistoryInfo []SubBalInfoHistoryDetails `json:"historyInfo,omitempty"`
}

type SubBalInfoHistoryDetails struct {
	ActivationDate   string        `json:"activationDate,omitempty"`
	DeactivationDate string        `json:"deactivationDate,omitempty"`
	LineBalDoc       *LineBalances `json:"lineBalInfo,omitempty"`
}

type LineBalances struct {
	Balances            *SubBalances                          `json:"balances"`
	LastUsageTimeStamp  string                                `json:"lastUsageTimeStamp"`
	SubscriberCounters  *SubscriberCounters                   `json:"subscriberCounters,omitempty"`
	ActiveReservations  map[string]interface{}                `json:"activeReservations,omitempty"`
	TransactionHistory  *TransactionHistoryCDB                `json:"transactionHistory,omitempty"`
	SubscriptionCharges map[string]*SubscriptionChargeDetails `json:"subscriptionCharges,omitempty"`
	NextBackupDate      int64                                 `json:"nextBackupDate,omitempty"`
}

//UsageDetailsStructure -
type UsageDetails struct {
	CumulativeUsage       int64  `json:"cumulativeUsage"`
	CumulativeReservation int64  `json:"cumulativeReservation"`
	TransactionCount      int64  `json:"transactionCount"`
	LastUpdateTime        string `json:"lastUpdateTime"`
	// NextResetDate         string `json:"nextResetDate"`
	// LastResetDate         string `json:"lastResetDate"`
}

type discountTime struct {
	startTime time.Time
	endTime   time.Time
}

/*type DiscountProfile struct {
	SubServiceTypes []string  `json:"subServiceTypes"`
	RatType         []string  `json:"ratType"`
	Discounts       Discounts `json:"discounts"`
}
*/

type DiscountType int32

const (
	DiscountType_Percentage DiscountType = 0
	DiscountType_Fixed      DiscountType = 1
)

//Discounts -
type Discounts struct {
	Dates              []string     `json:"dates"`
	Days               []int32      `json:"days"`
	DiscountPercentage int32        `json:"discountPercentage"`
	StartDate          string       `json:"startDate"`
	EndDate            string       `json:"endDate"`
	StartDay           int32        `json:"startDay"`
	EndDay             int32        `json:"endDay"`
	StartTime          string       `json:"startTime"`
	EndTime            string       `json:"endTime"`
	Discount           float64      `json:"discount"`
	DiscountMode       DiscountType `json:"discountMode"`
}

// billcycle
type BSBILLCYCLE struct {
	//	Ktab      string               `json:"ktab"`
	BillCycle map[string]BILLCYCLE `json:"billCycles"`
}
type BILLCYCLE struct {
	//BillCycleID   int64               `json:"billCycleId"`
	BillCycleName string              `json:"billCycleName"`
	BillGenData   []BillGenDataStruct `json:"billGenData"`
}

type BillGenDataStruct struct {
	BillDate           string `json:"billDate"`
	BillDueDate        string `json:"billDueDate"`
	BillDueDateType    int    `json:"billDueDateType"`
	PenaltyGracePeriod int    `json:"penaltyGracePeriod"`
	BillPeriod         struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	} `json:"billPeriod"`
	TestBillDate string `json:"testBillDate"`
}
type SubscriberState int32

const (
	SubscriberState_Active       SubscriberState = 0
	SubscriberState_Preactive    SubscriberState = 1
	SubscriberState_Barred       SubscriberState = 2
	SubscriberState_Deactivated  SubscriberState = 3
	SubscriberState_RechargeOnly SubscriberState = 4
	SubscriberState_NoCredit     SubscriberState = 5
)

var SubscriberState_name = map[int32]string{
	0: "Active",
	1: "Preactive",
	2: "Barred",
	3: "Deactivated",
	4: "RechargeOnly",
	5: "NoCredit",
}

type SubscriberAlias struct {
	SubscriberKey string `json:"subscriberKey"`
	Imsi          string `json:"imsi,omitempty"`
}

type RateProfileInfo struct {
	Status         int32      `json:"status"`
	KTab           string     `json:"KTAB,omitempty"`
	ActivationTime int64      `json:"activationTime"`
	ExpirationTime int64      `json:"expirationTime,omitempty"`
	SubServiceType []string   `json:"subServiceType"`
	RatType        []string   `json:"ratType,omitempty"`
	RatedUnit      int64      `json:"ratedUnit"`
	UnitSize       *int64     `json:"unitSize"`
	Rates          []Rates    `json:"rates"`
	TaxId          string     `json:"taxId,omitempty"`
	RatingType     RatingType `json:"ratingType"`
}

type TaxProfileInfo struct {
	TaxCategory    string    `json:"taxCategory"`
	TaxMode        int32     `json:"taxMode"`
	ActivationTime int64     `json:"activationTime"`
	ExpirationTime int64     `json:"expirationTime"`
	TaxItems       []TaxItem `json:"taxItems"`
}

type TaxItem struct {
	TaxCode string  `json:"taxCode"`
	TaxRate float64 `json:"taxRate"`
}

type RatingType int32

const (
	RatingType_TIERED     RatingType = 1
	RatingType_TELESCOPIC RatingType = 2
)

type DiscountProfile struct {
	Status          int32     `json:"status,omitempty"`
	ActivationTime  int64     `json:"activationTime,omitempty"`
	ExpirationTime  int64     `json:"expirationTime,omitempty"`
	SubServiceTypes []string  `json:"subServiceTypes"`
	RatType         []string  `json:"ratType"`
	Discounts       Discounts `json:"discounts"`
}

type RateProfileCDB struct {
	KTab     string                            `json:"KTAB,omitempty"`
	Profiles map[string]RateProfileVersionInfo `json:"profiles"`
}

type RateProfileVersionInfo struct {
	VersionMap map[string]RateProfileInfo `json:"versions"`
}

type DiscountProfileVersionInfo struct {
	VersionMap map[string]DiscountProfile `json:"versions"`
}

type TaxProfileVersionInfo struct {
	VersionMap map[string]TaxProfileInfo `json:"versions"`
}

type RPTABLES struct {
	RPTABLE_Details map[string]RPTABLE_Detail `json:"tables"`
}

type RPTABLE_Detail struct {
	DefaultRP        string `json:"defaultRP"`
	NetworkAttribute string `json:"mappingAttribute"`
	CreationDate     string `json:"creationDate"`
	ModifiedDate     string `json:"modifiedDate"`
}

type RPTABLE struct {
	NumberofParts int                                         `json:"numberofParts"`
	NumberDetails map[string][]*bulkratetrie.NumPrefixDetails `json:"numberDetails"`
}

type ZoneProfileVersionInfo struct {
	KTab    string                 `json:"KTAB"`
	Profile map[string]ZoneProfile `json:"profiles"`
}

type ZoneProfile struct {
	Description string                  `json:"description,omitempty"`
	ZoneType    string                  `json:"zoneType,omitempty"`
	VersionMap  map[string]ZoneVersions `json:"versions"`
}

type ZoneVersions struct {
	Status         int32    `json:"status,omitempty"`
	ActivationTime int64    `json:"activationTime,omitempty"`
	ExpireDate     int64    `json:"expireTime,omitempty"`
	CountryCode    []string `json:"countryCode"`
}

type MCCMNC struct {
	Mcc         string `json:"mcc,omitempty"`
	Mnc         string `json:"mnc,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type SpNumEvent struct {
	SpecialNumber     string `json:"specialNumber"`
	ServiceType       string `json:"serviceType"`
	ServiceProviderID string `json:"serviceProviderID"`
}

type RoamType int

const (
	RoamType_NoRoam            RoamType = 0
	RoamType_NationalRoam      RoamType = 1
	RoamType_InternationalRoam RoamType = 2
	RoamType_All               RoamType = -1
)

var RoamType_name = map[int]string{
	0:  "NoRoam",
	1:  "NationalRoam",
	2:  "InternationalRoam",
	-1: "All",
}

var RoamType_value = map[string]int{
	"NoRoam":            0,
	"NationalRoam":      1,
	"InternationalRoam": 2,
	"All":               -1,
}

const (
	RateProfileCacheType     string = "rateProfile"
	DiscountProfileCacheType string = "discountProfile"
	SpConfigCacheType        string = "spConfig"
	SubServiceCacheType      string = "subService"
	TaxProfileCacheType      string = "taxProfile"
)

const (
	RateProfileType         string = "rateProfile"
	QuotaProfileType        string = "quotaProfile"
	FupProfileType          string = "fupProfile"
	ConnectionProfileType   string = "connectionProfile"
	NetworkUsageProfileType string = "networkUsageProfile"
)

var ProfileTypeMap = map[string]string{
	"rateProfile":         "RP",
	"quotaProfile":        "QP",
	"fupProfile":          "FUP",
	"connectionProfile":   "Connection",
	"networkUsageProfile": "NetworkUsage",
}

type TimeSchemaProfileVersionInfo struct {
	KTab    string                       `json:"KTAB"`
	Profile map[string]TimeSchemaProfile `json:"profiles"`
}

type TimeSchemaProfile struct {
	Description string                               `json:"description,omitempty"`
	VersionMap  map[string]TimeSchemaProfileVersions `json:"versions"`
}

type TimeSchemaProfileVersions struct {
	Status         int32      `json:"status,omitempty"`
	ActivationTime int64      `json:"activationTime,omitempty"`
	ExpireDate     int64      `json:"expireTime,omitempty"`
	TimeSchema     TimeSchema `json:"timeSchema"`
}

type TimeSchema struct {
	Dates     []string `json:"dates"`
	Days      []int32  `json:"days"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	StartDay  int32    `json:"startDay"`
	EndDay    int32    `json:"endDay"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
}

// //ReservedBalance - structure to store reserved balances
// type ReservedBalance struct {
//      BalanceId   string            `json:"balanceId"`
//      Reservation []CashReservation `json:"reservation"`
// }

//CashReservation - structure to store each cashreservation
// type CashReservation struct {
//      UsedUnits                int64   `json:"usedUnits"`
//      GrantedUnits             int64   `json:"grantedUnits"`
//      CumulativeReservedAmount float64 `json:"cumulativeReservedAmount"`
//      RateProfileID            string  `json:"rateProfileId"`
//      PackageID                string  `json:"packageId"`
// }

/*
//ActiveReservation - active reservatiOn information
type RatingGroup struct {
	GrantedUnits    int64            `json:"grantedUnits"`
	UsedUnits       int64            `json:"usedUnits"`
	ServiceType     string           `json:"serviceType"`
	ReservedBundle  []ReservedBundle `json:"reservedProducts"`
	ReservedBalance *ReservedBalance `json:"reservedBalances"`
	ProductRenewal  bool             `json:"productRenewal"`
}

//ReservedBundle - structure to store reserved product
type ReservedBundle struct {
	BundleId    string `json:"bundleId"`
	Reservation int64  `json:"reservation"`
	PackageID   string `json:"packageId"`
	InstanceID  string `json:"instanceId"`
}

type ReservedBalance struct {
	RateProfileID     string            `json:"rateProfileId"`
	DiscountProfileID string            `json:"discountProfileId"`
	PackageID         string            `json:"packageId"`
	ProductID         string            `json:"productId"`
	InstanceID        string            `json:"instanceId"`
	Balances          []CashReservation `json:"balances"`
}

type CashReservation struct {
	BalanceId         string  `json:"balanceId"`
	ReservationAmount float64 `json:"reservedAmount"`
}*/

type SessionStructure struct {
	RatingGroups         map[string]RatingGroup `json:"ratingGroups"`
	SessionType          string                 `json:"sessionType"`
	SecondaryDocKey      string                 `json:"secondaryDocKey"`
	SessionUsageDetails  *SessionUsageDetails   `json:"sessionUsageDetails"`
	PremiumChargeDetails *PremiumChargeDetails  `json:"premiumChargeDetails,omitempty"`
}

type PRPInfo struct {
	ProfileId      string `json:"profileId"`
	ProfileVersion string `json:"profileVersion"`
}

type PremiumChargeDetails struct {
	ConnectionFeeColleted bool     `json:"connectionFeeColleted,omitempty"`
	NUCCommitTime         int64    `json:"nucCommitTime,omitempty"`
	NUCPRPInfo            *PRPInfo `json:"nucPRPInfo,omitempty"`
	FUPPRPInfo            *PRPInfo `json:"fupPRPInfo,omitempty"`
	PreviousNUCCommitTime int64    `json:"previousNucCommitTime,omitempty"`
}

type SessionUsageDetails struct {
	SessionSvcUsage         int64              `json:"sessionSvcUsage,omitempty"`
	SessionCumReserv        int64              `json:"sessionCumReserv,omitempty"`
	SessionCashBalanceUsage map[string]float64 `json:"sessionCashBalanceUsage,omitempty"`
}

func (state SubscriberState) String() string {
	return SubscriberState_name[int32(state)]
}

type PulseConfig struct {
	Size int64 `json:"size"`
	Unit int64 `json:"unit"`
}

type PulseType int32

const (
	PulseType_SECOND PulseType = 1
	PulseType_MINUTE PulseType = 2
	PulseType_BYTE   PulseType = 3
	PulseType_KIB    PulseType = 4
	PulseType_MIB    PulseType = 5
	PulseType_GIB    PulseType = 6
	PulseType_UNITS  PulseType = 7
	PulseType_KB     PulseType = 8
	PulseType_MB     PulseType = 9
	PulseType_GB     PulseType = 10
)

type Result int32

const (
	Result_NO_OP                 Result = 0
	Result_SUCCESS               Result = 1
	Result_FAILURE               Result = 2
	Result_TRANSACTION_NOT_FOUND Result = 3
)

var Result_name = map[int32]string{
	0: "NO_OP",
	1: "Success",
	2: "Failure",
	3: "Transaction_not_found",
}

type RequestedUnitType int32

const (
	RequestedUnit_VOLUME RequestedUnitType = 1
	RequestedUnit_TIME   RequestedUnitType = 2
	RequestedUnit_SSU    RequestedUnitType = 3
)

func (result Result) String() string {
	return Result_name[int32(result)]
}

type Release int32

const (
	Release_NO_OP   Release = 0
	Release_BALANCE Release = 1
	Release_ALL     Release = 2
	Release_COMMIT  Release = 3
)

var Release_name = map[int32]string{
	0: "NO_OP",
	1: "Balance",
	2: "All",
	3: "Commit",
}

func (release Release) String() string {
	return Release_name[int32(release)]
}

//ActiveReservation
//RatingGroup - rating group active reservatin information
type RatingGroup struct {
	GrantedUnits              int64                                 `json:"grantedUnits"`
	ServiceType               string                                `json:"serviceType"`
	UnlimitedBenefits         bool                                  `json:"unlimitedBenefits"`
	DiscountApplicability     comPb.DiscountApplicability           `json:"discountApplicability"`
	SPPulseConfig             *PulseConfig                          `json:"spPulseInfo"`
	ReservedBundle            []ReservedBundle                      `json:"reservedProducts"`
	ReservedBalance           *ReservedBalance                      `json:"reservedBalances"`
	ReservedBalancePRP        []*ReservedBalance                    `json:"reservedBalancesPRP"`
	LastTransactionDetails    string                                `json:"lastTransactionDetail"`
	LastTransactionResult     Result                                `json:"lastTransactionResult"`
	LastUsedUnitsInfo         []*comPb.UsedUnitInfo                 `json:"lastUsedUnitsInfo"`
	LastGrantISN              *int32                                `json:"grantISN,omitempty"`
	LastUsageReleaseOperation Release                               `json:"lastUsageReleaseOperation"`
	LastReservationDate       string                                `json:"lastReservationDate"`
	ProductRenewed            bool                                  `json:"productRenewed"`
	FinalUnitIndication       comPb.FinalUnitIndication             `json:"finalUnitIndication"`
	Counters                  map[string]*comPb.UsageCounterDetails `json:"counters,omitempty"`
	ExtraChargedAmount        int64                                 `json:"extraChargedAmount"`
	BeforeBTPSwitchUnits      int64                                 `json:"beforeBTPSwitchUnits"`
	AfterBTPSwitchUnits       int64                                 `json:"afterBTPSwitchUnits"`
	BTPSwitchTime             string                                `json:"btpSwitchTime"`
	ReservationTimeUnix       int64                                 `json:"reservationTimeUnix"`
	ReservationDurationInSec  int64                                 `json:"reservationDurationInSec"`
}

type ReservedBalance struct {
	RateProfileID          string                    `json:"rateProfileId`
	RateProfile            RateProfile               `json:"rateProfile"`
	DiscountProfile        *DiscountProfile          `json:"discountProfile,omitempty"`
	DiscountProfileID      string                    `json:"discountProfileId,omitempty"`
	PackageID              string                    `json:"packageId"`
	ProductID              string                    `json:"productId"`
	InstanceID             string                    `json:"instanceId"`
	Balances               []CashReservation         `json:"balances"`
	DiscountUnits          []DiscountUnit            `json:"discountUnits,omitempty"`
	RateProfileVersion     string                    `json:"rateProfileVersion,omitempty"`
	DiscountProfileVersion string                    `json:"discountProfileVersion,omitempty"`
	NestedRateProfMap      map[string]*RateProfile   `json:"nestedRateProfMap,omitempty"`
	ReservedTableInfo      TableInfo                 `json:"reservedTableInfo,omitempty"`
	PackageCountryInfo     *comPb.PackageCountryInfo `json:"packageCountryInfo,omitempty"`
}

type CashReservation struct {
	BalanceId         string  `json:"balanceId"`
	ReservationAmount float64 `json:"reservedAmount"`
}

type DiscountUnit struct {
	Units         int64 `json:"units"`
	ApplyDiscount bool  `json:"applyDiscount"`
}

//has to be deleted and generated structures should be used
//ReservedBundle - structure to store reserved bundles in active reservation block
type ReservedBundle struct {
	BundleId           string                       `json:"bundleId"`
	Reservation        int64                        `json:"reservation"`
	PackageID          string                       `json:"packageId"`
	InstanceID         string                       `json:"instanceId"`
	QuotaProfile       *comPb.QuotaProfile          `json:"quotaProfile,omitempty"`
	AssociatedProdInfo *comPb.AssociatedProductInfo `json:"associatedProductInfo,omitempty"`
	PackageCountryInfo *comPb.PackageCountryInfo    `json:"packageCountryInfo,omitempty"`
}

type ThresholdType int32

const (
	ThresholdType_Value      int32 = 1
	ThresholdType_Percentage int32 = 2
)

var ThresholdType_name = map[int32]string{
	ThresholdType_Value:      "Value",
	ThresholdType_Percentage: "Percentage",
}

var ThresholdType_value = map[string]int32{
	"Value":      ThresholdType_Value,
	"Percentage": ThresholdType_Percentage,
}

type AlertType int32

const (
	AlertType_Expiry  int32 = 1
	AlertType_Renewal int32 = 2
)

var AlertType_name = map[int32]string{
	AlertType_Expiry:  "Expiry",
	AlertType_Renewal: "Renewal",
}

var AlertType_value = map[string]int32{
	"Expiry":  AlertType_Expiry,
	"Renewal": AlertType_Renewal,
}

type AlertSchedule int32

const (
	AlertSchedule_DaysRelative int32 = 1
	AlertSchedule_DaysAbsolute int32 = 2
	AlertSchedule_Hours        int32 = 3
	AlertSchedule_OnExpiry     int32 = 4
	AlertSchedule_OnReset      int32 = 5
	//This is internal AlertSchedule defined for BTP expiry event
	AlertSchedule_OnBTPExpiry int32 = 6
)

var AlertSchedule_name = map[int32]string{
	AlertSchedule_DaysRelative: "DaysRelative",
	AlertSchedule_DaysAbsolute: "DaysAbsolute",
	AlertSchedule_Hours:        "Hours",
	AlertSchedule_OnExpiry:     "OnExpiry",
	AlertSchedule_OnReset:      "OnReset",
	AlertSchedule_OnBTPExpiry:  "OnBTPExpiry",
}

var AlertSchedule_value = map[string]int32{
	"DaysRelative": AlertSchedule_DaysRelative,
	"DaysAbsolute": AlertSchedule_DaysAbsolute,
	"Hours":        AlertSchedule_Hours,
	"OnExpiry":     AlertSchedule_OnExpiry,
	"OnReset":      AlertSchedule_OnReset,
	"OnBTPExpiry":  AlertSchedule_OnBTPExpiry,
}

//@Threshold Policy CDB structure
type ThresholdPolicies struct {
	KTab              string                      `json:"KTAB"`
	ThresholdPolicies map[string]*ThresholdPolicy `json:"policies"`
}

type ThresholdPolicy struct {
	CreationDate     string             `json:"creationDate"`
	ModifiedDate     string             `json:"modifiedDate"`
	ThresholdType    int32              `json:"thresholdType"`
	ThresholdDetails []ThresholdDetails `json:"thresholdDetails"`
}

type ThresholdDetails struct {
	ThresholdDisplayName string  `json:"thresholdDisplayName"`
	Thresholdvalue       float64 `json:"thresholdvalue"`
}

//@ Alert Policy CDB structure
type AlertPolicyDetails struct {
	KTab          string                  `json:"KTAB"`
	AlertPolicies map[string]*AlertPolicy `json:"policies"`
}

type AlertPolicy struct {
	CreationDate string  `json:"creationDate"`
	ModifiedDate string  `json:"modifiedDate"`
	AlertType    int32   `json:"alertType"`
	Alerts       []Alert `json:"alertDetails"`
}

type Alert struct {
	AlertDisplayName   string `json:"alertDisplayName"`
	AlertSchedule      int32  `json:"alertSchedule"`
	AlertScheduleValue *int64 `json:"alertScheduleValue"`
}

type AlrtStatus int32

const (
	AlrtStatus_NOT_SENT     AlrtStatus = 1
	AlrtStatus_SUCCESS_SENT AlrtStatus = 2
	AlrtStatus_FAILURE_SENT AlrtStatus = 3
	AlrtStatus_PAST         AlrtStatus = 4
)

type AlertStatus struct {
	AlertDisplayName string     `json:"alertDisplayName,omitempty"`
	Status           AlrtStatus `json:"status"`
}

// SpendingLimitControlPolicies - SLCPolicy CDB Structure
type SpendingLimitPoliciesCDB struct {
	KTab     string                             `json:"KTAB"`
	Policies map[string]*SpendingLimitPolicyCDB `json:"policies"`
}

type SpendingLimitPolicyCDB struct {
	CreationDate string        `json:"creationDate"`
	ModifiedDate string        `json:"modifiedDate"`
	Unit         string        `json:"unit"`
	Overridable  bool          `json:"overridable"`
	TierDetails  []TierDetails `json:"tiers"`
}

// RequestTierDetails - List of tiers for Spending Limit Control Policy request(SLP)
type TierDetailsReq struct {
	StartRange *float64 `json:"startRange" validate:"required"`
	EndRange   float64  `json:"endRange" validate:"required"`
	Status     string   `json:"status" validate:"required"`
}

// TierDetails - List of tiers for Spending Limit Control Policy (SLP)
type TierDetails struct {
	StartRange *int64 `json:"startRange" validate:"required"`
	EndRange   int64  `json:"endRange" validate:"required"`
	Status     string `json:"status" validate:"required"`
}

//ClassOfServiceDetails - Class Of Service CDB structure
type ClassOfServiceDetails struct {
	KTab           string                     `json:"KTAB"`
	CosDefinitions map[string]*ClassOfService `json:"cosDefinitions"`
}

type ClassOfService struct {
	CreationDate          string `json:"creationDate"`
	ModifiedDate          string `json:"modifiedDate"`
	RefCount              int32  `json:"refCount"`
	SubscriptionType      int32  `json:"subscriptionType"`
	SubscriberLifeCycleId string `json:"subscriberLifeCycleId"`
}

type ServiceConfig struct {
	Services map[string]*SubServices `json:"services"`
}

type SubServices struct {
	AllowedSubServices []string `json:"allowedSubServices"`
}

type SubInfo struct {
	ServiceProviderId string `json:"serviceProviderId"`
	CosId             string `json:"cosId"`
	Status            string `json:"status"`
	Cas               int64  `json:"cas"`
}

//const for subscriber type
const (
	PREPAID int32 = iota
	POSTPAID
	HYBRID
)

const (
	HttpStatus_400 = "400"
	HttpStatus_404 = "404"
)

const (
	ResponseCode_400000 = "400000"
	ResponseCode_400001 = "400001"
	ResponseCode_400002 = "400002"
	ResponseCode_404000 = "404000"
	ResponseCode_404001 = "404001"
)

const (
	PREPAID_STR  = "PRE"
	POSTPAID_STR = "POST"
	HYBRID_STR   = "HYD"
)

const (
	PREPAID_OPERATION  = "PREPAID_OPERATION"
	POSTPAID_OPERATION = "POSTPAID_OPERATION"
	HYBRID_OPERATION   = "HYBRID_OPERATION"
)

const (
	GC_BALANCE_ID = "BAL000000"
)

// Constants for EventName used in StateConfig.StateTransitions for Prepaid
const (
	PRE_SERVICE_USAGE         = "FIRST-SERVICE-USAGE"
	PRE_NO_CREDIT             = "NO-CREDIT"
	PRE_EXPLICIT_ACTIVATION   = "EXPLICIT-ACTIVATION"
	PRE_EXPLICIT_DEACTIVATION = "EXPLICIT-DEACTIVATION"
	PRE_STATE_TIMEOUT         = "STATE-TIMEOUT"
	PRE_CASH_BALANCE_EXPIRY   = "CASH-BALANCE-EXPIRY"
	PRE_EXPLICIT_BAR          = "EXPLICIT-BAR"
	PRE_EXPLICIT_UNBAR        = "EXPLICIT-UNBAR"
	PRE_RECHARGE              = "RECHARGE"
	PRE_BASE_PLAN_EXPIRED     = "BASE-PLAN-EXPIRED"
	PRE_BALANCE_TRANSFER      = "BALANCE-TRANSFER"
	PRE_BASE_PLAN_ACTIVATED   = "BASE-PLAN-ACTIVATED"
)

// Constants for EventName used in StateConfig.StateTransitions for Postpaid
const (
	POST_SERVICE_USAGE         = "FIRST-SERVICE-USAGE"
	POST_EXPLICIT_ACTIVATION   = "EXPLICIT-ACTIVATION"
	POST_EXPLICIT_DEACTIVATION = "EXPLICIT-DEACTIVATION"
	POST_STATE_TIMEOUT         = "STATE-TIMEOUT"
	POST_EXPLICIT_BAR          = "EXPLICIT-BAR"
	POST_EXPLICIT_UNBAR        = "EXPLICIT-UNBAR"
	POST_PAYMENT_RECEIVED      = "PAYMENT-RECEIVED"
	POST_NO_CREDIT             = "NO-CREDIT"
	POST_CREDIT_LIMIT_INCREASE = "CREDIT-LIMIT-INCREASE"
	//Extra event
	ADMIN_STATE_CHANGE = "ADMIN-STATE-CHANGE"
)

type ServiceList struct {
	Services map[string]*SubServiceList `json:"services"`
}

// SubscriberLifeCycleInfo - Subscriber LifeCycle List CDB structure
// LifeCycles - Map of various Subscriber LifeCycles with subscriberLifeCycleID as the key
type SubscriberLifeCycleInfo struct {
	KTab       string                          `json:"KTAB"`
	LifeCycles map[string]*SubscriberLifeCycle `json:"lifeCycles"`
}

// States - Map of various States with StateName as the key
type SubscriberLifeCycle struct {
	CreationDate     string                  `json:"creationDate"`
	ModifiedDate     string                  `json:"modifiedDate"`
	SubscriptionType int32                   `json:"subscriptionType"`
	RefCount         int32                   `json:"refCount"`
	States           map[string]*StateConfig `json:"states"`
}

type StateConfig struct {
	StateTransitions    []Transition               `json:"stateTransitions" validate:"required"`
	ServiceAvailability map[string]*SubServiceList `json:"serviceAvailability" validate:"required"`
	StateValidity       *StateValidity             `json:"stateValidity" validate:"required"`
}

type Transition struct {
	EventName string `json:"eventName"`
	NextState string `json:"nextState"`
}

type StateValidity struct {
	Period               int32 `json:"period" validate:"required"`
	EligibleForExtension bool  `json:"eligibleForExtension" validate:"required"`
}

type ServiceProvider struct {
	KTab                         string                        `json:"KTAB"`
	ServiceProviderID            string                        `json:"serviceProviderId"`
	ServiceProviderName          string                        `json:"serviceProviderName"`
	ServiceProviderType          string                        `json:"serviceProviderType"`
	ServiceProviderConfiguration *ServiceProviderConfiguration `json:"serviceProviderConfiguration,omitempty"`
	DatabaseIdentifier           string                        `json:"databaseIdentifier,omitempty"`
}

type ServiceProviderConfiguration struct {
	CurrencyUsage               *CurrencyUsage                          `json:"currencyUsage,omitempty"`
	ServiceUsage                map[string]*ServUsage                   `json:"serviceUsage,omitempty"`
	HomeNetworkDetails          *HomeNetworkDetails                     `json:"homeNetworkDetails,omitempty"`
	SubscriberTypeConfiguration map[string]*SubscriberTypeConfiguration `json:"subscriberTypeConfigurations,omitempty"`
	FeatureLevelConfiguration   map[string]map[string]interface{}       `json:"featureLevelConfiguration,omitempty"`
}

type HomeNetworkDetails struct {
	CountryCode  int32    `json:"countryCode,omitempty"`
	Mcc          string   `json:"mcc,omitempty"`
	Mnc          string   `json:"mnc,omitempty"`
	NetworkCodes []string `json:"networkCodes,omitempty"`
	MscAddress   string   `json:"mscAddress,omitempty"`
}

type CurrencyUsage struct {
	CurrencyConversionFactor    int64  `json:"currencyConversionFactor,omitempty"`
	CurrencyDecimalPlaces       int32  `json:"currencyDecimalPlaces,omitempty"`
	CurrencyRoundingRule        string `json:"currencyRoundingRule,omitempty"`
	CurrencyDecimalSeparator    string `json:"currencyDecimalSeparator,omitempty"`
	CurrencyDisplaySymbol       string `json:"currencyDisplaySymbol,omitempty"`
	UsageChgRoundingRule        string `json:"usageChgRoundingRule,omitempty"`
	RoundingPrecisionForPrepaid int32  `json:"roundingPrecisionForPrepaid,omitempty"`
}

type ServUsage struct {
	Size int64 `json:"size,omitempty"`
	Unit int64 `json:"unit,omitempty"`
}

type SubscriberTypeConfiguration struct {
	SupportedStates  []string `json:"supportedStates,omitempty"`
	GcWalletValidity int32    `json:"gcWalletValidity,omitempty"`
}

type SubServiceList struct {
	SubServices []string `json:"allowedSubServices"`
}

type SubscriberLifeCycleMaster struct {
	KTab        string                     `json:"KTAB"`
	EventList   map[string][]string        `json:"eventList"`
	ServiceList map[string]*SubServiceList `json:"services"`
	SubService  map[int32]string           `json:"subServiceMapping,omitempty"`
}

// constants for subscriber sub services
const (
	CALL_LEG                 = "CALL_LEG"
	ROAM_TYPE                = "ROAM_TYPE"
	CALL_DESTINATION_TYPE    = "CALL_DESTINATION_TYPE"
	FORWARDED_FLAG           = "FORWARDED_FLAG"
	SMS_TYPE                 = "SMS_TYPE"
	RENEWAL_SUPPORT          = "RENEWAL_SUPPORT"
	RECHARGE_SUPPORT         = "RECHARGE_SUPPORT"
	BALANCE_TRANSFER_SUPPORT = "BALANCE_TRANSFER_SUPPORT"
)

const (
	CallDestinationType_NATIONAL      float64 = 1.0
	CallDestinationType_INTERNATIONAL float64 = 2.0
)

const (
	SubServiceRoamType_NO_ROAM       float64 = 1.0
	SubServiceRoamType_NATIONAL      float64 = 2.0
	SubServiceRoamType_INTERNATIONAL float64 = 3.0
)

const (
	CallLeg_MO float64 = 1.0
	CallLeg_MT float64 = 2.0
)

const (
	ForwardedFlag_NON_FORWARDED float64 = 1.0
	ForwardedFlag_FORWARDED     float64 = 2.0
)

const (
	SmsType_MOMT float64 = 201.0
	SmsType_MOAT float64 = 202.0
	SmsType_AOMT float64 = 203.0
	SmsType_AOAT float64 = 204.0
)

const (
	TotalNoOfSecondsInADay int64 = 60 * 60 * 24 // 86400
)

const (
	Renewal_RENEWAL_SUPPORT = 1.0
)

const (
	Recharge_RECHARGE_SUPPORT = 1.0
)

const (
	Balance_Transfer_RECHARGE_SUPPORT = 1.0
)

const (
	CASH_BALANCE = "CASH_BALANCE"
	ALL_BALANCES = "ALL_BALANCES"
)

type SubscriptionCharges struct {
	SubscriptionTransactionMap map[string]*SubscriptionChargeDetails `json:"subscriptionCharges"`
}

type SubscriptionChargeDetails struct {
	PackageDetails       []PackageDetail        `json:"packageDetails"`
	FeeCollected         float64                `json:"feeCollected"`
	TransactionStatus    string                 `json:"transactionStatus"`
	RenewalDate          string                 `json:"renewalDate,omitempty"`
	ReservationTimeStamp string                 `json:"reservationTimeStamp"`
	OtherSubscriber      string                 `json:"otherSubscriber,omitempty"`
	BackUpData           map[string]interface{} `json:"backUpData,omitempty"`
}

type PackageDetail struct {
	PackageID    string `json:"packageId"`
	InstanceID   string `json:"instanceId"`
	UsageEndDate string `json:"usageEndDate,omitempty"`
}

type CounterProvisionDetails struct {
	CounterName     string `json:"counterName"`
	ThresholdPolicy string `json:"thresholdPolicy,omitempty"`
	MaximumLimit    int64  `json:"maximumLimit,omitempty"`
}

type CountersProvision struct {
	PackageLevelServiceCounter map[string]map[string]*CounterProvisionDetails `json:"serviceCounters,omitempty"`
	PackageLevelPackageCounter map[string]*CounterProvisionDetails            `json:"packageCounters,omitempty"`
	PackageLevelWalletCounter  map[string]map[string]*CounterProvisionDetails `json:"walletCounters,omitempty"`
}

type PackageLevelCounters struct {
	PackageLevelServiceCounter map[string]map[string]*CounterlistDetails `json:"serviceCounters,omitempty"`
	PackageLevelPackageCounter map[string]*CounterlistDetails            `json:"packageCounters,omitempty"`
	PackageLevelWalletCounter  map[string]map[string]*CounterlistDetails `json:"walletCounters,omitempty"`
}

type PackageCounterDetails struct {
	PackageLevelServiceCounter map[string]map[string]*SubscriberCounterDetails `json:"serviceCounters,omitempty"`
	PackageLevelPackageCounter map[string]*SubscriberCounterDetails            `json:"packageCounters,omitempty"`
	PackageLevelWalletCounter  map[string]map[string]*SubscriberCounterDetails `json:"walletCounters,omitempty"`
}

type SubscriberCounters struct {
	ServiceLevelCounter    map[string]map[string]*SubscriberCounterDetails `json:"serviceLevelCounter,omitempty"`
	PackageLevelCounter    map[string]map[string]*PackageCounterDetails    `json:"packageLevelCounter,omitempty"`
	SubscriberLevelCounter map[string]*SubscriberCounterDetails            `json:"subscriberLevelCounter,omitempty"`
	WalletLevelCounter     map[string]map[string]*SubscriberCounterDetails `json:"walletLevelCounter,omitempty"`
}
type SubscriberCounterDetails struct {
	Usage           int64                   `json:"usage"`
	Reserved        int64                   `json:"reserved"`
	UnitType        string                  `json:"unitType"`
	ResetDate       string                  `json:"resetDate,omitempty"`
	ProrationFactor float64                 `json:"prorationFactor,omitempty"`
	MaximumLimit    int64                   `json:"maximumLimit,omitempty"`
	RatedUsage      *RatedUsageDetails      `json:"ratedUsage,omitempty"`
	ThresholdPolicy *CounterThresholdPolicy `json:"thresholdPolicy,omitempty"`
}
type RatedUsageDetails struct {
	Volume int64 `json:"volume"`
	Time   int64 `json:"time"`
	Unit   int64 `json:"unit"`
}
type CounterDefinition struct {
	CounterType string `json:"counterType"`
	ServiceType string `json:"serviceType"`
	UnitType    string `json:"unitType"`
}
type RuleCounters struct {
	Counters map[string]*RuleCounterDetails `json:counters`
}
type RuleCounterDetails struct {
	CounterType string `json:"counterType"`
	UnitType    string `json:"unitType"`
	CounterPath string `json:"counterPath"`
}
type SubscriberCountersDefinations struct {
	KTab        string                              `json:"KTAB"`
	Counterlist *CounterListMasterDocDetails        `json:"counterList" validate:"required"`
	Rule        map[string]map[string]*RuleCounters `json:"rules" validate:"required"`
}

type CounterListMasterDocDetails struct {
	Service    map[string]map[string]*CounterlistDetails `json:"serviceLevelCounter,omitempty"`
	Subscriber map[string]*CounterlistDetails            `json:"subscriberLevelCounter,omitempty"`
	Wallet     map[string]map[string]*CounterlistDetails `json:"walletLevelCounter,omitempty"`
	Packages   *PackageLevelCounters                     `json:"packageLevelCounter,omitempty"`
}

type CounterlistDetails struct {
	CounterType           string `json:"counterType" validate:"required"`
	ServiceType           string `json:"serviceType,omitempty"`
	UnitType              string `json:"unitType" validate:"required"`
	DisplayUnit           string `json:"displayUnit" validate:"required"`
	MaximumLimit          int64  `json:"maximumLimit,omitempty"`
	ResetPolicy           string `json:"resetPolicy,omitempty"`
	ThresholdPolicy       string `json:"thresholdPolicy,omitempty"`
	Status                string `json:"status,omitempty"`
	IsProrationApplicable bool   `json:"isProrationApplicable,omitempty"`
	CreationDate          string `json:"creationDate,omitempty"`
	ModifiedDate          string `json:"modifiedDate,omitempty"`
}

const (
	TRANSACTION_STATUS_RESERVE = "RESERVE"
	TRANSACTION_STATUS_COMMIT  = "COMMIT"
	TRANSACTION_STATUS_REVERT  = "REVERT"
	TRANSACTION_STATUS_CREDIT  = "CREDIT"
	TRANSACTION_STATUS_DEBIT   = "DEBIT"
)

const (
	SUBSCRIPTION_CHARGES_SUB_DOC            = "subscriptionCharges"
	SUBSCRIPTION_CHARGES_RENEWALS           = "RENEWALS"
	SUBSCRIPTION_CHARGES_TRANSACTION_STATUS = "transactionStatus"
)

var SubTypeMap = map[int32]string{
	0: "Prepaid",
	1: "Postpaid",
	2: "Hybrid",
}

type SlcpAlias struct {
	SubscriberKey string `json:"subscriberKey"`
}

const (
	SRV_RECHARGE         = "RECHARGE"
	SRV_BALANCE_TRANSFER = "TRANSFERBALANCE"
	SRV_SMS              = "SMS"
	SRV_MMS              = "MMS"
	SRV_DATA             = "DATA"
	SRV_VOICE            = "VOICE"
)

type PackageCDB struct {
	KTab                  string                        `json:"KTAB,omitempty"`
	PackageId             string                        `json:"packageId,omitempty"`
	PackageName           string                        `json:"packageName,omitempty"`
	PackageVersion        string                        `json:"packageVersion,omitempty"`
	PackageType           string                        `json:"packageType,omitempty"`
	PackageUserType       string                        `json:"packageUserType,omitempty"`
	Description           string                        `json:"description,omitempty"`
	ServiceProviderId     string                        `json:"serviceProviderId,omitempty"`
	PurchaseFee           *float64                      `json:"purchaseFee,omitempty"`
	PurchasePolicy        string                        `json:"purchasePolicy,omitempty"`
	EffectiveStartDate    string                        `json:"effectiveStartDate,omitempty"`
	ProrateOnCommencement string                        `json:"prorateOnCommencement,omitempty"`
	ProrateOnTermination  string                        `json:"prorateOnTermination,omitempty"`
	RenewalPolicy         string                        `json:"renewalPolicy,omitempty"`
	RenewalFee            *float64                      `json:"renewalFee,omitempty"`
	ProrationFactor       float64                       `json:"ProrationFactor"`
	FeeCharged            *SubscriptionFeeCharged       `json:"feeCharged,omitempty"`
	RefundDetails         *RefundDetails                `json:"refundDetails,omitempty"`
	RenewalCycleDay       *int32                        `json:"renewalCycleDay,omitempty"`
	ActionOnRenewalFail   string                        `json:"actionOnRenewalFail,omitempty"`
	AlertNotiPolicy       string                        `json:"alertNotiPolicy,omitempty"`
	ActivationDate        string                        `json:"activationDate,omitempty"`
	ExpiryDate            string                        `json:"expiryDate,omitempty"`
	Validity              *int32                        `json:"validity,omitempty"`
	ValidityUnit          string                        `json:"validityUnit,omitempty"`
	State                 string                        `json:"state,omitempty"`
	AssociatedProducts    map[string]*ProductProfileCDB `json:"associatedComponents,omitempty"`
	Priority              *int32                        `json:"priority,omitempty"`
	TaxOnPurchaseFee      string                        `json:"taxOnPurchaseFee,omitempty"`
	TaxOnRenewalFee       string                        `json:"taxOnRenewalFee,omitempty"`
	CUGID                 []string                      `json:"cugId,omitempty"`
	FeatureType           string                        `json:"featureType,omitempty"`
	FnFeature             *FnFeature                    `json:"fnFeature,omitempty"`
	DiscountOnRenewalFee  float64                       `json:"discountOnRenewalFee,omitempty"`
	DiscountOnPurchaseFee float64                       `json:"discountOnPurchaseFee,omitempty"`
	RoamingCC             []string                      `json:"roamingCC,omitempty"`
	DestinationCC         []string                      `json:"destinationCC,omitempty"`
	RoamingZoneIds        []string                      `json:"roamingZoneIds,omitempty"`
	DestinationZoneIds    []string                      `json:"destinationZoneIds,omitempty"`
}

type ProductProfileCDB struct {
	KTab                       string   `json:"KTAB,omitempty"`
	Name                       string   `json:"name,omitempty"`
	Description                string   `json:"description,omitempty"`
	ProductType                string   `json:"productType,omitempty"`
	ServiceType                string   `json:"serviceType,omitempty"`
	SubServiceType             []string `json:"subServiceType,omitempty"`
	RatType                    []string `json:"ratType,omitempty"`
	Qos                        string   `json:"qos",omitempty"`
	BalanceRollover            *bool    `json:"balanceRollover,omitempty"`
	RolloverType               string   `json:"rolloverType,omitempty"`
	BalanceRolloverLimit       *int64   `json:"balanceRolloverLimit,omitempty"`
	BalRolloverValidity        *int32   `json:"balRolloverValidity,omitempty"`
	BalRolloverValidityUnit    string   `json:"balRolloverValidityUnit,omitempty"`
	AllowedIntervalForRollover *int32   `json:"allowedIntervalForRollover,omitempty"`
	State                      string   `json:"state,omitempty"`
	Quota                      *int64   `json:"quota,omitempty"`
	QuotaDisplayUnit           string   `json:"quotaDisplayUnit,omitempty"`
	RenewalPolicy              string   `json:"renewalPolicy,omitempty"`
	AlertNotiPolicy            string   `json:"alertNotiPolicy,omitempty"`
	ThresholdNotiPolicy        string   `json:"thresholdNotiPolicy,omitempty"`
	SlcPolicy                  string   `json:"slcPolicy,omitempty"`
	Validity                   *int32   `json:"validity,omitempty"`
	ValidityUnit               string   `json:"validityUnit,omitempty"`
}

type CreditStatusCDB struct {
	Credit                    float64 `json:"credit"`
	LastUpdateDate            string  `json:"lastUpdateDate,omitempty"`
	CurrentActiveReservations int64   `json:"currentActiveReservation,omitempty"`
}

type CashBalancesCDB struct {
	BalanceId              string                  `json:"balanceId,omitempty"`
	BalanceType            string                  `json:"balanceType,omitempty"`
	EffectiveStartDate     string                  `json:"effectiveStartDate,omitempty"`
	EffectiveExpiryDate    string                  `json:"effectiveExpiryDate,omitempty"`
	LastRechargeType       string                  `json:"lastRechargeType"`
	LastRechargeDate       string                  `json:"lastRechargeDate"`
	Priority               *int32                  `json:"priority"`
	State                  int32                   `json:"state"`
	ThresholdPolicy        *CashBalThresholdPolicy `json:"thresholdPolicy,omitempty"`
	CreditStatus           *CreditStatusCDB        `json:"creditStatus"`
	TransactionID          string                  `json:"transactionId"`
	TransactionDescription string                  `json:"transactionDescription"`
}

//CheckConnection - func for readiness and liveliness probe
func CheckConnection() bool {
	mlog.MavLog(mlog.INFO, "CB_CONN_CHECK", "Enter- CheckConnection")
	if CollectionObj != nil {
		if !IsCouchbaseHealthOkV2(gCDbConnectionList) {
			mlog.MavLog(mlog.ERROR, "CB_CONN_CHECK", "Exit- CheckConnection, Couchbase connection down")
			mlog.MavAlarm(ALARM_TOPIC, "CHFCouchbaseConnectionDown")
			CBDownAlarmRaised = true
			return false
		}
		mlog.MavLog(mlog.INFO, "CB_CONN_CHECK", "Exit- CheckConnection, Couchbase connection active")
		if CBDownAlarmRaised {
			mlog.MavAlarm(ALARM_TOPIC, "CHFCouchbaseConnectionUp")
		}
		return true
	} else {
		mlog.MavLog(mlog.INFO, "CB_CONN_CHECK", "Exit- CheckConnection, Couchbase connection object is nil")
		return false
	}
}

func InitializeCbConnection() error {
	initializeKTAB()
	mlog.MavLog(mlog.INFO, transID, "Enter- InitializeCbConnection")
	var err error
	if Cluster == nil {
		opts := gocb.ClusterOptions{Username: config.CommonStaticConf.Config.CouchbaseConfig.Username, Password: config.CommonStaticConf.Config.CouchbaseConfig.Password}
		Cluster, err = gocb.Connect(config.CommonStaticConf.Config.CouchbaseConfig.Hostname, opts)
		if err != nil {
			mlog.MavLog(mlog.ERROR, transID, "Exit- Error in connecting to couchbase: ", err)
			return err
		}
	}
	if CollectionObj == nil {
		mlog.MavLog(mlog.INFO, transID, "Couchbase connection is already established")
		bucket := Cluster.Bucket(config.CommonDynConf.Config.CouchbaseConfig.Bucket)
		if bucket == nil {
			mlog.MavLog(mlog.ERROR, transID, "Exit- InitializeCbConnection")
			return errors.New("Failed to get Bucket")
		}
		err = bucket.WaitUntilReady(time.Second*time.Duration(config.CommonDynConf.Config.CouchbaseConfig.ConnectTimeoutSec), nil)
		if err != nil {
			mlog.MavLog(mlog.ERROR, transID, "Exit- Timeout in connecting to bucket:", err)
			return err
		}
		CollectionObj = bucket.DefaultCollection()
		mlog.MavLog(mlog.INFO, transID, "Exit- InitializeCbConnection, connection ready")
	}

	initCouchBaseConfig(&gCommonCbCfgData)
	gCDbConnectionList, err = InitCouchDBConnectionv2(gCommonCbCfgData, false)
	if err != nil {
		//log.Println("InitCouchDBConnection failed", err)
		mlog.MavLog(mlog.ERROR, transID, "InitCouchDBConnection failed", err)

		//raise alarm if DB is down
		mlog.MavAlarm(ALARM_TOPIC, "CHFCouchbaseConnectionDown")
		return err
	}
	//log.Println("Init CouchDBConnection Success.")
	mlog.MavLog(mlog.ERROR, transID, "Init CouchDBConnection Success.")
	mlog.MavAlarm(ALARM_TOPIC, "CHFCouchbaseConnectionUp")

	mlog.MavLog(mlog.INFO, transID, "Exit- InitializeCbConnection")
	return nil
}

func initCouchBaseConfig(commonCbCfgData *CommonCbCfgStruct) {
	mlog.MavLog(mlog.INFO, transID, "Enter- InitializeCbConnection")
	var cfgObj CommCbCfgIntStruct
	var localDc bool = true
	cfgObj.ClusterName = "DC1"
	cfgObj.CouchBaseIP = config.CommonStaticConf.Config.CouchbaseConfig.Hostname
	cfgObj.CouchBasePort = "8091"
	cfgObj.CouchBaseUname = config.CommonStaticConf.Config.CouchbaseConfig.Username
	cfgObj.CouchBasePwd = config.CommonStaticConf.Config.CouchbaseConfig.Password
	cfgObj.CouchBaseBucketName = config.CommonDynConf.Config.CouchbaseConfig.Bucket
	cfgObj.CbEphBucketName = ""
	cfgObj.IsLocalDc = &localDc
	cfgObj.BucketOpenTimeout = config.CommonDynConf.Config.CouchbaseConfig.Timeout
	cfgObj.ConnectionTimeout = config.CommonDynConf.Config.CouchbaseConfig.ConnectTimeoutSec
	cfgObj.MaxCouchFailureThresholdCount = 0
	cfgObj.MaxNodeFailureThresholdPerc = 0

	commonCbCfgData.CommCbIntData.CommCbCfgData = append(commonCbCfgData.CommCbIntData.CommCbCfgData, cfgObj)
	mlog.MavLog(mlog.INFO, transID, "CommCbCfgData - ", len(commonCbCfgData.CommCbIntData.CommCbCfgData))
	mlog.MavLog(mlog.INFO, transID, "Exit- InitializeCbConnection")
}

/*
//GetBucketConnectionObj - Func to get the connection object for default collection in the bucket
func GetBucketConnectionObj(bucketName string) (*gocb.Collection, error) {
	mlog.MavTrace(transDetails.TraceFlag,mlog.INFO, transID, "Enter- GetBucketConnectionObj")
	if len(bucketName) > 0 {
		bucket := Cluster.Bucket(bucketName)
		if bucket == nil {
			mlog.MavTrace(transDetails.TraceFlag,mlog.ERROR, transID, "Exit- GetBucketConnectionObj")
			return nil, errors.New("Failed to get Bucket")
		}
		mlog.MavTrace(transDetails.TraceFlag,mlog.INFO, transID, "Exit- GetBucketConnectionObj")
		return bucket.DefaultCollection(), nil
	}
	mlog.MavTrace(transDetails.TraceFlag,mlog.ERROR, transID, "Exit- GetBucketConnectionObj")
	return nil, errors.New("Bucket name is invalid")
}
*/

/*func CreateThresholdStatus(transDetails *config.LoggerTransactionDetails, thresholdID string, serviceProviderID string, quota int64) ([]ThresholdStatus, error) {

	if thresholdID == "" {
		return nil, errors.New("EMPTY_THRESHOLD_ID")
	}
	docKey := THRESHOLD_POLICY_SID + serviceProviderID + "::ThresholdPolicies"
	thresholdPolicy := ThresholdPolicy{}

	getPaths := []string{}
	getTargetArr := []interface{}{}
	getPaths = append(getPaths, "policies."+thresholdID)
	getTargetArr = append(getTargetArr, &thresholdPolicy)

	_, err := GetSubDocuments(transDetails, docKey, getPaths, getTargetArr)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			mlog.MavTrace(transDetails.TraceFlag,mlog.ERROR, transDetails.TransID, "Error - createThresholdStatus", err)
			return nil, err
		} else if errors.Is(err, gocb.ErrPathNotFound) {
			mlog.MavTrace(transDetails.TraceFlag,mlog.ERROR, transDetails.TransID, "Error - createThresholdStatus", err)
			return nil, err
		}
	}

	tempTHstatus := make([]ThresholdStatus, len(thresholdPolicy.ThresholdDetails))
	for key, value := range thresholdPolicy.ThresholdDetails {

		if thresholdPolicy.ThresholdType == ThresholdType_Value {
			tempTHstatus[key].Value = int64(value.Thresholdvalue)

		} else if thresholdPolicy.ThresholdType == ThresholdType_Percentage {
			val := 100 - value.Thresholdvalue
			tempTHstatus[key].Value = int64((val * float64(quota)) / 100)
		}
		tempTHstatus[key].Status = false

	}
	return tempTHstatus, nil
}*/

//GetDocument - get a document that matches primaryKey from couchbase. This API returns CAS and error if occurred
func GetDocument(transDetails *config.LoggerTransactionDetails, key string, target []interface{}) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetDocument, Key: ", key)
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET, config.GetKtab(key), serviceProviderID})
	crdlResponse, err := CollectionObj.Get(key, &gocb.GetOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in GetDocument")
		return 0, err
	}

	err = crdlResponse.Content(target[0])
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in GetDocument: ", err)
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetDocument")
	return crdlResponse.Cas(), nil
}

//CreateDocument - Insert a document into couchbase with provided key. This API tries to insert a doc into couchbase with provided key. If the key already exists, then it returns error Otherwise it will insert a new doc. This API returns CAS of the doc after updation and error if occurred.
func CreateDocument(transDetails *config.LoggerTransactionDetails, key string, target []interface{}, ttlTimer *time.Duration) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- InsertDocument, Key: ", key)
	inOps := gocb.InsertOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			inOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			inOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.INSERT, config.GetKtab(key), serviceProviderID})
	InsertResponse, err := CollectionObj.Insert(key, target[0], &inOps)
	args := []string{config.CB_RESP, config.INSERT, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in InsertDocument: ", err)
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- InsertDocument")
	return InsertResponse.Cas(), nil
}

//UpsertDocument - Upsert a document into couchbase with provided key. This API tries to insert a doc into couchbase with provided key. If the key already exists, then the doc is updated with the doc passed as param to this API. Other wise it will insert a new doc. This API returns CAS of the doc after updation and error if occurred.
func UpsertDocument(transDetails *config.LoggerTransactionDetails, key string, target []interface{}, ttlTimer *time.Duration) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertDocument, Key: ", key)
	upOps := gocb.UpsertOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			upOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			upOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.UPSERT, config.GetKtab(key), serviceProviderID})
	InsertResponse, err := CollectionObj.Upsert(key, target[0], &upOps)
	args := []string{config.CB_RESP, config.UPSERT, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in UpsertDocument")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertDocument")
	return InsertResponse.Cas(), nil
}

func ReplaceDocumentWithCas(transDetails *config.LoggerTransactionDetails, key string, target []interface{}, cas gocb.Cas, ttlTimer *time.Duration) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- ReplaceDocumentWithCas, Key: ", key)
	upOps := gocb.ReplaceOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), Cas: cas}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			upOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			upOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REPLACE, config.GetKtab(key), serviceProviderID})
	ReplaceResponse, err := CollectionObj.Replace(key, target[0], &upOps)
	args := []string{config.CB_RESP, config.REPLACE, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in ReplaceDocumentWithCas")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- ReplaceDocumentWithCas")
	return ReplaceResponse.Cas(), nil
}

//DeletetDocument - Upsert a document into couchbase with provided key. This API tries to insert a doc into couchbase with provided key. If the key already exists, then the doc is updated with the doc passed as param to this API. Other wise it will insert a new doc. This API returns CAS of the doc after updation and error if occurred.
func DeletetDocument(transDetails *config.LoggerTransactionDetails, key string) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- DeletetDocument, Key: ", key)
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REMOVE, config.GetKtab(key), serviceProviderID})
	RemoveResponse, err := CollectionObj.Remove(key, &gocb.RemoveOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.REMOVE, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in DeletetDocument")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeletetDocument")
	return RemoveResponse.Cas(), nil
}

//GetSubDocuments - API fetches all the sub docs from the paths provided in subDocPaths array from the doc with key provided and fills the structures provided in target array. It return CAS of LookupIn op and error. This API returns midway if unmarshalling any of the subdoc into provided structures fails with corresponding error and 0 cas
func GetSubDocuments(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments, Target field list size mismatch: ", err)
		return 0, err
	}
	ops := make([]gocb.LookupInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.GetSpec(subDocPaths[i], nil)
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID})

	multiGetResult, err := CollectionObj.LookupIn(key, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments, Error in fetching sub docs: ", err)
		return 0, err
	}
	for i := 0; i < pathsArrayLen; i++ {
		decodeError := multiGetResult.ContentAt(uint(i), target[i])
		if decodeError != nil {
			mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments, Error decoding the value of subdoc - ", subDocPaths[i], " error:", decodeError)
			return 0, decodeError
		}
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments")
	return multiGetResult.Cas(), nil
}

//GetSubDocumentsWithFetchStatus - API fetches all the sub docs from the paths provided in subDocPaths array from the doc with key provided and fills the structures provided in target array. It returns array of boolean which represent fetch success/failure of corresponsing subDocs, CAS of LookupIn op and error. This API Does not return midway if unmarshalling any of the subdoc into provided structures fails.
func GetSubDocumentsWithFetchStatus(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}) (*[]bool, gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetSubDocumentsWithFetchStatus, Doc ID: ", key, " Paths: ", subDocPaths)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithFetchStatus, Target field list size mismatch: ", err)
		return nil, 0, err
	}
	ops := make([]gocb.LookupInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.GetSpec(subDocPaths[i], nil)
	}

	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID})
	multiGetResult, err := CollectionObj.LookupIn(key, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithFetchStatus, Error in fetching sub docs: ", err)
		return nil, 0, err
	}
	resultArr := make([]bool, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		decodeError := multiGetResult.ContentAt(uint(i), target[i])
		if decodeError != nil {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "decoding the value of subdoc - ", subDocPaths[i], " error:", decodeError)
			continue
		}
		resultArr[i] = true
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithFetchStatus")
	return &resultArr, multiGetResult.Cas(), nil
}

//UpsertSubDocumentsWithCas - This Func updates the specified subdoc with values provided in target. This also take CAS which will be passed while upserting to couchbase. This returns any error that is received
func UpsertSubDocumentsWithCas(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, cas gocb.Cas, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocumentsWithCas, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocumentsWithCas, Error: ", err)
		return err
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.UpsertSpec(subDocPaths[i], target[i], &gocb.UpsertSpecOptions{CreatePath: true})
	}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), Cas: cas}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocumentsWithCas, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocumentsWithCas")
	return nil
}

//UpsertSubDocuments - This Func updates the specified subdoc with values provided in target. This returns any error that is received
func UpsertSubDocuments(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocuments, Error: ", err)
		return err
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.UpsertSpec(subDocPaths[i], target[i], &gocb.UpsertSpecOptions{CreatePath: true})
	}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocuments, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocuments")
	return nil
}

//UpsertSubDocuments - This Func updates the specified subdoc with values provided in target, also creates parent document if it doesnt exist.This returns any error that is received
func UpsertSubDocumentsAutoCreateDoc(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocuments, Error: ", err)
		return err
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.UpsertSpec(subDocPaths[i], target[i], &gocb.UpsertSpecOptions{CreatePath: true})
	}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), StoreSemantic: gocb.StoreSemanticsUpsert}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocuments, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocuments")
	return nil
}

func InsertSubDocuments(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- InsertSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- InsertSubDocuments, Error: ", err)
		return err
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.InsertSpec(subDocPaths[i], target[i], &gocb.InsertSpecOptions{CreatePath: true})
	}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.INSERT_SUB_DOC, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.INSERT_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- InsertSubDocuments, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- InsertSubDocuments")
	return nil
}

//DeleteSubDocuments - This Func updates the specified subdoc with values provided in target. This returns any error that is received
func DeleteSubDocumentsWithTTL(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- DeleteSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths)
	pathsArrayLen := len(subDocPaths)
	if pathsArrayLen == 0 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments, empty subDocPaths")
		return nil
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.RemoveSpec(subDocPaths[i], nil)
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REMOVE_SUB_DOC, config.GetKtab(key), serviceProviderID})
	mutateInOpt := &gocb.MutateInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			mutateInOpt.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second

		} else {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry :", *ttlTimer)
			mutateInOpt.Expiry = *ttlTimer
		}
	}
	_, err := CollectionObj.MutateIn(key, ops, mutateInOpt)
	args := []string{config.CB_RESP, config.REMOVE_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments")
	return nil
}

//CheckDocumentExists - returns if the document exists or not for given key
func CheckDocumentExists(transDetails *config.LoggerTransactionDetails, docKey string) bool {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- CheckDocumentExists")
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(docKey)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.EXISTS, config.GetKtab(docKey), serviceProviderID})
	result, err := CollectionObj.Exists(docKey, &gocb.ExistsOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.EXISTS, config.GetKtab(docKey), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CheckDocumentExists, error in checking doc exists :", err)
		return false
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CheckDocumentExists : ", result.Exists())
	return result.Exists()
}

//CheckSubDocumentExists - returns if the sub document exists or not for given key
func CheckSubDocExists(transDetails *config.LoggerTransactionDetails, docKey string, subDocPaths []string, multiExists []bool) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- CheckSubDocExists")
	pathsArrayLen := len(subDocPaths)
	if pathsArrayLen == 0 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CheckSubDocExists, empty subDocPaths")
		return nil
	}
	ops := make([]gocb.LookupInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.ExistsSpec(subDocPaths[i], nil)
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(docKey)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.EXISTS_SUB_DOC, config.GetKtab(docKey), serviceProviderID})
	multiLookupResult, err := CollectionObj.LookupIn(docKey, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.EXISTS_SUB_DOC, config.GetKtab(docKey), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CheckSubDocExists, Error in LookupIn document: ", err)
		return err
	}

	for i := 0; i < pathsArrayLen; i++ {
		multiExists[i] = multiLookupResult.Exists(uint(i))
	}

	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CheckSubDocExists")
	return nil

}

func GetSubscriberKey(transDetails *config.LoggerTransactionDetails, supi string) string {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter GetSubscriberKey, subscriberId : ", supi)
	var subDocKey string
	subID, supiType := extractSubscriberKey(transDetails, supi)
	if supiType == SUPI_NO_OP {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "could not extract subscriber ID")
		return supi
	} else {
		if supiType == SUPI_IMSI {
			subDocKey = subID
		} else {
			subAlias := SubscriberAlias{}
			_, err := GetDocument(transDetails, SubAliasSid+subID, []interface{}{&subAlias})
			if err == nil {
				mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "fetched subscriber alias docuemtn : ", subAlias)
				subDocKey = subAlias.SubscriberKey
				if subAlias.Imsi != "" {
					subDocKey = subAlias.Imsi
				}
			} else {
				if errors.Is(err, gocb.ErrDocumentNotFound) {
					mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "subcriber alias document not found, check for subscriber main doc")
					if CheckDocumentExists(transDetails, SubInfoSid+subID) {
						mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "subscriber main doc exists with key :", subID)
						subDocKey = subID
					}
				} else {
					mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "error in fetching subscriber alias document : ", err)
				}
			}
		}
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit GetSubscriberKey, subscriberId : ", subDocKey)
	return subDocKey
}

func extractSubscriberKey(transDetails *config.LoggerTransactionDetails, subID string) (string, int) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter extractSubscriberKey, subscriberId : ", subID)
	var subDocKey string
	supiType := SUPI_NO_OP
	if strings.HasPrefix(subID, IMSI_PREFIX) {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "recieved subscriber with imsi- prefix")
		start := len(IMSI_PREFIX)
		end := len(subID)
		subDocKey = subID[start:end]
		if len(subDocKey) != 0 {
			supiType = SUPI_IMSI
		}
	} else if strings.HasPrefix(subID, NAI_PREFIX) {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "recieved subscriber with nai- prefix")
		start := len(NAI_PREFIX)
		end := strings.LastIndex(subID, "@")
		if end != -1 {
			subDocKey = subID[start:end]
			if len(subDocKey) != 0 {
				supiType = SUPI_NAI
			} else {
				mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "invalid NAI format")
			}
		}
	} else {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "unhandled subscriber identifier : ", subID)
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit extractSubscriberKey, subscriberKey : ", subDocKey, "supiType :", supiType)
	return subDocKey, supiType
}

func AppendSubDocumentsWithCas(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, cas gocb.Cas, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- AppendSubDocumentsWithCas, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	ops := []gocb.MutateInSpec{gocb.ArrayAppendSpec(subDocPaths[0], target, &gocb.ArrayAppendSpecOptions{HasMultiple: true})}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), Cas: cas}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.APPEND, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.APPEND, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- AppendSubDocumentsWithCas, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- AppendSubDocumentsWithCas")
	return nil
}

func AppendSubDocumentsReturnCas(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, cas *gocb.Cas, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- AppendSubDocumentsReturnCas, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	ops := []gocb.MutateInSpec{gocb.ArrayAppendSpec(subDocPaths[0], target, &gocb.ArrayAppendSpecOptions{HasMultiple: true})}
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), Cas: *cas}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.APPEND, config.GetKtab(key), serviceProviderID})
	updateResult, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.APPEND, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- AppendSubDocumentsReturnCas, Error in Mutating document: ", err)
		return err
	}
	*cas = updateResult.Cas()
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- AppendSubDocumentsReturnCas")
	return nil
}

//SetDocumentExpiry - set expiry for a document
func SetDocumentExpiry(transDetails *config.LoggerTransactionDetails, docKey string, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "traceid:", transDetails.TraceID, "Enter- SetDocumentExpiry")

	var ttl time.Duration
	touchOps := gocb.TouchOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "traceid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			ttl = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			ttl = *ttlTimer
		}
	}

	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(docKey)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.TOUCH, config.GetKtab(docKey), serviceProviderID})
	_, err := CollectionObj.Touch(docKey, ttl, &touchOps)
	args := []string{config.CB_RESP, config.TOUCH, config.GetKtab(docKey), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "traceid:", transDetails.TraceID, "Error - SetDocumentExpiry :", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "traceid:", transDetails.TraceID, "Exit - SetDocumentExpiry")
	return nil
}

//GetActiveReservation - func to get list of sub docs
func GetActiveReservation(subID string, arDoc *map[string]SessionStructure) error {
	mlog.MavLog(mlog.INFO, transID, "Enter- GetActiveReservation, subID: ", subID)
	docKey := SubBalSid + subID
	mlog.MavLog(mlog.INFO, transID, "document Key", docKey)

	ops := []gocb.LookupInSpec{gocb.GetSpec("activeReservations", &gocb.GetSpecOptions{})}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(docKey)
	config.UpdateCounterMetrics(&config.LoggerTransactionDetails{}, config.CB_REQ, []string{config.GET_SUB_DOC, config.GetKtab(docKey), serviceProviderID})
	result, err := CollectionObj.LookupIn(docKey, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET_SUB_DOC, config.GetKtab(docKey), serviceProviderID}
	config.UpdateRespMetrics(&config.LoggerTransactionDetails{}, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavLog(mlog.ERROR, transID, "Exit- GetActiveReservation, Error in fetching sub docs: ", err)
		return err
	}
	decodeError := result.ContentAt(0, arDoc)
	if decodeError != nil {
		mlog.MavLog(mlog.ERROR, transID, "Exit- GetSubDocuments, Error decoding the value of subdoc error:", decodeError)
		return decodeError
	}
	mlog.MavLog(mlog.INFO, transID, "Exit- GetSubDocuments activeReservation :", arDoc)
	return nil
}

//GetActiveReservation - func to get list of sub docs
func GetActiveReservationWithCas(transDetails *config.LoggerTransactionDetails, subID string, arDoc *map[string]SessionStructure) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetActiveReservation, subID: ", subID)
	docKey := SubBalSid + subID
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "document Key", docKey)

	ops := []gocb.LookupInSpec{gocb.GetSpec("activeReservations", &gocb.GetSpecOptions{})}

	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(docKey)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET_SUB_DOC, config.GetKtab(docKey), serviceProviderID})

	result, err := CollectionObj.LookupIn(docKey, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})

	args := []string{config.CB_RESP, config.GET_SUB_DOC, config.GetKtab(docKey), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetActiveReservation, Error in fetching sub docs: ", err)
		return 0, err
	}
	decodeError := result.ContentAt(0, arDoc)
	if decodeError != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments, Error decoding the value of subdoc error:", decodeError)
		return 0, decodeError
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocuments activeReservation :", arDoc)
	return result.Cas(), nil
}

func AddEscapeChar(input []byte) string {
	op := []byte("")
	op = append(op, '`')
	op = append(op, input...)
	op = append(op, '`')
	return string(op)
}

func GetServiceProviderList(transDetails *config.LoggerTransactionDetails) []string {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter - GetServiceProviderList")
	var spList []string
	spListMap := make(map[string]string)
	_, err := GetSubDocuments(transDetails, SpListKey, []string{SP_PATH}, []interface{}{&spListMap})
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Error in fetching the sp list document, error: ", err)
		mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetServiceProviderList")
		return spList
	}
	for key, _ := range spListMap {
		spList = append(spList, key)
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetServiceProviderList")
	return spList
}

func GetSubscriberLifeCycle(transDetails *config.LoggerTransactionDetails, serviceProviderID string, subLifeCycle *SubscriberLifeCycleInfo) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter - GetSubscriberLifeCycle")
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "serviceProviderID : ", serviceProviderID)
	key := SubLifeCycleSid + serviceProviderID + "::" + SUB_LIFE_CYCLE_LIST
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "lifeCycle Key: ", key)
	cas, err := GetDocument(transDetails, key, []interface{}{subLifeCycle})
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Error in fetching the subLifeCycleList document, error: ", err)
		mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetSubscriberLifeCycleList")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetServiceProviderList")
	return cas, nil
}

//GetCosInfo - return the cos definitions for a service provider ID
func GetCosInfo(transDetails *config.LoggerTransactionDetails, spId string, cosDefinition map[string]ClassOfService) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "Enter - GetCosInfo, spId: ", spId)
	key := CosInfoSid + spId + "::" + COS_INFO_SUFFIX
	_, err := GetSubDocuments(transDetails, key, []string{COS_INFO_PATH}, []interface{}{&cosDefinition})
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "Error in fetching sub doc path for key: ", key, " error: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "Exit - GetCosInfo")
	return nil
}

func UpdateSubscriberState(transDetails *config.LoggerTransactionDetails, subID string, state string, cas gocb.Cas) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "Enter - UpdateSubscriberState, subscriber: ", subID, " state: ", state)
	var key string = SubInfoSid + subID
	currTime := time.Now().Format("2006-01-02T15:04:05Z")
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "Exit - UpdateSubscriberState")
	return UpsertSubDocumentsWithCas(transDetails, key, []string{"state", "lastStatusChangeDate"}, []interface{}{state, currTime}, cas, nil)
}

//UpsertSubDocumentsReturnCas - This Func updates the specified subdoc with values provided in target. This also take CAS which will be passed while upserting to couchbase.   This returns any error that is received along with Cas
func UpsertSubDocumentsReturnCas(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, cas *gocb.Cas, ttlTimer *time.Duration) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocumentsReturnCas, Doc ID: ", key, " Paths: ", subDocPaths, " target: ", target)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- UpsertSubDocumentsReturnCas, Error: ", err)
		return err
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.UpsertSpec(subDocPaths[i], target[i], &gocb.UpsertSpecOptions{CreatePath: true})
	}
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID})
	muOps := gocb.MutateInOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout), Cas: *cas}
	if ttlTimer != nil {
		if ttlTimer.Minutes() > float64(60*24*30) {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Expiry more then 30 days:", *ttlTimer)
			muOps.Expiry = time.Duration(time.Now().Add(*ttlTimer).Unix()) * time.Second
		} else {
			muOps.Expiry = *ttlTimer
		}
	}
	startTime := config.GetTimeNs()
	updateResult, err := CollectionObj.MutateIn(key, ops, &muOps)
	args := []string{config.CB_RESP, config.UPSERT_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocumentsReturnCas, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- UpsertSubDocumentsReturnCas")
	*cas = updateResult.Cas()
	return nil
}

//GetDocumentInMultipleFormats - get a document that matches primaryKey from couchbase. This API returns CAS and error if occurred
func GetDocumentInMultipleFormats(transDetails *config.LoggerTransactionDetails, key string, target []interface{}) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetDocumentInMultipleFormats, Key: ", key)
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET, config.GetKtab(key), serviceProviderID})
	crdlResponse, err := CollectionObj.Get(key, &gocb.GetOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in GetDocumentInMultipleFormats")
		return 0, err
	}
	for _, targetValue := range target {
		err = crdlResponse.Content(targetValue)
		if err != nil {
			mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in GetDocumentInMultipleFormats: ", err)
			return 0, err
		}
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetDocumentInMultipleFormats")
	return crdlResponse.Cas(), nil
}

//GetLineDocument - fetch subscriber document
func GetLineDocument(transDetails *config.LoggerTransactionDetails, imsi string, target []interface{}) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetLineDocument")
	docKey := SubInfoSid + imsi
	_, err := GetDocument(transDetails, docKey, target)
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetLineDocument")
	return err
}

//GetLineBalDocument - fetch subscriber Line-Bal document
func GetLineBalDocument(transDetails *config.LoggerTransactionDetails, imsi string, target []interface{}) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetLineBalDocument")
	docKey := SubBalSid + imsi
	_, err := GetDocument(transDetails, docKey, target)
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetLineBalDocument")
	return err
}

//GetAliasDocument - fetch subscriber alias document
func GetAliasDocument(transDetails *config.LoggerTransactionDetails, imsi string, target []interface{}) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetAliasDocument")
	docKey := SubAliasSid + imsi
	_, err := GetDocument(transDetails, docKey, target)
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetAliasDocument")
	return err
}

func GetServiceProviderCurrencyUsage(transDetails *config.LoggerTransactionDetails, key string, target []interface{}) error {
	mlog.MavLog(mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter-GetServiceProviderCurrencyUsage, Key: ", key)
	errRetryCount := 0
retry:
	_, err := GetSubDocuments(transDetails, SpInfoSid+key,
		[]string{"serviceProviderConfiguration.currencyUsage"}, target)
	if err != nil {
		if errors.Is(err, gocb.ErrTemporaryFailure) {
			if errRetryCount < 2 {
				mlog.MavLog(mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "temporary couchbase error, Fetch service provider currency usage config : ", err)
				errRetryCount++
				goto retry
			}
		}
	}
	mlog.MavLog(mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit-GetServiceProviderCurrencyUsage, Key: ", key)
	return err
}

func GetServiceProviderFeatureLevelConfig(transDetails *config.LoggerTransactionDetails, key string, target map[string]map[string]interface{}) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetServiceProviderFeatureLevelConfig, Key: ", key)
	errRetryCount := 0
retry:
	_, err := GetSubDocuments(transDetails, SpInfoSid+key,
		[]string{"serviceProviderConfiguration.featureLevelConfiguration"}, []interface{}{&target})
	if err != nil {
		if errors.Is(err, gocb.ErrTemporaryFailure) {
			if errRetryCount < 2 {
				mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "temporary couchbase error, Fetch service provider feature config : ", err)
				errRetryCount++
				goto retry
			}
		}
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetServiceProviderFeatureLevelConfig, Key: ", key)
	return err
}

//DeleteSubDocuments - This Func updates the specified subdoc with values provided in target. This returns any error that is received
func DeleteSubDocuments(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- DeleteSubDocuments, Doc ID: ", key, " Paths: ", subDocPaths)
	pathsArrayLen := len(subDocPaths)
	if pathsArrayLen == 0 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments, empty subDocPaths")
		return nil
	}
	ops := make([]gocb.MutateInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		ops[i] = gocb.RemoveSpec(subDocPaths[i], nil)
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REMOVE_SUB_DOC, config.GetKtab(key), serviceProviderID})
	_, err := CollectionObj.MutateIn(key, ops, &gocb.MutateInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.REMOVE_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments, Error in Mutating document: ", err)
		return err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteSubDocuments")
	return nil
}

//DeleteDocument
func DeleteDocument(transDetails *config.LoggerTransactionDetails, key string) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- DeleteDocument, Key: ", key)
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REMOVE, config.GetKtab(key), serviceProviderID})
	RemoveResponse, err := CollectionObj.Remove(key, &gocb.RemoveOptions{Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.REMOVE, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in DeleteDocument")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteDocument")
	return RemoveResponse.Cas(), nil
}

func DeleteDocumentWithCas(transDetails *config.LoggerTransactionDetails, key string, cas gocb.Cas) (gocb.Cas, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- DeleteDocument, Key: ", key)
	startTime := config.GetTimeNs()
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.REMOVE, config.GetKtab(key), ""})
	RemoveResponse, err := CollectionObj.Remove(key, &gocb.RemoveOptions{Cas: cas, Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.REMOVE, config.GetKtab(key), ""}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- Error in DeleteDocument")
		return 0, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- DeleteDocument")
	return RemoveResponse.Cas(), nil
}

func GetSpendingLimitPolicies(transDetails *config.LoggerTransactionDetails, serviceProviderID string, policies map[string]*SpendingLimitPolicyDetail) error {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter - GetSpendingLimitPolicies")
	errRetryCount := 0

retry:
	_, err := GetSubDocuments(transDetails, SpendingLimitPolicySid+serviceProviderID+"::"+SPENDING_LIMIT_POLICY_LIST, []string{SPENDING_LIMIT_POLICY_PATH}, []interface{}{&policies})
	if err != nil {
		if errors.Is(err, gocb.ErrTemporaryFailure) {
			if errRetryCount < 2 {
				mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "temporary couchbase error, Fetch subscriber balance doc : ", err)
				errRetryCount++
				goto retry
			} else {
				mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetSpendingLimitPolicies, error in Fetch spending limit policy doc, Max retries done")
			}
		} else {
			mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetSpendingLimitPolicies, permanent couchbase failure, fetch spending limit policy doc : ", err)
		}
	}

	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit - GetSpendingLimitPolicies")
	return err
}

//UpdateEndDate -  usageEndDate should be last second of the day for renewal to happen as per the local time
func UpdateEndDate(usageEndDate string, transDetails *config.LoggerTransactionDetails) string {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Enter- updateEndDate")
	t, _ := time.Parse(time.RFC3339, usageEndDate)
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "usageEndDate: ", usageEndDate)
	localTime := t.Local()
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "LocalTime corresponding to usageEndDate: ", localTime)
	hour, min, sec := localTime.Clock()
	if hour != 23 || min != 59 || sec != 59 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "LocalTime is not 23.59.59hrs")
		localTimeWithAddition := localTime.AddDate(0, 0, 1)
		roundedLocalTime := time.Date(localTimeWithAddition.Year(), localTimeWithAddition.Month(), localTimeWithAddition.Day(), 0, 0, 0, 0, localTimeWithAddition.Location())
		lastSecondEndTime := roundedLocalTime.Add(-1 * time.Second)
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "LocalTime rounded to 23.59.59hrs: ", lastSecondEndTime)
		utcTime := lastSecondEndTime.UTC().Format("2006-01-02T15:04:05Z")
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "UTC time corresponding to local 23.59.59hrs: ", utcTime)
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Exit - updateEndDate")
		return utcTime
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Exit - updateEndDate", usageEndDate)
	return usageEndDate
}

//GetNorStartDate - This function will return the usageStartDate in 00hrs if its not 00hrs
func GetNorStartDate(usageStartDate string, transDetails *config.LoggerTransactionDetails) string {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Enter- GetNorStartDate", usageStartDate)
	t, _ := time.Parse(time.RFC3339, usageStartDate)
	localTime := t.Local()
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "LocalTime corresponding to usageStartDate: ", localTime)
	hour, min, sec := localTime.Clock()
	if hour != 00 || min != 00 || sec != 00 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "LocalTime is not 00hrs")
		roundedLocalTime := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, localTime.Location())
		mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "LocalTime rounded to 00hrs: ", roundedLocalTime)
		utcTime := roundedLocalTime.UTC().Format("2006-01-02T15:04:05Z")
		mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "UTC time corresponding to local 00hrs: ", utcTime)
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Exit - GetNorStartDate", utcTime)
		return utcTime
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "Exit - GetNorStartDate", usageStartDate)
	return usageStartDate
}

func GetServiceProviderID(key string) string {
	keySplit := strings.Split(key, "::")

	if len(keySplit) > 1 && len(keySplit[1]) > 2 && keySplit[1][0:2] == "SP" {
		return keySplit[1]
	}
	return "NA"
}

// AddDate - function to add month,day, year and day to date string
func AddDate(datestr string, unit string, delta int) string {
	// datestr = "2020-04-01T00:00:00.000Z"
	// dateArr := strings.Split(datestr, "T")
	// datestr = dateArr[0] + "T23:59:59Z"
	t, err := time.Parse(time.RFC3339, datestr)
	if err != nil {
		mlog.MavLog(mlog.ERROR, "", "Error - AddDate")
		return ""
	}
	switch unit {
	case "DAY":
		return t.AddDate(0, 0, delta).Format(time.RFC3339)
	case "WEEK":
		return t.AddDate(0, 0, 7*delta).Format(time.RFC3339)
	case "MON":
		return t.AddDate(0, delta, 0).Format(time.RFC3339)
	case "YEAR":
		return t.AddDate(delta, 0, 0).Format(time.RFC3339)
	default:
		mlog.MavLog(mlog.ERROR, "", "Error - AddDate - default")
		return ""
	}
}

func GetSubscribeProrationFeeOrQuota(transDetails *config.LoggerTransactionDetails, prorateOnCommencement string, isPkgProduct int32,
	usageStartDate string, usageEndDate string, renewalPolicy string, renewalCycleDay *int32, chargeCycle int64,
	validity *int32, validityUnit string, renewalFee float64, quota int64) (float64, float64, int64) {

	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "GetSubscribeProrationFeeOrQuota !!! prorateOnCommencement : '", prorateOnCommencement, "' isPkgProduct : ", isPkgProduct,
		"usageStartDate : ", usageStartDate, "usageEndDate : ", usageEndDate, "renewalPolicy : ", renewalPolicy,
		"validity : ", validity, "validityUnit : ", validityUnit, "renewalFee : ", renewalFee, "quota : ", quota)

	prorationFactor := 1.0
	proratedRenewalFee := renewalFee
	proratedQuota := quota
	applyproration := false

	if validityUnit == "WEEK" && renewalCycleDay != nil {
		applyproration = true
	}

	if applyproration == false && (renewalCycleDay == nil || *renewalCycleDay == 0) && chargeCycle == 0 {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "prorationFactor : 1 ")
		return prorationFactor, proratedRenewalFee, proratedQuota
	}

	if (validity != nil && *validity == 1 && validityUnit == "DAY") || prorateOnCommencement == "" || renewalPolicy == "ONEOFF" ||
		prorateOnCommencement == "BILL_CURRENT_CYCLE_FULLY" || (prorateOnCommencement == "PRORATE_PRICE_ONLY" && isPkgProduct == PRODUCT) {
		mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "prorationFactor: 1 ")
		return prorationFactor, proratedRenewalFee, proratedQuota
	}

	regularUsageEndDate := AddDate(usageStartDate, validityUnit, int(*validity))
	iRegularUsageEndDate, _ := time.Parse(time.RFC3339, regularUsageEndDate)
	iUsageEndDate, _ := time.Parse(time.RFC3339, usageEndDate)
	iUsageStartDate, _ := time.Parse(time.RFC3339, usageStartDate)
	regularUsagePeriod := int(iRegularUsageEndDate.Sub(iUsageStartDate).Hours()) / 24
	actualUsagePeriod := (iUsageEndDate.Sub(iUsageStartDate).Hours()) / 24
	actualUsagePeriod = math.Ceil(actualUsagePeriod)

	prorationFactor = float64(actualUsagePeriod) / float64(regularUsagePeriod)

	if prorationFactor > 1 {
		prorationFactor = 1
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "Something wrong with prorationFactor calculation ... need to check ...")
	}

	if renewalFee > 0.0 {
		proratedRenewalFee = float64(int64(prorationFactor * renewalFee))
	} else {
		proratedQuota = int64(prorationFactor * float64(quota))
	}

	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "GetSubscribeProrationFeeOrQuota !!! regularUsageEndDate : ", regularUsageEndDate, "regularUsagePeriod :", regularUsagePeriod, "actualUsagePeriod :", actualUsagePeriod,
		"prorationFactor : ", prorationFactor, "proratedRenewalFee : ", proratedRenewalFee, "proratedQuota : ", proratedQuota)
	return prorationFactor, proratedRenewalFee, proratedQuota
}

func GetUnSubscribeProrationFeeOrQuota(transDetails *config.LoggerTransactionDetails, prorateOnTermination string, subscriptionType SubscriptionType, isPkgProduct int32,
	usageStartDate string, usageEndDate string, unsubscribeDate string, renewalPolicy string,
	validity *int32, validityUnit string, renewalFeeCharged float64, quota int64, includeTerminationDay bool) (float64, float64, int64) {

	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "GetUnSubscribeProrationFeeOrQuota !!! prorateOnTermination : '", prorateOnTermination, "' usageStartDate : ", usageStartDate, "usageEndDate : ", usageEndDate, "unsubscribeDate : ", unsubscribeDate,
		"SubscriptionType", subscriptionType, "renewalPolicy : ", renewalPolicy, "validity : ", validity, "validityUnit : ", validityUnit, "renewalFeeCharged : ", renewalFeeCharged, "quota :", quota)

	prorationFactor := 0.0
	renewalFeeRefund := 0.0
	nonUsagePeriodQuota := int64(0)

	if renewalFeeCharged > 0 {
		prorationFactor = 0.0
	} else {
		prorationFactor = 1.0
	}

	if (validity != nil && *validity == 1 && validityUnit == "DAY") || prorateOnTermination == "" || (prorateOnTermination == "PRORATE" && subscriptionType == SubscriptionType_ADDON && isPkgProduct == PRODUCT) || prorateOnTermination == "BILL_CURRENT_CYCLE_FULLY" || (prorateOnTermination == "PRORATE_PRICE_ONLY" && isPkgProduct == PRODUCT) {
		return prorationFactor, renewalFeeRefund, nonUsagePeriodQuota
	}
	iUsageStartDate, _ := time.Parse(time.RFC3339, usageStartDate)
	iUsageEndDate, _ := time.Parse(time.RFC3339, usageEndDate)
	iUsageEndDate = iUsageEndDate.Add(1 * time.Second)

	regularUsagePeriod := (iUsageEndDate.Sub(iUsageStartDate).Hours()) / 24
	regularUsagePeriod = math.Ceil(regularUsagePeriod)
	iUnsubscribeDate, _ := time.Parse(time.RFC3339, unsubscribeDate)
	nonUsagePeriod := int(iUsageEndDate.Sub(iUnsubscribeDate).Hours()) / 24
	if includeTerminationDay == true {
		nonUsagePeriod = nonUsagePeriod + 1
	}

	prorationFactor = float64(nonUsagePeriod) / float64(regularUsagePeriod)

	if prorationFactor > 1 {
		prorationFactor = 1
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, "Something wrong with prorationFactor calculation ... need to check ...")
	}

	if renewalFeeCharged > 0.0 {
		renewalFeeRefund = prorationFactor * renewalFeeCharged
	} else {
		nonUsagePeriodQuota = int64(prorationFactor * float64(quota))
	}

	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, "GetUnSubscribeProrationFeeOrQuota !!! regularUsagePeriod : ", regularUsagePeriod, "nonUsagePeriod : ", nonUsagePeriod, "prorationFactor : ", prorationFactor,
		"renewalFeeRefund : ", renewalFeeRefund, "nonUsagePeriodQuota :", nonUsagePeriodQuota)
	return prorationFactor, renewalFeeRefund, nonUsagePeriodQuota
}

func GetSubscriptionStateStringToInt(stype string) int32 {

	switch stype {
	case "A":
		return int32(SubscriptionState_Active)
	case "D":
		return int32(SubscriptionState_Deactivated)
	case "S":
		return int32(SubscriptionState_Suspended)
	case "P":
		return int32(SubscriptionState_PendingActivation)
	default:
		return -1
	}
}

func GetSubscriberCounterListByID(transDetails *config.LoggerTransactionDetails, spID string) (*SubscriberCountersDefinations, error) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, spID, "Enter- GetSubscriberCounterListByID")
	docKey := SubCounterSid + spID + SEPARATOR + COUNTER_LIST
	subscriberCounters := &SubscriberCountersDefinations{}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, spID, "Enter- GetSubscriberCounterListByID-", docKey)
	_, err := GetDocument(transDetails, docKey, []interface{}{subscriberCounters})
	if err != nil {
		mlog.MavTrace(transDetails.TraceFlag, mlog.ERROR, transDetails.TransID, spID, "Error - GetSubscriberCounterListByID :", docKey, err)
		return nil, err
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.INFO, transDetails.TransID, spID, "Exit- GetSubscriberCounterListByID")
	return subscriberCounters, err

}

func isWalletExpired(expiryTime string) bool {
	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	var walletExpiry int64
	timeVal, _ := time.Parse(time.RFC3339, expiryTime)
	walletExpiry = timeVal.Unix()
	if walletExpiry < currentTime.Unix() {
		return true
	}
	return false
}

func getFeeTypeBasedOnWalletForType(walletForType string) FeeType {
	feeType := FeeType_Other
	if walletForType == "walletForPurchaseFee" {
		feeType = FeeType_Purchanse
	} else if walletForType == "walletForRenewalFee" {
		feeType = FeeType_Renewal
	}
	return feeType
}

func getFeeStringBasedOnWalletForType(walletForType string) string {
	feeType := "Other"
	if walletForType == "walletForPurchaseFee" {
		feeType = "PurchaseFee"
	} else if walletForType == "walletForRenewalFee" {
		feeType = "RenewalFee"
	}
	return feeType
}

func PurchaseOrRenewalFeeCreditUpdate(transID string, totalAmountToCredited float64, balances *map[string]*SubscriberAccountParams, debitedWalletDetails *[]*DebitedWalletDetails, purchaseOrRenewalType string) bool {
	mlog.MavLog(mlog.INFO, transID, "Enter- PurchaseOrRenewalFeeCredit, Purchase/Renewal type", purchaseOrRenewalType)
	totalFee := 0.0
	for _, BalIdFeeDetails := range *debitedWalletDetails {
		totalFee += BalIdFeeDetails.Fee
	}
	if totalFee != totalAmountToCredited {
		mlog.MavLog(mlog.ERROR, transID, "PurchaseOrRenewalFeeCredit :: something wrong.. input amount and LineDoc updated amount not matching !!")
		return false
	}
	for _, BalIdFeeDetails := range *debitedWalletDetails {
		totalFee += BalIdFeeDetails.Fee
		(*balances)[BalIdFeeDetails.BalanceId].CreditStatus.Credit += BalIdFeeDetails.Fee
	}
	mlog.MavLog(mlog.INFO, transID, "Exit- PurchaseOrRenewalFeeCredit")
	return true
}

func GetWalletChangeDetailsForEdrForRefund(transID string, cashBalance *map[string]*SubscriberAccountParams, debitedWalletDetails *[]*DebitedWalletDetails, walletChangeDetailForEdr *map[string]*WalletBalance, instanceId string, packageId string, FeeType string, modifyTempData bool, walletFeeTypeRecord []string, feeRefund float64, subscriberType int32) bool {

	mlog.MavLog(mlog.DEBUG, transID, "Enter- GetWalletChangeDetailsForEdrForRefund", feeRefund, subscriberType)

	if subscriberType == POSTPAID {
		walletChangeDetails := new(WalletBalance)
		walletChangeDetails.InstanceID = instanceId
		walletChangeDetails.PackageID = packageId
		walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType(FeeType)
		walletChangeDetails.ChangedBalance = &feeRefund
		(*walletChangeDetailForEdr)[FeeType+instanceId+packageId+strconv.Itoa(int(feeRefund))] = walletChangeDetails
		mlog.MavLog(mlog.DEBUG, transID, "GetWalletChangeDetailsForEdrForRefund for postpaid")
		return true
	}

	for index, BalIdFeeDetails := range *debitedWalletDetails {
		if cashBal, ok := (*cashBalance)[BalIdFeeDetails.BalanceId]; ok {
			walletChangeDetails := new(WalletBalance)
			walletChangeDetails.BalanceName = cashBal.BalanceId
			walletChangeDetails.BalanceType = cashBal.BalanceType
			walletChangeDetails.EffectiveExpiryDate = cashBal.EffectiveExpiryDate
			walletChangeDetails.EffectiveStartDate = cashBal.EffectiveStartDate
			keyFeeType := ""
			if len(walletFeeTypeRecord) == 0 {
				walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType("Other")
				keyFeeType = FeeType
			} else {
				walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType(walletFeeTypeRecord[index])
				keyFeeType = walletFeeTypeRecord[index]
			}
			walletChangeDetails.InstanceID = instanceId
			walletChangeDetails.PackageID = packageId
			prevBal := cashBal.CreditStatus.Credit
			walletChangeDetails.PreviousBalance = &prevBal
			changedBal := BalIdFeeDetails.Fee
			walletChangeDetails.ChangedBalance = &changedBal
			currBal := cashBal.CreditStatus.Credit + BalIdFeeDetails.Fee
			walletChangeDetails.CurrentBalance = &currBal
			if modifyTempData { // Updating only on the duplicate cashBalance struct to get correct EDR calculation for both purchaseFee & RenewalFee refund changes
				cashBal.CreditStatus.Credit = currBal
			}
			(*walletChangeDetailForEdr)[BalIdFeeDetails.BalanceId+keyFeeType+instanceId+packageId+strconv.Itoa(int(changedBal))] = walletChangeDetails
		}
	}

	mlog.MavLog(mlog.DEBUG, transID, "Exit- GetWalletChangeDetailsForEdrForRefund")
	return true
}

func addFeeTypeDetailsForEDRForSubscriberCreate(transID string, balId string, cashBal CashBalancesCDB, walletChangeDetailForEdr *map[string]*WalletBalance) {
	previousBal := 0.0
	_, ok := (*walletChangeDetailForEdr)[balId]
	if ok {
		mlog.MavLog(mlog.DEBUG, transID, "addFeeTypeDetailsForEDRForSubscriberCreate entry already exists in EDR for balId", balId)
	} else {
		changedBal := cashBal.CreditStatus.Credit
		walletChangeDetails := new(WalletBalance)
		walletChangeDetails.BalanceName = cashBal.BalanceId
		walletChangeDetails.BalanceType = cashBal.BalanceType
		walletChangeDetails.EffectiveExpiryDate = cashBal.EffectiveExpiryDate
		walletChangeDetails.EffectiveStartDate = cashBal.EffectiveStartDate
		walletChangeDetails.ChangedBalance = &changedBal
		walletChangeDetails.CurrentBalance = &changedBal
		walletChangeDetails.PreviousBalance = &previousBal
		walletChangeDetails.FeeType = "Recharge"
		(*walletChangeDetailForEdr)[balId] = walletChangeDetails
		mlog.MavLog(mlog.DEBUG, transID, "addFeeTypeDetailsForEDRForSubscriberCreate for balId", balId, "changedBal", changedBal)

	}
}

func CheckSubBalanceForPurchaseBasedOnBalanceID(transID string, cashBalance *map[string]*SubscriberAccountParams, Fee float64, FeeType string, balId string, previousWalletCreditTotal float64, balanceIdCreditMap *map[string]float64, walletChangeDetailForEdr *map[string]*WalletBalance, instanceId string, packageId string, subscriberType int32) (bool, float64) {
	mlog.MavLog(mlog.DEBUG, transID, "Enter- CheckSubBalanceForPurchaseBasedOnBalanceID , balId", balId, "InstanceId", instanceId, "packageId", packageId)
	mlog.MavLog(mlog.INFO, transID, "previousWalletCreditTotal", previousWalletCreditTotal, "totalFee", Fee, "FeeType", FeeType)
	totalCreditElems := 0
	tempChangedBal := 0.0
	if cashBal, ok := (*cashBalance)[balId]; ok {
		if isWalletExpired(cashBal.EffectiveExpiryDate) { //R.FRN-11293-US4
			mlog.MavLog(mlog.DEBUG, transID, "wallet with BalId", balId, "expired, check other passed wallets")
			return false, 0.0
		}
		walletChangeDetails := new(WalletBalance)
		walletChangeDetails.BalanceName = cashBal.BalanceId
		walletChangeDetails.BalanceType = cashBal.BalanceType
		walletChangeDetails.EffectiveExpiryDate = cashBal.EffectiveExpiryDate
		walletChangeDetails.EffectiveStartDate = cashBal.EffectiveStartDate
		walletChangeDetails.InstanceID = instanceId
		walletChangeDetails.PackageID = packageId
		walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType(FeeType)

		credit := (cashBal.CreditStatus.Credit + previousWalletCreditTotal) - Fee
		if (credit >= float64(config.CommonDynConf.Config.MinBalanceThreshold)) || subscriberType == POSTPAID {
			mlog.MavLog(mlog.INFO, transID, "MinBalanceThreshold", float64(config.CommonDynConf.Config.MinBalanceThreshold), "final credit", credit)
			for balanceID, CreditValue := range *balanceIdCreditMap {
				mlog.MavLog(mlog.INFO, transID, "Inside balanceIdCreditMap loop, balId", balanceID, "CreditValue", CreditValue)
				(*cashBalance)[balanceID].CreditStatus.Credit = (*cashBalance)[balanceID].CreditStatus.Credit - CreditValue
				(*cashBalance)[balanceID].CreditStatus.LastUpdateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
				(*cashBalance)[balanceID].FeeTypeDeducted = getFeeTypeBasedOnWalletForType(FeeType)
				totalCreditElems++
			}

			if previousWalletCreditTotal > 0 {
				previousBal := (*cashBalance)[balId].CreditStatus.Credit
				walletChangeDetails.PreviousBalance = &previousBal
				(*cashBalance)[balId].CreditStatus.Credit = (*cashBalance)[balId].CreditStatus.Credit + previousWalletCreditTotal - Fee
				(*balanceIdCreditMap)[balId] = Fee - previousWalletCreditTotal
				mlog.MavLog(mlog.DEBUG, transID, "balanceIdCreditMap len =  1", totalCreditElems)
				mlog.MavLog(mlog.DEBUG, transID, "balanceIdCreditMap len =  1", "BalID", balId, "Latest Balance", (*cashBalance)[balId].CreditStatus.Credit, "UsedBalance", (*balanceIdCreditMap)[balId])
				changeBal := Fee - previousWalletCreditTotal
				currentBal := previousBal - changeBal
				tempChangedBal -= changeBal
				walletChangeDetails.ChangedBalance = &tempChangedBal
				walletChangeDetails.CurrentBalance = &currentBal

			} else {
				previousBal := (*cashBalance)[balId].CreditStatus.Credit
				walletChangeDetails.PreviousBalance = &previousBal
				(*cashBalance)[balId].CreditStatus.Credit = credit
				(*balanceIdCreditMap)[balId] = Fee
				changeBal := Fee
				currentBal := previousBal - changeBal
				tempChangedBal -= changeBal
				walletChangeDetails.ChangedBalance = &tempChangedBal
				walletChangeDetails.CurrentBalance = &currentBal
			}

			(*cashBalance)[balId].CreditStatus.LastUpdateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			mlog.MavLog(mlog.INFO, transID, credit, "Exit- CheckSubscriberBalanceForPurchase  - true")
			(*cashBalance)[balId].FeeTypeDeducted = getFeeTypeBasedOnWalletForType(FeeType)

			mlog.MavLog(mlog.INFO, transID, credit, "walletChangeDetailForEdr true case KEY", balId+FeeType+instanceId+packageId)
			(*walletChangeDetailForEdr)[balId+FeeType+instanceId+packageId] = walletChangeDetails //update this map with BalanceId used and balance details that will be used in EDR

			return true, 0.0
		}
		mlog.MavLog(mlog.INFO, transID, credit, "Exit- CheckSubscriberBalanceForPurchase  - false")
		(*balanceIdCreditMap)[balId] = (*cashBalance)[balId].CreditStatus.Credit - float64(config.CommonDynConf.Config.MinBalanceThreshold)

		previousBal := (*cashBalance)[balId].CreditStatus.Credit
		walletChangeDetails.PreviousBalance = &previousBal
		changeBal := (*balanceIdCreditMap)[balId]
		currentBal := previousBal - changeBal
		tempChangedBal -= changeBal
		walletChangeDetails.ChangedBalance = &tempChangedBal
		walletChangeDetails.CurrentBalance = &currentBal
		(*walletChangeDetailForEdr)[balId+FeeType+instanceId+packageId+strconv.Itoa(int(tempChangedBal))] = walletChangeDetails
		mlog.MavLog(mlog.INFO, transID, credit, "walletChangeDetailForEdr false case KEY", balId+FeeType+instanceId+packageId)

		return false, (*cashBalance)[balId].CreditStatus.Credit - float64(config.CommonDynConf.Config.MinBalanceThreshold)
	}
	return false, 0.0
}

func CheckAndUpdateSubBalanceForPurchaseBasedOnBalanceIDCDB(transID string, cashBalance *map[string]*CashBalancesCDB, Fee float64, FeeType string, balId string, previousWalletCreditTotal float64, balanceIdCreditMap *map[string]float64, walletChangeDetailForEdr *map[string]*WalletBalance, instanceId string, packageId string, origWalletChangeDetailForEdr *map[string]*WalletBalance, subscriberType int32) (bool, float64) {
	mlog.MavLog(mlog.DEBUG, transID, "Enter- CheckAndUpdateSubBalanceForPurchaseBasedOnBalanceIDCDB , balId", balId, "instanceId", instanceId, "packageId", packageId)
	mlog.MavLog(mlog.INFO, transID, "previousWalletCreditTotal", previousWalletCreditTotal, "totalFee", Fee)
	totalCreditElems := 0
	tempChangedBal := 0.0
	if cashBal, ok := (*cashBalance)[balId]; ok {
		if isWalletExpired(cashBal.EffectiveExpiryDate) { //R.FRN-11293-US4
			mlog.MavLog(mlog.DEBUG, transID, "wallet with BalId", balId, "expired, check other passed wallets")
			return false, 0.0
		}
		//For Subscriber Create flow, EDR should have entry with FeeType as "Recharge" with changedBalance and currentBalance as original balance without debitted details.
		addFeeTypeDetailsForEDRForSubscriberCreate(transID, balId, *cashBal, origWalletChangeDetailForEdr)
		walletChangeDetails := new(WalletBalance)
		walletChangeDetails.BalanceName = cashBal.BalanceId
		walletChangeDetails.BalanceType = cashBal.BalanceType
		walletChangeDetails.EffectiveExpiryDate = cashBal.EffectiveExpiryDate
		walletChangeDetails.EffectiveStartDate = cashBal.EffectiveStartDate
		walletChangeDetails.InstanceID = instanceId
		walletChangeDetails.PackageID = packageId
		walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType(FeeType)

		credit := (cashBal.CreditStatus.Credit + previousWalletCreditTotal) - Fee
		if credit >= float64(config.CommonDynConf.Config.MinBalanceThreshold) || subscriberType == POSTPAID {
			mlog.MavLog(mlog.INFO, transID, "MinBalanceThreshold", float64(config.CommonDynConf.Config.MinBalanceThreshold), "final credit", credit)
			for balanceID, CreditValue := range *balanceIdCreditMap {
				mlog.MavLog(mlog.INFO, transID, "Inside balanceIdCreditMap loop, balId", balanceID, "CreditValue", CreditValue)
				(*cashBalance)[balanceID].CreditStatus.Credit = (*cashBalance)[balanceID].CreditStatus.Credit - CreditValue
				(*cashBalance)[balanceID].CreditStatus.LastUpdateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
				//(*cashBalance)[balanceID].FeeTypeDeducted = getFeeTypeBasedOnWalletForType(FeeType)
				totalCreditElems++
			}
			if previousWalletCreditTotal > 0 {
				previousBal := (*cashBalance)[balId].CreditStatus.Credit
				walletChangeDetails.PreviousBalance = &previousBal
				(*cashBalance)[balId].CreditStatus.Credit = (*cashBalance)[balId].CreditStatus.Credit + previousWalletCreditTotal - Fee
				(*balanceIdCreditMap)[balId] = Fee - previousWalletCreditTotal
				mlog.MavLog(mlog.DEBUG, transID, "balanceIdCreditMap len =  1", totalCreditElems)
				mlog.MavLog(mlog.DEBUG, transID, "balanceIdCreditMap len =  1", "BalID", balId, "Latest Balance", (*cashBalance)[balId].CreditStatus.Credit, "UsedBalance", (*balanceIdCreditMap)[balId])
				changeBal := Fee - previousWalletCreditTotal
				currentBal := previousBal - changeBal
				tempChangedBal -= changeBal
				walletChangeDetails.ChangedBalance = &tempChangedBal
				walletChangeDetails.CurrentBalance = &currentBal

			} else {
				previousBal := (*cashBalance)[balId].CreditStatus.Credit
				walletChangeDetails.PreviousBalance = &previousBal
				(*cashBalance)[balId].CreditStatus.Credit = credit
				(*balanceIdCreditMap)[balId] = Fee
				changeBal := Fee
				currentBal := previousBal - changeBal
				tempChangedBal -= changeBal
				walletChangeDetails.ChangedBalance = &tempChangedBal
				walletChangeDetails.CurrentBalance = &currentBal
			}

			(*cashBalance)[balId].CreditStatus.LastUpdateDate = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			mlog.MavLog(mlog.INFO, transID, credit, "Exit- CheckSubscriberBalanceForPurchase  - true")

			mlog.MavLog(mlog.INFO, transID, credit, "walletChangeDetailForEdr true case KEY", balId+FeeType+instanceId+packageId)
			(*walletChangeDetailForEdr)[balId+FeeType+instanceId+packageId] = walletChangeDetails //update this map with BalanceId used and balance details that will be used in EDR

			return true, 0.0
		}
		mlog.MavLog(mlog.INFO, transID, credit, "Exit- CheckSubscriberBalanceForPurchase  - false")
		(*balanceIdCreditMap)[balId] = (*cashBalance)[balId].CreditStatus.Credit - float64(config.CommonDynConf.Config.MinBalanceThreshold)

		previousBal := (*cashBalance)[balId].CreditStatus.Credit
		walletChangeDetails.PreviousBalance = &previousBal
		changeBal := (*balanceIdCreditMap)[balId]
		currentBal := previousBal - changeBal
		tempChangedBal -= changeBal
		walletChangeDetails.ChangedBalance = &tempChangedBal
		walletChangeDetails.CurrentBalance = &currentBal
		(*walletChangeDetailForEdr)[balId+FeeType+instanceId+packageId+strconv.Itoa(int(tempChangedBal))] = walletChangeDetails
		mlog.MavLog(mlog.INFO, transID, credit, "walletChangeDetailForEdr false case KEY", balId+FeeType+instanceId+packageId)

		return false, (*cashBalance)[balId].CreditStatus.Credit - float64(config.CommonDynConf.Config.MinBalanceThreshold)
	}
	return false, 0.0
}

func CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee(transID string, featureConfig *map[string]map[string]interface{}, walletForPurchaseOrRenewalFee *[]string, source interface{}, Fee float64, walletForType string, walletChangeDetailForEdr *map[string]*WalletBalance, purchaseAndRenewalBalIdCreditMap *map[string]map[string]float64, instanceId string, packageId string, subscriberType int32) (bool, error) {

	mlog.MavLog(mlog.INFO, transID, "Enter- CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee, walletForType", walletForType, "type : ", subscriberType)
	mlog.MavLog(mlog.DEBUG, transID, "Enter- CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee, walletForType", walletForType, "type : ", subscriberType)
	cashBalCDB, isCashBalCDB := source.(*map[string]*CashBalancesCDB)
	subAccParam, isSubscAccParam := source.(*map[string]*SubscriberAccountParams)
	var subscriberBalanceCDB *map[string]*CashBalancesCDB
	var subscriberBalance *map[string]*SubscriberAccountParams

	if subscriberType == POSTPAID {
		if isCashBalCDB {
			for balanceID, cashBal := range *cashBalCDB {
				if !isWalletExpired(cashBal.EffectiveExpiryDate) { //R.FRN-11293-US4
					addFeeTypeDetailsForEDRForSubscriberCreate(transID, balanceID, *cashBal, walletChangeDetailForEdr)
				}
			}
		}

		walletChangeDetails := new(WalletBalance)
		walletChangeDetails.InstanceID = instanceId
		walletChangeDetails.PackageID = packageId
		walletChangeDetails.FeeType = getFeeStringBasedOnWalletForType(walletForType)
		Fee = Fee * -1
		walletChangeDetails.ChangedBalance = &Fee
		(*walletChangeDetailForEdr)[walletForType+instanceId+packageId+strconv.Itoa(int(Fee))] = walletChangeDetails
		mlog.MavLog(mlog.DEBUG, transID, "CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee for postpaid")
		return true, nil
	}

	if featureConfig == nil || walletForPurchaseOrRenewalFee == nil || walletChangeDetailForEdr == nil || purchaseAndRenewalBalIdCreditMap == nil {
		mlog.MavLog(mlog.ERROR, transID, "CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee :: some of the input params are nil")
		return false, errors.New("INTERNAL_ERROR")
	}

	if isCashBalCDB {
		subscriberBalanceCDB = cashBalCDB
	} else if isSubscAccParam {
		subscriberBalance = subAccParam
	} else {
		mlog.MavLog(mlog.ERROR, transID, "CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee :: invalid input")
		return false, errors.New("INTERNAL_ERROR")
	}

	balanceIdCreditMap := make(map[string]float64)
	previousWalletCredit := 0.0
	previousWalletCreditTotal := 0.0
	balanceDebited := false
	consumedFromPurchaseOrRenewalFee := false
	consumedFromSPPurchaseOrRenewalFee := false

	tempWalletChangeDetailForEdr := make(map[string]*WalletBalance)

	for k := range balanceIdCreditMap {
		delete(balanceIdCreditMap, k)
	}

	for _, balanceID := range *walletForPurchaseOrRenewalFee {
		mlog.MavLog(mlog.INFO, transID, "Input BalanceID for Purchase/Renewal Fee", balanceID)
		previousWalletCreditTotal = previousWalletCreditTotal + previousWalletCredit //R.FRN-11293-US5
		if isCashBalCDB {
			balanceDebited, previousWalletCredit = CheckAndUpdateSubBalanceForPurchaseBasedOnBalanceIDCDB(transID, subscriberBalanceCDB, Fee, walletForType, balanceID, previousWalletCreditTotal, &balanceIdCreditMap, &tempWalletChangeDetailForEdr, instanceId, packageId, walletChangeDetailForEdr, subscriberType)
		} else if isSubscAccParam {
			balanceDebited, previousWalletCredit = CheckSubBalanceForPurchaseBasedOnBalanceID(transID, subscriberBalance, Fee, walletForType, balanceID, previousWalletCreditTotal, &balanceIdCreditMap, &tempWalletChangeDetailForEdr, instanceId, packageId, subscriberType)
		}
		if !balanceDebited {
			continue
		}
		consumedFromPurchaseOrRenewalFee = true
		break
	}

	if !consumedFromPurchaseOrRenewalFee {
		for k := range tempWalletChangeDetailForEdr {
			delete(tempWalletChangeDetailForEdr, k)
		}
		if SPLevelWalletForPurchaseFee, ok := (*featureConfig)["general"][walletForType]; ok {
			mlog.MavLog(mlog.INFO, transID, "SP config for walletForType:", walletForType, "=", SPLevelWalletForPurchaseFee)
			previousWalletCredit := 0.0
			previousWalletCreditTotal := 0.0
			for k := range balanceIdCreditMap {
				delete(balanceIdCreditMap, k)
			}
			SPWalletFeeBalIdSlice := ConvertToSlice(SPLevelWalletForPurchaseFee)
			for _, balIdValue := range SPWalletFeeBalIdSlice {
				mlog.MavLog(mlog.INFO, transID, "SPLevelWalletForPurchaseFee Type", balIdValue)
				previousWalletCreditTotal = previousWalletCreditTotal + previousWalletCredit //R.FRN-11293-US5
				if isCashBalCDB {
					balanceDebited, previousWalletCredit = CheckAndUpdateSubBalanceForPurchaseBasedOnBalanceIDCDB(transID, subscriberBalanceCDB, Fee, walletForType, balIdValue, previousWalletCreditTotal, &balanceIdCreditMap, &tempWalletChangeDetailForEdr, instanceId, packageId, walletChangeDetailForEdr, subscriberType)

				} else if isSubscAccParam {
					balanceDebited, previousWalletCredit = CheckSubBalanceForPurchaseBasedOnBalanceID(transID, subscriberBalance, Fee, walletForType, balIdValue, previousWalletCreditTotal, &balanceIdCreditMap, &tempWalletChangeDetailForEdr, instanceId, packageId, subscriberType)
				}
				if !balanceDebited {
					continue
				}
				consumedFromSPPurchaseOrRenewalFee = true
				break
			}
		} else {
			return false, errors.New("SP_CONFIG_walletForPurchaseFee_MISSING")
		}

	}

	mlog.MavLog(mlog.INFO, transID, "Exit- CheckAndUpdateBalancesBasedOnWalletForPurchaseOrRenewalFee, consumedFromSPPurchaseOrRenewalFee = ", consumedFromSPPurchaseOrRenewalFee, "consumedFromPurchaseOrRenewalFee", consumedFromPurchaseOrRenewalFee, "balanceIdCreditMap", balanceIdCreditMap)

	if consumedFromPurchaseOrRenewalFee || consumedFromSPPurchaseOrRenewalFee {
		for key, value := range tempWalletChangeDetailForEdr {
			(*walletChangeDetailForEdr)[key] = value
		}
		mlog.MavLog(mlog.DEBUG, transID, "walletChangeDetailForEdr, walletType", walletForType, "walletDetails", walletChangeDetailForEdr)
		(*purchaseAndRenewalBalIdCreditMap)[walletForType] = balanceIdCreditMap
		return true, nil
	}
	return false, nil
}

func ConvertToSlice(Events interface{}) []string {
	events := make([]string, 0)
	switch reflect.TypeOf(Events).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(Events)
		for i := 0; i < s.Len(); i++ {
			//            fmt.Println(s.Index(i))
			events = append(events, s.Index(i).Interface().(string))
		}
		return events
	}
	return nil
}

//return renewalFee and Discount deducted
func CalculateRenewalFeeAfterDiscountSlab(transDetails *config.LoggerTransactionDetails, rewnewalCycleCount int64, renewalFee float64, RenewalPolicy string, discountSlabDetails *DiscountInfo, currencyConversionFactor float64) (float64, float64, bool) {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- CalculateRenewalFeeAfterDiscountSlab renewalFee : ", renewalFee)

	changedFee := renewalFee
	isDiscountSlabApplied := false
	actualDiscount := 0.0
	if rewnewalCycleCount > 0 && discountSlabDetails.Discount != nil {
		sort.SliceStable(discountSlabDetails.Discount, func(p, q int) bool {
			return discountSlabDetails.Discount[p].Low < discountSlabDetails.Discount[q].Low
		})
		if len(discountSlabDetails.Discount) > 0 && rewnewalCycleCount < discountSlabDetails.Discount[0].Low {
			return changedFee, actualDiscount, false
		}

		for _, val := range discountSlabDetails.Discount {
			if rewnewalCycleCount >= val.Low && rewnewalCycleCount <= val.High {
				if val.DiscountVal != nil {
					isDiscountSlabApplied = true
					if val.DiscountType == "FIXED" {
						changedFee = renewalFee - *val.DiscountVal
						actualDiscount = *val.DiscountVal
					} else {
						changedFee = 0.01 * renewalFee * (100 - *val.DiscountVal)
						actualDiscount = *val.DiscountVal * renewalFee * 0.01
					}
					break
				}
			}
		}

	}
	if changedFee < 0 {
		changedFee = 0
	}
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CalculateRenewalFeeAfterDiscountSlab newRenewalFee : ", changedFee, " actualDiscountSlab: ", actualDiscount, " isDiscountSlabApplied: ", isDiscountSlabApplied)

	return changedFee, actualDiscount, isDiscountSlabApplied
}

func GetRenewalFeeForNormalPkg(transID string, renewalSlab []*RenewalFeeSlab, renewalFee *float64, RenewalPolicy string) float64 {
	mlog.MavLog(mlog.DEBUG, transID, "Enter- GetRenewalFeeForNormalPkg renewalFee : ", renewalFee, " RenewalPolicy:", RenewalPolicy)
	var newFee float64
	if renewalFee != nil {
		newFee = *renewalFee
	}
	if RenewalPolicy != "PERSLAB" {
		return newFee
	} else if RenewalPolicy == "PERSLAB" && renewalSlab != nil {
		for _, feeSlabData := range renewalSlab {
			if feeSlabData.Low == 1 {
				newFee = feeSlabData.Fee
			}

		}

	}
	mlog.MavLog(mlog.DEBUG, transID, "Exit- GetRenewalFeeForNormalPkg newFee : ", newFee)
	return newFee
}

// FRN-15255 (explained in FRN-14451 Phase2)
func GeneralRoundingLogic(transDetails *config.LoggerTransactionDetails, subscriberType int32, amount float64, currencyUsage *CurrencyUsage) float64 {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GeneralRoundingLogic, subscriberType: ", subscriberType, ", amount: ", amount)

	roundingRule := string(commconsts.RoundDwn) // Default behavior is Round DOWN

	if amount == 0 {
		goto FINAL
	}
	if currencyUsage == nil {
		amount = CurrencyRounding(transDetails, amount, roundingRule)
	} else {

		roundingRule = currencyUsage.CurrencyRoundingRule

		if subscriberType == POSTPAID {
			// For POSTPAID rounding needs to be done on amount with currency conversion Factor
			amount = CurrencyRounding(transDetails, amount, roundingRule)
		} else if subscriberType == PREPAID {
			// For PREPAID rounding needs to be done on amount without currency conversion Factor

			var currConvFactorFloat64 float64 = 1.0 // Default value, to avoid division by zero
			var roundingPrecision int32 = 0         // Default value

			if currencyUsage.CurrencyConversionFactor != 0 {
				currConvFactorFloat64 = float64(currencyUsage.CurrencyConversionFactor)
			}

			if currencyUsage.RoundingPrecisionForPrepaid != 0 {
				roundingPrecision = currencyUsage.RoundingPrecisionForPrepaid
			}

			// Default behavior for 0 value is same as POSTPAID flow, Rounding after cuurencyConversionFactor
			if roundingPrecision == 0 {
				amount = CurrencyRounding(transDetails, amount, roundingRule)
				goto FINAL
			}

			mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "GeneralRoundingLogic, currConvFactorFloat64: ", currConvFactorFloat64, ", CurrencyRoundingRule: ", roundingRule, ", RoundingPrecisionForPrepaid: ", roundingPrecision)

			roundingPrecisionFactor := math.Pow10(int(roundingPrecision))
			if roundingPrecisionFactor == currConvFactorFloat64 {
				// Since both are same, behavior will be same as POSTPAID
				amount = CurrencyRounding(transDetails, amount, roundingRule)
				goto FINAL
			} else if roundingPrecisionFactor > currConvFactorFloat64 {
				// Restrict RoundingPrecisionForPrepaid to Currency Conversion Factor
				roundingPrecisionFactor = currConvFactorFloat64
			}

			// Precautionary check (to avoid division by zero)
			if roundingPrecisionFactor <= 0 {
				roundingPrecisionFactor = 1.0
			}

			// First Divide the amount with Currency Conversion Factor
			amount = amount / currConvFactorFloat64

			// Rounding to be done on this amount
			amount = amount * roundingPrecisionFactor
			amount = CurrencyRounding(transDetails, amount, roundingRule)
			amount = amount / roundingPrecisionFactor

			// Multiply back the Currency Conversion Factor
			amount = amount * currConvFactorFloat64
		}
	}
FINAL:
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GeneralRoundingLogic, Final amount after Rounding : ", amount)
	return amount
}

func CurrencyRounding(transDetails *config.LoggerTransactionDetails, fee float64, roundingRule string) float64 {
	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- CurrencyRounding, Fee : ", fee, ", RoundingRule: ", roundingRule)

	if roundingRule == string(commconsts.RoundUp) {
		fee = math.Ceil(fee)
	} else if roundingRule == string(commconsts.NoRoundOff) {
		//No action as fee need to return as it is
	} else if roundingRule == string(commconsts.RoundOff) {
		fee = math.Round(fee)
	} else {
		//default behaviour is to return floor value
		fee = math.Floor(fee)
	}

	mlog.MavTrace(transDetails.TraceFlag, mlog.DEBUG, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- CurrencyRounding, newFee : ", fee)
	return fee
}

//GetQuotaValue - It will return quota value
func GetQuotaValue(transID string, quotaSlab *QuotaSlab, quota int64, renewalPolicy string, pkgRenewalCycleCount *int32) int64 {
	mlog.MavLog(mlog.DEBUG, transID, "Enter- GetQuotaValue, quota: ", quota, "renewalPolicy: ", renewalPolicy)
	var quotaValue int64
	if renewalPolicy != "PERSLAB" {
		quotaValue = quota
	} else if renewalPolicy == "PERSLAB" && quotaSlab != nil {
		//Added a variable so in future if prdRenewalCycleCount is supported we could add a case for that
		var renewalCycleCount *int32
		if quotaSlab.SlabInputType == PKGLEVELRENCOUNT {
			renewalCycleCount = pkgRenewalCycleCount
		}
		for _, slab := range quotaSlab.Counter {
			if slab.Low != nil && slab.High != nil && slab.Quota != nil {
				mlog.MavLog(mlog.DEBUG, transID, "slab Low:", slab.Low, ", slab High:", slab.High, ", slab Quota:", slab.Quota)
				if renewalCycleCount == nil || *renewalCycleCount == 0 {
					if *slab.Low == 1 {
						quotaValue = *slab.Quota
						break
					}
				} else if *renewalCycleCount >= *slab.Low && (*renewalCycleCount <= *slab.High || *slab.High == -1) {
					quotaValue = *slab.Quota
					mlog.MavLog(mlog.DEBUG, transID, "RenewalCycleCount: ", *renewalCycleCount, "quotaValue: ", quotaValue)
					break
				}
			} //Not adding else case since quotaSlab validations will be added in the validate.go for every request
		}
	}
	mlog.MavLog(mlog.DEBUG, transID, "Exit- GetQuotaValue quotaValue : ", quotaValue)
	return quotaValue
}

//GetSubDocumentsWithMulSpecSupport - API fetches all the sub docs from the paths provided in subDocPaths array based on the Spec provided from the doc with key provided and fills the structures provided in target array. It return CAS of LookupIn op and error. This API returns midway if unmarshalling any of the subdoc into provided structures fails with corresponding error and 0 cas
func GetSubDocumentsWithMulSpecSupport(transDetails *config.LoggerTransactionDetails, key string, subDocPaths []string, target []interface{}, pathSpec []string) (gocb.Cas, error) {
	mlog.MavLog(mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Enter- GetSubDocumentsWithMulSpecSupport, Doc ID: ", key, " Paths: ", subDocPaths)
	pathsArrayLen := len(subDocPaths)
	if len(target) != pathsArrayLen {
		err := errors.New("Target field list size mismatch")
		mlog.MavLog(mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithMulSpecSupport, Target field list size mismatch: ", err)
		return 0, err
	}
	ops := make([]gocb.LookupInSpec, pathsArrayLen)
	for i := 0; i < pathsArrayLen; i++ {
		if pathSpec[i] == GET_SPEC {
			ops[i] = gocb.GetSpec(subDocPaths[i], nil)
		} else if pathSpec[i] == COUNT_SPEC {
			ops[i] = gocb.CountSpec(subDocPaths[i], nil)
		} else if pathSpec[i] == EXISTS_SPEC {
			ops[i] = gocb.ExistsSpec(subDocPaths[i], nil)
		}
	}
	startTime := config.GetTimeNs()
	serviceProviderID := GetServiceProviderID(key)
	config.UpdateCounterMetrics(transDetails, config.CB_REQ, []string{config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID})

	multiGetResult, err := CollectionObj.LookupIn(key, ops, &gocb.LookupInOptions{
		Timeout: time.Second * time.Duration(config.CommonDynConf.Config.CouchbaseConfig.Timeout)})
	args := []string{config.CB_RESP, config.GET_SUB_DOC, config.GetKtab(key), serviceProviderID}
	config.UpdateRespMetrics(transDetails, []int64{startTime, config.GetTimeNs()}, err, args, nil)
	if err != nil {
		mlog.MavLog(mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithMulSpecSupport, Error in fetching sub docs: ", err)
		return 0, err
	}
	for i := 0; i < pathsArrayLen; i++ {
		decodeError := multiGetResult.ContentAt(uint(i), target[i])
		if decodeError != nil {
			mlog.MavLog(mlog.ERROR, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithMulSpecSupport, Error decoding the value of subdoc - ", subDocPaths[i], " error:", decodeError)
			return 0, decodeError
		}
	}
	mlog.MavLog(mlog.INFO, transDetails.TransID, "sessionid:", transDetails.TraceID, "Exit- GetSubDocumentsWithMulSpecSupport")
	return multiGetResult.Cas(), nil

}
