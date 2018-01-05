#!/usr/bin/env bash

YEAR=$(cat version | cut -d "." -f1)
if [[ ! ${YEAR} =~ ^[0-9][0-9][0-9][0-9] ]]; then
    exit 1
fi

OLD_VER=$(cat version | cut -d "." -f2)

CURRENT_YEAR=$(date -u +%Y)

if [[ ${CURRENT_YEAR} -ne ${YEAR} ]]; then
    NEW_VER=1
else
    NEW_VER=$((${OLD_VER} + 1))
fi

echo  $YEAR.$NEW_VER > version

cat version