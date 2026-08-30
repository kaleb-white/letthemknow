# Regexps for matching IPv4 addr parts
OCTET_SHARED='[0-2][0-5][0-9]|1?[0-9][0-9]|[0-9])'
OCTET="(${OCTET_SHARED}"
OCTET_UNCAPTURED="(?:${OCTET_SHARED}"

# <octet 1-4 of IPv4 addr>. Returns regexp to match value.
regexp_for_octet() {
	local i=1 regexp="^"
	while (( $i <= 4)); do
		# Append caputred or uncaptured
		if [[ $i == $1 ]]; then
			regexp="${regexp}${OCTET}{1}?"
		else
			regexp="${regexp}${OCTET_UNCAPTURED}{1}?"
		fi

		# Add dot if not in fourth octet
		if (( $i != 4 )); then
			regexp="${regexp}\."
		else
			regexp="${regexp}$"
		fi
		i=$(( i + 1 ))
	done
	printf $regexp
}

# <IPv4 value> <octet (0-255) 1-4 of IPv4 addr>. Returns value.
echo $(rg -o ''"$(regexp_for_octet "$2")"'' -r '$1' <<< "$1")

