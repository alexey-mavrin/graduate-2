package integration_test

import (
	"fmt"
	"os"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func storeAndGetAccount() {
	storeAndGet(common.AccountRecord,
		"-l http://localhost -u us1 -p pass1",
	)
}

func storeAndGetNote() {
	storeAndGet(common.NoteRecord,
		"-t text1",
	)
}

func storeAndGetCard() {
	storeAndGet(common.CardRecord,
		"-ch CardHolder -num 1111222233334444 -em 12 -ey 2027 -c 123",
	)
}

func storeAndDeleteAccount() {
	storeAndDelete(common.AccountRecord,
		"-l http://localhost -u us1 -p pass1",
	)
}

func storeAndDeleteNote() {
	storeAndDelete(common.NoteRecord,
		"-t text1",
	)
}

func storeAndDeleteCard() {
	storeAndDelete(common.CardRecord,
		"-ch CardHolder -num 1111222233334444 -em 12 -ey 2027 -c 123",
	)
}

func storeAndUpdateAccount() {
	storeAndUpdate(common.AccountRecord,
		"-l http://localhost -u us1 -p pass1",
	)
}

func storeAndUpdateNote() {
	storeAndUpdate(common.NoteRecord,
		"-t text1",
	)
}

func storeAndUpdateCard() {
	storeAndUpdate(common.CardRecord,
		"-ch CardHolder -num 1111222233334444 -em 12 -ey 2027 -c 123",
	)
}

func storeAndGet(rt common.RecordType,
	data string,
) {
	By(fmt.Sprintf("Running 'client %s get'", rt))
	_, _, err := runClient("user -a register")
	Expect(err).NotTo(HaveOccurred(), "Client should register")

	stdOut, stdErr, err := runClient(
		fmt.Sprintf("%s -a store -n name_1 %s", rt, data),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should store %s", rt),
	)
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a get -i 1", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should retrieve %s by ID", rt),
	)
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a get -n name_1", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should retrieve %s by name", rt),
	)

}

func storeAndUpdate(rt common.RecordType,
	data string,
) {
	By(fmt.Sprintf("Running 'client %s update'", rt))
	_, _, err := runClient("user -a register")
	Expect(err).NotTo(HaveOccurred(), "Client should register")

	// store initial record
	stdOut, stdErr, err := runClient(
		fmt.Sprintf("%s -a store -n name_1 %s", rt, data),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should store %s", rt),
	)

	// rename (update) record using ID
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a update -i 1 -n name_01", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should rename %s by ID", rt),
	)

	// update record using new name
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a update -n name_01 -m some_metainfo", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should update %s by name", rt),
	)

	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a get -n name_01", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(stdOut).To(ContainSubstring("some_metainfo"))
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should retrieve %s by name", rt),
	)

}

func storeAndDelete(rt common.RecordType,
	data string,
) {
	By(fmt.Sprintf("Running 'client %s delete'", rt))
	_, _, err := runClient("user -a register")
	Expect(err).NotTo(HaveOccurred(), "Client should register")

	// store record 1
	stdOut, stdErr, err := runClient(
		fmt.Sprintf("%s -a store -n name_1 %s", rt, data),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should store %s", rt),
	)

	// store record 2
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a store -n name_2 %s", rt, data),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should store %s", rt),
	)

	// delete record by ID
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a delete -n name_1", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should delete %s by name", rt),
	)

	// delete record by name
	stdOut, stdErr, err = runClient(
		fmt.Sprintf("%s -a delete -n name_2", rt),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		fmt.Sprintf("Client should delete %s by name", rt),
	)
}

func storeAndGetBinary() {
	By("Running 'client bin get'")
	_, _, err := runClient("user -a register")
	Expect(err).NotTo(HaveOccurred(), "Client should register")

	file, err := makeBinFile(256)
	Expect(err).NotTo(HaveOccurred(), "Binary is to be created")
	stdOut, stdErr, err := runClient(
		fmt.Sprintf("bin -a store -n bin_name_1 -f %s", file),
	)
	fmt.Print(stdOut, stdErr)

	fileOut := file + ".out"
	Expect(err).NotTo(HaveOccurred(),
		"Client should save binary data",
	)

	h1, err := getHash(file)
	Expect(err).NotTo(HaveOccurred(), "Get original file hash")

	stdOut, stdErr, err = runClient(
		fmt.Sprintf("bin -a get -i 1 -f %s", fileOut),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		"Client should get binary data by ID",
	)
	h2, err := getHash(fileOut)
	Expect(err).NotTo(HaveOccurred(), "Get saved file hash")
	Expect(h1).To(Equal(h2), "Files should be the same")
	os.Remove(fileOut)

	stdOut, stdErr, err = runClient(
		fmt.Sprintf("bin -a get -n bin_name_1 -f %s", fileOut),
	)
	fmt.Print(stdOut, stdErr)
	Expect(err).NotTo(HaveOccurred(),
		"Client should get binary data by name",
	)

	h2, err = getHash(fileOut)
	Expect(err).NotTo(HaveOccurred(), "Get saved file hash")
	Expect(h1).To(Equal(h2), "Files should be the same")
	os.Remove(file)
	os.Remove(fileOut)
}
