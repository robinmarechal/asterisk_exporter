# Asterisk Exporter

Prometheus exporter for Asterisk call center phone systems

Tested on Asterisk 1.13.

***Disclaimer**: To be honest I don't even know how Asterisk works. I just wrote an exporter based on my small observations with the one Asterisk system (v1.13) that was in place during this specific mission. Since i was not able to install python3 on the server, I decided to go for my own exporter, inspired by [tainguyenbp's Python one](https://github.com/tainguyenbp/asterisk_exporter). This version is written in Go and provides more features*.

## How does it work

In a nutshell, this exporter simply runs a bunch of commands to the asterisk agent (usually `/usr/bin/asterisk` or `/usr/sbin/asterisk`). It then parses the output and format it for Prometheus.

*Basically, this exporter is nothing else than a translator that transforms Asterisk outputs into Prometheus readable metrics.*

Some executed command examples : 
```bash
asterisk -rx 'core show uptime seconds'
asterisk -rx 'core show channels count'
asterisk -rx 'sip show peers'
asterisk -rx 'agent show all'
asterisk -rx 'iax2 show channels'
...
```

## Installation and Usage

The `node_exporter` listens on HTTP port 9795 by default. See the `--help` output for more options.

## Collectors

Metrics are splitted into multiple collectors, allowing users to enable/disabled them evetually depending on the plugins/modules they have installed (once again, *i don't know how Asterisk works*). This allows better portability between Asterisk systems, adaptability and eases possible future evolutions and contributions.

Collectors are enabled by providing a `--collector.<name>` flag.

Collectors are disabled by providing a `--collector.<name>=false` flag.

### Enabled by default

Name     | Description 
---------|-------------
agents | Gather metrics from `agents show ...` commands.
core | Gather metrics from `core show ...` commands.
sip | Gather metrics from `sip show ...` commands.


### Disabled by default

Name     | Description 
---------|-------------
bridge | Gather metrics from `bridge show ...` commands.
calendar | Gather metrics from `calendar show ...` commands.
confbridge | Gather metrics from `confbridge show ...` commands.
iax2 | Gather metrics from `iax2 show ...` commands.
module | Gather metrics from `module show ...` commands.


## Metrics

List of exposted metrics when all collectors are enabled :

```
# HELP asterisk_agents_defined Number of defined agents
# TYPE asterisk_agents_defined gauge
asterisk_agents_defined
# HELP asterisk_agents_logged Number of logged agents
# TYPE asterisk_agents_logged gauge
asterisk_agents_logged
# HELP asterisk_agents_online_logged Number of logged online agents
# TYPE asterisk_agents_online_logged gauge
asterisk_agents_online_logged
# HELP asterisk_agents_online_talking Number of talking online agents
# TYPE asterisk_agents_online_talking gauge
asterisk_agents_online_talking
# HELP asterisk_agents_talking Number of talking agents
# TYPE asterisk_agents_talking gauge
asterisk_agents_talking
# HELP asterisk_bridges_info Bridges info
# TYPE asterisk_bridges_info gauge
asterisk_bridges_info 
# HELP asterisk_bridges_technologies_info Bridge technologies info
# TYPE asterisk_bridges_technologies_info gauge
asterisk_bridges_technologies_info
# HELP asterisk_calendars_count Number of calendars
# TYPE asterisk_calendars_count gauge
asterisk_calendars_count 
# HELP asterisk_confbridges_info ConfBridge information
# TYPE asterisk_confbridges_info gauge
asterisk_confbridges_info
# HELP asterisk_core_active_calls Number of currently active calls
# TYPE asterisk_core_active_calls gauge
asterisk_core_active_calls 
# HELP asterisk_core_active_channels Number of currently active channels
# TYPE asterisk_core_active_channels gauge
asterisk_core_active_channels 
# HELP asterisk_core_images_registered Number of registered images
# TYPE asterisk_core_images_registered gauge
asterisk_core_images_registered
# HELP asterisk_core_last_reload_seconds Number of seconds since last reload
# TYPE asterisk_core_last_reload_seconds gauge
asterisk_core_last_reload_seconds
# HELP asterisk_core_processed_calls Number of processed calls
# TYPE asterisk_core_processed_calls gauge
asterisk_core_processed_calls
# HELP asterisk_core_system_memory_buffer_bytes System buffer memory in bytes
# TYPE asterisk_core_system_memory_buffer_bytes gauge
asterisk_core_system_memory_buffer_bytes
# HELP asterisk_core_system_memory_free_bytes System free memory in bytes
# TYPE asterisk_core_system_memory_free_bytes gauge
asterisk_core_system_memory_free_bytes
# HELP asterisk_core_system_memory_total_bytes System total memory in bytes
# TYPE asterisk_core_system_memory_total_bytes gauge
asterisk_core_system_memory_total_bytes
# HELP asterisk_core_system_processes Number of system processes
# TYPE asterisk_core_system_processes gauge
asterisk_core_system_processes
# HELP asterisk_core_system_swap_free_bytes System free swap in bytes
# TYPE asterisk_core_system_swap_free_bytes gauge
asterisk_core_system_swap_free_bytes
# HELP asterisk_core_system_swap_total_bytes System total swap in bytes
# TYPE asterisk_core_system_swap_total_bytes gauge
asterisk_core_system_swap_total_bytes
# HELP asterisk_core_system_uptime_seconds Number of seconds since system startup
# TYPE asterisk_core_system_uptime_seconds gauge
asterisk_core_system_uptime_seconds
# HELP asterisk_core_tasks_processed_total Number of processed tasks
# TYPE asterisk_core_tasks_processed_total counter
asterisk_core_tasks_processed_total
# HELP asterisk_core_tasks_processes_in_queue Task processes in queue
# TYPE asterisk_core_tasks_processes_in_queue gauge
asterisk_core_tasks_processes_in_queue
# HELP asterisk_core_tasks_processors Number of task processors
# TYPE asterisk_core_tasks_processors gauge
asterisk_core_tasks_processors
# HELP asterisk_core_version Version info
# TYPE asterisk_core_version gauge
asterisk_core_version
# HELP asterisk_exporter_collector_error Collector errors. 0 = no error, 1 = error occurred
# TYPE asterisk_exporter_collector_error gauge
asterisk_exporter_collector_error
# HELP asterisk_iax2_channels_active Number of IAX Active channels
# TYPE asterisk_iax2_channels_active gauge
asterisk_iax2_channels_active
# HELP asterisk_modules_count Number of installed modules
# TYPE asterisk_modules_count gauge
asterisk_modules_count
# HELP asterisk_sip_active_channels Number of active SIP channels
# TYPE asterisk_sip_active_channels gauge
asterisk_sip_active_channels
# HELP asterisk_sip_active_dialogs Number of active SIP dialogs
# TYPE asterisk_sip_active_dialogs gauge
asterisk_sip_active_dialogs
# HELP asterisk_sip_active_subscriptions Number of active SIP subscriptions
# TYPE asterisk_sip_active_subscriptions gauge
asterisk_sip_active_subscriptions
# HELP asterisk_sip_current_monitored_offline Number of currently monitored offline SIP
# TYPE asterisk_sip_current_monitored_offline gauge
asterisk_sip_current_monitored_offline
# HELP asterisk_sip_current_monitored_online Number of currently monitored online SIP
# TYPE asterisk_sip_current_monitored_online gauge
asterisk_sip_current_monitored_online
# HELP asterisk_sip_current_qualified Current number of qualified SIP
# TYPE asterisk_sip_current_qualified gauge
asterisk_sip_current_qualified
# HELP asterisk_sip_current_unknown Current number of unknown SIP
# TYPE asterisk_sip_current_unknown gauge
asterisk_sip_current_unknown
# HELP asterisk_sip_current_unmonitored_offline Number of currently unmonitored offline SIP
# TYPE asterisk_sip_current_unmonitored_offline gauge
asterisk_sip_current_unmonitored_offline
# HELP asterisk_sip_current_unmonitored_online Number of currently unmonitored online SIP
# TYPE asterisk_sip_current_unmonitored_online gauge
asterisk_sip_current_unmonitored_online
# HELP asterisk_sip_users Number of users
# TYPE asterisk_sip_users gauge
asterisk_sip_users
```

## Command line

```
usage: asterisk_exporter [<flags>]

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9795"
                              The address to listen on for HTTP requests.
      --asterisk.path="/usr/sbin/asterisk"
                              Path to Asterisk binary
      --metrics.prefix="asterisk"
                              Prefix of exposed metrics
      --web.telemetry-path="/metrics"
                              Path under which to expose metrics.
      --web.enable-exporter-metrics
                              Include metrics about the exporter itself (process_*, go_*).
      --web.enable-promhttp-metrics
                              Include metrics about the http server itself (promhttp_*)
      --web.max-requests=40   Maximum number of parallel scrape requests. Use 0 to disable.
      --collector.agents      Enable agents collector
      --collector.core        Enable core collector
      --collector.sip         Enable sip collector
      --collector.bridge      Enable bridge collector
      --collector.calendar    Enable calendar collector
      --collector.confbridge  Enable confbridge collector
      --collector.iax2        Enable iax2 collector
      --collector.module      Enable module collector
      --log.level=info        Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt     Output format of log messages. One of: [logfmt, json]
      --version               Show application version.
```

## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)

Building:

```bash
git clone https://github.com/robinmarechal/asterisk_exporter.git
cd asterisk_exporter
make
./asterisk_exporter <flags>
```

To see all available configuration flags:

```bash
./asterisk_exporter -h
```

## Evolutions

I'm most likely not going to do further work on this unless it's required at my work.

If you really need features or bufgixes, you can still create issues, and even try it out youself. Golang is not a hard language to learn (or, at least, this kind of exporter does not require a strong Go knowledge) and this exporter is a pretty basic one.