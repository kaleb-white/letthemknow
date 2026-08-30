# Description: utility for querying ipconfig for size and start of subnet.
# Usage: 
## Takes one of two arguments: 'size' or 'start'. 
## Start - returns the first IPv4 address in the subnet.
## Size - returns the size of the subnet.
## Takes an optional second argument for a specific network adapter.

# Options
## -s / --start: When set, returns the start of the subnet. 
## -# / --size: When set, returns the size of the subnet.
## -c / --current: When set, returns the current ip.
## -m / --mask-size: When set, returns the size of the bitmask.

START=
SIZE=
CURRENT=
MASK_SIZE=
NETWORK_ADAPTER=

HELP_STRING="Usage: [-s --start Return start of subnet] [-# --size Return size of subnet] [-c --current Return current ip address] [-n --network-adapter Specific network adapter to search for (otherwise uses first that specifies a subnet mask] \n" 
OPTS=$(getopt -o "s#cmn:" --long "start,size,current,mask-size,network-adapter:" -- "$@")

if [ $? -ne 0 ]; then
	printf "$HELP_STRING"
	exit 1
fi

eval set -- "$OPTS"
while true; do
	case $1 in
		-s | --start) 
			START=1
			shift
			;;
		-\# | --size) 
			SIZE=1
			shift
			;;
		-c | --current) 
			CURRENT=1
			shift
			;;
		-m | --mask-size) 
			MASK_SIZE=1
			shift
			;;
		-n | --network-adapter)
			shift
			NETWORK_ADAPTER="$1"
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

if 
	[[ -z "$START" || "$START" =~ ^- ]] &&\
	[[ -z "$SIZE" || "$SIZE" =~ ^- ]] &&\
	[[ -z "$MASK_SIZE" || "$MASK_SIZE" =~ ^- ]] &&\
	[[ -z "$CURRENT" || "$CURRENT" =~ ^- ]] ; then
	printf "Must have one of size, start, or current.\n"
	printf "$HELP_STRING"
	exit 1
fi

# Regexps for matching IPv4 addr parts
OCTET_SHARED='[0-2][0-5][0-9]|1?[0-9][0-9]|[0-9])'
OCTET="(${OCTET_SHARED}"
OCTET_UNCAPTURED="(?:${OCTET_SHARED}"

# Regexps for network adapter blocks in ipconfig result. Known adapter regexp allows prepending with adapter name.
NETWORK_ADAPTER_KNOWN='(?:(?!\n[^\n\ ])[\s\S])*?Subnet Mask(?:(?!\n[^\n\ ])[\s\S])*?\n(?=[^\n\ ])'
NETWORK_ADAPTER_UNKNOWN='\n[^\n\ ](?:(?!\n[^\n\ ])[\s\S])*?Subnet Mask(?:(?!\n[^\n\ ])[\s\S])*?\n(?=[^\n\ ])'

if [[ -z "$NETWORK_ADAPTER" ]]; then
	NETWORK_ADAPTER="$NETWORK_ADAPTER_UNKNOWN"
else
	NETWORK_ADAPTER="${NETWORK_ADAPTER}$NETWORK_ADAPTER_KNOWN"
fi

ADAPTER_BLOCK=$(ipconfig | grep -oPz ''"$NETWORK_ADAPTER"'' | tr -d '\0')

if [[ -z "$ADAPTER_BLOCK" ]]; then
	echo "Failed to find matching adapter block with subnet mask specification. Please select an existing adapter block with subnet mask specified."
	exit 1
fi

# <ip address> <octet #>
octet_val_adapter() {
	exec 3>&2
	exec 2> /dev/null
	local from_utils=$(./utils/octet-val.sh "$1" "$2")
	if [[ -z "$from_utils" ]]; then
		printf $(./octet-val.sh "$1" "$2")
	else
		printf "$from_utils"
	fi
	exec 2>&3
}

# <ipconfig value>
get_ip_like() {
	echo $( \
		rg "$1" <<< "$ADAPTER_BLOCK" |\
		rg -o ''"${OCTET}\.${OCTET}\.${OCTET}\.${OCTET}"''\
	)
}

# <n in log_2(n)>
log_base_2() {
	local n
	if (( $# != 0 )); then
		n="$1"
	else
		while read -r line; do
			n="$line"
		done
	fi

	printf $(python -c 'import math;print(round(math.log('"$n"', 2)))')
}

# <mask> Returns the number of bits in the mask.
calculate_mask_bits() {
	local mask="$1"

	# Define parts
	local mask_size=0 shft=1 bits octet
	local total_bits=0
	for octet_idx in 4 3 2 1; do
		octet=$(octet_val_adapter $mask $octet_idx)
		bits=$(log_base_2 "(255-$octet)+1" | rg -o '.*\.?')
		if (( $bits == 0 )); then
			break
		fi
		total_bits=$(( $total_bits + $bits ))
	done
	printf $(( 32 - $total_bits))
}

# <mask> Returns number of IPv4 Addresses to perform DNS lookup on in subnet.
calculate_mask_size() {
	local mask="$1"

	# Define parts
	local mask_size=0 shft=1 octet bits
	for octet_idx in 4 3 2 1; do
		octet=$(octet_val_adapter $mask $octet_idx)
		bits=$(log_base_2 "(255-$octet)+1" | rg -o '.*\.?')
		if (( $bits == 0 )); then
			continue
		fi
		mask_size=$(( (2 ** $bits) * $shft ))
		shft=$(( $shft * 256 ))
	done
	
	printf $mask_size
}

# <divisee> <divisor>
floor() {
	printf $(python -c 'import math;print(round(math.floor('"$1"'/'"$2"')))')
}

# <divisee> <divisor>
ceil() {
	printf $(python -c 'import math;print(round(math.ceil('"$1"'/'"$2"')))')
} 

# <mask> <curr IPv4> Calculates the lowest IPv4 address in the local network.
calculate_first_ipv4() {
	local mask="$1" curr_IPv4="$2"	

	local mask_size=$(calculate_mask_size "$mask")
	local subnet_top_bit=$(( $(log_base_2 "$mask_size") - 1 ))
	local top_octet_idx=$(( 5 - $(ceil "$subnet_top_bit" "8") ))
	local top_octet_bit=$(( $subnet_top_bit % 8 + 1 ))

	local first
	for octet_idx in 1 2 3 4; do
		if (( $octet_idx > $top_octet_idx )); then
			first="$first.0"	
		elif (( $octet_idx == $top_octet_idx )); then
			local curr_IPv4_octet_val=$(octet_val_adapter $curr_IPv4 $octet_idx)
			first="$first.$(( \
				$(floor \
					"$curr_IPv4_octet_val" \
					"$(( 2 ** $top_octet_bit ))" \
				) *  2 ** $top_octet_bit \
			))"
		else
			first="$first.$(octet_val_adapter $curr_IPv4 $octet_idx)"
		fi
	done
	printf $(rg -o '[0-9].*' <<< "$first")
} 

if [[ -n "$START" ]]; then
	echo "$(calculate_first_ipv4 $(get_ip_like "Subnet Mask") $(get_ip_like "IPv4 Address"))"
elif [[ -n "$SIZE" ]]; then
	echo "$(calculate_mask_size $(get_ip_like "Subnet Mask"))"
elif [[ -n "$CURRENT" ]]; then
	echo "$(get_ip_like  "IPv4 Address")"
elif [[ -n "$MASK_SIZE" ]]; then
	echo "$(calculate_mask_bits "$(get_ip_like  "Subnet Mask")")"
else
	echo "Command unrecognized. Enter 'start' for the first IPv4 address in the subnet or 'size' for the size of the subnet."
fi
