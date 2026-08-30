#!/bin/bash

# Description: Use to pair to a device
# Usage: 

log() { if [ $VERBOSE ]; then echo "[INFO] $1" >&1; fi; }
warn() { if [ $VERBOSE ]; then echo "[WARN] $1" >&1; fi; }
error() { if [ $VERBOSE ]; then echo "[ERROR] $1" >&1; exit 1; fi; }

# Options
HELP_STRING="""
	-s / --save-file: If set, the paired device's MAC address will be saved to this file.
	-c / --code: Device code given by pairing device 
	-v / --verbose: Print log messages to stdout
	-i / --save-ips: Save the list of ips to .up-addresses
	-l / --load-ips: Load from list of ips in .up-addresses"""

SAVE_IPS=
LOAD_IPS=
SAVE_FILE=
CODE=
VERBOSE=0

OPTS=$(getopt -o "vs:c:il" --long "verbose,save-file:,code:,save-ips,load-ips" -- "$@")

if [ $? -ne 0 ]; then
	printf "Failed to parse options.\n"
	printf "$HELP_STRING"
	exit 1
fi

eval set -- "$OPTS"
while true; do
	case $1 in
		-s | --save-file) 
			shift
			SAVE_FILE="$1"
			shift
			;;
		-c | --code) 
			shift
			CODE="$1"
			shift
			;;
		-v | --verbose) 
			VERBOSE=1
			shift
			;;
		-i | --save-ips) 
			SAVE_IPS=1
			shift
			;;
		-l | --load-ips) 
			LOAD_IPS=1
			shift
			;;

		--)
			shift
			break
			;;
		*) 
			printf "$HELP_STRING"
			exit 1;;
	esac
done

if [[ -z "$CODE" || "$CODE" =~ ^- ]]; then
	printf "Required argument -c | --code was missing or began with a - when read, indicating a missing argument.\n"
	print
	f "$HELP_STRING"
	exit 1
fi

# Constants
OCTET='([0-2][0-5][0-9]|1?[0-9][0-9]|[0-9])'
UP_ADDR_FILE=".up-addresses"
ANDROID_DNS_NAMES="pixel|android"
ANDROID_OS_NAMES="pixel|linux"

# Regular expressions (figured I'd put them here so if nmap changes outputs it's easy to update them). Sed expressions below
REGEXP_IP="${OCTET}\.${OCTET}\.${OCTET}\.${OCTET}"
REGEXP_NMAP_MAC_LINE="^MAC Address"
REGEXP_NMAP_TCP_PORTS="^.*/tcp"
REGEXP_NMAP_PORT_LINE="^[1-9]"
REGEXP_NMAP_HOSTNAME_LINE="^Nmap scan"
REGEXP_EMPTY_STRING="^$"
REGEXP_ONLY_SPACE="^\s+"
REGEXP_NMAP_OS_LINE="^Running:"

SED_EXTRACT_MAC='s/MAC Address: (.*) \(Unknown\)/\1/'
SED_EXTRACT_PORT='s:/tcp$::'
SED_EXTRACT_HOSTNAME='s/Nmap scan report for ([^\(0-9]*).*/\1/g'
SED_EXTRACT_OS="s/Running: //"

# Helpers
RANGES="""40000-45000
35000-39999
45001-50000
30000-34999
50001-55000
25000-29999
55001-60000
20000-24999
60001-65535
15000-19999
10000-14999
5000-9999
1024-4999
"""
## Output list
ADB_PAIR_PORT_LIST=()
ADB_PAIR_FOUND=1
ADB_PAIR_MAC_ADDR=
## <IP address>
scan_ip() {
	ADB_PAIR_FOUND=1
	while IFS= read -r line; do
		log "Checking range $line..."
		local RESULT=$(nmap "$1" -T5 -n -Pn -p "$line" --open)
		# Sample line indicating open port: 40203/tcp	open	unknown
		# This if means there's at least one port
		if rg -q '^[1-9]' <<< "$RESULT"; then
			# Extract MAC address if we ned to save it
			ADB_PAIR_MAC_ADDR=$(rg ''"$REGEXP_NMAP_MAC_LINE"'' <<< "$RESULT" | sed -E ''"$SED_EXTRACT_MAC"'')

			# Extract actual ports
			while IFS= read -r line; do
				ADB_PAIR_PORT_LIST+=($(rg -o ''"$REGEXP_NMAP_TCP_PORTS"'' | sed ''"$SED_EXTRACT_PORT"''))
			done <<< $(rg ''"$REGEXP_NMAP_PORT_LINE"'')

			# Attempt to connect to each at once
			PIPES=()
			for port in "${ADB_PAIR_PORT_LIST[@]}"; do
				log "Attempting to pair with $1 on port $port..."
				PIPES+=(''"$port"'_PIPES')
				coproc ''"$port"'_PIPES' { adb pair "$1":"$port"; }
#				echo "$CODE" >&"${$PIPES[1]}"
			done

			# Parse outputs for valid response
			for PORT in "${PIPES[@]}"; do
				log "a"					
			done

		fi
	done <<< "$RANGES"
	log $(nmap "$1" -T4 -n -Pn -p 1024-65535 --open)
}

if [[ -z "$LOAD_IPS" ]]; then
	# Parse ipconfig for subnet
	log "Parsing ipconfig result..."
	IP=$(./utils/subnet-info.sh -c)
	MASK_SIZE=$(./utils/subnet-info.sh -m)

	log "IP: $IP"
	log "MASK_SIZE: $MASK_SIZE"

	# nmap subnet for hosts that are up
	log "Running nmap for hosts against $IP/$MASK_SIZE..."
	## -sn = just ping, -n = don't perform reverse DNS lookup
	UP_HOSTS=$(nmap "$IP/$MASK_SIZE" -sn | rg -o ''"$REGEXP_IP"'')
	log "Found IP addresses $(tr '\n' ' ' <<< "$UP_HOSTS")."

	echo "$UP_HOSTS" > "$UP_ADDR_FILE"
fi

# nmap for hostnames 
log "Running nmap for host names..."
HOSTNAME_REPORT=$(nmap -sn -iL "$UP_ADDR_FILE") 
## Extract hostname lines
HOSTNAMES=$(rg ''"$REGEXP_NMAP_HOSTNAME_LINE"'' <<< "$HOSTNAME_REPORT") 

## Scan list: list of IPs to port scan
SCAN_LIST=()

## Parse for hostnames 
MATCH=0
TOTAL_HOSTS=0
while IFS= read -r line; do
	hostname=$(sed --r ''"$SED_EXTRACT_HOSTNAME"'' <<< "$line")
	TOTAL_HOSTS=$((TOTAL_HOSTS + 1))
	if [[ $? -ne 0 ]] || rg -q '^\s+|^[0-9]' <<< "$hostname"; then
		log "Hostname not found or error on line $line."
		log "Continuing."
		HOSTNAME_FOUND=1
		break
	fi

	if rg -q -i ''"$ANDROID_DNS_NAMES"'' <<< "$hostname"; then
		MATCH=0
		IP=$(rg -o ''"$REGEXP_IP"''  <<< "$line")
		SCAN_LIST+=("$IP")
		log "Matching hostname found (matched portion $hostname) with IP $IP."
	fi
done <<< "$HOSTNAMES"

if [[ $MATCH -eq 0 ]]; then
	for IP in "${SCAN_LIST[@]}"; do
		log "Beginning scan on $IP..."
		scan_ip "$IP"
	done
fi

log "Could not determine android IP from pinging $TOTAL_HOSTS hosts. Either no matching hostname (adjust hostnames in ANDROID_DNS_NAMES), DHCP server not configured to put records in DNS server, DNS server not configured to serve reverse DNS records, or see above error."
for i in 1 2 3; do
	log "Running OS detection against up hosts (attempt $i)/3..."
	TOTAL_MATCHES=0
	while IFS= read -r ip; do
		log "Trying host $ip..."
		RESULT=$(nmap "$ip" -O -F -T4)
		OS=$(rg ''"$REGEXP_NMAP_OS_LINE"'' <<< "$RESULT" | sed ''"$SED_EXTRACT_OS"'')
		IP=$(rg ''"$REGEXP_IP"'' <<< "$RESULT")
		if rg -q ''"$REGEXP_ONLY_SPACE"'' <<< "$OS" || rg -q ''"$REGEXP_EMPTY_STRING"'' <<< "$OS"; then
			log "No OS identified. Trying next..."
			continue
		fi
		log "Device is running any of: $OS."
		if rg -q -i ''"$ANDROID_OS_NAMES"'' <<< "$OS"; then
			MATCH=0
			log "IP matched possible Android OS names. Adding IP to scan list: $IP."
			SCAN_LIST+=($IP)
			TOTAL_MATCHES=$(($TOTAL_MATCHES + 1))
		fi
	done <<< "$UP_HOSTS"
	if [[ $MATCH -eq 0 ]]; then
		break
	fi
done

if [[ $MATCH -eq 0]]; then
	log "$TOTAL_MATCHES OS matches found. Scanning for open ports..."
	for IP in "${SCAN_LIST[@]}"; do
		log "Beginning scan on $IP..."
		scan_ip "$IP"
		if [[ $ADB_PAIR_FOUND -eq 0 ]]; then
			log "Found open port!
		fi
	done	
else
	error "Could not find a network device + port to attempt to connect to. Double check the device is connected and has wireless debugging enabled."
fi


if [[ -z "$SAVE_IPS" && -z "$LOAD_IPS" ]]; then
	log "Removing up hosts file $UP_ADDR_FILE..."
	rm "$UP_ADDR_FILE"
fi
