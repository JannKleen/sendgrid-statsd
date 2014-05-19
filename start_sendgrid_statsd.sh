#!/bin/bash

./sendgrid-statsd --statsd_host=`basename "$STATSD_PORT"`
