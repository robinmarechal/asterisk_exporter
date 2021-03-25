#!/bin/bash

lines=(
    "agent show all"
    "agent show online"
    "ari show status"
    "bridge show all"
    "bridge technology show"
    "calendar show calendars"
    "calendar show types"
    "confbridge show menus"
    "confbridge show profile bridges"
    "confbridge show profile users"
    "core show applications"
    "core show calls + seconds"
    "core show channels count"
    "core show channeltypes"
    "core show hanguphandlers all"
    "core show image formats"
    "core show sysinfo"
    "core show taskprocessors"
    "core show threads"
    "core show uptime"
    "core show version"
    "iax2 show channels"
    "module show"
    "sip show channels"
    "sip show subscriptions"
    "sip show channelstats"
    "sip show peers"
    "sip show users"
)

IFS=$'\n' 
for CMD in ${lines[@]}
do 
    echo "CMD: asterisk -rx '$CMD'"
    asterisk2 -rx "$CMD"
    echo "########################################################"
done
