#!/usr/bin/env python3

"""Regression tests for nightly source selection; no network access required."""

import contextlib
import importlib.util
import io
import json
from pathlib import Path
import runpy
import subprocess
import unittest
from unittest.mock import patch


SCRIPT = Path(__file__).with_name("resolve-nightly-versions.py")
SPEC = importlib.util.spec_from_file_location("resolve_nightly_versions", SCRIPT)
resolver = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(resolver)


def kubernetes_refs(minors=(97, 98, 99, 100, 101)):
    refs = {f"refs/heads/release-1.{minor}": f"{minor:040x}" for minor in minors}
    refs["refs/heads/master"] = "f" * 40
    return refs


def release_refs(*versions):
    return {f"refs/tags/{version}": "a" * 40 for version in versions}


class KubernetesSelectionTests(unittest.TestCase):
    def test_numeric_support_window_and_master_keep_exact_source_shas(self):
        refs = kubernetes_refs()
        refs["refs/heads/release-1.999-experimental"] = "b" * 40
        refs["refs/tags/release-1.999"] = "c" * 40
        self.assertEqual(resolver.kubernetes_matrix(refs), {"include": [
            {"k8s_branch": f"release-1.{minor}", "k8s_sha": f"{minor:040x}"}
            for minor in (98, 99, 100, 101)
        ] + [{"k8s_branch": "master", "k8s_sha": "f" * 40}]})

    def test_empty_or_unrelated_refs_fail(self):
        for refs in ({}, {"refs/heads/master": "f" * 40}):
            with self.subTest(refs=refs), self.assertRaisesRegex(ValueError, "release branches"):
                resolver.kubernetes_matrix(refs)

    def test_missing_minor_is_not_replaced_by_older_release(self):
        with self.assertRaisesRegex(ValueError, "release-1.99"):
            resolver.kubernetes_matrix(kubernetes_refs((97, 98, 100, 101)))

    def test_incomplete_window_fails(self):
        with self.assertRaisesRegex(ValueError, "support window"):
            resolver.kubernetes_matrix(kubernetes_refs((100, 101)))

    def test_missing_master_fails(self):
        refs = kubernetes_refs()
        del refs["refs/heads/master"]
        with self.assertRaisesRegex(ValueError, "master"):
            resolver.kubernetes_matrix(refs)


class ControllerSelectionTests(unittest.TestCase):
    def test_latest_patch_of_two_latest_minors_excludes_prereleases(self):
        refs = release_refs(
            "v0.9.99", "v0.10.9", "v0.10.10", "v0.11.2", "v0.11.1",
            "v0.12.0-alpha.1", "v1.0.0-rc.1", "v99.0", "v0.11.3^{}",
        )
        self.assertEqual(resolver.controller_refs(refs), ["develop", "v0.10.10", "v0.11.2"])

    def test_major_version_transition(self):
        refs = release_refs("v0.99.10", "v1.0.2", "v1.0.10", "v1.1.0")
        self.assertEqual(resolver.controller_refs(refs), ["develop", "v1.0.10", "v1.1.0"])

    def test_requires_two_stable_minors(self):
        for refs in ({}, release_refs("v0.10.0", "v0.10.1", "v0.11.0-rc.1")):
            with self.subTest(refs=refs), self.assertRaisesRegex(ValueError, "two kapp-controller"):
                resolver.controller_refs(refs)


class RemoteResolutionTests(unittest.TestCase):
    def test_keeps_full_sha(self):
        sha = "1234567890abcdef1234567890abcdef12345678"
        with patch.object(resolver.subprocess, "run", return_value=subprocess.CompletedProcess(
            [], 0, stdout=f"{sha}\trefs/heads/master\n",
        )):
            self.assertEqual(resolver.remote_refs("kubernetes/kubernetes"), {"refs/heads/master": sha})

    def test_empty_and_invalid_sha_fail(self):
        for output in ("", "1234567\trefs/heads/master\n", "z" * 40 + "\trefs/heads/master\n"):
            with self.subTest(output=output), patch.object(
                resolver.subprocess, "run", return_value=subprocess.CompletedProcess([], 0, stdout=output),
            ), self.assertRaises(ValueError):
                resolver.remote_refs("kubernetes/kubernetes")

    def test_git_failure_propagates(self):
        with patch.object(resolver.subprocess, "run", side_effect=subprocess.CalledProcessError(128, "git")):
            with self.assertRaises(subprocess.CalledProcessError):
                resolver.remote_refs("kubernetes/kubernetes")

    def test_main_outputs_machine_readable_values(self):
        output = io.StringIO()
        with patch.object(resolver, "remote_refs", side_effect=[
            kubernetes_refs(), {"refs/heads/main": "d" * 40}, release_refs("v0.10.0", "v0.11.0"),
        ]), contextlib.redirect_stdout(output):
            resolver.main()
        values = dict(line.split("=", 1) for line in output.getvalue().splitlines())
        self.assertEqual(json.loads(values["matrix"]), resolver.kubernetes_matrix(kubernetes_refs()))
        self.assertEqual(values["kind_sha"], "d" * 40)
        self.assertEqual(json.loads(values["controller_refs"]), ["develop", "v0.10.0", "v0.11.0"])

    def test_failed_resolution_publishes_no_partial_matrix(self):
        output = io.StringIO()
        with patch.object(resolver, "remote_refs", side_effect=[
            kubernetes_refs(), {"refs/heads/main": "d" * 40}, {},
        ]), contextlib.redirect_stdout(output), self.assertRaises(ValueError):
            resolver.main()
        self.assertEqual(output.getvalue(), "")

    def test_cli_reports_network_failure_and_exits_nonzero(self):
        output, errors = io.StringIO(), io.StringIO()
        with patch("subprocess.run", side_effect=subprocess.CalledProcessError(128, "git")), \
                contextlib.redirect_stdout(output), contextlib.redirect_stderr(errors), \
                self.assertRaises(SystemExit) as exit_status:
            runpy.run_path(str(SCRIPT), run_name="__main__")
        self.assertEqual(exit_status.exception.code, 1)
        self.assertEqual(output.getvalue(), "")
        self.assertIn("::error::Nightly version resolution failed", errors.getvalue())


if __name__ == "__main__":
    unittest.main()
