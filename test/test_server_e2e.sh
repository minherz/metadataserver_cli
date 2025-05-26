#!/bin/bash

echo "End-to-end testing of metadata server CLI"

# find open port
busy_ports=$(ss -Htan | awk '{print $4}' | cut -d':' -f2 | grep -v "^$")
read lower_port upper_port < /proc/sys/net/ipv4/ip_local_port_range
for (( port = lower_port ; port <= upper_port ; port++ )); do
    if ! echo "${busy_ports}" | grep -q "${port}"; then
        break;
    fi
done

# NOTE: in super rare case when ALL ports are busy it fails
echo "Launching metadata server on port ${port}"

./metadataserver_cli -a '127.0.0.1' -p "${port}" &
sleep 1
pid=$(pgrep metadataserver | tail -1)
if [[ "|${pid}|" == "||" ]]; then
    echo "Failed to start metadata server. Exiting..."
    exit 1
fi

test=$(curl -sX GET "http://127.0.0.1:${port}/computeMetadata/v1")

echo "Test completed."
echo "Sending Ctrl-C to stop metadata server (PID:${pid})"
kill -2 "${pid}"

if [[ "${test}" != "ok" ]]; then
    echo "Failed to get valid response from metadata server"
    exit 2
fi

