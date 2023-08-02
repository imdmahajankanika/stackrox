#!/usr/bin/env -S python3 -u

"""
Run qa-tests-backend in a IBM CLOUD POWERVS Openshift cluster provided via automation-flavors/powervs.
"""
import os
from base_qa_e2e_test import make_qa_e2e_test_runner_midstream
from clusters import AutomationFlavorsCluster

# set required test parameters
os.environ["DEPLOY_STACKROX_VIA_OPERATOR"] = "true"
os.environ["ORCHESTRATOR_FLAVOR"] = "openshift"
os.environ["ROX_POSTGRES_DATASTORE"] = "true"
os.environ["USE_MIDSTREAM_IMAGES"] = "true"
os.environ["REMOTE_CLUSTER_ARCH"] = "ppc64le"

make_qa_e2e_test_runner_midstream(cluster=AutomationFlavorsCluster()).run()
