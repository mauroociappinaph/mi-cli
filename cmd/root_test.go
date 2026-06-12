package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCommand() *cobra.Command {
	// Reset output on root and subcommands to avoid state leakage
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
	greetCmd.SetOut(nil)
	greetCmd.SetErr(nil)
	versionCmd.SetOut(nil)
	versionCmd.SetErr(nil)
	return rootCmd
}

func TestRootCmd_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "mi-cli")
	assert.Contains(t, output, "greet")
	assert.Contains(t, output, "version")
}

func TestRootCmd_VersionFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "commit")
	assert.Contains(t, output, "built")
}

func TestGreetCmd_DefaultName(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"greet"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "¡Hola, mundo! 👋")
}

func TestGreetCmd_CustomName(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"greet", "Mauro"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "¡Hola, Mauro! 👋")
}

func TestGreetCmd_JSONOutput(t *testing.T) {
	viper.Set("output", "json")
	defer viper.Set("output", "text")

	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"greet", "Test"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := strings.TrimSpace(buf.String())
	assert.Contains(t, output, `{"greeting": "¡Hola, Test! 👋"}`)
	assert.True(t, strings.HasSuffix(output, "}"))
	assert.JSONEq(t, `{"greeting": "¡Hola, Test! 👋"}`, output)
}

func TestGreetCmd_InvalidArgs(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestCommand()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"greet", "a", "b", "c"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts at most 1 arg")
}

func TestRootCmd_PersistentFlags(t *testing.T) {
	viper.Reset()
	// Re-bind flags after reset
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	cmd := newTestCommand()
	cmd.SetArgs([]string{"--verbose", "version"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	require.NoError(t, err)

	assert.True(t, viper.GetBool("verbose"), "verbose flag should be bound to viper")
}

func TestInitConfig_WithConfigFile(t *testing.T) {
	t.Skip("Requiere archivo de config temporal")
}