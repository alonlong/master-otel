#!/bin/bash

SERVER_IP=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1)

# start services
function start_services() {
    nohup bin/stored --grpc-addr $SERVER_IP:9082 &
    sleep 1
    nohup bin/ctld --stored-addr $SERVER_IP:9082 --grpc-addr $SERVER_IP:9083 &
    sleep 1
    nohup bin/apid --ctld-addr $SERVER_IP:9083 --http-addr $SERVER_IP:9084 &
    sleep 1

    ps -ef | grep $SERVER_IP | grep bin
    echo "Services are started"
}

# stop services
function stop_services() {
    ps -ef | grep $SERVER_IP | grep bin | awk '{print $2}' | xargs kill

    ps -ef | grep $SERVER_IP | grep bin
    echo "Services are stopped"
}

# status services
function status_services() {
    ps -ef | grep $SERVER_IP | grep bin
}

# test add user
function test_add_user() {
    curl -XPOST -H "Content-Type: application/json" -d '{"email":"test@live.cn", "username":"test"}' http://$SERVER_IP:9084/user
}

# test delete user
function test_delete_user() {
    curl -XDELETE http://$SERVER_IP:9084/user/$1
}

# if the `logs` directory does not exist, create it
if [ ! -d "logs" ]; then
    mkdir logs
fi

# check command
if [ "$1" == "start" ]; then
    start_services
elif [ "$1" == "stop" ]; then
    stop_services
elif [ "$1" == "restart" ]; then
    stop_services
    start_services
elif [ "$1" == "status" ]; then
    status_services
elif [ "$1" == "test" ]; then
    if [ "$2" == "create" ]; then
        test_add_user
    elif [ "$2" == "delete" ]; then
        test_delete_user $3
    else
        echo "Usage: $0 test {create|delete}"
        exit 1
    fi
else
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
fi

