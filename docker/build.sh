#!/bin/bash
docker build --no-cache -t us.gcr.io/mchirico/auth:pig .
docker push us.gcr.io/mchirico/auth:pig
