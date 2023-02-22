#!/usr/bin/env bash
set -eou pipefail


deployment_value=NA
namespace_value=NA
clustername_value=NA
clusterid_value=NA
format_value=table

process_arg() {
    arg=$1

    key="$(echo "$arg" | cut -d "=" -f 1)"
    value="$(echo "$arg" | cut -d "=" -f 2)"
     
    if [[ "$key" == "deployment" ]]; then
        deployment_value="$value"
    elif [[ "$key" == "namespace" ]]; then
	namespace_value="$value"
    elif [[ "$key" == "clustername" ]]; then
	clustername_value="$value"
    elif [[ "$key" == "clusterid" ]]; then
	clusterid_value="$value"
    elif [[ "$key" == "format" ]]; then
	format_value="$value"
    fi
}

process_args() {
     echo "In process_arguments"
     for arg in "$@"; do
         echo "$arg"
	 process_arg "$arg"
     done
}

pad_string() {
    string="$1"
    length="$2"

    current_length=${#string}

    #echo "string= $string"
    #echo "length= $length"
    #echo "current_length= $current_length"

    if [[ $current_length -ge $length ]]; then
        echo "$string"
    else
        num_spaces=$((length - current_length))
        #echo "num_spaces= $num_spaces"

        #padded="$(printf "%s%*s\n" "$string" "$num_spaces" "")"
        #echo $padded
        for ((i = 0; i < num_spaces; i = i + 1)); do
            string="$string "
        done

        echo "$string"
    fi

}

process_args $@

port=8443
port=8000
export OPEN_BROWSER=false
#export OPEN_BROWSER=true
logmein localhost:$port &> token_file.txt
token="$(cat token_file.txt | sed 's|.*token=||' | sed 's|&type.*||')"

password="$(cat ./deploy/k8s/central-deploy/password)"
curl -sSkf -u "admin:$password" -o /dev/null -w '%{redirect_url}' "https://localhost:$port/sso/providers/basic/4df1b98c-24ed-4073-a9ad-356aec6bb62d/challenge?micro_ts=0"

if [[ "$deployment_value" == "NA" ]]; then
    json_deployments="$(curl --location --silent --request GET "https://localhost:$port/v1/deployments" -k -H "Authorization: Bearer $token")"

    if [[ "$namespace_value" != "NA" ]]; then
	json_deployments="$(echo "$json_deployments" | jq --arg namespace "$namespace_value" '{deployments: [.deployments[] | select(.namespace == $namespace)]}')"
    fi

    if [[ "$clustername_value" != "NA" ]]; then
	json_deployments="$(echo "$json_deployments" | jq --arg clustername "$clustername_value" '{deployments: [.deployments[] | select(.cluster == $clustername)]}')"
    fi
    
    if [[ "$clusterid_value" != "NA" ]]; then
	json_deployments="$(echo "$json_deployments" | jq --arg clusterid "$clusterid_value" '{deployments: [.deployments[] | select(.clusterId == $clusterid)]}')"
    fi

    ndeployment="$(echo $json_deployments | jq '.deployments | length')"
    deployments=()
    for ((i = 0; i < ndeployment; i = i + 1)); do
        deployments+=("$(echo "$json_deployments" | jq .deployments[$i].id | tr -d '"')")
    done
else
    deployments=($deployment_value)
fi


netstat_lines=""

for deployment in ${deployments[@]}; do
    listening_endpoints="$(curl --location --silent --request GET "https://localhost:$port/v1/listening_endpoints/deployment/$deployment" -k --header "Authorization: Bearer $token")" || true
    if [[ "$listening_endpoints" != "" ]]; then
        nlistening_endpoints="$(echo $listening_endpoints | jq '.listeningEndpoints | length')"
	if [[ "$nlistening_endpoints" > 0 ]]; then
	    if [[ "$format_value" == "json" ]]; then
                echo "deployment= $deployment"
                echo $listening_endpoints | jq
                echo
            fi	
	fi

        for ((j = 0; j < nlistening_endpoints; j = j + 1)); do
            l4_proto="$(echo $listening_endpoints | jq .listeningEndpoints[$j].endpoint.protocol | tr -d '"')"
            if [[ "$l4_proto" == L4_PROTOCOL_TCP ]]; then
                proto=tcp
            elif [[ "$l4_proto" == L4_PROTOCOL_UDP ]]; then
                proto=udp
            else
               proto=unkown
            fi

            name="$(echo $listening_endpoints | jq .listeningEndpoints[$j].signal.name | tr -d '"')"
            plop_port="$(echo $listening_endpoints | jq .listeningEndpoints[$j].endpoint.port | tr -d '"')"
            namespace="$(echo $listening_endpoints | jq .listeningEndpoints[$j].namespace | tr -d '"')"
            clusterId="$(echo $listening_endpoints | jq .listeningEndpoints[$j].clusterId | tr -d '"')"
            podId="$(echo $listening_endpoints | jq .listeningEndpoints[$j].podId | tr -d '"')"
            containerName="$(echo $listening_endpoints | jq .listeningEndpoints[$j].containerName | tr -d '"')"
            pid="$(echo $listening_endpoints | jq .listeningEndpoints[$j].signal.pid | tr -d '"')"

	    name="$(pad_string $name 20)"
	    pid="$(pad_string $pid 9)"
	    plop_port="$(pad_string $plop_port 7)"
	    proto="$(pad_string $proto 7)"
	    namespace="$(pad_string $namespace 15)"
	    clusterId="$(pad_string $clusterId 40)"
	    podId="$(pad_string $podId 55)"
	    containerName="$(pad_string $containerName 20)"

	    netstat_line="${name}${pid}${plop_port}${proto}${namespace}${clusterId}${podId}${containerName}\n"
            netstat_lines="${netstat_lines}${netstat_line}"
        done
    fi
done

echo
if [[ "$format_value" == "table" ]]; then
    header="Program name\tPID\tPort\tProto\tNamespace\tClusterId\t\t\t\tpodId\t\t\tcontainerName"

    name="$(pad_string "Program name" 20)"
    pid="$(pad_string "PID" 9)"
    plop_port="$(pad_string "Port" 7)"
    proto="$(pad_string "Proto" 7)"
    namespace="$(pad_string "Namespace" 15)"
    clusterId="$(pad_string "ClusterId" 40)"
    podId="$(pad_string "podId" 55)"
    containerName="$(pad_string "containerName" 20)"

    header="${name}${pid}${plop_port}${proto}${namespace}${clusterId}${podId}${containerName}"

    echo -e "$header"
    echo -e "$netstat_lines"
fi
