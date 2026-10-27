import tempfile
import unittest
from pathlib import Path

from prototype import install_no_replace, quote_sieve_string


class QuoteSieveStringTests(unittest.TestCase):
    def test_escapes_quotes_and_backslashes(self):
        self.assertEqual(
            r'"Hello \"world\" \\ mailbox"',
            quote_sieve_string('Hello "world" \\ mailbox'),
        )

    def test_rejects_newline(self):
        with self.assertRaisesRegex(ValueError, "control characters"):
            quote_sieve_string("line one\nline two")

    def test_allows_unicode(self):
        self.assertEqual('"Grüße 世界"', quote_sieve_string("Grüße 世界"))


class InstallNoReplaceTests(unittest.TestCase):
    def setUp(self):
        self.temporary_directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary_directory.cleanup)
        self.root = Path(self.temporary_directory.name)
        self.source = self.root / "source.sieve"
        self.target = self.root / "vacation.sieve"
        self.source.write_text('require ["vacation"];\n', encoding="utf-8")
        self.source.chmod(0o640)

    def test_installs_content_and_mode(self):
        install_no_replace(self.source, self.target)
        self.assertEqual(self.source.read_bytes(), self.target.read_bytes())
        self.assertEqual(0o640, self.target.stat().st_mode & 0o777)

    def test_does_not_replace_existing_target(self):
        self.target.write_text("keep me\n", encoding="utf-8")
        with self.assertRaises(FileExistsError):
            install_no_replace(self.source, self.target)
        self.assertEqual("keep me\n", self.target.read_text(encoding="utf-8"))

    def test_does_not_replace_symlink(self):
        other = self.root / "other"
        other.write_text("keep me\n", encoding="utf-8")
        self.target.symlink_to(other)
        with self.assertRaises(FileExistsError):
            install_no_replace(self.source, self.target)
        self.assertEqual("keep me\n", other.read_text(encoding="utf-8"))

    def test_missing_directory_fails(self):
        with self.assertRaises(FileNotFoundError):
            install_no_replace(self.source, self.root / "missing" / "target.sieve")


if __name__ == "__main__":
    unittest.main()
