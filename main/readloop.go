package main

func ReadLoop(loopFullPath string) bool {
	var err error
	loopData, err = ReadCSV(loopFullPath, true)
	if err != nil {
		return false
	}

	if len(loopData.Header) == 0 {
		return false
	}

	if len(loopData.Data) == 0 {
		return false
	}

	return true
}
