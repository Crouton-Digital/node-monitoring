#!/bin/bash
image_version_line=false
while read -r; do
    if [[ $REPLY == *"_version"* ]]; then
        image_version_line=true
        echo "$REPLY"
    elif [ "$image_version_line" = true ]; then
		echo "${REPLY//\"*\"/\"$1\"}"
        image_version_line=false
    else
        echo "$REPLY"
    fi
done < vars.tf > vars.tf.t
mv vars.tf{.t,}

