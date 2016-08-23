set -u -e
set -x
# main.lua
echo ${MQTT_USERNAME} set.
echo ${MQTT_PASSWORD} set.
echo ${MQTT_SERVER_HOSTNAME} set.
echo ${MQTT_PORT} set.

# apsetup.lua

function decorate_variable() {
   local variable_name=$1
   echo "@@${variable_name}@@"
}

function replace_variable_with_file() {
    local variable_name=$(decorate_variable $1)
    local file_source=$2
    local file_destination=$3

    local tmp=$(mktemp)
    cp $file_destination $tmp

    echo "reading from ${file_source} to set ${variable_name}"
    sed -i "/${variable_name}/ r ${file_source}" $tmp
    sed -i "/${variable_name}/d" $tmp
    cp $tmp build/$file_destination
}

function replace_variable_with_variable() {
    local variable_name=$(decorate_variable $1)
    local variable_content=$2
    local file_destination=$3
    
    local tmp=$(mktemp)
    cp $file_destination $tmp

    sed -i "s/$variable_name/$variable_content/" $tmp
    cp $tmp build/$file_destination
}

mkdir build || true
replace_variable_with_file "AP_PAYLOAD_PAGE" "html/apsetup.html" "apsetup.lua"
replace_variable_with_variable "MQTT_USERNAME" ${MQTT_USERNAME} "main.lua"
replace_variable_with_variable "MQTT_PASSWORD" ${MQTT_PASSWORD} "main.lua" 
replace_variable_with_variable "MQTT_SERVER_HOSTNAME" ${MQTT_SERVER_HOSTNAME} "main.lua"
replace_variable_with_variable "MQTT_PORT" ${MQTT_PORT} "main.lua"
cp init.lua build/ 
