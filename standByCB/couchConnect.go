package standByCB

import (
	//cc "commconsts"

	//cs "commonstructs"
	"errors"
	"sync"

	"encoding/json"
	ml "mwp-appcommon/mavcimclient"
	"time"

	"github.com/couchbase/gocb/v2"
)

type CommCbCfgIntStruct struct {
	ClusterName                   string `json:"clusterName" validate:"required"`
	CouchBaseIP                   string `json:"hostAddr" validate:"required"`
	CouchBasePort                 string `json:"connPort"`
	CouchBaseUname                string `json:"userName" validate:"required"`
	CouchBasePwd                  string `json:"password" validate:"required"`
	CouchBaseBucketName           string `json:"main_bucketName" validate:"required"`
	CbEphBucketName               string `json:"ephemeral_bucketName" validate:"required"`
	IsLocalDc                     *bool  `json:"is_local_dc" validate:"required"`
	BucketOpenTimeout             int    `json:"bucketOpenTimeout_sec" validate:"required"`
	ConnectionTimeout             int    `json:"connectionTimeout_sec" validate:"required"`
	MaxCouchFailureThresholdCount int    `json:"maxCouchFailureThreshold_count" validate:"required"`
	MaxNodeFailureThresholdPerc   int    `json:"maxNodeFailureThreshold_perc" validate:"required"`
}

type CommCbTimerIntStruct struct {
	CacheTTLForChatRoomData *int `json:"cacheTTLForChatRoomData_sec" validate:"required"`
	CacheTTLForUserDoc      *int `json:"cacheTTLForUserDoc_sec" validate:"required"`
	BcidDocTtl              *int `json:"businessChatDocTTL_sec" validate:"required"`
	DelBizzDocsTtl          *int `json:"delBizzDocsTTL_sec" validate:"required"`
	FeedbackfileTtl         *int `json:"cacheTTLForFeedbackFile_sec" validate:"required"`
}

//CommonCbCfgStruct - defined for all microservices-COMMAS_COMMON_CB_CFG_FILE
type CommonCbCfgStruct struct {
	CommCbIntData CommCbIntStruct `json:"config" validate:"required"`
}

type CommCbIntStruct struct {
	CommCbCfgData      []CommCbCfgIntStruct `json:"couchDb_Cfg" validate:"required,dive"`
	CommCbTimerCfgData CommCbTimerIntStruct `json:"timer_Cfg" validate:"required"`
}

var (
	gCDbConnectionMap   map[string]*CouchConnectionStruct
	gCouchRetryCountCfg int
	gLocalClusterName   string
)
var (
	gCDbEphConnectionMap   map[string]*CouchConnectionStruct
	gCouchEphRetryCountCfg int
	gLocalEphClusterName   string
)

func SetCouchRetryCfg(counter int) {
	gCouchRetryCountCfg = counter
}

//CouchConnectionStruct struct to handle all couch connection details.
type CouchConnectionStruct struct {
	ClusterName                   string
	ConnectionStr                 string
	UserName                      string
	Password                      string
	ConnectionTimeout             int
	BucketOpenTimeout             int
	BucketName                    string
	BucketPtr                     *gocb.Bucket
	Collectionptr                 *gocb.Collection
	ClusterPtr                    *gocb.Cluster
	ConnectionStatus              bool
	FailureCount                  int
	MaxCouchFailureThresholdCount int
	MaxNodeFailureThresholdCount  int //-> Per cluster config
	ConnectionStatusLock          sync.Mutex
	IsEphemeralBkt                bool
}

//InitCouchDBConnectionv2 Connect to couch and retun local and remote collection
func InitCouchDBConnectionv2(commonCbCfgData CommonCbCfgStruct, useEphemeralBkt bool) (map[string]*CouchConnectionStruct, error) {

	var lTransID string
	localClusterName := ""
	ml.MavLog(ml.INFO, lTransID, "Entering InitCouchDBConnection")

	if !useEphemeralBkt {
		gCDbConnectionMap = make(map[string]*CouchConnectionStruct)
	} else {
		gCDbEphConnectionMap = make(map[string]*CouchConnectionStruct)
	}

	couchConfigList := commonCbCfgData.CommCbIntData.CommCbCfgData

	for _, cdbCfg := range couchConfigList {

		connectionStr := "couchbase://" + cdbCfg.CouchBaseIP
		ml.MavLog(ml.INFO, lTransID, "Connect to CouchBase with: ", connectionStr)

		if *cdbCfg.IsLocalDc == true {
			localClusterName = cdbCfg.ClusterName
		}

		ml.MavLog(ml.INFO, lTransID, "Authenticate with Username: ", cdbCfg.CouchBaseUname)
		ml.MavLog(ml.INFO, lTransID, "Authenticate with Passwd: ", cdbCfg.CouchBasePwd)
		bktName := cdbCfg.CouchBaseBucketName
		if useEphemeralBkt {
			bktName = cdbCfg.CbEphBucketName
		}

		connectionData := CouchConnectionStruct{
			ClusterName:                   cdbCfg.ClusterName,
			ConnectionStr:                 connectionStr,
			ConnectionTimeout:             cdbCfg.ConnectionTimeout,
			BucketOpenTimeout:             cdbCfg.BucketOpenTimeout,
			UserName:                      cdbCfg.CouchBaseUname,
			Password:                      cdbCfg.CouchBasePwd,
			BucketName:                    bktName,
			MaxCouchFailureThresholdCount: cdbCfg.MaxCouchFailureThresholdCount,
			IsEphemeralBkt:                useEphemeralBkt,
		}

		err := connectToCouchDb(&connectionData, lTransID)
		if err != nil {
			ml.MavLog(ml.ERROR, lTransID, "Connect error:", err)
		}
		if connectionData.BucketPtr != nil {
			//Ping Bucket and get total Node count.
			pings, err := connectionData.BucketPtr.Ping(&gocb.PingOptions{
				ReportID:     "couchbase_healthcheck_report",
				ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
			})
			if err != nil {
				ml.MavLog(ml.ERROR, lTransID, "InitCouchDBConnectionv2: Cluster Ping failed During ", err)
				return nil, err
			}
			pingReport := pings.Services[gocb.ServiceTypeKeyValue]
			ml.MavLog(ml.INFO, lTransID, "InitCouchDBConnectionv2: total Node count in cluster ", len(pingReport))

			//Update max node failure Threshold count as per the configured percent
			connectionData.MaxNodeFailureThresholdCount = int(len(pingReport) * (cdbCfg.MaxNodeFailureThresholdPerc / 100))

			//Single node configuration. set to 1
			if len(pingReport) == 1 || connectionData.MaxNodeFailureThresholdCount <= 0 {
				connectionData.MaxNodeFailureThresholdCount = 1
			}
		}

		//Add connection struct to map even if connecting fails. Reconnect happens during health check
		if !useEphemeralBkt {
			gCDbConnectionMap[cdbCfg.ClusterName] = &connectionData
			gLocalClusterName = localClusterName
		} else {
			gCDbEphConnectionMap[cdbCfg.ClusterName] = &connectionData
			gLocalEphClusterName = localClusterName
		}
	}

	ml.MavLog(ml.ERROR, lTransID, "gCDbConnectionMap - ", len(gCDbConnectionMap))
	ml.MavLog(ml.ERROR, lTransID, "useEphemeralBkt - ", useEphemeralBkt)
	ml.MavLog(ml.ERROR, lTransID, "gCDbEphConnectionMap - ", len(gCDbEphConnectionMap))

	if (len(gCDbConnectionMap) < 1 && !useEphemeralBkt) || (len(gCDbEphConnectionMap) < 1 && useEphemeralBkt) {
		//return error
		ml.MavLog(ml.ERROR, lTransID, "CouchDBConnection failed for all the dc.")
		errStr := "CouchDBConnection failed for all the dc"
		return nil, errors.New(errStr)
	} else if useEphemeralBkt {
		ml.MavLog(ml.INFO, lTransID, "InitCouchDBConnection is Success for Ephemeral bucket. gCDbEphConnectionMap:", gCDbEphConnectionMap)
		return gCDbEphConnectionMap, nil
	}

	ml.MavLog(ml.INFO, lTransID, "InitCouchDBConnection is Success. gCDbConnectionMap:", gCDbConnectionMap)
	return gCDbConnectionMap, nil
}

//connectToCouchDb establish a couch connection with given information
func connectToCouchDb(connStruct *CouchConnectionStruct, lTransID string) error {
	cDbCluster, err := gocb.Connect(connStruct.ConnectionStr, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: connStruct.UserName,
			Password: connStruct.Password,
		},
		TimeoutsConfig: gocb.TimeoutsConfig{ConnectTimeout: time.Duration(connStruct.ConnectionTimeout) * time.Second},
	})
	if err != nil {
		ml.MavLog(ml.ERROR, lTransID, "Connect error:", err)

		connStruct.MarkCouchConnectionInactive()
		return err
	}
	cDbBucket := cDbCluster.Bucket(connStruct.BucketName)
	//Wait until bucket is ready to use
	err = cDbBucket.WaitUntilReady(time.Duration(connStruct.BucketOpenTimeout)*time.Second, nil)
	if err != nil || cDbBucket == nil {
		ml.MavLog(ml.ERROR, lTransID, "OpenBucket error: ", err)

		connStruct.MarkCouchConnectionInactive()
		RaiseCouchbaseDownAlarm(connStruct)
		return err
	}
	//update bucket and collection pointers
	connStruct.Collectionptr = cDbBucket.DefaultCollection()
	connStruct.BucketPtr = cDbBucket
	connStruct.ClusterPtr = cDbCluster

	connStruct.MarkCouchConnectionActive()
	return err
}
func IsEphCouchbaseHealthOkV2() bool {
	return IsCouchbaseHealthOkV2(gCDbEphConnectionMap)
}

//IsCouchbaseHealthOkV2 Health check for couch db connection
func IsCouchbaseHealthOkV2(connectionMap map[string]*CouchConnectionStruct) bool {
	/*
	   1.if bucket is nil,then connect to couch
	   2.if system is not nil, then ping and check for connection status
	   3. If connection ping is negative, Mark connection status as false
	*/
	noOfClustersDown := 0
	lTransID := "READINESS_PROBE"
	for connetionIndex, connectionData := range connectionMap {
		bucket := connectionData.BucketPtr
		if bucket != nil {
			pings, err := bucket.Ping(&gocb.PingOptions{
				ReportID:     "couchbase_healthcheck_report",
				ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
			})
			if err != nil {
				ml.MavLog(ml.ERROR, lTransID, "Healthcheck probe to couchbase failed with err:", err)
				return false
			}
			pingReport, _ := json.Marshal(pings)
			ml.MavLog(ml.INFO, lTransID, "Couchbase ping report in JSON format:", string(pingReport))

			for service, pingReports := range pings.Services {
				clusterDown := false
				noOfNodesDown := 0

				for _, pingReport := range pingReports {
					if pingReport.State != gocb.PingStateOk {
						ml.MavLog(ml.ERROR, lTransID, "Ping not ok for service:", service)

						RaiseCouchbaseNodeDownAlarm(connectionData, pingReport.Remote)
						//Increment nodes down count.
						noOfNodesDown++
					} else {
						RaiseCouchbaseNodeUpAlarm(connectionData, pingReport.Remote)
					}

				}
				//mark Cluster down if number of nodes down is more than or equal to threshold value
				if noOfNodesDown >= connectionData.MaxNodeFailureThresholdCount {
					ml.MavLog(ml.ERROR, lTransID, "Service down for cluster:", pings.ID, "Bucket:", connetionIndex)
					//mark cluster as Down
					func() {
						connectionData.ConnectionStatusLock.Lock()
						connectionData.ConnectionStatus = false
						connectionData.ConnectionStatusLock.Unlock()
					}()
					RaiseCouchbaseDownAlarm(connectionData)
					clusterDown = true
					//Increment cluster down count
					noOfClustersDown++
				}
				if !clusterDown {
					//Mark service status as up
					connectionData.MarkCouchConnectionActive()
				}
			}
		} else {
			//Bucket is Nil. let's connect to couch
			err := connectToCouchDb(connectionData, lTransID)
			if err != nil {
				ml.MavLog(ml.ERROR, lTransID, "Service down for cluster: ", connetionIndex)
				noOfClustersDown++
			} else {
				//Mark service status as up
				connectionData.MarkCouchConnectionActive()
			}
		}
	}
	if noOfClustersDown == len(connectionMap) {
		ml.MavLog(ml.ERROR, lTransID, "All couchbase clusters are down.")
		return false
	}
	return true
}

//getActivePrimaryCouchConnection return active couch connection to be used for couch operations
func getActivePrimaryCouchConnection(connectionMap map[string]*CouchConnectionStruct, localClusterName string) (*CouchConnectionStruct, bool) {
	/*
		1. Check for Local DC connection. Return if active
	*/
	primaryConnection := connectionMap[localClusterName]
	if checkIfConnectionIsActive(primaryConnection) && primaryConnection.ClusterPtr != nil {
		return primaryConnection, true
	}
	ml.MavLog(ml.WARN, "", "No active connection found")
	return primaryConnection, false
}

//getActiveSecondaryCouchConnection return active couch connection to be used for couch operations
func getActiveSecondaryCouchConnection(connectionMap map[string]*CouchConnectionStruct, localClusterName string) (*CouchConnectionStruct, bool) {
	/*
		1. Loop over given connections
		3. Check for secondary connection. Return if active
	*/
	for index, value := range connectionMap {
		if index == localClusterName {
			//Skip local DC
			continue
		}
		if checkIfConnectionIsActive(value) && value.ClusterPtr != nil {
			return value, true
		}
	}
	ml.MavLog(ml.WARN, "", "No active connection found")
	return nil, false
}

func checkIfConnectionIsActive(connection *CouchConnectionStruct) bool {
	return connection.ConnectionStatus
}

//MarkCouchConnectionActive mark given connection as active
func (conn *CouchConnectionStruct) MarkCouchConnectionActive() {
	if conn == nil {
		return
	}
	conn.ConnectionStatusLock.Lock()
	//Reset cluser failure count to Zero if cluster is Up
	conn.FailureCount = 0
	conn.ConnectionStatus = true
	conn.ConnectionStatusLock.Unlock()
	RaiseCouchbaseUpAlarm(conn)
}

//MarkCouchConnectionInactive mark current connection as inactive
func (conn *CouchConnectionStruct) MarkCouchConnectionInactive() {
	if conn == nil {
		return
	}
	conn.ConnectionStatusLock.Lock()
	conn.FailureCount++
	//mark Cluster down if number of nodes down is more than or equal to threshold value
	if conn.FailureCount >= conn.MaxCouchFailureThresholdCount {
		conn.ConnectionStatus = false
		RaiseCouchbaseDownAlarm(conn)
	}
	conn.ConnectionStatusLock.Unlock()
}

//closeConnection Close connection to given cluster
func (conn *CouchConnectionStruct) closeConnection() {
	if conn == nil {
		return
	}
	conn.ClusterPtr.Close(&gocb.ClusterCloseOptions{})
}

//CloseCouchDbConnection Close all connections to couch
func CloseCouchDbConnection(connectionMap map[string]*CouchConnectionStruct) {
	for _, connection := range connectionMap {
		connection.closeConnection()
	}
}

//CheckForActiveCouchConnection Check if current connection is active. else switch to Remote
//Default: Local DC is set as primary active
func CheckForActiveCouchConnection(connectionList map[string]*CouchConnectionStruct, useEphemeralBkt bool) (activeConnection *CouchConnectionStruct) {
	/*
		1. Check if local dc is up.
		2. If local is down switch to remote.
		3. If remote is down, return primary, even if it is down
	*/
	ml.MavLog(ml.INFO, "NULL", "Entering CheckForActiveCouchConnection.")
	localClusterName := gLocalClusterName
	if useEphemeralBkt {
		localClusterName = gLocalEphClusterName
	}
	//Get primary(local). Use if active
	primary, ok := getActivePrimaryCouchConnection(connectionList, localClusterName)
	if ok {
		activeConnection = primary
		RaiseCouchbaseUpAlarm(activeConnection)
		ml.MavLog(ml.INFO, "NULL", "CheckForActiveCouchConnection, Using connection ", activeConnection.ConnectionStr)
		return
	}

	activeConnection, ok = getActiveSecondaryCouchConnection(connectionList, localClusterName)
	//Active connection found. Return
	if ok {
		RaiseCouchbaseUpAlarm(activeConnection)
		ml.MavLog(ml.INFO, "NULL", "CheckForActiveCouchConnection, Using connection ", activeConnection.ConnectionStr)
		return
	}

	ml.MavLog(ml.WARN, "NULL", "CheckForActiveCouchConnection, No active connection. Returning Primary")
	//Return non nil connection here. worst case nil
	return returnNonNilConnection(connectionList)
}

//returnNonNilConnection return a non nil connection
func returnNonNilConnection(connectionList map[string]*CouchConnectionStruct) (activeConnection *CouchConnectionStruct) {
	for _, connection := range connectionList {
		if connection.ClusterPtr != nil {
			activeConnection = connection
			break
		}
	}
	return
}

//RaiseCouchbaseDownAlarm raise alarm for couch cluster down
func RaiseCouchbaseDownAlarm(couchConnData *CouchConnectionStruct) {

	managedObjMap := make(map[string]string)
	managedObjMap["ClusterName"] = couchConnData.ClusterName
	addlInfoMap := make(map[string]string)
	isLocalClusterStr := "false"
	isEphemeralBktStr := "false"

	localClusterName := gLocalClusterName
	if couchConnData.IsEphemeralBkt {
		localClusterName = gLocalEphClusterName
		isEphemeralBktStr = "true"
	}
	if localClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}

	addlInfoMap["IsLocalCluster"] = isLocalClusterStr
	managedObjMap["IsEphemeralBucket"] = isEphemeralBktStr
	ml.MavLog(ml.ERROR, "NULL", "Raising CouchbaseClusterDown alarm.", "managedObjMap:", managedObjMap, "addlInfoMap:", addlInfoMap)
	ml.MavAlarmWithManagedObjNAddlInfo("EVENT", "CouchbaseClusterDown", managedObjMap, addlInfoMap)

}

//RaiseCouchbaseUpAlarm raise alarm for couch cluster up
func RaiseCouchbaseUpAlarm(couchConnData *CouchConnectionStruct) {

	managedObjMap := make(map[string]string)
	managedObjMap["ClusterName"] = couchConnData.ClusterName
	addlInfoMap := make(map[string]string)
	isLocalClusterStr := "false"
	isEphemeralBktStr := "false"

	localClusterName := gLocalClusterName
	if couchConnData.IsEphemeralBkt {
		localClusterName = gLocalEphClusterName
		isEphemeralBktStr = "true"
	}
	if localClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}
	if gLocalClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}
	addlInfoMap["IsLocalCluster"] = isLocalClusterStr
	managedObjMap["IsEphemeralBucket"] = isEphemeralBktStr
	ml.MavLog(ml.INFO, "NULL", "Raising CouchbaseClusterUp alarm.", "managedObjMap:", managedObjMap, "addlInfoMap:", addlInfoMap)
	ml.MavAlarmWithManagedObjNAddlInfo("EVENT", "CouchbaseClusterUp", managedObjMap, addlInfoMap)

}

//RaiseCouchbaseNodeDownAlarm raise alarm for couch node down
func RaiseCouchbaseNodeDownAlarm(couchConnData *CouchConnectionStruct, nodeIp string) {

	managedObjMap := make(map[string]string)
	managedObjMap["ClusterName"] = couchConnData.ClusterName
	managedObjMap["NodeIp"] = nodeIp
	addlInfoMap := make(map[string]string)
	isLocalClusterStr := "false"
	isEphemeralBktStr := "false"

	localClusterName := gLocalClusterName
	if couchConnData.IsEphemeralBkt {
		localClusterName = gLocalEphClusterName
		isEphemeralBktStr = "true"
	}
	if localClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}
	if gLocalClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}
	addlInfoMap["IsLocalCluster"] = isLocalClusterStr
	managedObjMap["IsEphemeralBucket"] = isEphemeralBktStr
	ml.MavLog(ml.ERROR, "NULL", "Raising CouchbaseNodeDown alarm.", "managedObjMap:", managedObjMap, "addlInfoMap:", addlInfoMap)
	ml.MavAlarmWithManagedObjNAddlInfo("EVENT", "CouchbaseNodeDown", managedObjMap, addlInfoMap)

}

//RaiseCouchbaseNodeUpAlarm raise alarm for couch node up
func RaiseCouchbaseNodeUpAlarm(couchConnData *CouchConnectionStruct, nodeIp string) {

	managedObjMap := make(map[string]string)
	managedObjMap["ClusterName"] = couchConnData.ClusterName
	managedObjMap["NodeIp"] = nodeIp
	addlInfoMap := make(map[string]string)
	isLocalClusterStr := "false"
	isEphemeralBktStr := "false"

	localClusterName := gLocalClusterName
	if couchConnData.IsEphemeralBkt {
		localClusterName = gLocalEphClusterName
		isEphemeralBktStr = "true"
	}
	if localClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}
	if gLocalClusterName == couchConnData.ClusterName {
		isLocalClusterStr = "true"
	}

	addlInfoMap["IsLocalCluster"] = isLocalClusterStr
	managedObjMap["IsEphemeralBucket"] = isEphemeralBktStr

	ml.MavLog(ml.INFO, "NULL", "Raising CouchbaseNodeUp alarm.", "managedObjMap:", managedObjMap, "addlInfoMap:", addlInfoMap)
	ml.MavAlarmWithManagedObjNAddlInfo("EVENT", "CouchbaseNodeUp", managedObjMap, addlInfoMap)

}

//GetConnectionMap expose connection map outside package
func GetConnectionMap(useEphemeralBkt bool) map[string]*CouchConnectionStruct {
	if !useEphemeralBkt {
		return gCDbConnectionMap
	}
	return gCDbEphConnectionMap
}
