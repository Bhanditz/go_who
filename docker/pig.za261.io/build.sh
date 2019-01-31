#!/bin/bash
docker build --no-cache -t us.gcr.io/mchirico/auth:za261 .
docker push us.gcr.io/mchirico/auth:za261
