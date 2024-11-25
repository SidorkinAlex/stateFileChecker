# **State File Checker**
## **Overview**
State File Checker is a Go-based command-line tool designed to verify the consistency of files within a specified directory. It checks for changes in file hashes, identifies new files, and detects missing or deleted files. This tool is particularly useful for ensuring the integrity of application states and managing file versions.
## **Features**
- **Hash Comparison**: Computes and compares SHA-1 hashes of files to detect changes.
- **File Status Reporting**: Generates a report of new, changed, and deleted files.
- **Ignore List**: Supports an ignore list to skip specified files or directories during the check.
- **Custom Commands**: Allows users to specify commands to execute upon success or failure of the consistency check.
  ## **Installation**
  To install the State File Checker, clone the repository and build the project:

  ```bash

  git clone https://github.com/SidorkinAlex/stateFileChecker.git

  cd stateFileChecker

  make build
  ```
  ## **Usage**
  Run the tool from the command line with the required parameters:

  bash

  ./stateFileChecker -s <source\_directory>
  ### **Parameters**
- -s: Specifies the source directory to check.
- --success--run: Command to execute if the consistency check is successful.
- --failed--run: Command to execute if the consistency check fails.
  ### **Example**
  bash

  ./stateFileChecker -s /path/to/your/directory --success--run "echo Success!" --failed--run "echo Failure!"
  ## **Output**
  The results of the consistency check will be logged to a file named hash\_diff\_output.txt, which includes:

- **New Files**: List of files that were added.
- **Changed Files**: List of files that have been modified.
- **Deleted Files**: List of files that are missing.
  ## **Dependencies**
  This project relies on the following Go packages:

- flag: For command-line argument parsing.
- crypto/sha1: For generating file hashes.
- encoding/csv: For reading CSV files containing hash data.
- encoding/json: For parsing JSON manifest files.
  ## **Contributing**
  Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.
  ## **License**
  This project is licensed under the MIT License. See the LICENSE file for more details.

  This README provides a comprehensive overview of the State File Checker, including its purpose, features, installation instructions, usage examples, and contribution guidelines.