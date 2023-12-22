#!/bin/bash

sed -i "s|\"API_URL\": \"default\"|\"API_URL\": \"$API_URL\"|g" src/config.json
