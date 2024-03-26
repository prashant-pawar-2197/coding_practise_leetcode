package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
	"github.com/couchbase/gocb/v2"
)



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

var couchbaseDetail *CouchbaseDetails = &CouchbaseDetails{}

type CouchbaseDetails struct {
	Cluster            *gocb.Cluster
	CollectionObj      *gocb.Collection
	BucketPtr			*gocb.Bucket
	BucketName          string
	CouchbaseServer		string
	CouchBaseUname		string
	CouchBasePwd		string
	BillingBucket      *gocb.Collection
	connflag           bool
	transID            string
	CBDownAlarmRaised  bool
	gCDbConnectionList map[string]*CouchConnectionStruct
	gCommonCbCfgData   CommCbIntStruct
}

type CdrConfig struct {
    Config struct {
        CdrSettings []struct {
	//		MsgInterface	  ms.MsgInterfaceConfig `json:"msgInterface"`
            HttpServer        string `json:"httpServer"`
            Http2Server        string `json:"http2Server"`
            CouchbaseServer   string `json:"couchbaseServer"`
            CouchbaseBucket   string `json:"couchbaseBucket"`
            CouchbaseUsername string `json:"couchbaseUsername"`
            CouchbasePassword string `json:"couchbasePassword"`
            DumpCdr           bool   `json:"dumpCdr"`
            MetricCfgFile     string `json:"metricCfgFile"`
            MetricPort        string `json:"metricPort"`
            AlarmsCfgFile     string `json:"alarmsCfgFile"`
            SctFile           string `json:"sctFile"`
            ProfileFile       string `json:"profileFile"`
            ActionFile        string `json:"actionFile"`
			LogTxnSctFile     string `json:"logTxnSctFile"`
            TmplFile          string `json:"tmplFile"`
			CustomSetFile     string `json:"customSetFile"`
			ReadCustomSetFromCB bool `json:readCustomSetFromCB,omitempty`
			CustomSetDocID	  string `json:customSetDocId,omitempty`
            ComponentName     string `json:"componentName"`
            PollInterval      int    `json:"pollInterval"`
			TimerServiceURI  string `json:"TimerServiceURI"`
			NotifyServerURI   string `json:"NotifyServerURI"`
			UDFSAPItimeout    int `json:"UDFSAPItimeout"`
			UDSFMaxTimerTimeOut int `json:"UDSFMaxTimerTimeOut"`
			UDSFMinTimerTimeOut int `json:"UDSFMinTimerTimeOut"`
			KafkaBrokers      string `json:"kafkaBrokers"`
			KafkaTopic        string `json:"kafkaTopic"`
			OpenDistro        bool   `json:"openDistro"`
			ElUsername        string `json:"elUsername"`
			ElPassword        string `json:"elPassword"`
			ElCerts           string `json:"caCert"`
			LogLevel          int `json:"logLevel"`
			MaxNumBerOfUdsfCallBack int `json:"maxNumBerOfUdsfCallBack"`
            ConsoleLogging    bool   `json:"consoleLogging"`
            TickerInterval    int    `json:"tickerInterval"`
            WriteToFile       bool   `json:"writeToFile"`
			ReadconfigModeType  string `json:"readconfigModeType"`
    		VaultConfigFilePath   string   `json:"vaultConfigFilePath"`
        } `json:"processorCfg"`
		StandByCbConfig		[]CommCbCfgIntStruct	`json:"couchDb_Cfg"`
    } `json:"config"`
}

func InitializeCbConnection(config CommCbIntStruct) error {
	//initializeKTAB()
	fmt.Println("Enter- InitializeCbConnection")
	var err error
	// if couchbaseDetail.Cluster == nil {
	// 	opts := gocb.ClusterOptions{Username: config.CommCbCfgData[0].CouchBaseUname, Password: config.CommCbCfgData[0].CouchBasePwd}
	// 	couchbaseDetail.Cluster, err = gocb.Connect(config.CommCbCfgData[0].CouchBaseIP, opts)
	// 	if err != nil {
	// 		fmt.Println("Exit- Error in connecting to couchbase: ", err)
	// 		// return err
	// 	}
	// }
	// if couchbaseDetail.CollectionObj == nil {
	// 	fmt.Println("Couchbase connection is already established")
	// 	bucket := couchbaseDetail.Cluster.Bucket(config.CommCbCfgData[0].CouchBaseBucketName)
	// 	if bucket == nil {
	// 		fmt.Println("Exit- InitializeCbConnection")
	// 		//return errors.New("Failed to get Bucket")
	// 	}
	// 	err = bucket.WaitUntilReady(time.Second*time.Duration(config.CommCbCfgData[0].ConnectionTimeout), nil)
	// 	if err != nil {
	// 		fmt.Println("Exit- Timeout in connecting to bucket:", err)
	// 		//return err
	// 	}
	// 	couchbaseDetail.CollectionObj = bucket.DefaultCollection()
	// 	fmt.Println("Exit- InitializeCbConnection, connection ready")
	// }

//	initCouchBaseConfig(&gCommonCbCfgData)
	couchbaseDetail.gCDbConnectionList, err = InitCouchDBConnectionv2(config, false)
	fmt.Println(couchbaseDetail.gCDbConnectionList)
	if err != nil {
		//log.Println("InitCouchDBConnection failed", err)
		fmt.Println("InitCouchDBConnection failed", err)

		//raise alarm if DB is down
		fmt.Println("CHFCouchbaseConnectionDown")
		return err
	}
	//log.Println("Init CouchDBConnection Success.")
	fmt.Println("Init CouchDBConnection Success.")
	//fmt.Println(ALARM_TOPIC, "CHFCouchbaseConnectionUp")

	fmt.Println("Exit- InitializeCbConnection")
	fmt.Println(couchbaseDetail)
	return nil
}

//InitCouchDBConnectionv2 Connect to couch and retun local and remote collection
func InitCouchDBConnectionv2(commonCbCfgData CommCbIntStruct, useEphemeralBkt bool) (map[string]*CouchConnectionStruct, error) {

	var lTransID string
	localClusterName := ""
	fmt.Println("Entering InitCouchDBConnection")

	if !useEphemeralBkt {
		gCDbConnectionMap = make(map[string]*CouchConnectionStruct)
	} else {
		gCDbEphConnectionMap = make(map[string]*CouchConnectionStruct)
	}

	couchConfigList := commonCbCfgData.CommCbCfgData

	for _, cdbCfg := range couchConfigList {

		connectionStr := "couchbase://" + cdbCfg.CouchBaseIP
		fmt.Println("Connect to CouchBase with: ", connectionStr)

		if *cdbCfg.IsLocalDc == true {
			localClusterName = cdbCfg.ClusterName
		}

		fmt.Println("Authenticate with Username: ", cdbCfg.CouchBaseUname)
		fmt.Println("Authenticate with Passwd: ", cdbCfg.CouchBasePwd)
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
			fmt.Println("Connect error:", err)
		}
		if connectionData.BucketPtr != nil {
			//Ping Bucket and get total Node count.
			pings, err := connectionData.BucketPtr.Ping(&gocb.PingOptions{
				ReportID:     "couchbase_healthcheck_report",
				ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
			})
			if err != nil {
				fmt.Println("InitCouchDBConnectionv2: Cluster Ping failed During ", err)
				return nil, err
			}
			pingReport := pings.Services[gocb.ServiceTypeKeyValue]
			fmt.Println("InitCouchDBConnectionv2: total Node count in cluster ", len(pingReport))

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

	fmt.Println("gCDbConnectionMap - ", len(gCDbConnectionMap))
	fmt.Println("useEphemeralBkt - ", useEphemeralBkt)
	fmt.Println("gCDbEphConnectionMap - ", len(gCDbEphConnectionMap))

	if (len(gCDbConnectionMap) < 1 && !useEphemeralBkt) || (len(gCDbEphConnectionMap) < 1 && useEphemeralBkt) {
		//return error
		fmt.Println("CouchDBConnection failed for all the dc.")
		errStr := "CouchDBConnection failed for all the dc"
		return nil, errors.New(errStr)
	} else if useEphemeralBkt {
		fmt.Println("InitCouchDBConnection is Success for Ephemeral bucket. gCDbEphConnectionMap:", gCDbEphConnectionMap)
		return gCDbEphConnectionMap, nil
	}

	fmt.Println("InitCouchDBConnection is Success. gCDbConnectionMap:", gCDbConnectionMap)
	return gCDbConnectionMap, nil
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
		//Added by Prashant from here
		if !connectionData.ConnectionStatus{
			fmt.Println("Making Bucket ptr nil")
			connectionData.BucketPtr = nil
		}
		// till here 
		bucket := connectionData.BucketPtr
		if bucket != nil {
			pings, err := bucket.Ping(&gocb.PingOptions{
				ReportID:     "couchbase_healthcheck_report",
				ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
			})
			if err != nil {
				fmt.Println("Healthcheck probe to couchbase failed with err:", err)
				return false
			}
			pingReport, _ := json.Marshal(pings)
			fmt.Println("Couchbase ping report in JSON format:", string(pingReport))

			for service, pingReports := range pings.Services {
				clusterDown := false
				noOfNodesDown := 0

				for _, pingReport := range pingReports {
					if pingReport.State != gocb.PingStateOk {
						fmt.Println("Ping not ok for service:", service)

						//RaiseCouchbaseNodeDownAlarm(connectionData, pingReport.Remote)
						//Increment nodes down count.
						noOfNodesDown++
					} else {
						//RaiseCouchbaseNodeUpAlarm(connectionData, pingReport.Remote)
					}

				}
				//mark Cluster down if number of nodes down is more than or equal to threshold value
				if noOfNodesDown >= connectionData.MaxNodeFailureThresholdCount {
					fmt.Println("Service down for cluster:", pings.ID, "Bucket:", connetionIndex)
					//mark cluster as Down
					func() {
						connectionData.ConnectionStatusLock.Lock()
						connectionData.ConnectionStatus = false
						connectionData.ConnectionStatusLock.Unlock()
					}()
					//RaiseCouchbaseDownAlarm(connectionData)
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
				fmt.Println("Service down for cluster: ", connetionIndex)
				noOfClustersDown++
			} else {
				//Mark service status as up
				connectionData.MarkCouchConnectionActive()
			}
		}
	}
	if noOfClustersDown == len(connectionMap) {
		fmt.Println("All couchbase clusters are down.")
		return false
	}
	return true
}

func connectToCouchDb(connStruct *CouchConnectionStruct, lTransID string) error {
	fmt.Println("Initialising connection --", connStruct.ConnectionStr)
	cDbCluster, err := gocb.Connect(connStruct.ConnectionStr, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: connStruct.UserName,
			Password: connStruct.Password,
		},
		TimeoutsConfig: gocb.TimeoutsConfig{ConnectTimeout: time.Duration(connStruct.ConnectionTimeout) * time.Second},
	})
	if err != nil {
		fmt.Println("Connect error:", err)

		connStruct.MarkCouchConnectionInactive()
		return err
	}
	cDbBucket := cDbCluster.Bucket(connStruct.BucketName)
	//Wait until bucket is ready to use
	err = cDbBucket.WaitUntilReady(time.Duration(connStruct.BucketOpenTimeout)*time.Second, nil)
	 if err != nil || cDbBucket == nil {
	 	fmt.Println("OpenBucket error: ", err)

	 	connStruct.MarkCouchConnectionInactive()
	// //	RaiseCouchbaseDownAlarm(connStruct)
	 	return err
	}
	//update bucket and collection pointers
	connStruct.Collectionptr = cDbBucket.DefaultCollection()
	connStruct.BucketPtr = cDbBucket
	connStruct.ClusterPtr = cDbCluster

	connStruct.MarkCouchConnectionActive()
	return err
}

//CheckForActiveCouchConnection Check if current connection is active. else switch to Remote
//Default: Local DC is set as primary active
func CheckForActiveCouchConnection(connectionList map[string]*CouchConnectionStruct, useEphemeralBkt bool) (activeConnection *CouchConnectionStruct) {
	/*
		1. Check if local dc is up.
		2. If local is down switch to remote.
		3. If remote is down, return primary, even if it is down
	*/
	fmt.Println("NULL", "Entering CheckForActiveCouchConnection.")
	localClusterName := gLocalClusterName
	if useEphemeralBkt {
		localClusterName = gLocalEphClusterName
	}
	//Get primary(local). Use if active
	primary, ok := getActivePrimaryCouchConnection(connectionList, localClusterName)
	if ok {
		activeConnection = primary
		//RaiseCouchbaseUpAlarm(activeConnection)
		fmt.Println("NULL", "CheckForActiveCouchConnection, Using connection ", activeConnection.ConnectionStr)
		return
	}

	activeConnection, ok = getActiveSecondaryCouchConnection(connectionList, localClusterName)
	//Active connection found. Return
	if ok {
		//RaiseCouchbaseUpAlarm(activeConnection)
		fmt.Println("Fetching Secondary Active Connection", localClusterName)
		fmt.Println("NULL", "CheckForActiveCouchConnection, Using connection ", activeConnection.ConnectionStr)
		return
	}

	fmt.Println("NULL", "CheckForActiveCouchConnection, No active connection. Returning Primary")
	//Return non nil connection here. worst case nil
	return returnNonNilConnection(connectionList)
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
	fmt.Println("No active connection found")
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
	fmt.Println("No active connection found")
	return nil, false
}

func checkIfConnectionIsActive(connection *CouchConnectionStruct) bool {
	return connection.ConnectionStatus
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

//CommonCbCfgStruct - defined for all microservices-COMMAS_COMMON_CB_CFG_FILE
// type CommonCbCfgStruct struct {
// 	CommCbIntData CommCbIntStruct `json:"config" validate:"required"`
// }

type CommCbIntStruct struct {
	CommCbCfgData      []CommCbCfgIntStruct `json:"couchDb_Cfg" validate:"required,dive"`
//	CommCbTimerCfgData CommCbTimerIntStruct `json:"timer_Cfg" validate:"required"`
}

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
	//RaiseCouchbaseUpAlarm(conn)
}

func (conn *CouchConnectionStruct) MarkCouchConnectionInactive() {
	if conn == nil {
		return
	}
	conn.ConnectionStatusLock.Lock()
	conn.FailureCount++
	//mark Cluster down if number of nodes down is more than or equal to threshold value
	if conn.FailureCount >= conn.MaxCouchFailureThresholdCount {
		conn.ConnectionStatus = false
		//RaiseCouchbaseDownAlarm(conn)
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

// func switchCBDC(abcd func(), parameters ...interface{}){
// 	activeConnection := CheckForActiveCouchConnection(gCDbConnectionList, false)
// 	Cluster = activeConnection.ClusterPtr       
// 	CollectionObj = activeConnection.BucketPtr.DefaultCollection()
// 	abcd()
// }
const SLEEPCONST time.Duration = 2 * time.Second
const RETRYCOUNTER int = 3


func fetchDocFromCB() error {
	defer wg.Done()
	ticker := time.NewTicker(3 * time.Second)
	for {
        select {
        case _ = <-ticker.C:
			var data map[string]interface{}
			elem, err := couchbaseDetail.CollectionObj.Get("doc", nil)
			if errors.Is(err, gocb.ErrInternalServerFailure) || errors.Is(err, gocb.ErrTimeout) {
				fmt.Println("---------------Switching to secondary DB------------")
				//counter which will take care of retry attempts
				//sleep 
				//time.Sleep(SLEEPCONST)
				//switchCBDC(fetchDocFromCB, data)
			}else{
				elem.Content(&data)
				elemData, _ := json.Marshal(data)
				fmt.Println(string(elemData))
			}
        }
    }

}
func CBClusterHeatlCheck(){
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	for {
        select {
        case _ = <-ticker.C:
			healthFlag := IsCouchbaseHealthOkV2(couchbaseDetail.gCDbConnectionList)
			if healthFlag {
				activeConnection := CheckForActiveCouchConnection(couchbaseDetail.gCDbConnectionList, false)
				couchbaseDetail.Cluster = activeConnection.ClusterPtr       
				couchbaseDetail.CollectionObj = activeConnection.BucketPtr.DefaultCollection()
				couchbaseDetail.BucketName = activeConnection.BucketName
				couchbaseDetail.CouchbaseServer = activeConnection.ConnectionStr
				couchbaseDetail.CouchBaseUname = activeConnection.UserName
				couchbaseDetail.CouchBasePwd = activeConnection.Password
				couchbaseDetail.BucketPtr = activeConnection.BucketPtr
				fmt.Println(activeConnection.ClusterName)  
			} else{
				fmt.Println("ClusterDown")
				os.Exit(1)
			}
        }
    }
}
type myDoc struct {
	Discount float64 
	RetailCharge float64 
	DB string
}

var wg sync.WaitGroup
func main()  {
	wg.Add(1)
	var cdrProconfig CdrConfig
	readFileData, err := ioutil.ReadFile("cbConfig.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdrProconfig)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	var config CommCbIntStruct
	config.CommCbCfgData = cdrProconfig.Config.StandByCbConfig
	
// initialise both connections(primary + Secondary)
	InitializeCbConnection(config)	
// doing health check on both connections
	go CBClusterHeatlCheck()

	// for i := 0; i < RETRYCOUNTER; i++ {
	// 	time.Sleep(SLEEPCONST)
	// 	err := fetchDocFromCB() // query To be performed
	// 	if err == nil{
	// 		break
	// 	}
	// 	if err != nil && i == RETRYCOUNTER-1 {
	// 		fmt.Println("Failed to perform the query, number of retry attempts reached")
	// 		return
	// 	}
	// }
//	go fetchDocFromCB()
	wg.Wait()
}


/*	document := myDoc{Discount: 1050000000,RetailCharge: 1200000000}
	_, errr := CollectionObj.Insert("doc", &document, nil)
	if errr != nil {
		fmt.Println(err)
	}else{
		fmt.Println("doc inserted")
	}
		document := myDoc{Discount: 1050000000,RetailCharge: 1200000000, DB: "secondary"}
	_, errr := CollectionObj.Upsert("doc", &document, &gocb.UpsertOptions{Timeout: 3 * time.Second})
	if errr != nil{
		fmt.Println("Error occured", err)
	}
	*/
