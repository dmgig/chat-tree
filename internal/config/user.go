package config

const (
	DefaultOutputDir = "output"
	BuildPrompt      = "" +
		"\nYou are helping to generate documentation for a software project.\n\n" +
		"You will receive code **one file at a time**.\n\n" +
		"For each file:\n" +
		"- Build on the documentation you've already written.\n" +
		"- You will receive your **previous output** along with the new file content.\n" +
		"- Use the previous documentation as a base and **add to or revise** it only if the new file provides relevant information.\n\n" +
		"⚠️ Do not speculate or guess. If a file does not provide enough information, simply return the **same documentation** without change.\n\n" +
		"Your goal is to produce **complete project documentation**, progressively updated as more code is received.\n\n" +
		"The final documentation should include:\n" +
		"- A list of available **highest level commands or features available to someone running the program in the cli,\n\n" +
		"we don't need to give info about under the hood commands** Even if there is no update to the deocumentation, you should respond again with the documentation. Your final response will be stored to be reviewed.\n" +
		"- **Sample usage**\n" +
		"- A **high-level explanation** of what the program does\n\n" +
		"Example:\n\n" +
		"# Project Name\n\n" +
		"## Overview\n" +
		"...\n\n" +
		"## Available Commands\n" +
		"- `command-name`: What it does\n" +
		"  - Example usage...\n\n" +
		"## Modules / Features\n" +
		"- Description of key modules\n" +
		"- Their roles / interactions\n"
	ReviewPrompt = "" +
		"You are now reviewing the final documentation for accuracy.\n\n" +
		"You will receive:\n" +
		"- The full documentation that was previously generated\n" +
		"- A list of files that were used\n\n" +
		"Review the documentation for accuracy, completeness, and clarity. If something seems wrong, unclear, or missing, revise it. Otherwise, you may return the same result.\n\n" +
		"You may request specific files by name if needed, and they will be provided in follow-up messages. However, try to work with the files listed unless essential."
	DefaultModel = "gpt-4-turbo"
)
