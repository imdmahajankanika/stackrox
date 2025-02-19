#!/usr/bin/env -S python3 -u

"""
Run QA e2e tests against a given cluster.
"""
from pre_tests import PreSystemTests
from ci_tests import QaE2eTestPart1, QaE2eTestPart2, QaE2eDBBackupRestoreTest
from post_tests import PostClusterTest, CheckStackroxLogs, FinalPost
from runners import ClusterTestSetsRunner


def make_qa_e2e_test_runner(cluster):
    return ClusterTestSetsRunner(
        cluster=cluster,
        sets=[
            {
                "name": "QA tests part I",
                "pre_test": PreSystemTests(),
                "test": QaE2eTestPart1(),
                "post_test": PostClusterTest(
                    check_stackrox_logs=True,
                    artifact_destination_prefix="part-1",
                ),
            },
            {
                "name": "QA tests part II",
                "test": QaE2eTestPart2(),
                "post_test": PostClusterTest(
                    check_stackrox_logs=True,
                    artifact_destination_prefix="part-2",
                ),
                "always_run": False,
            },
            {
                "name": "DB backup and restore",
                "test": QaE2eDBBackupRestoreTest(),
                "post_test": CheckStackroxLogs(
                    check_for_errors_in_stackrox_logs=True,
                    artifact_destination_prefix="db-test",
                ),
                "always_run": False,
            },
        ],
        final_post=FinalPost(
            store_qa_tests_data=True,
        ),
    )
