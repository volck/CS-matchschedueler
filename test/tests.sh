#!/bin/bash
curl -s localhost:1337/makeMatch -d @justmatch.json | jq
curl -s localhost:1337/makeMatch -d @justmatch_now.json | jq
curl -s localhost:1337/makeMatch -d @justmatch_1h.json | jq

curl -s localhost:1337/scheduleMatches | jq
