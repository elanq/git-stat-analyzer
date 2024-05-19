package helper_test

import (
	"testing"

	"github.com/elanq/git-stat-analyzer/helper"
)

func TestParseAuthor(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Author: devopswizards <devops@ajaib.co.id>",
			expected: "devops@ajaib.co.id",
		},
		{
			input: "Author: Elan Qisthi <elan.aji@ajaib.co.id>",

			expected: "elan.aji@ajaib.co.id",
		},
		{
			input: "commit 569eec0afc4f844ce5be48b8b18c145352583e3e",

			expected: "",
		},
		{
			input:    "Date:   Thu Feb 22 12:07:44 2024 +0700",
			expected: "",
		},
	}

	for _, c := range cases {
		actual := helper.ParseAuthor(c.input)
		if actual != c.expected {
			t.Log("Test ERROR. expected: ", c.expected, " actual:", actual)
			t.Fail()
		}
	}
}

func TestParseDate(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Author: devopswizards <devops@ajaib.co.id>",
			expected: "Mon Jan 1 00:00:00 0001 +0000",
		},
		{
			input: "Author: Elan Qisthi <elan.aji@ajaib.co.id>",

			expected: "Mon Jan 1 00:00:00 0001 +0000",
		},
		{
			input: "commit 569eec0afc4f844ce5be48b8b18c145352583e3e",

			expected: "Mon Jan 1 00:00:00 0001 +0000",
		},
		{
			input:    "Date:   Thu Feb 1 12:07:44 2024 +0700",
			expected: "Thu Feb 1 12:07:44 2024 +0700",
		},
	}

	for _, c := range cases {
		actual := helper.ParseDate(c.input)
		actualString := actual.Format(helper.TimeFormat)
		if actualString != c.expected {
			t.Log("Test ERROR. expected: ", c.expected, " actual:", actualString)
			t.Fail()
		}
	}
}

func TestParseCommit(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Author: devopswizards <devops@ajaib.co.id>",
			expected: "",
		},
		{
			input: "Author: Elan Qisthi <elan.aji@ajaib.co.id>",

			expected: "",
		},
		{
			input: "commit 569eec0afc4f844ce5be48b8b18c145352583e3e",

			expected: "569eec0afc4f844ce5be48b8b18c145352583e3e",
		},
		{
			input:    "Date:   Thu Feb 22 12:07:44 2024 +0700",
			expected: "",
		},
	}

	for _, c := range cases {
		actual := helper.ParseCommit(c.input)
		if actual != c.expected {
			t.Log("Test ERROR. expected: ", c.expected, " actual:", actual)
			t.Fail()
		}
	}
}

func TestParseCommitMessage(t *testing.T) {
	cases := []struct {
		input               string
		expectedFileChanges []string
		expectedAddedLine   int
		expectedRemovedLine int
	}{
		{
			input: ` 3 files changed, 29 insertions(+), 29 deletions(-)
	 pom.xml                |  4 ++--
	 odt-web-app/pom.xml    |  8 ++++----
	 odt-service/pom.xml    | 18 +++++++++---------
			`,
			expectedFileChanges: []string{
				"pom.xml",
				"odt-web-app/pom.xml",
				"odt-service/pom.xml",
			},
			expectedAddedLine:   29,
			expectedRemovedLine: 29,
		},
		{
			input: ` 3 files changed, 29 deletions(-)
	 pom.xml                |  4 ++--
	 odt-web-app/pom.xml    |  8 ++++----
	 odt-service/pom.xml    | 18 +++++++++---------
			`,
			expectedFileChanges: []string{
				"pom.xml",
				"odt-web-app/pom.xml",
				"odt-service/pom.xml",
			},
			expectedAddedLine:   0,
			expectedRemovedLine: 29,
		},
		{
			input: ` 3 files changed, 29 insertions(+)
	 pom.xml                |  4 ++--
	 odt-web-app/pom.xml    |  8 ++++----
	 odt-service/pom.xml    | 18 +++++++++---------
			`,
			expectedFileChanges: []string{
				"pom.xml",
				"odt-web-app/pom.xml",
				"odt-service/pom.xml",
			},
			expectedAddedLine:   29,
			expectedRemovedLine: 0,
		},
	}

	for _, c := range cases {
		actualFileChanges, actualAddedLine, actualRemovedLine := helper.ParseCommitMessage(c.input)
		if len(actualFileChanges) != len(c.expectedFileChanges) {
			t.Log("Testfile changes ERROR. expected: ", c.expectedFileChanges, " actual:", actualFileChanges)
			t.FailNow()
		}
		for i, f := range c.expectedFileChanges {
			if f != actualFileChanges[i] {
				t.Log("Testfile changes ERROR. expected: ", f, " actual:", actualFileChanges[i])
				t.FailNow()
			}
		}
		if actualAddedLine != c.expectedAddedLine {
			t.Log("Test added line ERROR. expected: ", c.expectedAddedLine, " actual:", actualAddedLine)
			t.FailNow()
		}
		if actualRemovedLine != c.expectedRemovedLine {
			t.Log("Test removed line ERROR. expected: ", c.expectedRemovedLine, " actual:", c.expectedRemovedLine)
			t.FailNow()
		}
	}
}
