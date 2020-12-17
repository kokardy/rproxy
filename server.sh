#!/bin/bash

./rproxy \
    -scheme $SCHEME \
    -rhost $RHOST \
    -addr $ADDR \
    -ori $ORI -dest $DEST \
    -ori $ORI2 -dest $DEST2 \
