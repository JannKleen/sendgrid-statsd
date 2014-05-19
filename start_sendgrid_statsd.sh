#!/bin/bash

./sendgrid-statsd --statsd_host=`basename "$STATSD_PORT_8125_UDP"`
