package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type ByName []os.FileInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {

	result, err := recursiveFolderDown(path, "", printFiles)

	if err != nil {
		return fmt.Errorf("Recursive Function return error: " + err.Error())
	}

	out.WriteString(result)

	return nil
}

func recursiveFolderDown(path, prefix string, printFiles bool) (string, error) {
	var result bytes.Buffer

	file, err := os.Open(path)

	if err != nil {
		return "", fmt.Errorf("Folder not open! " + err.Error())
	}

	names, err := file.Readdir(0)

	if err != nil {
		return "", fmt.Errorf("Can't readdir!")
	}

	if !printFiles {
		names = deleteFilesFromArrayNames(names)
	}

	sort.Sort(ByName(names))

	for idx := 0; idx < len(names); idx++ {
		if names[idx].IsDir() {
			result.WriteString(concatPrefixAndName(idx, len(names), prefix, printFolderName(names[idx])))
			recursiveString, err := recursiveFolderDown(pathWithSeparator(file.Name(), names[idx].Name()), printFolderPrefix(idx, len(names), prefix), printFiles)
			if err == nil {
				result.WriteString(recursiveString)
			} else {
				return "", err
			}
		} else {
			result.WriteString(concatPrefixAndName(idx, len(names), prefix, printFileNameAndSize(names[idx])))
		}
	}

	return result.String(), nil
}

func deleteFilesFromArrayNames(names []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0, len(names))

	for idx := 0; idx < len(names); idx++ {
		if names[idx].IsDir() {
			result = append(result, names[idx])
		}
	}
	return result
}

func pathWithSeparator(folderName, file string) string {
	return (folderName + string(os.PathSeparator)) + file
}

func concatPrefixAndName(idx, size int, prefix, name string) string {
	return printElementPrefix(idx, size, prefix) + name
}

func printElementPrefix(idx, size int, prefix string) string {
	if size-idx != 1 {
		return prefix + "├───"
	} else {
		return prefix + "└───"
	}
}
func printFolderPrefix(idx, size int, prefix string) string {
	if size-idx != 1 {
		return prefix + "│	"
	} else {
		return prefix + "	"
	}
}

func printFolderName(file os.FileInfo) string {
	return file.Name() + "\n"
}

func printFileNameAndSize(file os.FileInfo) string {
	return (file.Name() + printAwesomeFileSize(file.Size())) + "\n"
}

func printAwesomeFileSize(size int64) string {
	var result bytes.Buffer

	if size > 0 {
		result.WriteString(" (")
		result.WriteString(strconv.FormatInt(size, 10))
		result.WriteString("b)")
	} else {
		result.WriteString(" (empty)")
	}

	return result.String()
}

func main() {
	stdOut := os.Stdout
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	stdOut.WriteString(out.String())
}
