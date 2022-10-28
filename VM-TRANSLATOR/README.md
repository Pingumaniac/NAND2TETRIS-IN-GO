# VM-TRANSLATOR-IN-GO

## Summary
The Elements of Computing Systems: Building a Modern Computer from First Principles (nand2tetris) Chapter 7-8 Project completed using Go (1.18.4).

The aim of this project is to translate a program written in intermediate code which is designed to run on a Virtual Machine (VM) into Hack assembly language.

## About the author

Name: Youngjae Moon
* Bachelor of Science in Computer Science, and Applied Mathematics and Statistics at Stony Brook University (SUNY Korea).

## Instructions on checking out the latest stable version of this assembler

#### Method 1:
1. Open the main page for our GitHub repository: https://github.com/Pingumaniac/VM-Translator-IN-GO
2. Click the following button: <img src = "https://user-images.githubusercontent.com/63883314/115416097-69ade280-a232-11eb-8401-8c41362ab4c2.png" width="44" height="14">
3. Click 'Download ZIP' option.
4. Unzip the folder.

#### Method 2:
1.  Copy the web URL to your clipboard.
2.  Open 'Git Bash' from your local computer. You must have installed Git from the following page to do this: https://git-scm.com/downloads
3.  Move to the preferred directory.
4.  Next, type the following:
```
git clone
```

## Instructions on executing this software

### 1.Install Go

Please download Go from the following website: https://go.dev/dl/

It is recommended to download the latest version.

2. Open Terminal and then move to the corresponding folder. Next, enter the following to test each exemplar file for Chapter 7 and folder for Chapter 8:
```
go run VMTranslator.go ./vm_files_ch7/BasicTest.vm
go run VMTranslator.go ./vm_files_ch7/PointerTest.vm
go run VMTranslator.go ./vm_files_ch7/SimpleAdd.vm
go run VMTranslator.go ./vm_files_ch7/StackTest.vm
go run VMTranslator.go ./vm_files_ch7/StaticTest.vm
```
```
go run VMTranslator.go ./vm_files_ch8/BasicLoop
go run VMTranslator.go ./vm_files_ch8/FibonacciElement
go run VMTranslator.go ./vm_files_ch8/NestedCall
go run VMTranslator.go ./vm_files_ch8/SimpleFunction
go run VMTranslator.go ./vm_files_ch8/StaticsTest
```
## Bug tracking

* All users can view and report a bug in "GitHub Issues" of our repository. 
* Here is the URL for viewing and reporting a list of bugs: https://github.com/Pingumaniac/VM-Translator-IN-GO/issues
