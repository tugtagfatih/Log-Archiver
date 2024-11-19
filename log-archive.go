package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("-----------------TUGTAGFATIH'S LOG ARCHIVER-----------------")

	logDir := flag.String("logdir", "", "Path to the log directory to archive.")
	flag.Parse()

	if *logDir == "" {
		fmt.Println("Usage: log-archive -logdir=<log-directory> \nSince you have not entered a log directory, the \"/var/log\" directory will be used by default")
		*logDir = "/var/log"
	}

	// Timestamp for archive
	timestamp := time.Now().Format("20060102_150405")
	archiveName := fmt.Sprintf("logs_archive_%s.tar.gz", timestamp)
	archivePath := filepath.Join(*logDir, "archived_logs")

	//Create archive Path
	err := os.MkdirAll(archivePath, 0755)
	if err != nil {
		fmt.Printf("Error: Archive directory could not be created: %v\n", err)
		return
	}

	// Full path of achive
	archiveFullPath := filepath.Join(archivePath, archiveName)
	err = compressLogs(*logDir, archiveFullPath, archiveName)

	if err != nil {
		noPermErr := fmt.Sprintf("open %s/%s: permission denied", archivePath, archiveName)
		errstr := fmt.Sprintf("%s", err)
		if errstr == noPermErr {
			fmt.Println("Required permissions were not found, please run as root so the files can be created")
		} else {
			fmt.Printf("Error: Logs could not be compressed: %v\n", err)
			return
		}

	}

	//Logging process details
	logFilePath := filepath.Join(*logDir, "archive_log.txt")
	logEntry := fmt.Sprintf("%s: Archived logs to %s\n", timestamp, archiveName)
	err = appendToFile(logFilePath, logEntry)
	if err != nil {
		fmt.Printf("Error: Could not write to log file: %v\n", err)
		return
	}
	fmt.Printf("Archiving completed: %s\n", archiveFullPath)
}

func compressLogs(logDir, archivePath, archiveName string) error {
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Exclude the archive file and archive_log.txt
		if info.IsDir() || filepath.Base(path) == archiveName || filepath.Base(path) == "archive_log.txt" {
			return nil
		}

		return addFileToTar(tarWriter, path, logDir)
	})

	return err
}

func addFileToTar(tw *tar.Writer, filePath, baseDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return err
	}
	header.Name = relPath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

func appendToFile(filePath, text string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}
