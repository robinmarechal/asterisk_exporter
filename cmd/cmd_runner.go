package cmd

import (
	"bytes"
	"os/exec"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/robinmarechal/asteriskk_exporter/util"
)

//////////////////////////////////////////////////////////////////////////
///////////////////////// STRUCTS
//////////////////////////////////////////////////////////////////////////

// CmdRunner command struct
type CmdRunner struct {
	Logger log.Logger
	Cmd    string
}

// ChannelsInfo Channels and calls infos
type ChannelsInfo struct {
	ActiveChannels int64
	ActiveCalls    int64
	ProcessedCalls int64
}

// UptimeInfo uptime and reload time infos
type UptimeInfo struct {
	SystemUptimeSeconds int64
	LastReloadSeconds   int64
}

// PeersInfo peers infos
type PeersInfo struct {
	// asterisk -rx 'sip show peers' | grep 'sip peers' | grep 'Monitored' | grep 'Unmonitored'"
	// [sip_peers, monitored_online, monitored_offline, unmonitored_online, unmonitored_offline] = re.findall("\d+", sip_show_peers)
	SipPeers           int64
	MonitoredOnline    int64
	MonitoredOffline   int64
	UnmonitoredOnline  int64
	UnmonitoredOffline int64
	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*UNKNOWN\s' | wc -l"
	PeersStatusUnknown int64
	// asterisk -rx 'sip show peers' | grep -P '^\d{3,}.*OK\s\(\d+' | wc -l"
	PeersStatusQualified int64
}

// ThreadsInfo threads infos
type ThreadsInfo struct {
	ThreadCount int64
}

type AgentsInfo struct {
	//agent show all
	DefinedAgents int64
	LoggedAgents  int64
	TalkingAgents int64
}
type OnlineAgentsInfo struct {
	//agent show online
	OnlineDefinedAgents int64
	OnlineLoggedAgents  int64
	OnlineTalkingAgents int64
}

type BridgesInfo struct {
	// bridge show all
	Count int64
}

type BridgeTechnologiesInfo struct {
	// bridge technology show
	BridgeTechnologies []BridgeTechnology
}

type BridgeTechnology struct {
	Name      string
	Type      string
	Priority  string
	Suspended string
}

type CalendarsInfo struct {
	// calendar show calendars
	// calendar show types
	Count int64
}

type ConfBridgeInfo struct {
	// confbridge show menus
	Menus []string
	// confbridge show profile bridges
	Profiles []string
	// confbridge show profile users
	Users []string
}

type ChannelTypesInfo struct {
	// core show channeltypes
	ChannelTypes []ChannelType
}

type ChannelType struct {
	Type        string
	DeviceState bool
	Indications bool
	Transfer    bool
}

type ImagesInfo struct {
	// core show image formats
	Registered int64
}

type SystemInfo struct {
	// core show sysinfo
	TotalMemory  int64
	FreeMemory   int64
	BufferMemory int64
	TotalSwap    int64
	FreeSwap     int64
	ProcessCount int64
}

type TaskProcessorsInfo struct {
	// core show taskprocessors
	ProcessorCounter    int64
	ProcessedTasksTotal int64
	InQueue             int64
}

type VersionInfo struct {
	// core show version
	Version string
}

type IaxChannelsInfo struct {
	// iax2 show channels
	ActiveCount int64
}

type ModulesInfo struct {
	// module show
	ModuleCount int64
}

type SipChannelsInfo struct {
	// sip show channels
	// sip show subscriptions
	// sip show channelstats
	ActiveSipDialogs       int64
	ActiveSipSubscriptions int64
	ActiveSipChannels      int64
}

type UsersInfo struct {
	// sip show users
	Users int64
}

//////////////////////////////////////////////////////////////////////////
///////////////////////// DEFAULTS
//////////////////////////////////////////////////////////////////////////

var (
	DefaultUptimeInfo = UptimeInfo{
		SystemUptimeSeconds: -1,
		LastReloadSeconds:   -1,
	}

	DefaultChannelsInfo = ChannelsInfo{
		ActiveChannels: -1,
		ActiveCalls:    -1,
		ProcessedCalls: -1,
	}

	DefaultPeersInfo = PeersInfo{
		SipPeers:             -1,
		MonitoredOnline:      -1,
		MonitoredOffline:     -1,
		UnmonitoredOnline:    -1,
		UnmonitoredOffline:   -1,
		PeersStatusUnknown:   -1,
		PeersStatusQualified: -1,
	}

	DefaultThreadsInfo = ThreadsInfo{
		ThreadCount: -1,
	}

	DefaultAgentsInfo = AgentsInfo{
		DefinedAgents: -1,
		LoggedAgents:  -1,
		TalkingAgents: -1,
	}

	DefaultOnlineAgentsInfo = OnlineAgentsInfo{
		OnlineDefinedAgents: -1,
		OnlineLoggedAgents:  -1,
		OnlineTalkingAgents: -1,
	}

	DefaultBridgesInfo = BridgesInfo{
		Count: -1,
	}

	DefaultBridgeTechnologiesInfo = BridgeTechnologiesInfo{
		BridgeTechnologies: []BridgeTechnology{},
	}

	DefaultCalendarsInfo = CalendarsInfo{
		Count: -1,
	}

	DefaultConfBridgeInfo = ConfBridgeInfo{
		Menus:    []string{},
		Profiles: []string{},
		Users:    []string{},
	}

	DefaultChannelTypesInfo = ChannelTypesInfo{
		ChannelTypes: []ChannelType{},
	}

	DefaultImagesInfo = ImagesInfo{
		Registered: -1,
	}

	DefaultSystemInfo = SystemInfo{
		TotalMemory:  -1,
		FreeMemory:   -1,
		BufferMemory: -1,
		TotalSwap:    -1,
		FreeSwap:     -1,
		ProcessCount: -1,
	}

	DefaultTaskProcessorsInfo = TaskProcessorsInfo{
		ProcessorCounter:    -1,
		ProcessedTasksTotal: -1,
		InQueue:             -1,
	}

	DefaultVersionInfo = VersionInfo{
		Version: "",
	}

	DefaultIaxChannelsInfo = IaxChannelsInfo{
		ActiveCount: -1,
	}

	DefaultModulesInfo = ModulesInfo{
		ModuleCount: -1,
	}

	DefaultActiveSipDialogs       = int64(-1)
	DefaultActiveSipSubscriptions = int64(-1)
	DefaultActiveSipChannels      = int64(-1)

	DefaultUsersInfo = UsersInfo{
		Users: -1,
	}

	// Regexps

	AllNumbersRegexp              = regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*`)
	AllIntegersRegexp             = regexp.MustCompile(`\d+`)
	StringWithoutWhitespaceRegexp = regexp.MustCompile(`[^\s]+`)
	YesNoRegexp                   = regexp.MustCompile(`no|yes`)
)

//////////////////////////////////////////////////////////////////////////
///////////////////////// FACTORIES
//////////////////////////////////////////////////////////////////////////

// NewCmdRunner build cmdRunner instance
func NewCmdRunner(asteriskPath string, logger log.Logger) *CmdRunner {
	return &CmdRunner{
		Logger: logger,
		Cmd:    asteriskPath,
	}
}

//////////////////////////////////////////////////////////////////////////
///////////////////////// HELPERS
//////////////////////////////////////////////////////////////////////////

func (c *CmdRunner) run(asteriskCommand string) (string, error) {
	cmd := exec.Command(c.Cmd, "-rx", asteriskCommand)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	level.Debug(c.Logger).Log("msg", "Running command", "cmd", cmd.String())
	outBytes, err := cmd.Output()

	if err != nil {
		level.Error(c.Logger).Log("err", err, "cmd", cmd.String(), "stderr", stderr.String())
		return "", err
	}

	return util.SanitizeString(string(outBytes)), nil

}

//////////////////////////////////////////////////////////////////////////
///////////////////////// COMMANDS
//////////////////////////////////////////////////////////////////////////

// UptimeInfos get uptime and reload infos
func (c *CmdRunner) UptimeInfos() *UptimeInfo {
	out, err := c.run("core show uptime seconds")
	return c.newUptimeInfo(out, err)
}

// ChannelsInfo get channels and calls info
func (c *CmdRunner) ChannelsInfo() *ChannelsInfo {
	out, err := c.run("core show channels count")
	return c.newChannelsInfo(out, err)
}

// PeersInfo get peers infos
func (c *CmdRunner) PeersInfo() *PeersInfo {
	out, err := c.run("sip show peers")
	return c.newPeersInfo(out, err)
}

// ThreadsInfo get threads infos
func (c *CmdRunner) ThreadsInfo() *ThreadsInfo {
	// asterisk -rx 'core show threads' | tail -1 | cut -d' ' -f1"
	// Or
	// asterisk -rx 'core show threads'
	// => split \n, length -1

	out, err := c.run("core show threads")
	return c.newThreadsInfo(out, err)
}

func (c *CmdRunner) AgentsInfo() *AgentsInfo {
	out, err := c.run("agent show all")
	return c.newAgentsInfo(out, err)
}

func (c *CmdRunner) OnlineAgentsInfo() *OnlineAgentsInfo {
	out, err := c.run("agent show online")
	return c.newOnlineAgentsInfo(out, err)
}

func (c *CmdRunner) BridgesInfo() *BridgesInfo {
	out, err := c.run("bridge show all")
	return c.newBridgesInfo(out, err)
}

func (c *CmdRunner) BridgeTechnologiesInfo() *BridgeTechnologiesInfo {
	out, err := c.run("bridge technology show")
	return c.newBridgeTechnologiesInfo(out, err)
}

func (c *CmdRunner) CalendarsInfo() *CalendarsInfo {
	out, err := c.run("calendar show calendars")
	return c.newCalendarsInfo(out, err)
}

func (c *CmdRunner) ConfBridgeInfo() *ConfBridgeInfo {
	return &ConfBridgeInfo{
		Menus:    c.newConfBridgeMenus(c.run("confbridge show menus")),
		Profiles: c.newConfBridgeProfiles(c.run("confbridge show profile bridges")),
		Users:    c.newConfBridgeUsers(c.run("confbridge show profile users")),
	}
}

func (c *CmdRunner) ChannelTypesInfo() *ChannelTypesInfo {
	out, err := c.run("core show channeltypes")
	return c.newChannelTypesInfo(out, err)
}

func (c *CmdRunner) ImagesInfo() *ImagesInfo {
	out, err := c.run("core show image formats")
	return c.newImagesInfo(out, err)
}

func (c *CmdRunner) SystemInfo() *SystemInfo {
	out, err := c.run("core show sysinfo")
	return c.newSystemInfo(out, err)
}

func (c *CmdRunner) TaskProcessorsInfo() *TaskProcessorsInfo {
	out, err := c.run("core show taskprocessors")
	return c.newTaskProcessorsInfo(out, err)
}

func (c *CmdRunner) VersionInfo() *VersionInfo {
	out, err := c.run("core show version")
	return c.newVersionInfo(out, err)
}

func (c *CmdRunner) IaxChannelsInfo() *IaxChannelsInfo {
	out, err := c.run("iax2 show channels")
	return c.newIaxChannelsInfo(out, err)
}

func (c *CmdRunner) ModulesInfo() *ModulesInfo {
	out, err := c.run("module show")
	return c.newModulesInfo(out, err)
}

func (c *CmdRunner) SipChannelsInfo() *SipChannelsInfo {
	return &SipChannelsInfo{
		ActiveSipDialogs:       c.newActiveSipDialogs(c.run("sip show channels")),
		ActiveSipSubscriptions: c.newActiveSipSubscriptions(c.run("sip show subscriptions")),
		ActiveSipChannels:      c.newActiveSipChannels(c.run("sip show channelstats")),
	}
}

func (c *CmdRunner) UsersInfo() *UsersInfo {
	out, err := c.run("sip show users")
	return c.newUsersInfo(out, err)
}
