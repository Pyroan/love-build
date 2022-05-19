package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var doWinBuild bool = false
var doWebBuild bool = false
var clean bool = false

var inputDirectory string
var outputName string
var outputDirectory string
var version string

var lovePath string
var loveJSPath string

var verbose bool = false

func main() {
	var err error

	flag.BoolVar(&doWinBuild, "w", false, "create a windows build")
	flag.BoolVar(&doWebBuild, "b", false, "create an html5 build using love.js")
	flag.BoolVar(&clean, "clean", false, "delete .love file when finished")
	flag.StringVar(&outputDirectory, "d", "", "custom output directory (defaults to cwd, will create the directory if it doesn't exist.)")
	flag.StringVar(&outputName, "o", "", "custom output name (defaults to name of output directory)")
	flag.StringVar(&version, "version", "", "(optional) version name for this release")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")

	flag.Parse()

	inputDirectory = flag.Arg(0)

	if err = validateAndProcessArgs(); err != nil {
		fmt.Println("Argument Error: " + err.Error())
		return
	}
	fmt.Printf("doWinBuild: %v\ndoWebBuild: %v\nclean: %v\ninputDirectory: "+
		"%v\noutputName: %v\noutputDirectory: %v\nversion: %v\n", doWinBuild, doWebBuild, clean, inputDirectory, outputName, outputDirectory, version)
	fmt.Printf("lovePath: %v\nloveJSPath: %v\n", lovePath, loveJSPath)
	if err = generateLoveFile(); err != nil {
		fmt.Println("Error generating .love file: " + err.Error())
		return
	}
	if doWinBuild {
		if err = makeWinBuild(); err != nil {
			fmt.Println("Error generating Windows build: " + err.Error())
			return
		}
	}
	// if doWebBuild {
	// 	if err = makeWebBuild(); err != nil {
	// 		fmt.Println("Error generating Web build: " + err.Error())
	// 		return
	// 	}
	// }
	if clean {
		if err = cleanup(); err != nil {
			fmt.Println("Error during cleanup: " + err.Error())
			return
		}
	}
	fmt.Println("All done! Your builds are in " + outputDirectory)
	fmt.Println("Have a nice day :)")
}

func validateAndProcessArgs() error {
	var err error
	// first, make sure we're able to find the `love` executable and its related dlls/license(!)
	// TODO make this OS-independent...?
	lovePath, err = exec.LookPath("love")
	if err != nil {
		return errors.New("Couldn't find `love.exe`! Is it missing from $PATH?")
	}
	// LookPath() isn't guaranteed to return an absolute path for some reason.
	lovePath, err = filepath.Abs(lovePath)
	if err != nil {
		return err
	}
	// make sure Love's license can be found.
	fi, err := os.Stat(filepath.Join(filepath.Dir(lovePath), "license.txt"))
	if err != nil {
		return errors.New("Couldn't find LÖVE license!")
	}
	// if we're building for web, we need to make sure `love.js` is available too.
	if doWebBuild {
		loveJSPath, err = exec.LookPath("love.js")
		if err != nil {
			return errors.New("Couldn't find `love.js`, which is required for web builds. Is it missing from $PATH?")
		}
		loveJSPath, err = filepath.Abs(loveJSPath)
		if err != nil {
			return err
		}
	}

	// make sure inputDirectory is a real, extant directory.
	if inputDirectory == "" {
		return errors.New("Missing Input Directory!")
	}
	inputDirectory, err := filepath.Abs(inputDirectory)
	if err != nil {
		return err
	}
	fi, err = os.Stat(inputDirectory)
	if err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New(inputDirectory + " is not a directory!")
	}
	// also check to make sure inputDirectory at least contains a `main.lua`
	// (otherwise we know this can't be a Löve game)
	_, err = os.Stat(filepath.Join(inputDirectory, "main.lua"))
	if err != nil {
		return errors.New("Couldn't find main.lua in input directory!")
	}

	if outputName == "" {
		outputName = fi.Name()
	}
	if version != "" {
		outputName += "-" + version
	}

	if outputDirectory == "" {
		if outputDirectory, err = os.Getwd(); err != nil {
			return err
		}
	}
	fi, err = os.Stat(outputDirectory)
	if err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New(outputDirectory + " is not a directory!")
	}
	return nil
}

func generateLoveFile() error {
	loveFileName := getLoveFileName()
	archive, err := os.Create(loveFileName)
	if err != nil {
		return err
	}
	vPrint("Created " + loveFileName)
	vPrint("Copying files to archive...")
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	// TODO this doesn't handle symlinks
	err = filepath.Walk(inputDirectory, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip hidden (dot-prefixed) directories like .git
			if info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}

		// open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		// we want to strip the parent dir (inputDirectory) from the path, or else the .love file won't work.
		cutPath, err := filepath.Rel(inputDirectory, path)
		if err != nil {
			return err
		}
		// copy the file into our .love archive
		vPrint("\t" + cutPath)

		writer, err := zipWriter.Create(cutPath)
		if _, err = io.Copy(writer, file); err != nil {
			return err
		}

		file.Close()
		return nil
	})
	if err != nil {
		// delete the faulty archive.
		archive.Close()
		os.Remove(loveFileName)
		return err
	}
	vPrint("Done")
	return nil
}

func makeWebBuild() error {
	return nil
}

func makeWinBuild() error {
	dirPath := filepath.Join(outputDirectory, outputName+"_win")
	os.Mkdir(dirPath, 0755)
	gamePath := filepath.Join(dirPath, outputName+".exe")
	// grab love.exe and append our .love to it
	loveExecutable, err := os.Open(lovePath)
	if err != nil {
		return err
	}

	loveFile, err := os.Open(getLoveFileName())

	gameExec, err := os.Create(gamePath)
	if err != nil {
		return err
	}
	vPrint("Created " + gamePath)
	vPrint("Copying love.exe")
	_, err = io.Copy(gameExec, loveExecutable)
	if err != nil {
		return err
	}
	vPrint("Copying .love file")
	_, err = io.Copy(gameExec, loveFile)
	if err != nil {
		return err
	}

	loveExecutable.Close()
	loveFile.Close()
	gameExec.Close()
	vPrint("Finished generating " + outputName + ".exe")
	vPrint("Copying LÖVE License")
	license, err := os.Open(filepath.Join(filepath.Dir(lovePath), "license.txt"))
	if err != nil {
		return err
	}
	licenseCopy, err := os.Create(filepath.Join(dirPath, "license.txt"))
	if err != nil {
		return err
	}
	if _, err = io.Copy(licenseCopy, license); err != nil {

	}
	licenseCopy.Close()
	vPrint("Copying LÖVE `.dll`s")
	err = filepath.Walk(filepath.Dir(lovePath), func(path string, info os.FileInfo, err error) error {

		if filepath.Ext(path) == ".dll" {
			cutPath, err := filepath.Rel(filepath.Dir(lovePath), path)
			if err != nil {
				return err
			}
			vPrint("\t" + cutPath)
			copy, err := os.Create(filepath.Join(dirPath, filepath.Base(path)))
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err = io.Copy(copy, file); err != nil {
				return err
			}
			file.Close()
			copy.Close()
		}
		return nil
	})

	vPrint("Finished building for Windows")
	return nil
}

// Delete the `.love` file we created
func cleanup() error {
	err := os.Remove(getLoveFileName())
	if err != nil {
		return err
	}
	return nil
}

func vPrint(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}

func getLoveFileName() string {
	return filepath.Join(outputDirectory, outputName+".love")
}
