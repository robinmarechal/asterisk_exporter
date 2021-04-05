package cmd

import (
	"errors"
	"testing"

	"github.com/prometheus/common/promlog"
)

var (
	logCfg = &promlog.Config{}
	logger = promlog.New(logCfg)

	cmdRunner = CmdRunner{
		Logger: logger,
		Cmd:    "",
	}
)

//////////////////////////////////////////////////////////////////////////
///////////////////////// NewUptimeInfo
//////////////////////////////////////////////////////////////////////////

func TestNewUptimeInfo_ValidCommandOutput(t *testing.T) {
	sample :=
		`System uptime: 36520
Last reload: 12345`

	upInfo := cmdRunner.newUptimeInfo(sample, nil)

	if upInfo.SystemUptimeSeconds != 36520 {
		t.Errorf("System uptime has not been parsed correctly.\nExpected: %d\nActual: %d", 36520, upInfo.SystemUptimeSeconds)
	}

	if upInfo.LastReloadSeconds != 12345 {
		t.Errorf("Last reload has not been parsed correctly.\nExpected: %d\nActual: %d", 12345, upInfo.LastReloadSeconds)
	}
}

func TestNewUptimeInfo_InvalidCommandOutput(t *testing.T) {
	sample := `System uptime: 36520`

	upInfo := cmdRunner.newUptimeInfo(sample, nil)

	if *upInfo != DefaultUptimeInfo {
		t.Errorf("Uptime info should take default values since the output is not well formatted.")
	}
}

func TestNewUptimeInfo_WhenCommandError(t *testing.T) {
	err := errors.New("default error")

	upInfo := cmdRunner.newUptimeInfo("", err)

	if *upInfo != DefaultUptimeInfo {
		t.Errorf("Uptime info should take default values since the command resulted in error.")
	}
}

//////////////////////////////////////////////////////////////////////////
///////////////////////// NewChannelsInfo
//////////////////////////////////////////////////////////////////////////

func TestNewChannelsInfo_ValidCommandOutput(t *testing.T) {
	sample :=
		`12 active channels
25 active calls
789 calls processed`

	result := cmdRunner.newChannelsInfo(sample, nil)

	if result.ActiveChannels != 12 {
		t.Errorf("ActiveChannels has not been parsed correctly.\nExpected: %d\nActual: %d", 12, result.ActiveChannels)
	}

	if result.ActiveCalls != 25 {
		t.Errorf("ActiveCalls has not been parsed correctly.\nExpected: %d\nActual: %d", 25, result.ActiveCalls)
	}

	if result.ProcessedCalls != 789 {
		t.Errorf("ProcessedCalls has not been parsed correctly.\nExpected: %d\nActual: %d", 789, result.ProcessedCalls)
	}
}

func TestNewChannelsInfo_InvalidCommandOutput(t *testing.T) {
	sample := `12 active channels`

	result := cmdRunner.newChannelsInfo(sample, nil)

	if *result != DefaultChannelsInfo {
		t.Errorf("Channels info should take default values since the output is not well formatted.")
	}
}

func TestNewChannelsInfo_WhenCommandError(t *testing.T) {
	err := errors.New("default error")

	result := cmdRunner.newChannelsInfo("", err)

	if *result != DefaultChannelsInfo {
		t.Errorf("Channels info should take default values since the command resulted in error.")
	}
}

//////////////////////////////////////////////////////////////////////////
///////////////////////// NewPeersInfo
//////////////////////////////////////////////////////////////////////////

func TestNewPeersInfo_ValidCommandOutput_Monitored(t *testing.T) {
	sample :=
		`Name/username             Host                                    Dyn Forcerport Comedia    ACL Port     Status      Description                      
5 sip peers [Monitored: 2 online, 3 offline Unmonitored: 4 online, 5 offline]`

	result := cmdRunner.newPeersInfo(sample, nil)

	if result.SipPeers != 5 {
		t.Errorf("SipPeers has not been parsed correctly.\nExpected: %d\nActual: %d", 5, result.SipPeers)
	}

	if result.MonitoredOnline != 2 {
		t.Errorf("MonitoredOnline has not been parsed correctly.\nExpected: %d\nActual: %d", 2, result.MonitoredOnline)
	}

	if result.MonitoredOffline != 3 {
		t.Errorf("MonitoredOffline has not been parsed correctly.\nExpected: %d\nActual: %d", 3, result.MonitoredOffline)
	}

	if result.UnmonitoredOnline != 4 {
		t.Errorf("UnmonitoredOnline has not been parsed correctly.\nExpected: %d\nActual: %d", 4, result.UnmonitoredOnline)
	}

	if result.UnmonitoredOffline != 5 {
		t.Errorf("UnmonitoredOffline has not been parsed correctly.\nExpected: %d\nActual: %d", 5, result.UnmonitoredOffline)
	}
}

func TestNewPeersInfo_ValidCommandOutput_UnknownAndOkPeers_None(t *testing.T) {
	sample :=
		`Name/username             Host                                    Dyn Forcerport Comedia    ACL Port     Status      Description                      
5 sip peers [Monitored: 2 online, 3 offline Unmonitored: 4 online, 5 offline]`

	result := cmdRunner.newPeersInfo(sample, nil)

	if result.PeersStatusUnknown != 0 {
		t.Errorf("PeersStatusUnknown has not been parsed correctly.\nExpected: %d\nActual: %d", 789, result.PeersStatusUnknown)
	}

	if result.PeersStatusQualified != 0 {
		t.Errorf("PeersStatusQualified has not been parsed correctly.\nExpected: %d\nActual: %d", 789, result.PeersStatusQualified)
	}
}

func TestNewPeersInfo_ValidCommandOutput_UnknownAndOkPeers_SomeLines(t *testing.T) {
	sample :=
		`Name/username             Host                                    Dyn Forcerport Comedia    ACL Port     Status      Description  
12368 dfghjk UNKNOWN fzfze
1259 azerty UNKNOWN jthgezfaz
83348 ltejnoer OK (9875 zetrgefd
76198 zjebz OK (9875) hytgrfed
036 fzf OK (9875) hytgfrde
9763 zefzgezegrz UNKNOWN gege
16962 fzfze OK (9875) fzfbrhtgrfd
5 sip peers [Monitored: 2 online, 3 offline Unmonitored: 4 online, 5 offline]`

	result := cmdRunner.newPeersInfo(sample, nil)

	if result.PeersStatusUnknown != 3 {
		t.Errorf("PeersStatusUnknown has not been parsed correctly.\nExpected: %d\nActual: %d", 789, result.PeersStatusUnknown)
	}

	if result.PeersStatusQualified != 4 {
		t.Errorf("PeersStatusQualified has not been parsed correctly.\nExpected: %d\nActual: %d", 789, result.PeersStatusQualified)
	}
}

func TestNewPeersInfo_InvalidCommandOutput(t *testing.T) {
	sample := `???`

	result := cmdRunner.newPeersInfo(sample, nil)

	expected := PeersInfo{
		SipPeers:             DefaultPeersInfo.SipPeers,
		MonitoredOnline:      DefaultPeersInfo.MonitoredOnline,
		MonitoredOffline:     DefaultPeersInfo.MonitoredOffline,
		UnmonitoredOnline:    DefaultPeersInfo.UnmonitoredOnline,
		UnmonitoredOffline:   DefaultPeersInfo.UnmonitoredOffline,
		PeersStatusUnknown:   0,
		PeersStatusQualified: 0,
	}

	if *result != expected {
		t.Errorf("Peers info should take default values since the output is not well formatted.")
	}
}

func TestNewPeersInfo_WhenCommandError(t *testing.T) {
	err := errors.New("default error")

	result := cmdRunner.newPeersInfo("", err)

	if *result != DefaultPeersInfo {
		t.Errorf("Peers info should take default values since the command resulted in error.")
	}
}

//////////////////////////////////////////////////////////////////////////
///////////////////////// NewThreadsInfo
//////////////////////////////////////////////////////////////////////////

func TestNewThreadsInfo_ValidCommandOutput(t *testing.T) {
	sample :=
		`0x7f67583ae700 3695 netconsole           started at [ 1639] asterisk.c listener()
0x7f67b0713700 18 default_tps_processing_function started at [  202] taskprocessor.c default_listener_start()
0x7f67b0697700 19 bridge_manager_thread started at [ 4869] bridge.c bridge_manager_create()
0x7f67b078f700 17 db_sync_thread       started at [ 1022] db.c astdb_init()
0x7f67b080b700 16 default_tps_processing_function started at [  202] taskprocessor.c default_listener_start()
0x7f67b09fb700 12 logger_thread        started at [ 1595] logger.c init_logger()
0x7f67b0af3700 9 listener             started at [ 1699] asterisk.c ast_makesocket()
0x7f67b4dce700 8 default_tps_processing_function started at [  202] taskprocessor.c default_listener_start()
0x7f67b4e4a700 7 default_tps_processing_function started at [  202] taskprocessor.c default_listener_start()
9 threads listed.`

	result := cmdRunner.newThreadsInfo(sample, nil)

	if result.ThreadCount != 9 {
		t.Errorf("ThreadCount has not been computed correctly.\nExpected: %d\nActual: %d", 9, result.ThreadCount)
	}
}

func TestNewThreadsInfo_InvalidCommandOutput_0Thread(t *testing.T) {
	sample :=
		`0 threads listed.`

	result := cmdRunner.newThreadsInfo(sample, nil)

	if result.ThreadCount != 0 {
		t.Errorf("ThreadCount has not been computed correctly.\nExpected: %d\nActual: %d", 0, result.ThreadCount)
	}
}

func TestNewThreadsInfo_InvalidCommandOutput_1Thread(t *testing.T) {
	sample :=
		`0x7f67b09fb700 12 logger_thread        started at [ 1595] logger.c init_logger()
1 threads listed.`

	result := cmdRunner.newThreadsInfo(sample, nil)

	if result.ThreadCount != 1 {
		t.Errorf("ThreadCount has not been computed correctly.\nExpected: %d\nActual: %d", 1, result.ThreadCount)
	}
}

func TestNewThreadsInfo_WhenCommandError(t *testing.T) {
	err := errors.New("default error")

	result := cmdRunner.newThreadsInfo("", err)

	if *result != DefaultThreadsInfo {
		t.Errorf("Thread info should take default values since the command resulted in error.")
	}
}

func TestNewAgentsInfo(t *testing.T) {
	// agent show all
	sample := `Agent-ID Name                 State       Channel                        Talking with
Defined agents: 5, Logged in: 3, Talking: 1`
	var expected int64

	result := cmdRunner.newAgentsInfo(sample, nil)

	expected = 5
	if result.DefinedAgents != expected {
		t.Errorf("DefinedAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.DefinedAgents)
	}

	expected = 3
	if result.LoggedAgents != expected {
		t.Errorf("LoggedAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.LoggedAgents)
	}

	expected = 1
	if result.TalkingAgents != expected {
		t.Errorf("TalkingAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.TalkingAgents)
	}
}

func TestNewOnlineAgentsInfo(t *testing.T) {
	// agent show online

	sample := `Agent-ID Name                 State       Channel                        Talking with
Defined agents: 5, Logged in: 3, Talking: 1`
	var expected int64

	result := cmdRunner.newOnlineAgentsInfo(sample, nil)

	expected = 5
	if result.OnlineDefinedAgents != expected {
		t.Errorf("OnlineDefinedAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.OnlineDefinedAgents)
	}

	expected = 3
	if result.OnlineLoggedAgents != expected {
		t.Errorf("OnlineLoggedAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.OnlineLoggedAgents)
	}

	expected = 1
	if result.OnlineTalkingAgents != expected {
		t.Errorf("OnlineTalkingAgents has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.OnlineTalkingAgents)
	}
}

func TestNewBridgesInfo(t *testing.T) {
	// bridge show all
	sample := `Bridge-ID                            Chans Type            Technology
foo									 bar   1			   ertyu
foo									 bar   1			   ertyu
foo									 bar   1			   ertyu`
	var expected int64

	result := cmdRunner.newBridgesInfo(sample, nil)

	expected = 3
	if result.Count != expected {
		t.Errorf("BridgesInfo has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.Count)
	}
}

func TestNewBridgeTechnologiesInfo(t *testing.T) {
	// bridge show all
	sample := `Name                 Type                 Priority Suspended
softmix              MultiMix                   10 No
holding_bridge       Holding                    50 No
simple_bridge        1to1Mix                    50 Yes
native_rtp           Native                     90 No`
	var expected string

	result := cmdRunner.newBridgeTechnologiesInfo(sample, nil)

	if len(result.BridgeTechnologies) != 4 {
		t.Errorf("Function newBridgeTechnologiesInfo has not parsed output correctly.\nExpected size: 4\nActual: %d", len(result.BridgeTechnologies))
	}

	// 0

	expected = "softmix"
	if result.BridgeTechnologies[0].Name != expected {
		t.Errorf("BridgeTechnologies[0].Name has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[0].Name)
	}

	expected = "MultiMix"
	if result.BridgeTechnologies[0].Type != expected {
		t.Errorf("BridgeTechnologies[0].Type has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[0].Type)
	}

	expected = "10"
	if result.BridgeTechnologies[0].Priority != expected {
		t.Errorf("BridgeTechnologies[0].Priority has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[0].Priority)
	}

	expected = "No"
	if result.BridgeTechnologies[0].Suspended != expected {
		t.Errorf("BridgeTechnologies[0].Suspended has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[0].Suspended)
	}

	// 2

	expected = "simple_bridge"
	if result.BridgeTechnologies[2].Name != expected {
		t.Errorf("BridgeTechnologies[2].Name has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[2].Name)
	}

	expected = "1to1Mix"
	if result.BridgeTechnologies[2].Type != expected {
		t.Errorf("BridgeTechnologies[2].Type has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[2].Type)
	}

	expected = "50"
	if result.BridgeTechnologies[2].Priority != expected {
		t.Errorf("BridgeTechnologies[2].Priority has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[2].Priority)
	}

	expected = "Yes"
	if result.BridgeTechnologies[2].Suspended != expected {
		t.Errorf("BridgeTechnologies[2].Suspended has not been computed correctly.\nExpected: %s\nActual: %s", expected, result.BridgeTechnologies[2].Suspended)
	}
}

func TestNewCalendarsInfo(t *testing.T) {
	// calendar show calendars
	sample := `Calendar             Type       Status
	--------             ----       ------
	cal1				 typ		0
	cal2				 typ2       2`
	var expected int64

	result := cmdRunner.newCalendarsInfo(sample, nil)

	expected = 2
	if result.Count != expected {
		t.Errorf("Count has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.Count)
	}
}

func TestNewConfBridgeMenus(t *testing.T) {
	// confbridge show menus
	sample := `--------- Menus -----------
sample_admin_menu
default_menu
sample_user_menu`
	var expected string

	result := cmdRunner.newConfBridgeMenus(sample, nil)

	if len(result) != 3 {
		t.Errorf("Function newConfBridgeMenus has not parsed output correctly.\nExpected size: 3\nActual: %d", len(result))
	}

	expected = "sample_admin_menu"
	if result[0] != expected {
		t.Errorf("ConfBridgeMenus has not been computed correctly.\nExpected: %s\nActual: %s", expected, result[0])
	}

	expected = "default_menu"
	if result[1] != expected {
		t.Errorf("ConfBridgeMenus has not been computed correctly.\nExpected: %s\nActual: %s", expected, result[1])
	}

	expected = "sample_user_menu"
	if result[2] != expected {
		t.Errorf("ConfBridgeMenus has not been computed correctly.\nExpected: %s\nActual: %s", expected, result[2])
	}
}

func TestNewConfBridgeProfiles(t *testing.T) {
	// confbridge show profile bridges
	sample := `--------- Bridge Profiles -----------
default_bridge`
	var expected string

	result := cmdRunner.newConfBridgeProfiles(sample, nil)

	if len(result) != 1 {
		t.Errorf("Function newConfBridgeProfiles has not parsed output correctly.\nExpected size: 1\nActual: %d", len(result))
	}

	expected = "default_bridge"
	if result[0] != expected {
		t.Errorf("ConfBridgeProfiles has not been computed correctly.\nExpected: %s\nActual: %s", expected, result[0])
	}
}

func TestNewConfBridgeUsers(t *testing.T) {
	// confbridge show profile users
	sample := `--------- User Profiles -----------
default_user`
	var expected string

	result := cmdRunner.newConfBridgeUsers(sample, nil)

	if len(result) != 1 {
		t.Errorf("Function newConfBridgeUsers has not parsed output correctly.\nExpected size: 1\nActual: %d", len(result))
	}

	expected = "default_user"
	if result[0] != expected {
		t.Errorf("ConfBridgeUsers has not been computed correctly.\nExpected: %s\nActual: %s", expected, result[0])
	}
}

func TestNewChannelTypesInfo(t *testing.T) {
	// core show channeltypes
	sample := `Type             Description                              Devicestate  Indications  Transfer    
-----------      -----------                              -----------  -----------  ----------- 
Recorder         Bridge Media Recording Channel Driver    no           yes          no          
Announcer        Bridge Media Announcing Channel Driver   yes           yes          no          
CBAnn            Conference Bridge Announcing Channel     no           no          yes          
----------
3 channel drivers registered.`
	var expected ChannelTypesInfo

	result := cmdRunner.newChannelTypesInfo(sample, nil)

	if len(result.ChannelTypes) != 3 {
		t.Errorf("Function newChannelTypesInfo has not parsed output correctly.\nExpected size: 3\nActual: %d", len(result.ChannelTypes))
	}

	// 0

	expected = ChannelTypesInfo{
		ChannelTypes: []ChannelType{
			{
				Type:        "Recorder",
				DeviceState: false,
				Indications: true,
				Transfer:    false,
			},
			{
				Type:        "Announcer",
				DeviceState: true,
				Indications: true,
				Transfer:    false,
			},
			{
				Type:        "CBAnn",
				DeviceState: false,
				Indications: false,
				Transfer:    true,
			},
		},
	}

	for i, exp := range expected.ChannelTypes {
		if result.ChannelTypes[i] != exp {
			if result.ChannelTypes[i].Type != exp.Type {
				t.Errorf("ChannelTypes[%d].Type has not been computed correctly.\nExpected: %s\nActual: %s", i, exp.Type, result.ChannelTypes[i].Type)
			}
			if result.ChannelTypes[i].DeviceState != exp.DeviceState {
				t.Errorf("ChannelTypes[%d].DeviceState has not been computed correctly.\nExpected: %t\nActual: %t", i, exp.DeviceState, result.ChannelTypes[i].DeviceState)
			}
			if result.ChannelTypes[i].Indications != exp.Indications {
				t.Errorf("ChannelTypes[%d].Indications has not been computed correctly.\nExpected: %t\nActual: %t", i, exp.Indications, result.ChannelTypes[i].Indications)
			}
			if result.ChannelTypes[i].Transfer != exp.Transfer {
				t.Errorf("ChannelTypes[%d].Transfer has not been computed correctly.\nExpected: %t\nActual: %t", i, exp.Transfer, result.ChannelTypes[i].Transfer)
			}
		}
	}

}

func TestNewImagesInfo(t *testing.T) {
	// core show image formats
	sample := `      Name Extensions                                        Description     Format
	---- ----------                                        -----------     ------
3 image formats registered.`
	var expected int64

	result := cmdRunner.newImagesInfo(sample, nil)

	expected = 3
	if result.Registered != expected {
		t.Errorf("ImagesInfo has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.Registered)
	}
}

func TestNewSystemInfo(t *testing.T) {
	// core show sysinfo
	sample := `
	System Statistics
	-----------------
	  System Uptime:             61 hours
	  Total RAM:                 5915724 KiB
	  Free RAM:                  148876 KiB
	  Buffer RAM:                228980 KiB
	  Total Swap Space:          786428 KiB
	  Free Swap Space:           786428 KiB
	
	  Number of Processes:       294
	
	`

	result := cmdRunner.newSystemInfo(sample, nil)

	expected := SystemInfo{
		TotalMemory:  5915724 * 1024,
		FreeMemory:   148876 * 1024,
		BufferMemory: 228980 * 1024,
		TotalSwap:    786428 * 1024,
		FreeSwap:     786428 * 1024,
		ProcessCount: 294,
	}

	if *result != expected {
		t.Errorf("SystemInfo has not been computed correctly.\nExpected: %d\nActual: %d", expected, *result)
	}
}

func TestNewTaskProcessorsInfo(t *testing.T) {
	// core show taskprocessors
	sample := `
	Processor                                      Processed   In Queue  Max Depth  Low water High water
app_voicemail                                          0          0          0        450        500
ast_msg_queue                                          0          3          0        450        500
CCSS_core                                              1          0          1        450        500
hep_queue_tp                                           0          0          0        450        500
subm:ast_system-00000006                              5          2         15        450        500
subm:ast_system-00000041                              6          7          5        450        500
subm:ast_system-00000043                              7          0          5        450        500

7 taskprocessors
`

	result := cmdRunner.newTaskProcessorsInfo(sample, nil)

	expected := TaskProcessorsInfo{
		ProcessorCounter:    7,
		ProcessedTasksTotal: 19,
		InQueue:             12,
	}

	if *result != expected {
		t.Errorf("TaskProcessorsInfo has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.ProcessorCounter)
	}
}

func TestNewVersionInfo(t *testing.T) {
	// core show version
	sample := `Asterisk certified/13.8-cert4 built by root @ 1b0d6163fdc2 on a x86_64 running Linux on 2017-09-01 18:37:56 UTC`

	result := cmdRunner.newVersionInfo(sample, nil)

	if result.Version != sample {
		t.Errorf("VersionInfo has not been computed correctly.\nExpected: %s\nActual: %s", sample, result.Version)
	}
}

func TestNewIaxChannelsInfo(t *testing.T) {
	// iax2 show channels
	sample := `Channel               Peer                                      Username    ID (Lo/Rem)  Seq (Tx/Rx)  Lag      Jitter  JitBuf  Format  FirstMsg    LastMsg
7 active IAX channels`

	result := cmdRunner.newIaxChannelsInfo(sample, nil)

	expected := int64(7)
	if result.ActiveCount != expected {
		t.Errorf("ActiveCount has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.ActiveCount)
	}
}

func TestNewModulesInfo(t *testing.T) {
	// module show
	sample := `Module                         Description                              Use Count  Status      Support Level
app_agent_pool.so              Call center agent pool applications      0          Running              core
app_authenticate.so            Authentication Application               0          Running              core
app_bridgewait.so              Place the channel into a holding bridge  0          Running              core
app_cdr.so                     Tell Asterisk to not maintain a CDR for  0          Running              core
app_celgenuserevent.so         Generate an User-Defined CEL event       0          Running              core
res_timing_timerfd.so          Timerfd Timing Interface                 1          Running              core
6 modules loaded`

	result := cmdRunner.newModulesInfo(sample, nil)

	expected := int64(6)
	if result.ModuleCount != expected {
		t.Errorf("ModuleCount has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.ModuleCount)
	}
}

func TestNewActiveSipDialogs(t *testing.T) {
	// sip show channels
	sample := `Peer             User/ANR         Call ID          Format           Hold     Last Message    Expiry     Peer      
7 active SIP dialogs`

	result := cmdRunner.newActiveSipDialogs(sample, nil)

	expected := int64(7)
	if result != expected {
		t.Errorf("ActiveSipDialogs has not been computed correctly.\nExpected: %d\nActual: %d", expected, result)
	}
}

func TestNewActiveSipSubscriptions(t *testing.T) {
	// sip show subscriptions
	sample := `Peer             User             Call ID          Extension        Last state     Type            Mailbox    Expiry
7 active SIP subscriptions`

	result := cmdRunner.newActiveSipSubscriptions(sample, nil)

	expected := int64(7)
	if result != expected {
		t.Errorf("ActiveSipSubscriptions has not been computed correctly.\nExpected: %d\nActual: %d", expected, result)
	}
}

func TestNewActiveSipChannels(t *testing.T) {
	// sip show channelstats
	sample := `Peer             Call ID      Duration Recv: Pack  Lost       (     %) Jitter Send: Pack  Lost       (     %) Jitter
7 active SIP channels`

	result := cmdRunner.newActiveSipChannels(sample, nil)

	expected := int64(7)
	if result != expected {
		t.Errorf("ActiveSipChannels has not been computed correctly.\nExpected: %d\nActual: %d", expected, result)
	}
}

func TestNewUsersInfo(t *testing.T) {
	// sip show users
	sample := `Username                   Secret           Accountcode      Def.Context      ACL  Forcerport
foo						   bar				qwe				 abc  			  1	   zef
foo						   bar				qwe				 abc  			  1	   zef
foo						   bar				qwe				 abc  			  1	   zef`

	result := cmdRunner.newUsersInfo(sample, nil)

	expected := int64(3)
	if result.Users != expected {
		t.Errorf("Users has not been computed correctly.\nExpected: %d\nActual: %d", expected, result.Users)
	}
}
