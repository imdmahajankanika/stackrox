#!/usr/bin/env -S python3 -u

"""
Run version compatibility tests
"""
import os
from lookup_previous_helm_releases import get_previously_released_helm_chart_versions
from compatibility_test import make_compatibility_test_runner
from clusters import GKECluster

class SensorVersionsFailure(Exception):
    pass

# set required test parameters
os.environ["ORCHESTRATOR_FLAVOR"] = "k8s"
os.environ["ROX_POSTGRES_DATASTORE"] = "true"

chart_versions = get_previously_released_helm_chart_versions("stackrox-secured-cluster-services", 4)

gkecluster=GKECluster("compat-test")

failing_sensor_versions = []
for version in chart_versions:
    os.environ["SENSOR_CHART_VERSION"] = version
    try:
        make_compatibility_test_runner(cluster=gkecluster).run()
    except Exception:
        print(f"Exception \"{Exception}\" raised in compatibility test for sensor version {version}")
        failing_sensor_versions += version

if len(failing_sensor_versions) > 0:
    raise SensorVersionsFailure(f"Compatibility tests failed for Sensor versions " + ', '.join(failing_sensor_versions))

