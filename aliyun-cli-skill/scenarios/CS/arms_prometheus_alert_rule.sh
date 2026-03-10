#!/bin/bash
# CS Scenario: arms_prometheus_alert_rule
# Generated from: arms_prometheus_alert_rule.json

set -e

REGION="cn-hangzhou"

echo "==================================="
echo "CS - arms_prometheus_alert_rule"
echo "==================================="
echo ""

# Step 1: describe-clusters-v1
# Type: dependency_check
echo "[1] Running cs describe-clusters-v1..."
aliyun cs describe-clusters-v1 \
  --region-id "$REGION" \
  --page_size 10 \
  --output json > /tmp/cs001.json
CS001_RESULT=$(cat /tmp/cs001.json)
echo "✓ Step completed"

# Step 2: describe-vpcs
# Type: dependency_check
echo "[2] Running vpc describe-vpcs..."
aliyun vpc describe-vpcs \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --page-size 1 \
  --is-default true \
  --output json > /tmp/check_vpc.json
CHECK_VPC_RESULT=$(cat /tmp/check_vpc.json)
echo "✓ Step completed"

# Step 3: create-default-vpc
# Type: dependency
echo "[3] Running vpc create-default-vpc..."
aliyun vpc create-default-vpc \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --output json > /tmp/create_vpc.json
CREATE_VPC_RESULT=$(cat /tmp/create_vpc.json)
echo "  Waiting for Status to be Available..."
# TODO: Add waiter for Status == Available
echo "✓ Step completed"

# Step 4: create-v-switch
# Type: dependency
echo "[4] Running vpc create-v-switch..."
aliyun vpc create-v-switch \
  --region-id "$REGION" \
  --v-switch-name $(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 10) \
  --zone-id "cn-hangzhou-h" \
  --cidr-block $((RANDOM % 10 + 100)) \
  --vpc-id "$VPC_ASSIGN_VPCID" \
  --output json > /tmp/create_vsw.json
CREATE_VSW_RESULT=$(cat /tmp/create_vsw.json)
echo "  Waiting for Status to be Available..."
# TODO: Add waiter for Status == Available
echo "✓ Step completed"

# Step 5: create-cluster
# Type: dependency
echo "[5] Running cs create-cluster..."
aliyun cs create-cluster \
  --region-id "$REGION" \
  --body "{'name': "#EXPRESSION{'tf-test-' + GetRandomString(10)}", 'region_id': 'cn-hangzhou', 'cluster_type': 'ManagedKubernetes', 'cluster_spec': 'ack.pro.small', 'vpcid': '$.vpc_assign.VpcId', 'vswitch_ids': ['$.vsw_assign.VSwitchId'], 'service_cidr': '172.21.0.0/20', 'deletion_protection': False, 'profile': 'Serverless'}" \
  --output json > /tmp/create_cluster.json
CREATE_CLUSTER_RESULT=$(cat /tmp/create_cluster.json)
echo "  Waiting for state to be running..."
# TODO: Add waiter for state == running
echo "✓ Step completed"

# Step 6: create-prometheus-alert-rule
echo "[6] Running arms create-prometheus-alert-rule..."
aliyun arms create-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "1" \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 10%" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --notify-type "ALERT_MANAGER" \
  --type $((RANDOM % 89999 + 10000)) \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --output json > /tmp/y0yxd.json
Y0YXD_RESULT=$(cat /tmp/y0yxd.json)
echo "✓ Step completed"

# Step 7: describe-prometheus-alert-rule
echo "[7] Running arms describe-prometheus-alert-rule..."
aliyun arms describe-prometheus-alert-rule \
  --region-id "$REGION" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/rd001.json
RD001_RESULT=$(cat /tmp/rd001.json)
echo "✓ Step completed"

# Step 8: update-prometheus-alert-rule
echo "[8] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --duration "1" \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 10%" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/n02tu.json
N02TU_RESULT=$(cat /tmp/n02tu.json)
echo "✓ Step completed"

# Step 9: update-prometheus-alert-rule
echo "[9] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "2" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 10%" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/mzurq.json
MZURQ_RESULT=$(cat /tmp/mzurq.json)
echo "✓ Step completed"

# Step 10: update-prometheus-alert-rule
echo "[10] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "2" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 10%" \
  --labels "[{"name":"TF","value":"test1"}]" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/9i45j.json
9I45J_RESULT=$(cat /tmp/9i45j.json)
echo "✓ Step completed"

# Step 11: update-prometheus-alert-rule
echo "[11] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "2" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 10%" \
  --labels "[{"name":"TF","value":"test1"}]" \
  --annotations "[{"name":"TF","value":"test1"}]" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/of9ia.json
OF9IA_RESULT=$(cat /tmp/of9ia.json)
echo "✓ Step completed"

# Step 12: update-prometheus-alert-rule
echo "[12] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "2" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10" \
  --message "node available memory is less than 20%" \
  --labels "[{"name":"TF","value":"test1"}]" \
  --annotations "[{"name":"TF","value":"test1"}]" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/nyisx.json
NYISX_RESULT=$(cat /tmp/nyisx.json)
echo "✓ Step completed"

# Step 13: update-prometheus-alert-rule
echo "[13] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "2" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20" \
  --message "node available memory is less than 20%" \
  --labels "[{"name":"TF","value":"test1"}]" \
  --annotations "[{"name":"TF","value":"test1"}]" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/hrmrr.json
HRMRR_RESULT=$(cat /tmp/hrmrr.json)
echo "✓ Step completed"

# Step 14: update-prometheus-alert-rule
echo "[14] Running arms update-prometheus-alert-rule..."
aliyun arms update-prometheus-alert-rule \
  --region-id "$REGION" \
  --region-id "cn-hangzhou" \
  --duration "1" \
  --alert-name $((RANDOM % 89999 + 10000)) \
  --expression "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20" \
  --message "node available memory is less than 20%" \
  --type $((RANDOM % 89999 + 10000)) \
  --labels "[{"name":"TF2","value":"test2"}]" \
  --annotations "[{"name":"TF2","value":"test2"}]" \
  --cluster-id "$CLUSTER_ASSIGN_CLUSTERID" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/6sea2.json
6SEA2_RESULT=$(cat /tmp/6sea2.json)
echo "✓ Step completed"

# Step 15: delete-prometheus-alert-rule
echo "[15] Running arms delete-prometheus-alert-rule..."
aliyun arms delete-prometheus-alert-rule \
  --region-id "$REGION" \
  --alert-id "$Y0YXD_OUTPUT_PROMETHEUSALERTRULE_ALERTID" \
  --output json > /tmp/v6jr7.json
V6JR7_RESULT=$(cat /tmp/v6jr7.json)
echo "✓ Step completed"

# Step 16: delete-cluster
# Type: dependency_delete
echo "[16] Running cs delete-cluster..."
aliyun cs delete-cluster \
  --region-id "$REGION" \
  --cluster-id "$CREATE_CLUSTER_OUTPUT_CLUSTER_ID" \
  --output json > /tmp/del_cluster.json
DEL_CLUSTER_RESULT=$(cat /tmp/del_cluster.json)
echo "✓ Step completed"

# Step 17: delete-v-switch
# Type: dependency_delete
echo "[17] Running vpc delete-v-switch..."
aliyun vpc delete-v-switch \
  --region-id "$REGION" \
  --v-switch-id "$CREATE_VSW_OUTPUT_VSWITCHID" \
  --output json > /tmp/del_vsw.json
DEL_VSW_RESULT=$(cat /tmp/del_vsw.json)
echo "✓ Step completed"

# Step 18: delete-vpc
# Type: dependency_delete
echo "[18] Running vpc delete-vpc..."
aliyun vpc delete-vpc \
  --region-id "$REGION" \
  --vpc-id "$CREATE_VPC_OUTPUT_VPCID" \
  --output json > /tmp/del_vpc.json
DEL_VPC_RESULT=$(cat /tmp/del_vpc.json)
echo "✓ Step completed"


echo ""
echo "==================================="
echo "Cleanup (if needed)"
echo "==================================="
echo "Please clean up resources manually or implement cleanup logic."
