package cmd

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var contents = `# This code is free software; you can redistribute it and/or modify it under
# the terms of the new BSD License.
#
# Copyright (c) 2010, Sebastian Staudt

# A nano configuration file to enable syntax highlighting of some Git specific
# files with the GNU nano text editor (http://www.nano-editor.org)

# This syntax format is used for editing commit and tag messages
syntax "git commit/tag messages" "COMMIT_EDITMSG|TAG_EDITMSG"

# overlong lines
color brightred ".{72}(.*)"
color white "^.{72}"

# comment
color blue "^#.*$"

# JIRA ticket id
color yellow "[A-Z]+-[0-9]+"

# special comment lines
color green      "^# Changes to be committed:"
color red        "^# Changes not staged for commit:"
color brightblue "^# Untracked files:"
color brightblue "^# On branch .+$"
color brightblue "^# Your branch is ahead of .+$"
# No need to highlight this.
# color brightblue  "Your branch is up-to-date with .+$"
color brightred   "^# Your branch and '[^']+' have diverged"
color white       "#[[:space:]](deleted|modified|new file|renamed):[[:space:]].*"
color red         "#[[:space:]]deleted:"
color green       "#[[:space:]]modified:"
color brightgreen "#[[:space:]]new file:"
color cyan        "#[[:space:]]renamed:"

# Recolor hash symbols
color blue "#"

# This syntax format is used for interactive rebasing
syntax "git rebase todo" "git-rebase-todo"

# Default
color white ".*"

# Comments
color blue "^#.*"

# Commit IDs
color brightwhite "[0-9a-f]{7,40}"

# Rebase commands
color green       "^(e|edit) [0-9a-f]{7,40}"
color green       "^#  (e, edit)"
color brightgreen "^(f|fixup) [0-9a-f]{7,40}"
color brightgreen "^#  (f, fixup)"
color brightwhite "^(p|pick) [0-9a-f]{7,40}"
color brightwhite "^#  (p, pick)"
color cyan        "^(r|reword) [0-9a-f]{7,40}"
color cyan        "^#  (r, reword)"
color brightred   "^(s|squash) [0-9a-f]{7,40}"
color brightred   "^#  (s, squash)"
color yellow      "^(x|exec) [^ ]+ [0-9a-f]{7,40}"
color yellow      "^#  (x <cmd>, exec <cmd>)"

# Recolor hash symbols
color blue "#"
`

// gitForNanoCmd represents the gitForNano command
var gitForNanoCmd = &cobra.Command{
	Use:   "git-nano-highlight",
	Short: "Installs git commit message syntax highlighting for nano",
	Run: func(cmd *cobra.Command, args []string) {
		f, _ := cmd.Flags().GetBool("force")
		s := spinner.New(spinner.CharSets[11], time.Second / 10)
		s.Suffix = " Installing highlighting file"
		s.Start();
		file, err := installHighlighting(f)
		if err != nil {
			s.Stop()
			rootCmd.PrintErrln(err)
			os.Exit(1)
		}

		s.Suffix = " Updating .nanorc file"
		s.Restart()

		err = installNanorc(file)
		if err != nil {
			s.Stop()
			rootCmd.PrintErrln(err)
		}

		s.Stop()
		rootCmd.Println("Successfully installed git highlighting for nano")
	},
}

// Installs a reference to the newly created file in the
// local nanorc file
func installNanorc(rf string) error {
	file, err := nanorcFile()

	if err != nil {
		return err
	}

	exists, err := isPathExists(file)
	if err != nil {
		return err
	}

	// Output array
	var out []string

	// Remove occurrences of the include line if they exist
	if exists {
		fc, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// Any occurrences
		re := regexp.MustCompile(`(?mi)include.*\.gitcommit\.nanorc.*`)

		for _, l := range strings.Split(string(fc), "\n") {
			res := re.FindAllString(l, -1)
			if len(res) == 0 {
				out = append(out, l)
			}
		}
	}

	incl := fmt.Sprintf(`include "%s" # automatically added by dev`, rf)

	last := len(out)
	if (len(out) > 0) && (out[last - 1] == "") {
		out[last - 1] = incl
	} else {
		out = append(out, incl)
	}
	out = append(out, "")


	err = ioutil.WriteFile(file, []byte(strings.Join(out, "\n")), 755)
	return err
}

func highlightingFile() (string, error) {
	var path string

	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	file := ".gitcommit.nanorc"
	path = filepath.Join(home, file)

	return path, nil
}


func nanorcFile() (string, error) {
	var path string

	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	file := ".nanorc"
	path = filepath.Join(home, file)

	return path, nil
}

// Installs the syntax highlighting definition file,
// returning the path to the file
func installHighlighting(f bool) (string, error) {
	file, err := highlightingFile()
	if err != nil {
		return file, err
	}

	exists, err := isPathExists(file)
	if err != nil {
		return file, err
	}
	if exists && !f {
		return file, errors.New("file already exists")
	}

	err = ioutil.WriteFile(file, []byte(contents), 755)
	if err != nil {
		return file, err
	}

	return file, nil
}

func init() {
	rootCmd.AddCommand(gitForNanoCmd)

	gitForNanoCmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
}
