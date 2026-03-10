#!/bin/bash
# CS Scenario: Addon
# Generated from: Addon.json

set -e

REGION="cn-hangzhou"

echo "==================================="
echo "CS - Addon"
echo "==================================="
echo ""

# Step 1: describe-addons
echo "[1] Running cs describe-addons..."
aliyun cs describe-addons \
  --region "cn-hangzhou" \
  --cluster_type "ManagedKubernetes" \
  --output json > /tmp/list_addons.json
LIST_ADDONS_RESULT=$(cat /tmp/list_addons.json)
echo "✓ Step completed"

# Step 2: describe-clusters-v1
echo "[2] Running cs describe-clusters-v1..."
aliyun cs describe-clusters-v1 \
  --page_size 10 \
  --page_number 1 \
  --output json > /tmp/list_clusters.json
LIST_CLUSTERS_RESULT=$(cat /tmp/list_clusters.json)
echo "✓ Step completed"

# Step 3: list-addons
echo "[3] Running cs list-addons..."
aliyun cs list-addons \
  --region_id "cn-hangzhou" \
  --cluster_type "ManagedKubernetes" \
  --profile "Default" \
  --cluster_spec "ack.pro.small" \
  --cluster_version "1.34.3-aliyun.1" \
  --output json > /tmp/list_all_addons.json
LIST_ALL_ADDONS_RESULT=$(cat /tmp/list_all_addons.json)
echo "✓ Step completed"


echo ""
echo "==================================="
echo "Cleanup (if needed)"
echo "==================================="
echo "Please clean up resources manually or implement cleanup logic."
