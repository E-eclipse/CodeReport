# CodeReport

### CodeReport is a tool for generating reports from code.

![](result.png)
### Structure:
```text
│   go.mod
│   go.sum
│   main.go
│   result.png
|
├───docx
│       DocxGeneration.go
│       TableGeneration.go
│
├───interface
│       GenerationInterface.go
│       TableGenerationInterface.go
│
├───models
│       FileInfo.go
│
└───utils
        FileFunctions.go
        SupportFunctions.go
        welcome.go
```

### Ignore files

Create `.codereportignore` in the scanned project root to exclude files and folders from the report.

```gitignore
cache/
*.log
/root-only-folder/
!keep.log
```

CodeReport also ignores common dependency, cache, build, binary and media files by default, including `__pycache__`, `node_modules`, `.git`, `dist`, `build`, `*.pyc`, archives, images, fonts and executables.
