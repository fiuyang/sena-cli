package main

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

// Fungsi untuk mengenkripsi password
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func processRow(row *xlsx.Row, headerRow *xlsx.Row, targetColumn string) {
	cells := row.Cells
	headerCells := headerRow.Cells

	// Iterasi melalui sel-sel dalam baris data
	for i, cell := range cells {
		columnName := headerCells[i].String() // Ambil nama kolom dari baris header
		if columnName == targetColumn {
			value := cell.String()
			hashedPassword := hashPassword(value)
			mutex.Lock()
			cell.SetString(hashedPassword)
			mutex.Unlock()
		}
	}
}

func main() {
	startTime := time.Now()
	var excelFileName string
	var targetColumn string
	reader := bufio.NewReader(os.Stdin)

	lines := []string{
		"\033[34m   _____                    _    _           _     \n  / ____|                  | |  | |         | |    \n | (___   ___ _ __   __ _  | |__| | __ _ ___| |__  \n  \\___ \\ / _ \\ '_ \\ / _` | |  __  |/ _` / __| '_ \\ \n  ____) |  __/ | | | (_| | | |  | | (_| \\__ \\ | | |\n |_____/ \\___|_| |_|\\__,_| |_|  |_|\\__,_|___/_| |_|\n                                                   \n                                                   \033[0m \n" + "\x1B[38;2;66;211;146mA\x1B[39m \x1B[38;2;67;209;149mC\x1B[39m\x1B[38;2;68;206;152mL\x1B[39m\x1B[38;2;69;204;155mI\x1B[39m \x1B[38;2;70;201;158mt\x1B[39m\x1B[38;2;71;199;162mo\x1B[39m\x1B[38;2;72;196;165mo\x1B[39m\x1B[38;2;73;194;168ml\x1B[39m \x1B[38;2;74;192;171mf\x1B[39m\x1B[38;2;75;189;174mo\x1B[39m\x1B[38;2;76;187;177mr\x1B[39m \x1B[38;2;77;184;180mh\x1B[39m\x1B[38;2;78;182;183ma\x1B[39m\x1B[38;2;79;179;186ms\x1B[39m\x1B[38;2;80;177;190mh\x1B[39m\x1B[38;2;81;175;193mi\x1B[39m\x1B[38;2;82;172;196mn\x1B[39m\x1B[38;2;83;170;199mg\x1B[39m \x1B[38;2;84;165;205mc\x1B[39m\x1B[38;2;85;162;208mo\x1B[39m\x1B[38;2;86;160;211ml\x1B[39m\x1B[38;2;87;158;215mu\x1B[39m\x1B[38;2;88;155;218mm\x1B[39m\x1B[38;2;89;153;221mn\x1B[39m \x1B[38;2;90;150;224me\x1B[39m\x1B[38;2;91;148;227mx\x1B[39m\x1B[38;2;92;145;230mc\x1B[39m\x1B[38;2;93;143;233me\x1B[39m\x1B[38;2;94;141;236ml\x1B[39m\x1B[38;2;95;138;239m \x1B[39m",
	}

	for _, line := range lines {
		fmt.Println(line)
	}

	for {
		fmt.Print("Enter Excel file path: ")
		excelFileName, _ = reader.ReadString('\n')
		excelFileName = strings.TrimSpace(excelFileName)

		excelFileName = strings.Trim(excelFileName, "\"")

		if excelFileName == "" {
			fmt.Println("Excel file path cannot be empty")
			continue
		}
		break
	}

	for {
		fmt.Print("Enter column name to hash: ")
		targetColumn, _ = reader.ReadString('\n')
		targetColumn = strings.TrimSpace(targetColumn)

		if targetColumn == "" {
			fmt.Println("Column name cannot be empty")
			continue
		}
		break
	}

	// Baca file Excel
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %v\n", err)
	}

	// Dapatkan sheet pertama dari file Excel
	sheet := xlFile.Sheets[0]

	// Ambil baris pertama (header)
	headerRow := sheet.Rows[0]

	// Buat worker pool
	workerCount := 4
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for _, row := range sheet.Rows[1:] { // Mulai dari baris kedua untuk menghindari header
				if row != nil {
					processRow(row, headerRow, targetColumn)
				}
			}
		}(i)
	}

	wg.Wait()

	outputFileName := excelFileName[:len(excelFileName)-5] + "_hashed.xlsx"

	err = xlFile.Save(outputFileName)
	elapsedTime := time.Since(startTime)
	fmt.Printf("elapsed time %v\n", elapsedTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hash column successfully.")
}
