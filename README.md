# Data Generator with Progress Bar

This Go-based program generates mock data for a database and exports it as CSV files. It supports multiple tables with parent-child relationships, simulating data generation with progress updates in the terminal.

## Features
- Generates mock data for tables and their child tables.
- Progress bar displayed in the terminal showing the status of data generation.
- Supports configurable table settings via `/config/*.json`.
- Allows users to choose different configuration files and output folders.

## Prerequisites

- Go 1.16+ (Install Go from [here](https://golang.org/dl/))
- Git (Install Git from [here](https://git-scm.com/))

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yrmsa/data-generator-go.git
   cd data-generator-go
   ```

2. Initialize the Go module (if not already initialized):
   ```bash
   go mod init data-generator-go
   ```

3. Install dependencies:
   ```bash
   go get github.com/schollz/progressbar/v3
   ```

## Configuration

The program uses a `config.json` file to configure the data generation process. You can create your own `config.json` file or use one of the example configuration files.

### Example `config.json`:
```json
{
  "tables": [
    {
      "name": "cust",
      "rows": 10,
      "columns": [
        {
          "name": "cust_id",
          "generator": { "increment": true }
        },
        {
          "name": "name",
          "generator": {
            "random": {
              "prefix": "user_",
              "length": 8
            }
          }
        }
      ]
    },
    {
      "name": "cust_personal",
      "parent": "cust",
      "rows_per_parent": 1,
      "columns": [
        {
          "name": "cust_personal_id",
          "generator": { "table_increment": true }
        },
        {
          "name": "cust_id",
          "generator": { "parent_key": "cust_id" }
        },
        {
          "name": "gender",
          "generator": { "hardcoded": "L" }
        },
        {
          "name": "nickname",
          "generator": {
            "random": { 
              "prefix": "nickname_",
              "length": 7
            }
          }
        }
      ]
    },
    {
      "name": "cust_addr",
      "parent": "cust_personal",
      "rows_per_parent": 2,
      "columns": [
        {
          "name": "cust_addr_id",
          "generator": { "table_increment": true }
        },
        {
          "name": "cust_personal_id",
          "generator": { "parent_key": "cust_personal_id" }
        },
        {
          "name": "address",
          "generator": {
            "random": {
              "suffix": " Street",
              "length": 12
            }
          }
        }
      ]
    },
    {
      "name": "cust_attr",
      "parent": "cust_personal",
      "rows_per_parent": 10,
      "columns": [
        {
          "name": "cust_attrl_id",
          "generator": { "table_increment": true }
        },
        {
          "name": "cust_personal_id",
          "generator": { "parent_key": "cust_personal_id" }
        },
        {
          "name": "attr_key",
          "generator": { 
            "predefined_list": ["preference", "setting", "option", "config"]
          }
        },
        {
          "name": "attr_value",
          "generator": {
            "random": { "length": 6 }
          }
        }
      ]
    }
  ]
}
```

## Running the Program

1. Compile and run the program:
   ```bash
   go run main.go
   ```

2. The program will ask for the configuration file and folder path to save the CSV files.

### Example Output
```plaintext
Generating data for table: users (1000 rows)
Table: users [==========>              ] 300/1000
Table: users [=======================>] 1000/1000
Completed data generation for table: users
```

## Build Process

### Windows Systems
```batch
build.bat
```

### Linux/macOS Systems
```bash
# Set execute permissions
chmod +x build.sh

# Run build script
./build.sh
```

### Build Output Structure
```text
dist/
├── data-generator-arm64        # arm64 executable
├── data-generator-linux        # Linux executable
└── data-generator-windows.exe  # Windows executable
```


## Contributing

1. Fork the repository.
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/data-generator.git
   ```
3. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Make your changes and commit:
   ```bash
   git commit -m "Add new feature"
   ```
5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
6. Create a pull request from your fork to the main repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
