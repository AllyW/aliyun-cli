#!/bin/bash
# CS Scenario: EnvFeature
# Generated from: EnvFeature.json

set -e

REGION="cn-hangzhou"

echo "==================================="
echo "CS - EnvFeature"
echo "==================================="
echo ""

# Step 1: describe-clusters-v1
echo "[1] Running cs describe-clusters-v1..."
aliyun cs describe-clusters-v1 --page_size 1 --output json > /tmp/sly89.json
SLY89_RESULT=$(cat /tmp/sly89.json)
echo "Step completed"

# Variable assignment: clsid
CLSID_CLUSTER_ID="$.sly89.output.clusters[0].cluster_id"

# Step 2: create-environment
echo "[2] Running arms create-environment..."
aliyun arms create-environment \
  --environment-name ""$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 10)-arms-hz"" \
  --region-id "cn-hangzhou" \
  --environment-type "CS" \
  --environment-sub-type "ACK" \
  --bind-resource-id "$CLSID_CLUSTER_ID" \
  --output json > /tmp/qjtuy.json
QJTUY_RESULT=$(cat /tmp/qjtuy.json)
echo "Step completed"

# Step 3: list-environment-features
echo "[3] Running arms list-environment-features..."
aliyun arms list-environment-features \
  --region-id "cn-hangzhou" \
  --environment-id "$QJTUY_OUTPUT_DATA" \
  --output json > /tmp/y13ll.json
Y13LL_RESULT=$(cat /tmp/y13ll.json)
echo "Step completed"

# Step 4: delete-environment
echo "[4] Running arms delete-environment..."
aliyun arms delete-environment \
  --region-id "cn-hangzhou" \
  --environment-id "$QJTUY_OUTPUT_DATA" \
  --output json > /tmp/q8yc4.json
Q8YC4_RESULT=$(cat /tmp/q8yc4.json)
echo "Step completed"


echo ""
echo "==================================="
echo "Cleanup (if needed)"
echo "==================================="
echo "Please clean up resources manually or implement cleanup logic."
