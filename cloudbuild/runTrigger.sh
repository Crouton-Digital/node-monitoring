#!/bin/bash


function runCloudBuildTrigger() {
  # Deploy-Infra-Test-Env
  local triggerName=${1}
  # '{"branchName":"master", "substitutions": { "_DEPLOYMENT_PATH": "infra-test/deployment/botfun" } }'
  local jsonSubstitutionsVar=${2}

  curl -d "$jsonSubstitutionsVar" -X POST -H "Content-type: application/json" -H "Authorization: Bearer $(gcloud config config-helper --format='value(credential.access_token)')" https://cloudbuild.googleapis.com/v1/projects/asset-management-ci-cd/triggers/$triggerName:run

}

trigger_name=$1
deployment_patch=$2
branchName=$3

echo $trigger_name
echo $deployment_patch
echo $branchName

runCloudBuildTrigger $trigger_name '{"branchName":"'$branchName'", "substitutions": { "_DEPLOYMENT_PATH": "'$deployment_patch'" } }'
