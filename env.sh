#!/bin/bash
# APP settings
export API_ADDR=localhost:2727
export DB_URI="postgres://dev:CgqFv9wR1UePt6Aq@82.202.247.144:5432/license?sslmode=disable"

# Trace settings
export UPTRACE_DSN="http://pDTThrJirasdfqrxDzAer562XkTX24DQ@up.benefy.ru:14317/2"
export OTEL_RESOURCE_ATTRIBUTES="service.name=lts-service,service.version=0.1.x"

# Logger settings
# -1 = Debug; 0 = Info; 1 = Warn; 2 = Error ...
export LOG_LEVEL=-1
export BUN_VERBOSE=false

# Token settings
export SESSION_DURATION="600h"
export TOKEN_DURATION="6m"
#export TOKENEXP="60h"
#export TOKENKEY="LCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY2OTQ3NDcxMywiaWF0IjoxNjY5NDc0NzEzfQ.kD1xZbOm33yD4f0Q0Prlk"
export SIGNED_TOKEN="LCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY2OTQ3NDcxMywiaWF0IjoxNjY5NDc0NzEzfQ.kD1xZbOm33yD4f0Q0Prlk"