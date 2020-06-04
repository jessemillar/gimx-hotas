#!/usr/bin/env bash

ALL_MACRO_COMMENTS=$(grep "^#" < gimx/.gimx/macros/vkb-gladiator-mk-ii.txt)

echo "$ALL_MACRO_COMMENTS" > mapping-reference.txt
