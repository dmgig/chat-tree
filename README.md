# Chat Tree

## Overview
A software project for generating documentation through incremental code updates. This project involves managing structures, organizing documentation, and handling code updates efficiently. It interacts with external packages to handle token counts and model specifications for chat models. The project also includes functions for managing tokens and encoding for chat models, a review session function for enhancing the quality and accuracy of generated documentation, and a session package for creating and processing session files with OpenAI.

## Available Commands
- `openai`: Initiates prompt-based interactions with the OpenAI chat model.
  - Example usage: `chat-tree openai --prompt "your message here"`
- `document`: Generates documentation based on provided paths.
  - Example usage: `chat-tree document path/to/files --exclude pattern ...`
- `list`: Lists files based on provided paths and exclusion patterns.
  - Example usage: `chat-tree list path/to/files --exclude pattern ...`
- `list-models`: Lists available chat models with their maximum token limits.
  - Example usage: `chat-tree list-models`

## Modules / Features
- Incremental documentation generation
- Code update handling
- CLI interface for viewing documentation
- Structure management
- Organization of documentation
- File listing based on paths and exclusion patterns
- Retrieval and display of available chat models with their maximum token limits
- External module dependencies management
- Integration with OpenAI chat model for prompt-based interactions
- Functions for managing tokens and encoding for chat models
- Review session function to perform a review pass on the generated documentation
- Session package for creating and processing session files with OpenAI

## Sample Usage
- `chat-tree openai --prompt "your message here"`: Initiates a prompt-based interaction with the OpenAI chat model.
- `chat-tree document path/to/files --exclude pattern ...`: Generates documentation based on the specified paths while excluding certain patterns.
- `chat-tree list path/to/files --exclude pattern ...`: Lists files based on the specified paths while excluding certain patterns.
- `chat-tree list-models`: Lists available chat models with their maximum token limits.

## High-Level Explanation
The program automatically generates documentation for a software project through incremental code updates. It manages structures, organizes documentation, and handles code updates efficiently. Additionally, the program interacts with external packages to handle token counts and model specifications for chat models. The review session function allows for a detailed review of the generated documentation, enhancing the project's quality and accuracy. The session package, with the addition of the `Save` function, contributes by creating session directories, writing prompt/response files, processing them with OpenAI, and saving the responses. The latest addition of commands includes interacting with the OpenAI chat model, generating documentation, listing files, and displaying available chat models.