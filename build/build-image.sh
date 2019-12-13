#!/bin/bash

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
NAME="${1}"
NS="knative-sample"
DEFAULT_REG="registry.cn-hangzhou.aliyuncs.com"

GIT_COMMIT="$(git rev-parse --verify HEAD)"
GIT_BRANCH=`git branch | grep \* | cut -d ' ' -f2`
if [[ -z "${GIT_BRANCH}" || -z "${GIT_COMMIT}" ]]; then
  TAG="latest-$(date +%Y''%m''%d''%H''%M''%S)"
else
  TAG="${GIT_BRANCH}_${GIT_COMMIT:0:8}-$(date +%Y''%m''%d''%H''%M''%S)"
fi

docker build -t "${DEFAULT_REG}/${NS}/${NAME}:${TAG}" -f ${ROOTDIR}/Dockerfile-${NAME} ${ROOTDIR}/../
docker push "${DEFAULT_REG}/${NS}/${NAME}:${TAG}"

