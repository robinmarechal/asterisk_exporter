package cmd

import (
	"errors"
	"strings"

	"github.com/docker/go-units"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/robinmarechal/asteriskk_exporter/util"
)

func (c *CmdRunner) logError(value int64, err error) int64 {
	if err != nil {
		level.Error(*c.Logger).Log("err", err)
		return -1
	}

	return value
}

////////// UTILS

func (c *CmdRunner) newUptimeInfo(out string, err error) *UptimeInfo {
	if err != nil {
		return &DefaultUptimeInfo
	}

	lines := strings.Split(out, "\n")

	length := len(lines)
	if length != 2 {
		level.Error(*c.Logger).Log("err", "Uptime command is not well formatted.", "output", out)
		return &DefaultUptimeInfo
	}

	return &UptimeInfo{
		SystemUptimeSeconds: util.ExtractTrailingValueAfterColon(lines[0], c.Logger),
		LastReloadSeconds:   util.ExtractTrailingValueAfterColon(lines[1], c.Logger),
	}
}

func (c *CmdRunner) newChannelsInfo(out string, err error) *ChannelsInfo {
	if err != nil {
		return &DefaultChannelsInfo
	}

	lines := strings.Split(out, "\n")

	length := len(lines)
	if length != 3 {
		level.Error(*c.Logger).Log("err", "Channels command is not well formatted.", "output", out)
		return &DefaultChannelsInfo
	}

	return &ChannelsInfo{
		ActiveChannels: util.ExtractLeadingInteger(lines[0], c.Logger),
		ActiveCalls:    util.ExtractLeadingInteger(lines[1], c.Logger),
		ProcessedCalls: util.ExtractLeadingInteger(lines[2], c.Logger),
	}
}

func (c *CmdRunner) newPeersInfo(out string, err error) *PeersInfo {
	// asterisk -rx 'sip show peers' | grep 'sip peers' | grep 'Monitored' | grep 'Unmonitored'"
	// [sip_peers, monitored_online, monitored_offline, unmonitored_online, unmonitored_offline] = re.findall("\d+", sip_show_peers)

	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*UNKNOWN\s' | wc -l"
	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*OK\s\(\d+' | wc -l"
	// 123 zefzgezegrz UNKNOWN
	// 642 OK (9875

	if err != nil {
		level.Error(*c.Logger).Log("err", err)
		return &DefaultPeersInfo
	}

	obj := PeersInfo{}

	lines := strings.Split(out, "\n")
	peersInfoLine, err2 := extractPeersInfoLine(lines)

	if peersInfoLine == "" || err2 != nil {
		if peersInfoLine == "" {
			level.Error(*c.Logger).Log("err", "Error, command output is empty", "cmd", "show sip peers")
		} else {
			level.Error(*c.Logger).Log("err", err2, "cmd", "show sip peers")
		}
		setPeersMonitoringInfoToDefault(&obj)
	} else {
		// asterisk -rx 'sip show peers' | grep 'sip peers' | grep 'Monitored' | grep 'Unmonitored'"
		// [sip_peers, monitored_online, monitored_offline, unmonitored_online, unmonitored_offline] = re.findall("\d+", sip_show_peers)
		errors := setPeersInfoFromMonitoringInfoLine(c.Logger, &obj, peersInfoLine)

		for _, err := range *errors {
			if err != nil {
				level.Error(*c.Logger).Log("err", err)
			}
		}
	}

	setUnknownAndOkPeersCount(c.Logger, &obj, lines)

	return &obj
}

func setPeersInfoFromMonitoringInfoLine(logger *log.Logger, obj *PeersInfo, line string) *[]error {
	errors := make([]error, 5)

	submatchall := AllNumbersRegexp.FindAllString(line, 5)

	obj.SipPeers, errors[0] = util.StrToInt(submatchall[0])
	obj.MonitoredOnline, errors[1] = util.StrToInt(submatchall[1])
	obj.MonitoredOffline, errors[2] = util.StrToInt(submatchall[2])
	obj.UnmonitoredOnline, errors[3] = util.StrToInt(submatchall[3])
	obj.UnmonitoredOffline, errors[4] = util.StrToInt(submatchall[4])

	return &errors
}

func setPeersMonitoringInfoToDefault(obj *PeersInfo) {
	obj.SipPeers = DefaultPeersInfo.SipPeers
	obj.MonitoredOnline = DefaultPeersInfo.MonitoredOnline
	obj.MonitoredOffline = DefaultPeersInfo.MonitoredOffline
	obj.UnmonitoredOnline = DefaultPeersInfo.UnmonitoredOnline
	obj.UnmonitoredOffline = DefaultPeersInfo.UnmonitoredOffline
}

func setUnknownAndOkPeersCount(logger *log.Logger, obj *PeersInfo, lines []string) {
	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*UNKNOWN\s' | wc -l"
	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*OK\s\(\d+' | wc -l"
	obj.PeersStatusQualified = 0
	obj.PeersStatusUnknown = 0

	for _, line := range lines {
		if strings.Contains(line, "UNKNOWN") {
			obj.PeersStatusUnknown++
		} else if strings.Contains(line, "OK") {
			obj.PeersStatusQualified++
		}
	}
}

func extractPeersInfoLine(lines []string) (string, error) {
	for _, line := range lines {
		if strings.Contains(line, "sip peers") &&
			strings.Contains(line, "Monitored:") &&
			strings.Contains(line, "Unmonitored:") {
			return line, nil
		}
	}

	return "", errors.New("not found peers info line in provided lines")
}

func (c *CmdRunner) newThreadsInfo(out string, err error) *ThreadsInfo {
	// asterisk -rx 'core show threads' | tail -1 | cut -d' ' -f1"
	// Or
	// asterisk -rx 'core show threads'
	// => split \n, length -1

	if err != nil {
		return &DefaultThreadsInfo
	}

	lines := strings.Split(out, "\n")
	lastLine := lines[len(lines)-1]

	parts := strings.Split(lastLine, " ")
	valueStr := parts[0]

	intValue, err := util.StrToInt(valueStr)

	if err != nil {
		level.Error(*c.Logger).Log("err", err, "line", lastLine, "value", valueStr)
		return &DefaultThreadsInfo
	}

	return &ThreadsInfo{
		ThreadCount: intValue,
	}
}

func (c *CmdRunner) newAgentsInfo(out string, err error) *AgentsInfo {
	if err != nil {
		level.Error(*c.Logger).Log("err", err)
		return &DefaultAgentsInfo
	}

	lastLine := util.ExtractLastLine(out)

	if lastLine == "" {
		level.Error(*c.Logger).Log("err", "Error, command output is empty")
		return &DefaultAgentsInfo
	}

	// Defined agents: 5, Logged in: 3, Talking: 1
	submatchall := AllIntegersRegexp.FindAllString(lastLine, 3)

	return &AgentsInfo{
		DefinedAgents: util.StrToIntOrDefault(c.Logger, submatchall[0], -1),
		LoggedAgents:  util.StrToIntOrDefault(c.Logger, submatchall[1], -1),
		TalkingAgents: util.StrToIntOrDefault(c.Logger, submatchall[2], -1),
	}
}

func (c *CmdRunner) newOnlineAgentsInfo(out string, err error) *OnlineAgentsInfo {
	if err != nil {
		level.Error(*c.Logger).Log("err", err)
		return &DefaultOnlineAgentsInfo
	}

	lastLine := util.ExtractLastLine(out)

	if lastLine == "" {
		level.Error(*c.Logger).Log("err", "Error, command output is empty")
		return &DefaultOnlineAgentsInfo
	}

	// Defined agents: 5, Logged in: 3, Talking: 1
	submatchall := AllIntegersRegexp.FindAllString(lastLine, 3)

	return &OnlineAgentsInfo{
		OnlineDefinedAgents: util.StrToIntOrDefault(c.Logger, submatchall[0], -1),
		OnlineLoggedAgents:  util.StrToIntOrDefault(c.Logger, submatchall[1], -1),
		OnlineTalkingAgents: util.StrToIntOrDefault(c.Logger, submatchall[2], -1),
	}
}

func (c *CmdRunner) newBridgesInfo(out string, err error) *BridgesInfo {
	if err != nil {
		level.Error(*c.Logger).Log("err", err)
		return &DefaultBridgesInfo
	}

	count := util.CountLines(out) - 1

	return &BridgesInfo{
		Count: int64(count),
	}
}

func (c *CmdRunner) newBridgeTechnologiesInfo(out string, err error) *BridgeTechnologiesInfo {
	if err != nil {
		return &DefaultBridgeTechnologiesInfo
	}

	out = strings.TrimSuffix(out, "\n")
	lines := strings.Split(out, "\n")

	results := BridgeTechnologiesInfo{
		BridgeTechnologies: make([]BridgeTechnology, len(lines)-1),
	}

	for i := 1; i < len(lines); i++ {
		matches := StringWithoutWhitespaceRegexp.FindAllString(lines[i], 4)

		results.BridgeTechnologies[i-1] = BridgeTechnology{
			Name:      matches[0],
			Type:      matches[1],
			Priority:  matches[2],
			Suspended: matches[3],
		}
	}

	return &results
}

func (c *CmdRunner) newCalendarsInfo(out string, err error) *CalendarsInfo {
	if err != nil {
		return &DefaultCalendarsInfo
	}

	// Calendar             Type       Status
	// --------             ----       ------
	// cal1				 typ		0
	// cal2				 typ2       2

	return &CalendarsInfo{
		Count: int64(util.CountLines(out)) - 2,
	}
}

func (c *CmdRunner) newConfBridgeMenus(out string, err error) []string {
	if err != nil {
		return []string{}
	}

	// --------- Menus -----------
	// sample_admin_menu
	// default_menu
	// sample_user_menu

	lines := strings.Split(out, "\n")
	if len(lines) <= 1 {
		return []string{}
	}

	return lines[1:]
}

func (c *CmdRunner) newConfBridgeProfiles(out string, err error) []string {
	return c.newConfBridgeMenus(out, err)
}

func (c *CmdRunner) newConfBridgeUsers(out string, err error) []string {
	return c.newConfBridgeMenus(out, err)
}

func (c *CmdRunner) newChannelTypesInfo(out string, err error) *ChannelTypesInfo {
	if err != nil {
		return &DefaultChannelTypesInfo
	}

	// Type             Description                              Devicestate  Indications  Transfer
	// -----------      -----------                              -----------  -----------  -----------
	// Recorder         Bridge Media Recording Channel Driver    no           yes          no
	// Announcer        Bridge Media Announcing Channel Driver   yes           yes          no
	// CBAnn            Conference Bridge Announcing Channel     no           no          yes
	// ----------
	// 3 channel drivers registered.

	nbLines := util.CountLines(out)
	nbChannelTypes := nbLines - 4

	results := ChannelTypesInfo{
		ChannelTypes: make([]ChannelType, nbChannelTypes),
	}

	lines := strings.Split(out, "\n")
	header := lines[0]

	stateIdx := strings.Index(header, "Devicestate")

	for i := 2; i < len(lines)-2; i++ {
		line := lines[i]
		matches := YesNoRegexp.FindAllString(line[stateIdx:], 3)

		results.ChannelTypes[i-2] = ChannelType{
			Type:        line[:strings.Index(line, " ")], // start until first space
			DeviceState: matches[0] == "yes",
			Indications: matches[1] == "yes",
			Transfer:    matches[2] == "yes",
		}
	}

	return &results
}

func (c *CmdRunner) newImagesInfo(out string, err error) *ImagesInfo {
	if err != nil {
		return &DefaultImagesInfo
	}

	lastLine := util.ExtractLastLine(out)
	v := util.ExtractLeadingInteger(lastLine, c.Logger)

	return &ImagesInfo{
		Registered: v,
	}
}

func (c *CmdRunner) newSystemInfo(out string, err error) *SystemInfo {
	if err != nil {
		return &DefaultSystemInfo
	}

	//
	// System Statistics
	// -----------------
	// System Uptime:             9 hours
	// Total RAM:                 12959084 KiB
	// Free RAM:                  8689472 KiB
	// Buffer RAM:                184724 KiB
	// Total Swap Space:          4194304 KiB
	// Free Swap Space:           4194304 KiB

	// Number of Processes:       672
	//

	result := SystemInfo{
		TotalMemory:  -1,
		FreeMemory:   -1,
		BufferMemory: -1,
		TotalSwap:    -1,
		FreeSwap:     -1,
		ProcessCount: -1,
	}

	out = strings.TrimSuffix(out, "\n")
	lines := strings.Split(out, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		kv := strings.Split(line, ":")

		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		value := strings.TrimSpace(kv[1])

		switch key {
		case "Total RAM":
			result.TotalMemory = c.logError(units.RAMInBytes(value))
		case "Free RAM":
			result.FreeMemory = c.logError(units.RAMInBytes(value))
		case "Buffer RAM":
			result.BufferMemory = c.logError(units.RAMInBytes(value))
		case "Total Swap Space":
			result.TotalSwap = c.logError(units.RAMInBytes(value))
		case "Free Swap Space":
			result.FreeSwap = c.logError(units.RAMInBytes(value))
		case "Number of Processes":
			result.ProcessCount = util.StrToIntOrDefault(c.Logger, value, 0)
		}
	}

	return &result
}

func (c *CmdRunner) newTaskProcessorsInfo(out string, err error) *TaskProcessorsInfo {
	if err != nil {
		return &DefaultTaskProcessorsInfo
	}

	out = strings.TrimSuffix(out, "\n")
	lines := strings.Split(out, "\n")

	var sumProcessed int64 = 0
	var sumInQueue int64 = 0
	var count int64 = 0

	for i := 1; i < len(lines)-1; i++ {
		matches := StringWithoutWhitespaceRegexp.FindAllString(lines[i], 3)

		sumProcessed += util.StrToIntOrDefault(c.Logger, matches[1], 0)
		sumInQueue += util.StrToIntOrDefault(c.Logger, matches[2], 0)
		count++
	}

	return &TaskProcessorsInfo{
		ProcessorCounter:    count,
		ProcessedTasksTotal: sumProcessed,
		InQueue:             sumInQueue,
	}
}

func (c *CmdRunner) newVersionInfo(out string, err error) *VersionInfo {
	if err != nil {
		return &DefaultVersionInfo
	}

	// Asterisk certified/13.8-cert4 built by root @ 1b0d6163fdc2 on a x86_64 running Linux on 2017-09-01 18:37:56 UTC

	return &VersionInfo{
		Version: out,
	}
}

func (c *CmdRunner) newIaxChannelsInfo(out string, err error) *IaxChannelsInfo {
	if err != nil {
		return &DefaultIaxChannelsInfo
	}

	// Channel               Peer                                      Username    ID (Lo/Rem)  Seq (Tx/Rx)  Lag      Jitter  JitBuf  Format  FirstMsg    LastMsg
	// 7 active IAX channels

	lastLine := util.ExtractLastLine(out)
	v := util.ExtractLeadingInteger(lastLine, c.Logger)

	return &IaxChannelsInfo{
		ActiveCount: v,
	}
}

func (c *CmdRunner) newModulesInfo(out string, err error) *ModulesInfo {
	if err != nil {
		return &DefaultModulesInfo
	}

	// Module                         Description                              Use Count  Status      Support Level
	// app_agent_pool.so              Call center agent pool applications      0          Running              core
	// app_authenticate.so            Authentication Application               0          Running              core
	// app_bridgewait.so              Place the channel into a holding bridge  0          Running              core
	// app_cdr.so                     Tell Asterisk to not maintain a CDR for  0          Running              core
	// app_celgenuserevent.so         Generate an User-Defined CEL event       0          Running              core
	// res_timing_timerfd.so          Timerfd Timing Interface                 1          Running              core
	// 6 modules loaded

	lastLine := util.ExtractLastLine(out)
	v := util.ExtractLeadingInteger(lastLine, c.Logger)

	return &ModulesInfo{
		ModuleCount: v,
	}
}

func (c *CmdRunner) newActiveSipDialogs(out string, err error) int64 {
	if err != nil {
		return DefaultActiveSipDialogs
	}

	lastLine := util.ExtractLastLine(out)
	return util.ExtractLeadingInteger(lastLine, c.Logger)

}

func (c *CmdRunner) newActiveSipSubscriptions(out string, err error) int64 {
	if err != nil {
		return DefaultActiveSipSubscriptions
	}

	lastLine := util.ExtractLastLine(out)
	return util.ExtractLeadingInteger(lastLine, c.Logger)

}

func (c *CmdRunner) newActiveSipChannels(out string, err error) int64 {
	if err != nil {
		return DefaultActiveSipChannels
	}

	lastLine := util.ExtractLastLine(out)
	return util.ExtractLeadingInteger(lastLine, c.Logger)

}

func (c *CmdRunner) newUsersInfo(out string, err error) *UsersInfo {
	if err != nil {
		return &DefaultUsersInfo
	}

	return &UsersInfo{
		Users: int64(util.CountLines(out)) - 1,
	}
}

// core show calls + seconds ???
