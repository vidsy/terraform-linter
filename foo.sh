#!/bin/bash

STACKS=$(find /Users/stevenjack/Projects/Vidsy/infrastructure -name "*.tf" -not -path "*.terraform/*" -type f -exec dirname {} \; | sort | uniq)

for STACK in $STACKS; do
	echo "Linting $STACK"
	./terraform-linter -tf-directory="$STACK"
done
